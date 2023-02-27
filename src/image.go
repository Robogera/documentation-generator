package main

type Image struct {
	Reference  string       `"!img":Command Whitespace @Ident`
	Path       *Path        `Whitespace @@`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (image Image) Serve() ([]byte, error) {
    return nil, nil
}
