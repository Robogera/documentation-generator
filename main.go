package main

import (
	"log"
	// "fmt"
)

func main() {
	var err error
	var ast *INI
	ast, err = getAST("status = 401")
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	log.Printf("Err %s", ast)
}
