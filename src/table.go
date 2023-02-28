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

	processed_rows := make([]struct {
		Paths      [][]byte
		Paragraphs [][]byte
	}, len(table.Rows))

	for i, row := range table.Rows {
		processed_rows[i].Paths = make([][]byte, len(row.Paths))
		processed_rows[i].Paragraphs = make([][]byte, len(row.Paragraphs))
		for j, path := range row.Paths {
			content, err := path.Serve()
			if err != nil {
				return nil, err
			}
			processed_rows[i].Paths[j] = content
		}
		for j, paragraph := range row.Paragraphs {
			content, err := paragraph.Serve()
			if err != nil {
				return nil, err
			}
			processed_rows[i].Paragraphs[j] = content
		}
	}

	processed_title, err := serve(table.Title.Text, `{{ .Text }}`)

	processed_data := struct {
		Title     []byte
		Reference string
		Rows      []struct {
			Paths      [][]byte
			Paragraphs [][]byte
		}
	}{
		Title:     processed_title,
		Reference: table.Reference,
		Rows:      processed_rows,
	}

	tmpl, err := template.New("table.html").ParseFiles("templates/table.html")
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

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths      []*Path      `"img":Ident (Whitespace @@)`
	Paragraphs []*Paragraph `(EOL "txt":Ident Whitespace @@)+`
}
