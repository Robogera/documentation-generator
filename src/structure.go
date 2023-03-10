package main

import (
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
	Header    *Header    `| @@`
	Table     *Table     `| @@`
	Box       *Box       `| @@`
	List      *List      `| @@`
	Paragraph *Paragraph `| @@ ) EOL (EOL | EOF)`
}

type Header struct {
	Level string `@("*":Special "*":Special? "*":Special?)`
	Text  *Text  `Whitespace @@`
	ID    string `("#":Special @Ident)?`
}

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths      []*Path      `"img":Ident (Whitespace @@)+`
	Paragraphs []*Paragraph `(EOL "txt":Ident Whitespace @@)+`
}

type Table struct {
	// TODO: add colors
	Title *Text  `"!table":Command Whitespace @@`
	ID    string `("#":Special @Ident)?`
	Rows  []*Row `(EOL @@)+`
}

type Image struct {
	Path       *Path        `"!img":Command Whitespace @@`
	ID         string       `(Whitespace "#":Special @Ident)?`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Path struct {
	// Used by multiple file elements
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}

type List struct {
	// TODO: add unordered lists
	// (would require some syntax change to distinguish them)
	ID         string       `"!list":Command (Whitespace "#":Special @Ident)?`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Box struct {
	Type       string       `@( "!info":Command | "!warn":Command )`
	ID         string       `(Whitespace "#":Special @Ident)?`
	Paragraphs []*Paragraph `(EOL @@)+`
}

// Stub function for union types to satisfy the ParserUnionType interface
func (entry Entry) Union() {}

type Paragraph struct {
	Element []*ParagraphElement `@@+`
}

type ParagraphElement struct {
	Pos lexer.Position
	// Normal text just gets lumped together into big chunks
	// stored in "NormalText" string
	// And the special elements like bold text or urls are stored separately
	Link *Link `( @@`
	Bold *Bold `| @@`
	Code *Code `| @@`
	Text *Text `| @@ )`
}

// Stub function for union types to satisfy the ParserUnionType interface
func (elem ParagraphElement) Union() {}

type Link struct {
	Text *Text `"[" @@ "]"`
	// Urls can only be local references for now
	Url string `"(" @Ident ")"`
}

type Bold struct {
	Text *Text `"*":Special @@ "*":Special`
}

type Text struct {
	Text string `@( Ident | RusWord | Number | Whitespace | OpenParen | CloseParen | Punct )+`
}

type Code struct {
	Text string `Code @( Ident | RusWord | Number | Whitespace | OpenParen | CloseParen | Punct | Special )+ Code`
}
