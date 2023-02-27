package main

type Table struct {
	// TODO: add colors
	Reference string `"!table":Command Whitespace @Ident`
	Title     *Text  `Whitespace @@`
	Rows      []*Row `(EOL @@)+`
}

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths     []*Path    `"img":Ident (Whitespace @@)`
	Paragraph *Paragraph `EOL "txt":Ident Whitespace @@`
}

type Path struct {
	// Used by multiple file elements
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}

func (table Table) Serve() ([]byte, error) {
    return nil, nil
}

