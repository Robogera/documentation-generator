package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"log"
	"os"
)

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

	var SMLParser *participle.Parser[FILE] = participle.MustBuild[FILE](
		participle.Lexer(SMLLexer),
		participle.CaseInsensitive("Indent"),
	)

	var err error
	var ast *FILE

	ast, err = parseSML(SMLParser, "input/syntax.test")
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	log.Printf("Err %+v", ast)

	for _, entry := range ast.Entries {
		log.Printf("%+v\n", entry)
	}
}
