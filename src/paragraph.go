package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"text/template"
)

type Paragraph struct {
	Element []*ParagraphElement `@@+`
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

type ParagraphElement struct {
	Pos lexer.Position
	// Normal text just gets lumped together into big chunks
	// stored in "NormalText" string
	// And the special elements like bold text or urls are stored separately
	Link *Link `( @@`
	Bold *Bold `| @@`
	Text *Text `| @@ )`
}

type Link struct {
	Text *Text `"[" @@ "]"`
	// Urls can only be local references for now
	Url string `"(" @Ident ")"`
}

type Bold struct {
	Text *Text `"*":Special @@ "*":Special`
}

type Text struct {
	Text string `@( Ident | RusWord | Number | Whitespace | OpenParen | CloseParen | Punct )+`
}

// Stub function for union types to satisfy the ParserUnionType interface
func (elem ParagraphElement) Union() {}

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

func serve(data_struct any, template_string string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("new").Parse(template_string)
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
