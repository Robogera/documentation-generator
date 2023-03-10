package main

import (
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

	processed_body, err := document.Serve(headers, images, ids, links)
	if err != nil {
		log.Fatalf("Serving body failed with error %s\n", err)
	}

	stylesheets, err := os.ReadDir("./css")

	processed_css := make([][]byte, len(stylesheets))

	for i, stylesheet := range stylesheets {
		contents, err := os.ReadFile(filepath.Join("./css", stylesheet.Name()))
		if err != nil {
			log.Printf("Error: %s. Stylesheet %s invalid. Skipping...\n", err, stylesheet.Name())
			continue
		}
		processed_css[i] = contents
	}

	processed_toc, err := headers.generateTOC()
	if err != nil {
		log.Fatalf("Error generating TOC: %s\n", err)
	}

	processed_data := struct {
		Style [][]byte
		TOC   []byte
		Body  []byte
	}{
		Style: processed_css,
		TOC:   processed_toc,
		Body:  processed_body,
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

	err = tmpl.Execute(gohtml.NewWriter(dest), processed_data)
	if err != nil {
		log.Fatalf("Error writing to ouput: %s\n", err)
	}
}
