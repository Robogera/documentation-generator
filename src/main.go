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
	Infobox   *Infobox   ` | @@`
	Warnbox   *Warnbox   ` | @@`
	List      *List      ` | @@`
	Image     *Image     ` | @@ ) EOL (EOL | EOF)`
}

type Paragraph struct {
    // A "hardcore" regexp filter that enforces correct grammars on the parsing stage
    // was overkill and produced obscure error messages that are hard to interpret
    // TODO: check for basic mistakes like double whitespace, no punctuation, no capital letter, etc
    // after the parsing stage
	// Text string `@( Ident ( Punct? Whitespace OpenParen? Ident CloseParen? )* Punct )`
	Text string `@( Ident | Whitespace | OpenParen | CloseParen | Punct )+`
}

type Image struct {
    Path        string     `"!image":Command Whitespace @(Ident "." Ident) EOL`
	Description *Paragraph `@@`
}

type Infobox struct {
	Reference string     `"!info":Command Whitespace @Ident EOL`
	Text      *Paragraph `@@`
}

type Warnbox struct {
	Reference string     `"!warn":Command Whitespace @Ident EOL`
	Text      *Paragraph `@@`
}

type List struct {
	Reference string       `"!list":Command Whitespace @Ident`
	Entries   []*Paragraph `(EOL @@)*`
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
        // TODO: not sure if having a separate lexer token for command is a great idea
        // Maybe find a way to distinguish basic paragraphs from command-prefixed structures
        // in the parsing stage instead?
		{"Command", `![a-z]+`},
		{"Color", `#[0-9a-fA-F]{6}`},
		{"Ident", `[a-zA-Zа-яА-Я][a-zA-Zа-яА-Я'_]*`},
        {"OpenParen", `[\(\[]{1}`},
        {"CloseParen", `[\)\]]{1}`},
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

	ast, err = parseSML(SMLParser, "input/syntax.test")
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	log.Printf("Err %+v", ast)

	for _, entry := range ast.Entries {
		log.Printf("%+v\n", entry)
	}
}
