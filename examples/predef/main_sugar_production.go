//go:build sugar_production

package main

import "fmt"

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
