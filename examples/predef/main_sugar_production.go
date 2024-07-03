//go:build sugar_production

package main

import (
	"fmt"

	. "github.com/arcane-craft/sugar/syntax/predef"
)

func main() {
	TestIdentifier()
}

func TestIdentifier() {
	fmt.Println("TestIdentifier")
	fmt.Println("func TestIdentifier()")
	fmt.Println("github.com/arcane-craft/sugar/examples/predef")
	fmt.Println("/home/jinzhao/Documents/sugar/examples/predef/main.go")
	fmt.Println("20")
}
