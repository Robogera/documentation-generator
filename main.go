package main

import (
    "log"
    "os"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type FILE struct {
	// Properties []*Property `@@*`
    Entries []*Entry `(( @@ ) Newline Newline)*`
}

type Entry struct {
    Pos lexer.Position

    Command string ` "!" @Ident`
    Paragraph string `| @Ident`
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
        {"Ident", `[a-zA-Z_ ]+`},
        {"Command", `!`},
        {"Terminator", `;`},
        {"Newline", `\n`},
    })

    var SMLParser *participle.Parser[FILE] = participle.MustBuild[FILE](
        participle.Lexer(SMLLexer),
        participle.CaseInsensitive("Indent"),
    )

	var err error
	var ast *FILE

	ast, err = parseSML(SMLParser, "test.sml")
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	log.Printf("Err %+v", ast)

    for _, entry := range ast.Entries {
        log.Printf("%+v\n", entry)
    }
}

