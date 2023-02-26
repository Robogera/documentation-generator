package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// =========================//
// Basic structure of a file
// =========================//
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

// ======================//
// Warn/info boxes
// ======================//
type Box struct {
	//Reference string     `"!info":Command Whitespace @Ident EOL`
	Type       string       `"!":Punct @( "info":Ident | "warn":Ident )`
	Reference  string       `Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)*`
}

type Image struct {
	Reference  string       `"!":Punct "img":Ident Whitespace @Ident`
	Path       string       `Whitespace @("/":Special? (Ident "/":Special)* Ident "." Ident)`
	Paragraphs []*Paragraph `(EOL @@)*`
}

type List struct {
	Reference  string       `"!list":Command Whitespace @Ident`
	Paragraphs []*Paragraph `(EOL @@)*`
}

// ======================//
// Text paragraphs
// ======================//
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
	Link     *Link       `( @@`
	BoldText *BoldText   `| @@`
	Text     *NormalText `| @@ )`
}

type Link struct {
	// TODO: find way to not repeat this token construct in every elemnt declaration
	// that just wraps around the text chunks
	Text string `"[":OpenParen @( Ident | Whitespace | OpenParen | CloseParen | Punct )+ "]":CloseParen`
	Url  string `"(":OpenParen @( Ident | Whitespace | OpenParen | CloseParen | Punct )+ ")":CloseParen`
}

type NormalText struct {
	Text string `@( Ident | Whitespace | OpenParen | CloseParen | Punct )+`
}

type BoldText struct {
	Text string `"*":Special @( Ident | Whitespace | OpenParen | CloseParen | Punct )+ "*":Special`
}
