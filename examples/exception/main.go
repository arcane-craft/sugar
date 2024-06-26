//go:build !sugar_production

package main

import (
	"fmt"
	"io"
	"os"

	. "github.com/arcane-craft/sugar/syntax/exception"
)

func main() {
	Run()
}

func Run() ([]byte, error) {
	Try(func() {
		file, _ := os.Open("hello.txt")
		defer file.Close()
		content, _ := io.ReadAll(file)
		Return(content)
	}).Catch(Type[*os.PathError](), func(err error) {
		fmt.Println("error occured:", err)
		Throw(err)
	})
	return nil, nil
}
