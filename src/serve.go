package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
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

func (document Document) Serve(headers *HeaderStorage, images *ImageStorage, ids *IdStorage, links *LinkStorage) ([]byte, error) {

	processed_entries := make([][]byte, len(document.Entries))

	for i, entry := range document.Entries {

		// passing all the storages to
		contents, err := entry.Serve(headers, images, ids, links)
		if err != nil {
			return nil, fmt.Errorf("Entry at %s invalid, skipping: %s\n", entry.Pos.String(), err)
		}

		processed_entries[i] = contents
	}

	return serve(processed_entries,
		`{{ range . }}{{ printf "%s" . }}
        {{ end }}
    `)
}

func (entry Entry) Serve(headers *HeaderStorage, images *ImageStorage, ids *IdStorage, links *LinkStorage) ([]byte, error) {

	entry_type, err := unionType(&entry)
	if err != nil {
		return nil, fmt.Errorf("Serving entry: Error: %s at %s", err, entry.Pos.String())
	}

	switch entry_type {
	case "*main.Paragraph":
		log.Printf("Paragraph found at %s\n", entry.Pos.String())
		return entry.Paragraph.Serve(links)
	case "*main.Header":
		log.Printf("Header found at %s\n", entry.Pos.String())
		return entry.Header.Serve(headers, ids)
	case "*main.List":
		log.Printf("List found at %s\n", entry.Pos.String())
		return entry.List.Serve(links)
	case "*main.Image":
		log.Printf("Image found at %s\n", entry.Pos.String())
		return entry.Image.Serve(links, images, ids)
	case "*main.Table":
		log.Printf("Table found at %s\n", entry.Pos.String())
		return entry.Table.Serve(images, ids, links)
	case "*main.Box":
		log.Printf("Box found at %s\n", entry.Pos.String())
		return entry.Box.Serve(links)
	default:
		return nil, fmt.Errorf("Serving entry: Entry type %s at %s not defined", entry_type, entry.Pos.String())
	}
}

func (table Table) Serve(images *ImageStorage, ids *IdStorage, links *LinkStorage) ([]byte, error) {

	err := ids.push(table.ID)
	if err != nil {
		return nil, fmt.Errorf("Table serving error: %s", err)
	}

	processed_rows := make([]struct {
		Paths      [][]byte
		Paragraphs [][]byte
	}, len(table.Rows))

	for i, row := range table.Rows {

		processed_rows[i].Paths = make([][]byte, len(row.Paths))

		for j, path := range row.Paths {
			content, err := path.Serve(images)
			if err != nil {
				return nil, err
			}
			// TODO: check path validity
			processed_rows[i].Paths[j] = content
		}

		processed_rows[i].Paragraphs = make([][]byte, len(row.Paragraphs))

		for j, paragraph := range row.Paragraphs {
			content, err := paragraph.Serve(links)
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
		Title []byte
		ID    string
		Rows  []struct {
			Paths      [][]byte
			Paragraphs [][]byte
		}
	}{
		Title: processed_title,
		ID:    table.ID,
		Rows:  processed_rows,
	}

	return serve(processed_data, `
    <div class="wrapper">
        <table class="buttons" id="{{ .ID }}">
             <colgroup>
                <col style="width:30%;">
                <col style="width:70%;">
            </colgroup>
            <thead>
                <tr>
                    <th colspan="2">{{ printf "%s" .Title }}</th>
                </tr>
                <tr>
                    <th>Элемент</th>
                    <th>Функция</th>
                </tr>
            </thead>
            <tbody>
            {{ range .Rows }}<tr>
                <th>
                    <div class="grid">{{ range .Paths }}
                        <div {{ if false }}class="wide-element"{{ end }}>
                            <img src="{{ printf "%s" . }}" class="grid-image">
                        </div>
                    {{ end }}
                    </div>
                </th>
                <th>{{ range .Paragraphs }}<p>{{ printf "%s" . }}</p>{{ end }}</th>
            </tr>{{ end }}
            </tbody>
        </table>
    </div>
    `)
}

func (paragraph Paragraph) Serve(links *LinkStorage) ([]byte, error) {

	processed_elements := make([][]byte, len(paragraph.Element))

	for i, elem := range paragraph.Element {
		contents, err := elem.Serve(links)
		if err != nil {
			return nil, err
		}
		processed_elements[i] = contents
	}

	return serve(processed_elements, "<p>{{ range . }}{{ printf \"%s\" . }}{{ end }}</p>")
}

func (elem ParagraphElement) Serve(links *LinkStorage) ([]byte, error) {
	elem_type, err := unionType(&elem)
	if err != nil {
		return nil, fmt.Errorf("Serving paragraph element: Error: %s at %s", err, elem.Pos.String())
	}

	switch elem_type {
	case "*main.Link":
		links.push(elem.Link.Url)
		return serve(elem.Link, `<a href="#{{ .Url }}">{{ .Text.Text }}</a>`)
	case "*main.Bold":
		return serve(elem.Bold, `<b>{{ .Text.Text }}</b>`)
	case "*main.Text":
		return serve(elem.Text, `{{ .Text }}`)
	default:
		return nil, fmt.Errorf("Serving paragraph element: Element type %s at %s not defined", elem_type, elem.Pos.String())
	}
}

func (list List) Serve(links *LinkStorage) ([]byte, error) {

	paragraphs := make([][]byte, len(list.Paragraphs))

	for i, paragraph := range list.Paragraphs {
		contents, err := paragraph.Serve(links)
		if err != nil {
			continue
		}
		paragraphs[i] = contents
	}

	return serve(paragraphs, `
    <ol>
    {{ range . }}<li class="listelement">{{ printf "%s" . }}</li>
    {{ end }}</ol>
    `)
}

func (image Image) Serve(links *LinkStorage, images *ImageStorage, ids *IdStorage) ([]byte, error) {

	err := ids.push(image.ID)
	if err != nil {
		return nil, fmt.Errorf("Serving image error: %s", err)
	}
	// Underlying paragraphs are processes by the paragraph.Serve() method and
	// returned as a slice to avoid repetition
	processed_paragraphs := make([][]byte, len(image.Paragraphs))

	for i, paragraph := range image.Paragraphs {
		contents, err := paragraph.Serve(links)
		if err != nil {
			return nil, err
		}
		processed_paragraphs[i] = contents
	}

	processed_path, err := image.Path.Serve(images)
	if err != nil {
		return nil, err
	}

	processed_data := struct {
		ID         string
		Path       []byte
		Paragraphs [][]byte
	}{
		ID:         image.ID,
		Path:       processed_path,
		Paragraphs: processed_paragraphs,
	}

	return serve(processed_data, `
    <figure name={{ .ID }}>
        <img class="grid-image" src="{{ printf "%s" .Path }}"></img>
        <figcaption>
            {{ range .Paragraphs }}{{ printf "%s" . }}
            {{ end }}
        </figcaption>
    </figure>
    `)
}

func (path Path) Serve(images *ImageStorage) ([]byte, error) {
	images.push(path.Path)
	return serve(path.Path, "img/{{ . }}")
}

func (header Header) Serve(headers *HeaderStorage, ids *IdStorage) ([]byte, error) {

	processed_text, err := serve(header.Text, `{{ .Text }}`)
	if err != nil {
		return nil, err
	}
	processed_id := header.ID

	if processed_id == "" {
		processed_id = uuid.New().String()
	}

	headers.push(len(header.Level), processed_text, processed_id)
	err = ids.push(processed_id)
	if err != nil {
		return nil, err
	}

	// <h1> is already reserved by a page title so we construct h2 and upwards
	processed_level := 1 + len(header.Level)

	processed_data := struct {
		Level int
		Text  []byte
		ID    string
	}{
		Level: processed_level,
		Text:  processed_text,
		ID:    processed_id,
	}

	return serve(processed_data, `
    <h{{ printf "%d" .Level }} id="{{ .ID }}">{{ printf "%s" .Text }}</h{{ printf "%d" .Level }}>
    `)
}

func (box Box) Serve(links *LinkStorage) ([]byte, error) {
	processed_paragraphs := make([][]byte, len(box.Paragraphs))

	for i, paragraph := range box.Paragraphs {
		contents, err := paragraph.Serve(links)
		if err != nil {
			return nil, err
		}
		processed_paragraphs[i] = contents
	}

	processed_data := struct {
		Type       string
		ID         string
		Paragraphs [][]byte
	}{
		Type:       box.Type[1:],
		ID:         box.ID,
		Paragraphs: processed_paragraphs,
	}

	return serve(processed_data, `
    <div class="wrapper">
        <div class="{{ .Type }}" id="{{ .ID }}">
            {{ range .Paragraphs }}{{ printf "%s" . }}
            {{ end }}
        </div>
    </div>
    `)
}
