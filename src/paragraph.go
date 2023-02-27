package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Paragraph struct {
	Element []*ParagraphElement `@@+`
}

func (paragraph Paragraph) Serve() ([]byte, error) {
    return nil, nil
}

type ParagraphElement struct {
	Pos lexer.Position
	// Normal text just gets lumped together into big chunks
	// stored in "NormalText" string
	// And the special elements like bold text or urls are stored separately
	Link *Link `( @@`
	Bold *Bold `| @@`
	Text *Text `| @@ )`
}

// Stub function for union types to satisfy the ParserUnionType interface
func (elem ParagraphElement) Union() {}

type Link struct {
	Text *Text `"[":OpenParen @@ "]":CloseParen`
	// Urls can only be local references for now
	Url string `"(":OpenParen @Ident ")":CloseParen`
}

type Bold struct {
	Text *Text `"*":Special @@ "*":Special`
}

type Text struct {
	Text string `@( Ident | RusWord | Number | Whitespace | OpenParen | CloseParen | Punct )+`
}
