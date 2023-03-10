package main

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/yosssi/gohtml"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func main() {

	document_lexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Command", `![a-z]+`},
		{"Number", `[0-9]+`},
		{"OpenSquare", `\[`},
		{"CloseSquare", `\]`},
		{"Ident", `[a-zA-Z][a-zA-Z0-9_-]*`},
		{"RusWord", `[а-яА-ЯёЁ]+`},
		{"OpenParen", `\(`},
		{"CloseParen", `\)`},
		{"EOL", `[\n\r]{1}`},
		{"Whitespace", `[ ]+`},
		{"Punct", `[!\?.,;:\-+%'"—«»]+`},
		{"Special", `[\*\\/#<>&=]`},
		{"Code", "`"},
	})

	document_parser := participle.MustBuild[Document](
		participle.Lexer(document_lexer),
	)

	file_contents, err := os.ReadFile("input/syntax.test")
	if err != nil {
		log.Fatalf("Init: Error: File could not be read: %s\n", err)
	}

	document, err := document_parser.ParseBytes("", file_contents)
	if err != nil {
		log.Fatalf("Parser: Error: File could not be parsed: %s\n", err)
	}

	headers := new(HeaderStorage)
	images := newImageStorage()
	ids := newIdStorage()
	links := newLinkStorage()

	stylesheets, err := os.ReadDir(filepath.Join(".", "css"))
	if err != nil {
		log.Fatalf("Opening css folder error: %s\n", err)
	}

	processed_html, err := newProcessedHtml(len(stylesheets))
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	processed_html.Body, err = document.Serve(headers, images, ids, links)
	if err != nil {
		log.Fatalf("Serving body failed with error %s\n", err)
	}

	for i, stylesheet := range stylesheets {
		contents, err := os.ReadFile(filepath.Join("css", stylesheet.Name()))
		if err != nil {
			log.Printf("Error: %s. Stylesheet %s invalid. Skipping...\n", err, stylesheet.Name())
			continue
		}
		processed_html.Style[i] = contents
	}

	processed_html.TOC, err = headers.generateTOC()
	if err != nil {
		log.Fatalf("Error generating TOC: %s\n", err)
	}

	tmpl, err := template.New("main.html").ParseFiles("templates/main.html")
	if err != nil {
		log.Fatalf("Template reading error: %s\n", err)
	}

	dest, err := os.Create("output/index.html")
	if err != nil {
		log.Fatalf("Error creating file: %s\n", err)
	}
	defer dest.Close()

	err = tmpl.Execute(gohtml.NewWriter(dest), processed_html)
	if err != nil {
		log.Fatalf("Error writing to ouput: %s\n", err)
	}
}

type processedHtml struct {
	Style [][]byte
	TOC   []byte
	Body  []byte
}

func newProcessedHtml(styles_count int) (*processedHtml, error) {
	if styles_count < 1 {
		return nil, fmt.Errorf("Can't have less than one css")
	}
	return &processedHtml{
		Style: make([][]byte, styles_count),
	}, nil
}
