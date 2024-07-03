//go:build !sugar_production

package main

import (
	"fmt"

	. "github.com/arcane-craft/sugar/syntax/predef"
)

func main() {
	TestIdentifier()
}

func TestIdentifier() {
	fmt.Println(Function__)
	fmt.Println(PrettyFunction__)
	fmt.Println(Package__)
	fmt.Println(File__)
	fmt.Println(Line__)
}
