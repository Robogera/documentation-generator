package main

import (
	"fmt"
	"github.com/conflowio/parsley/ast/interpreter"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	terminal "github.com/conflowio/parsley/text/terminal"
	"os"
)

// NewParser returns with a new JSON parser
func NewParser() parsley.Parser {
	var value parser.Func

	array := combinator.SeqOf(
		terminal.Rune('['),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Array()),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	file := combinator.SeqOf(
		terminal.Word(nil, "!file", nil),
		combinator.SepBy(
			text.LeftTrim(&value, text.WsSpaces),
			text.LeftTrim(terminal.Word(nil, "\n\n", nil), text.WsSpaces),
		).Bind(interpreter.Array()),
	).Bind(interpreter.Select(1))

	keyValue := combinator.SeqOf(
		terminal.String("string", false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(&value, text.WsSpacesNl),
	)

	object := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.SepBy(
			text.LeftTrim(keyValue, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(','), text.WsSpaces),
		).Bind(interpreter.Object()),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Bind(interpreter.Select(1))

	value = combinator.Choice(
		terminal.String("string", false),
		terminal.Float("number"),
		terminal.Integer("integer"),
		array,
		object,
		terminal.Bool("boolean", "true", "false"),
		terminal.Nil("null", "null"),
	).Name("value")

	return value
}

func main() {
	jsonFilePath := "example.sml"
	if len(os.Args) > 1 {
		jsonFilePath = os.Args[1]
	}
	file, err := text.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}
	fs := parsley.NewFileSet(file)

	reader := text.NewReader(file)
	ctx := parsley.NewContext(fs, reader)
	s := combinator.Sentence(text.Trim(NewParser()))

	res, evalErr := parsley.Evaluate(ctx, s)
	if evalErr != nil {
		panic(evalErr)
	}
	fmt.Printf("Parser calls: %d\n", ctx.CallCount())
	fmt.Printf("%v\n", res)
}
