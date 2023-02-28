package main

import (
    "bytes"
    "text/template"
)

type Table struct {
	// TODO: add colors
	Reference string `"!table":Command Whitespace @Ident`
	Title     *Text  `Whitespace @@`
	Rows      []*Row `(EOL @@)+`
}

func (table Table) Serve() ([]byte, error) {

    processed_rows := make([]struct{
        Paths [][]byte
        Paragraph []byte
    }, len(table.Rows))
    
    for i, row := range table.Rows {
        processed_rows[i].Paths = make([][]byte, len(row.Paths))
        for j, path := range row.Paths {
            content, err := path.Serve()
            if err != nil {
                return nil, err
            }
            processed_rows[i].Paths[j] = content
        }
        content, err := row.Paragraph.Serve()
        if err != nil {
            return nil, err
        }
        processed_rows[i].Paragraph = content
    }

    tmpl, err := template.New("table").Parse(
        ""
    )
    return nil, nil
}

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths     []*Path    `"img":Ident (Whitespace @@)`
	Paragraph *Paragraph `EOL "txt":Ident Whitespace @@`
}

