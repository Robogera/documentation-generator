package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

// generic serve() that returns a byte slice of template string executed against the interface{}
// use with caution inside the elements' methods to make sure
// the correct data struct/template combo is used
func serve(data_struct any, template_string string) ([]byte, error) {

	buf := new(bytes.Buffer)

	// the name in New() doesn't really matter unless we do a ParseFiles()
	// which will error if name doesn't match the file basename (??)
	tmpl, err := template.New("generic-template").Parse(template_string)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(buf, data_struct)
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

// ================================================================
// Entry-specific Serve() methods that should be called recursively
// until the whole thing is constructed into a byte slice
// ================================================================

func (entry Entry) Serve() ([]byte, error) {

	entry_type, err := unionType(&entry)

	if err != nil {
		return nil, fmt.Errorf("Serving entry: Error: %s at %s", err, entry.Pos.String())
	}

	switch entry_type {
	case "*main.Paragraph":
		log.Printf("Paragraph found at %s\n", entry.Pos.String())
		return entry.Paragraph.Serve()
	case "*main.Header":
		log.Printf("Header found at %s\n", entry.Pos.String())
		return entry.Header.Serve()
	case "*main.List":
		log.Printf("List found at %s\n", entry.Pos.String())
		return entry.List.Serve()
	case "*main.Image":
		log.Printf("Image found at %s\n", entry.Pos.String())
		return entry.Image.Serve()
	case "*main.Table":
		log.Printf("Table found at %s\n", entry.Pos.String())
		return entry.Table.Serve()
	case "*main.Box":
		log.Printf("Box found at %s\n", entry.Pos.String())
		return entry.Box.Serve()
	default:
		return nil, fmt.Errorf("Serving entry: Entry type %s at %s not defined", entry_type, entry.Pos.String())
	}
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

	processed_title, err := serve(table.Title, `{{ .Text }}`)
	if err != nil {
		return nil, err
	}

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

func (paragraph Paragraph) Serve() ([]byte, error) {
	processed_elements := make([][]byte, len(paragraph.Element))
	for i, elem := range paragraph.Element {
		contents, err := elem.Serve()
		if err != nil {
			return nil, err
		}
		processed_elements[i] = contents
	}

	tmpl, err := template.New("paragraph").Parse(
		"<p>{{ range . }}{{ printf \"%s\" . }}{{ end }}</p>",
	)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, processed_elements)
	if err != nil {
		return nil, err
	}

	result := make([]byte, buf.Len())

	_, err = buf.Read(result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (elem ParagraphElement) Serve() ([]byte, error) {
	elem_type, err := unionType(&elem)
	if err != nil {
		return nil, fmt.Errorf("Serving paragraph element: Error: %s at %s", err, elem.Pos.String())
	}

	switch elem_type {
	case "*main.Link":
		return serve(elem.Link, `<a href="{{ .Url }}">{{ .Text.Text }}</a>`)
	case "*main.Bold":
		return serve(elem.Bold, `<b>{{ .Text.Text }}</b>`)
	case "*main.Text":
		return serve(elem.Text, `{{ .Text }}`)
	default:
		return nil, fmt.Errorf("Serving paragraph element: Element type %s at %s not defined", elem_type, elem.Pos.String())
	}
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

	processed_data := struct {
		Reference  string
		Path       []byte
		Paragraphs [][]byte
	}{
		Reference:  image.Reference,
		Path:       processed_path,
		Paragraphs: processed_paragraphs,
	}

	return serve(processed_data, "<figure name={{ .Reference }}>\n<img src=\"{{ printf \"%s\" .Path }}\">\n<figcaption>{{ range .Paragraphs }}{{ printf \"%s\" . }}{{ end }}</figcaption>\n</figure>")
}

func (path Path) Serve() ([]byte, error) {
	return serve(path.Path, "{{ . }}")
}

func (header Header) Serve() ([]byte, error) {
	processed_text, err := serve(header.Text, `{{ .Text }}`)
	if err != nil {
		return nil, err
	}

	// <h1> is already reserved by a page title so we construct h2 and upwards
	processed_level := 1 + len(header.Level)

	processed_data := struct {
		Level int
		Text  []byte
	}{
		Level: processed_level,
		Text:  processed_text,
	}
	return serve(processed_data, `<h{{ printf "%d" .Level }}>{{ printf "%s" .Text }}</h{{ printf "%d" .Level }}>`)
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

    return serve(processed_data, "<div class=\"wrapper\"><div name=\"{{ .Reference }}\" class=\"{{ .Type }}\">\n{{ range .Paragraphs }}{{ printf \"%s\" . }}{{ end }}\n</div></div>")
}
