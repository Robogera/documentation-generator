package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// ==============================================
// All declarative parsing rules go here
// Use verbose struct syntax insead of interface
// syntax whenever possible as it is more readable
//================================================

// =========================
// Basic structure of a file
// =========================
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
	Table     *Table     `| @@`
	Box       *Box       `| @@`
	List      *List      `| @@`
	Paragraph *Paragraph `| @@ ) EOL (EOL | EOF)`
}

// ===============
// Warn/info boxes
// ===============
type Box struct {
	Type       string       `@( "!info":Command | "!warn":Command )`
	Reference  string       `Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

// =============
// Single images
// =============
type Image struct {
	Reference  string       `"!img":Command Whitespace @Ident`
	Path       *Path        `Whitespace @@`
	Paragraphs []*Paragraph `(EOL @@)*`
}

// =============
// Ordered lists
// =============
type List struct {
	// TODO: add unordered lists
	// (would require some syntax change to distinguish them)
	Reference  string       `"!list":Command Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)+`
}

// ==================
// Tables with images
// ==================
type Table struct {
	// TODO: add colors
	Reference string `"!table":Command Whitespace @Ident`
	Title     *Text  `Whitespace @@`
	Rows      []*Row `(EOL @@)+`
}

type Path struct {
	// Used by multiple file elements
	Path string `@("/":Special? (Ident "/":Special)* Ident "." Ident)`
}

type Row struct {
	// TODO: (maybe) make it possible to have multiple paragraphs per row
	Paths     []*Path    `"img":Ident (Whitespace @@)`
	Paragraph *Paragraph `EOL "txt":Ident Whitespace @@`
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
