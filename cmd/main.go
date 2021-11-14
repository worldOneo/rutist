package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/worldOneo/rutist/ast"
	"github.com/worldOneo/rutist/interpreter"
	"github.com/worldOneo/rutist/tokens"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "main.rut", "Defines the file to execute")
	flag.Parse()
	abs, err := filepath.Abs(file)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	code := string(content)
	tokens, err := tokens.Lexer(code)
	if err != nil {
		log.Fatal(err)
	}
	parsed, err := ast.Parse(tokens, abs)
	if err != nil {
		log.Fatal(err)
	}
	_, err = interpreter.Run(abs, parsed)
	if err != nil {
		log.Fatal(err)
	}
}
