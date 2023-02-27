package main

type Box struct {
	Type       string       `@( "!info":Command | "!warn":Command )`
	Reference  string       `Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

func (box Box) Serve() ([]byte, error) {
    return nil, nil
}
