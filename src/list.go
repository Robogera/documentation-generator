package main

type List struct {
	// TODO: add unordered lists
	// (would require some syntax change to distinguish them)
	Reference  string       `"!list":Command Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (list List) Serve() ([]byte, error) {
    paragraphs := make([][]byte, len(list.Paragraphs))
    for i, paragraph := range list.Paragraphs {
        contents, err := paragraph.Serve()
        if err != nil {
            continue
        }
        paragraphs[i] = contents
    }
    return serve(paragraphs, `<ol>{{ range . }}<li>{{ printf "%s" . }}</li>{{ end }}</ol>`)
}

