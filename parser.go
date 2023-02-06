package main

import (
	"github.com/alecthomas/participle/v2"
)

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	String *string  `  @String`
	Float  *float64 `| @Float`
	Int    *int     `| @Int`
}

func getAST(input string) (*INI, error) {
	var parser *participle.Parser[INI]
	var err error

	parser, err = participle.Build[INI]()
	if err != nil {
		return nil, err
	}

	var ast *INI

	ast, err = parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
