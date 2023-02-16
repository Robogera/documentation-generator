package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"log"
	"os"
)

type FILE struct {
	Entries []*Entry `@@*`
}

type Entry struct {
	Pos lexer.Position

    Paragraph *Paragraph ` ( @@`
    Image     *Image     ` | @@ ) EOL (EOL | EOF)`
}

type Paragraph struct {
	Text string `@( Ident ( Punct? Whitespace Ident )* Punct )`
}

type Image struct {
    Path string `"!image":Command Whitespace @Filename EOL`
    Description *Paragraph `@@`
}

func parseSML(parser *participle.Parser[FILE], filename string) (*FILE, error) {
	var err error

	var ast *FILE
	var contents []byte

	contents, err = os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ast, err = parser.ParseBytes("", contents)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func main() {

	var SMLLexer *lexer.StatefulDefinition = lexer.MustSimple([]lexer.SimpleRule{
        {"Filename", `[a-zA-Z0-9-_]+\.[a-zA-Z]+`},
        {"Command", `![a-z]+`},
        {"Color", `#[0-9a-fA-F]{6}`},
		{"Ident", `[a-zA-Z][a-zA-Z_]*`},
		{"EOL", `[\n\r]{1}`},
		{"Whitespace", `[ ]+`},
		{"Punct", `[!\?.,;:]+`},
	})

	var SMLParser *participle.Parser[FILE] = participle.MustBuild[FILE](
		participle.Lexer(SMLLexer),
		participle.CaseInsensitive("Indent"),
		participle.CaseInsensitive("Line"),
	)

	var err error
	var ast *FILE

	ast, err = parseSML(SMLParser, "syntax.test")
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	log.Printf("Err %+v", ast)

	for _, entry := range ast.Entries {
		log.Printf("%+v\n", entry)
	}
}
