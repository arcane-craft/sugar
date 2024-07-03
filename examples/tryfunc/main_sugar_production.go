//go:build sugar_production

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

func TryOnly() (_ []byte, err9EIH73OK48 error) {
	file, errCEGNQ31MMK := os.Open("hello.txt")
	if errCEGNQ31MMK != nil {
		err9EIH73OK48 = errCEGNQ31MMK
		return
	}
	defer file.Close()
	content, errT07OPFVK5O := io.ReadAll(file)
	if errT07OPFVK5O != nil {
		err9EIH73OK48 = errT07OPFVK5O
		return
	}
	return content, nil
}

func TryWithDefer() (c []byte, e error) {
	defer func() {
		if errors.As(e, new(*os.PathError)) {
			fmt.Println("invalid file path")
			e = nil
		}
	}()

	file, errFENAO9VRCC := os.Open("hello.txt")
	if errFENAO9VRCC != nil {
		e = errFENAO9VRCC
		return
	}
	defer file.Close()
	content, err8T86277SQK := io.ReadAll(file)
	if err8T86277SQK != nil {
		e = err8T86277SQK
		return
	}
	return content, nil
}
