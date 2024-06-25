//go:build !sugar_production

package main

import (
	"fmt"
	"os"

	. "github.com/arcane-craft/sugar"
	. "github.com/arcane-craft/sugar/option"
	. "github.com/arcane-craft/sugar/result"
)

func main() {
	RunResultQuestion()
	RunOptionQuestion()
}

type File struct{}

func (File) Close() Result[Empty] {
	fmt.Println("Close()")
	return Ok(Empty{})
}

func OpenFile(name string) Result[*File] {
	fmt.Println("OpenFile()")
	return Ok(&File{})
}

func ReadFile(name string) Result[[]byte] {
	fmt.Println("ReadFile()")
	return Ok([]byte("Hello, World!"))
}

func WriteFile(name string, data []byte, perm os.FileMode) Result[Empty] {
	fmt.Println("WriteFile()")
	return Ok(Empty{})
}

func Mkdir(name string, perm os.FileMode) Result[Empty] {
	fmt.Println("Mkdir()")
	return Ok(Empty{})
}

func Rename(oldpath, newpath string) Result[Empty] {
	fmt.Println("Rename()")
	return Ok(Empty{})
}

func RunResultQuestion() Result[Empty] {
	WriteFile("example.txt", []byte("Hello, World!"), 0644).Q()
	data := ReadFile("example.txt").Q()
	fmt.Println("content1:", string(data))
	fmt.Println("content2:", string(ReadFile("example.txt").Q()))
	if data := ReadFile("example.txt").Q(); len(data) > 0 {
		fmt.Println("content1:", string(data))
		WriteFile("example.txt", []byte("Hello, World!"), 0644).Q()
	}
	WriteFile("example.txt", ReadFile("example.txt").Q(), 0644).Q()
	OpenFile("example.txt").Q().Close().Q()
	func() Result[[]byte] {
		data := ReadFile("example.txt").Q()
		return Ok(data)
	}().Q()
	Mkdir("example_dir", 0755).Q()
	return Rename("example.txt", "example_dir/example_renamed.txt")
}

func DeocdeJSON(content []byte) Option[map[string]any] {
	return Some(map[string]any{})
}

func Get[T any](key string, obj map[string]any) Option[T] {
	fmt.Println("Get()")
	var t T
	return Some(t)
}

func RunOptionQuestion() Option[string] {
	obj := DeocdeJSON([]byte(`{"hello":"world"}`)).Q()
	val := Get[string]("hello", obj).Q()
	fmt.Println(val)
	return Some(val)
}
