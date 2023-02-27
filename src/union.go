package main

// ===========================
// Generic function that returns the type information for the
// "union types" (see participle documentation)
// All union types should have lexer.Position as the first field
// and only one other field is allowed to be non-nil
// ===========================

import (
	"fmt"
	"reflect"
)

type ParserUnionType interface {
	*Entry | *ParagraphElement
}

func unionType[T ParserUnionType](data_structure T) (entry_type string, err error) {
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

	reflected_value := reflect.ValueOf(data_structure).Elem()
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