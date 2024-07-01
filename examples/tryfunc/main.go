//go:build !sugar_production

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	. "github.com/arcane-craft/sugar/syntax/tryfunc"
)

func main() {
	TryOnly()
}

func TryOnly() ([]byte, error) {
	file := Try(os.Open("hello.txt"))
	defer file.Close()
	content := Try(io.ReadAll(file))
	return content, nil
}

func TryWithDefer() (c []byte, e error) {
	defer func() {
		if errors.As(e, new(*os.PathError)) {
			fmt.Println("invalid file path")
			e = nil
		}
	}()

	file := Try(os.Open("hello.txt"))
	defer file.Close()
	content := Try(io.ReadAll(file))
	return content, nil
}
