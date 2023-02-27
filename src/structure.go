package main

import (
    "fmt"
    "log"
	"github.com/alecthomas/participle/v2/lexer"
)

type Document struct {
	// Entries are the basic components of the file
	// Like paragraphs of text, tables, images, blocks, etc.
	// They are separated by two EOL tokens as specified in
	// The Entry struct tags
	Entries []*Entry `@@+`
}

type Entry struct {
	// Position saved for better error messaging past the lexing and parsing
	// stages
	Pos lexer.Position

	Image     *Image     `( @@`
	Table     *Table     `| @@`
	Box       *Box       `| @@`
	List      *List      `| @@`
	Paragraph *Paragraph `| @@ ) EOL (EOL | EOF)`
}

// Stub function for union types to satisfy the ParserUnionType interface
func (entry Entry) Union() {}

func (entry Entry) Serve() ([]byte, error){
    entry_type, err := unionType(&entry)
    if err != nil {
        return nil, fmt.Errorf("Serving entry: Error: %s at %s", err, entry.Pos.String())
    }

    switch entry_type {
        case "*main.Paragraph":
            log.Printf("Paragraph found at %s\n", entry.Pos.String())
            return entry.Paragraph.Serve()
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

