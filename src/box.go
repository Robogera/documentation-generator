package main

import (
    "text/template"
    "bytes"
)

type Box struct {
	Type       string       `@( "!info":Command | "!warn":Command )`
	Reference  string       `Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (box Box) Serve() ([]byte, error) {
    processed_paragraphs := make([][]byte, len(box.Paragraphs))

    for i, paragraph := range box.Paragraphs {
        contents, err := paragraph.Serve()
        if err != nil {
            return nil, err
        }
        processed_paragraphs[i] = contents
    }

    processed_data := struct{
        Type string
        Reference string
        Paragraphs [][]byte
    }{
        Type: box.Type[1:],
        Reference: box.Reference,
        Paragraphs: processed_paragraphs,
    }

    tmpl, err := template.New("").Parse(
        "<div name=\"{{ .Reference }}\" class=\"{{ .Type }}\">\n{{ range .Paragraphs }}{{ printf \"%s\" . }}{{ end }}\n</div>",
    )
    if err != nil {
        return nil, err
    }

    buf := new(bytes.Buffer)

    err = tmpl.Execute(buf, processed_data)
    if err != nil {
        return nil, err
    }

    result := make([]byte, buf.Len())

    _, err = buf.Read(result)
    if err != nil {
        return nil, err
    }

    return result, nil
}
