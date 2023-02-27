package main

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"reflect"
)

// ==============================================
// All declarative parsing rules go here
// Use verbose struct syntax insead of interface
// syntax whenever possible as it is more readable
//================================================

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
	// Position saved for better error messaging past the lexing and parsing
	// stages
	Pos lexer.Position

	Image     *Image     `( @@`
	Table     *Table     `| @@`
	Box       *Box       `| @@`
	List      *List      `| @@`
	Paragraph *Paragraph `| @@ ) EOL (EOL | EOF)`
}

func (entry Entry) Type() (entry_type string, err error) {
	// Returns entry type by searching for the first non-nil Field
	// that is not Pos (which is always non-nil op successful parse)
	// Returns error if all fields are nil or reflect panics
	// Which really shouldn't happen regardless of user input but
	// can happen if ENtry struct is implemented wrong
	defer func() {
		if recovered := recover(); recovered != nil {
			entry_type = ""
			err = fmt.Errorf("Reflect paniced, Entry struct implemented incorrectly?")
		}
	}()

	reflected_value := reflect.ValueOf(entry)
	var field reflect.Value

	// Iterating from 1 because
	for i := 1; i < reflected_value.NumField(); i++ {
		field = reflected_value.Field(i)
		if !reflect.ValueOf(field.Interface()).IsNil() {
			return field.Type().String(), nil
		}
	}

	return "", fmt.Errorf("All fields are nil, parser failed?")

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
	Text string `@( Ident | RusWord | Number | Whitespace | OpenParen | CloseParen | Punct )+`
}
