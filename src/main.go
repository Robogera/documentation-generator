package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"log"
	"os"
)

func main() {

	document_lexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Command", `![a-z]+`},
		{"Color", `#[0-9a-fA-F]{6}`},
		{"Number", `[0-9]+`},
		{"Ident", `[a-zA-Z][a-zA-Z0-9_-]*`},
		{"RusWord", `[а-яА-ЯёЁ]+`},
		{"OpenParen", `[\(\[]{1}`},
		{"CloseParen", `[\)\]]{1}`},
		{"EOL", `[\n\r]{1}`},
		{"Whitespace", `[ ]+`},
		{"Punct", `[!\?.,;:\-']+`},
		{"Special", `[\*\\/]`},
	})

	document_parser := participle.MustBuild[FILE](
		participle.Lexer(document_lexer),
	)

	file_contents, err := os.ReadFile("input/syntax.test")
	if err != nil {
		log.Fatalf("Init: Error: File could not be read: %s\n", err)
	}

	syntax_tree, err := document_parser.ParseBytes("", file_contents)
	if err != nil {
		log.Fatalf("Parser: Error: File could not be parsed: %s\n", err)
	}

	for _, entry := range syntax_tree.Entries {
        entry_type, err := entry.Type()
        if err != nil {
            log.Printf("Reader: invalid Entry: %s", err)
        }
        switch entry_type {
        case "*main.Paragraph":
            log.Printf("Paragraph found at %s\n", entry.Pos.String())
        case "*main.List":
            log.Printf("List found at %s\n", entry.Pos.String())
        case "*main.Image":
            log.Printf("Image found at %s\n", entry.Pos.String())
        case "*main.Table":
            log.Printf("Table found at %s\n", entry.Pos.String())
        case "*main.Box":
            log.Printf("Box found at %s\n", entry.Pos.String())
        }
	}
}
