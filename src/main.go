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
        {"OpenSquare", `\[`},
        {"CloseSquare", `\]`},
		{"Ident", `[a-zA-Z][a-zA-Z0-9_-]*`},
		{"RusWord", `[а-яА-ЯёЁ]+`},
		{"OpenParen", `\(`},
		{"CloseParen", `\)`},
		{"EOL", `[\n\r]{1}`},
		{"Whitespace", `[ ]+`},
		{"Punct", `[!\?.,;:\-']+`},
		{"Special", `[\*\\/]`},
	})

	document_parser := participle.MustBuild[Document](
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
        contents, err := entry.Serve()
        if err != nil {
            log.Printf("Entry invalid, skipping: %s\n", err)
            continue
        }
        log.Printf("%s", string(contents))
	}
}
