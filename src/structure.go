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
}

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths      []*Path      `"img":Ident (Whitespace @@)+`
	Paragraphs []*Paragraph `(EOL "txt":Ident Whitespace @@)+`
}

type Table struct {
	// TODO: add colors
	Reference string `"!table":Command Whitespace @Ident`
	Title     *Text  `Whitespace @@`
	Rows      []*Row `(EOL @@)+`
}

type Image struct {
	Reference  string       `"!img":Command Whitespace @Ident`
	Path       *Path        `Whitespace @@`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Path struct {
	// Used by multiple file elements
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}

type List struct {
	// TODO: add unordered lists
	// (would require some syntax change to distinguish them)
	Reference  string       `"!list":Command Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Box struct {
	Type       string       `@( "!info":Command | "!warn":Command )`
	Reference  string       `Whitespace @Ident`
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
