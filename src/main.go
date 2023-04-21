package main

import (
	"bytes"
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

	input_files, err := os.ReadDir(filepath.Join("input"))
	if err != nil {
		log.Fatalf("Opening source files failed with error: %s", err)
	}

	buf := new(bytes.Buffer)

	for _, input_file := range input_files {

		if input_file.IsDir() {
			log.Printf("Skipping a subdir %s in the input dir\n", input_file.Name())
			continue
		} else {
			log.Printf("Reading from file %s\n", input_file.Name())
		}

		contents, err := os.ReadFile(filepath.Join("input", input_file.Name()))
		if err != nil {
			log.Printf("Reading file %s failed with error %s, skipping...\n", input_file.Name(), err)
			continue
		}
		buf.Write(contents)
	}

	combined_file_contents := make([]byte, buf.Len())

	_, err = buf.Read(combined_file_contents)
	if err != nil {
		log.Fatalf("Reading from combined input file buffer failed with error %s\n", err)
	}

	document, err := document_parser.ParseBytes("", combined_file_contents)
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

	scripts, err := os.ReadDir(filepath.Join(".", "js"))
	if err != nil {
		log.Fatalf("Opening js folder error: %s\n", err)
	}

	processed_html, err := newProcessedHtml(len(stylesheets), len(scripts))
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
		processed_html.Styles[i] = contents
	}

	for i, script := range scripts {
		contents, err := os.ReadFile(filepath.Join("js", script.Name()))
		if err != nil {
			log.Printf("Error: %s. Stylesheet %s invalid. Skipping...\n", err, script.Name())
			continue
		}
		processed_html.Scripts[i] = contents
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

	err = os.Mkdir(filepath.Join("output", "img"), 0755)
	if err != nil {
		log.Printf("Creating directory failed with error %s. Attempting to copy files anyway...\n", err)
	}

	for _, v := range images.dump() {
		input, err := os.ReadFile(filepath.Join("input", "img", v))
		if err != nil {
			log.Printf("Reading image %s failed with error %s, skipping...\n", v, err)
			continue
		}
		err = os.WriteFile(filepath.Join("output", "img", v), input, 0755)
		if err != nil {
			log.Printf("Writing to file %s failed with error %s, skipping...\n", v, err)
			continue
		}
	}
}

type processedHtml struct {
	Scripts [][]byte
	Styles  [][]byte
	TOC     []byte
	Body    []byte
}

func newProcessedHtml(styles_count, scripts_count int) (*processedHtml, error) {
	if styles_count < 1 {
		return nil, fmt.Errorf("Can't have less than one css")
	}
	return &processedHtml{
		Styles:  make([][]byte, styles_count),
		Scripts: make([][]byte, scripts_count),
	}, nil
}
