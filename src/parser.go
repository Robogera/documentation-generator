package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// =========================
// Basic structure of a file
// =========================
type FILE struct {
	// Entries are the basic components of the file
	// Like paragraphs of text, tables, images, blocks, etc.
	// They are separated by two EOL tokens as specified in
	// The Entry struct tags
	Entries []*Entry `@@+`
}

type Entry struct {
	Pos lexer.Position

	Image     *Image     `( @@`
	Box       *Box       `| @@`
	List      *List      `| @@`
	Paragraph *Paragraph `| @@ ) EOL (EOL | EOF)`
}

// ===============
// Warn/info boxes
// ===============
type Box struct {
	Type       string       `"!":Punct @( "info":Ident | "warn":Ident )`
	Reference  string       `Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Image struct {
	Reference  string       `"!":Punct "img":Ident Whitespace @Ident`
	Path       *Path        `Whitespace @@`
	Paragraphs []*Paragraph `(EOL @@)*`
}

type List struct {
	Reference  string       `"!":Punct "list":Ident Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

type Table struct {
	Reference string `"!":Punct "table":Ident Whitespace @Ident`
	Title     *Text  `Whitespace @@`
	Rows      []*Row `@@+`
}

type Row struct {
    // TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths      []*Path      `EOL "img":Ident (Whitespace @@)`
	Paragraph *Paragraph `EOL @@`
}

// ===============
// Text paragraphs
// ===============
type Paragraph struct {
	// A "hardcore" regexp filter that enforces correct grammars on the parsing stage
	// was overkill and produced obscure error messages that are hard to interpret
	// TODO: check for basic mistakes like double whitespace, no punctuation, no capital letter, etc
	// after the parsing stage
	// Text string `@( Ident ( Punct? Whitespace OpenParen? Ident CloseParen? )* Punct )`
	Element *ParagraphElement `@@+`
}

type ParagraphElement struct {
	// Normal text just gets lumped together into big chunks
	// stored in "NormalText" string
	// And the special elements like bold text or urls are stored separately
	Link *Link `( @@`
	Bold *Bold `| @@`
	Text *Text `| @@ )`
}

type Link struct {
	Text *Text `"[":OpenParen @@ "]":CloseParen`
	// Urls can only be local references for now
	Url string `"(":OpenParen @Ident ")":CloseParen`
}

type Bold struct {
	Text *Text `"*":Special @@ "*":Special`
}

type Text struct {
	Text string `@( Ident | Number | Whitespace | OpenParen | CloseParen | Punct )+`
}

type Path struct {
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}
