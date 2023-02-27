package main

type List struct {
	// TODO: add unordered lists
	// (would require some syntax change to distinguish them)
	Reference  string       `"!list":Command Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (list List) Serve() ([]byte, error) {
    return nil, nil
}

