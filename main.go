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
	Image     *Image     ` | @@`
	Table     *Table     ` | @@`
	List      *List      ` | @@ ) EOL`
}

type Paragraph struct {
	Text string `@( Ident ( WordSeparator Ident )* PunctuationMark )`
}

type Image struct {
	Path string `ImageDeclaration @( Ident "." ("png" | "jpg" | "jpeg") )`
}

type Table struct {
	Name string `TableDeclaration @Ident`
}

type List struct {
	Name string `ListDeclaration @Ident`
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
		{"Ident", `[a-zA-Z_]+`},
		{"EOL", `[\n\r]+`},
		{"TableDeclaration", `!Table +`},
		{"ListDeclaration", `!List +`},
		{"ImageDeclaration", `!Image +`},
		{"WordSeparator", `[!\?.,;:]?[ ]`},
		{"PunctuationMark", `[!\?.,;:]+`},
	})

	var SMLParser *participle.Parser[FILE] = participle.MustBuild[FILE](
		participle.Lexer(SMLLexer),
		participle.CaseInsensitive("Indent"),
		participle.CaseInsensitive("Line"),
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
