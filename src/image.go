package main

import (
    "bytes"
    "text/template"
)

type Image struct {
	Reference  string       `"!img":Command Whitespace @Ident`
	Path       *Path        `Whitespace @@`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (image Image) Serve() ([]byte, error) {

    // Underlying paragraphs are processes by the paragraph.Serve() method and
    // returned as a slice to avoid repetition
    processed_paragraphs := make([][]byte, len(image.Paragraphs))

    for i, paragraph := range image.Paragraphs {
        contents, err := paragraph.Serve()
        if err != nil {
            return nil, err
        }
        processed_paragraphs[i] = contents
    }

    processed_path, err := image.Path.Serve()
    if err != nil {
        return nil, err
    }

    processed_data := struct{
        Reference string
        Path []byte
        Paragraphs [][]byte
    }{
        Reference: image.Reference,
        Path: processed_path,
        Paragraphs: processed_paragraphs,
    }

    tmpl, err := template.New("image").Parse(
        "<figure name={{ .Reference }}>\n<img src=\"{{ printf \"%s\" .Path }}\">\n<figcaption>{{ range .Paragraphs }}<p>{{ printf \"%s\" . }}</p>{{ end }}</figcaption>\n/figure>")
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

type Path struct {
	// Used by multiple file elements
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}

func (path Path) Serve() ([]byte, error) {
    return serve(path.Path, "{{ . }}")
}

