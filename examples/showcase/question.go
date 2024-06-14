//go:build !sugar_production

package main

import (
	"fmt"
	"os"

	. "github.com/arcane-craft/sugar"
	. "github.com/arcane-craft/sugar/option"
	. "github.com/arcane-craft/sugar/result"
)

type File struct{}

func (File) Close() Result[Empty] {
	return Ok(Empty{})
}

func OpenFile(name string) Result[*File] {
	panic("TODO")
}

func ReadFile2(name string) Result[[]byte] {
	return FromR1Err(ReadFile(name))
}

func WriteFile2(name string, data []byte, perm os.FileMode) Result[Empty] {
	return FromErr(WriteFile(name, data, perm))
}

func Mkdir2(name string, perm os.FileMode) Result[Empty] {
	return FromErr(Mkdir(name, perm))
}

func Rename2(oldpath, newpath string) Result[Empty] {
	return FromErr(Rename(oldpath, newpath))
}

func RunQuestion() Result[Empty] {
	WriteFile2("example.txt", []byte("Hello, World!"), 0644).Q()
	data := ReadFile2("example.txt").Q()
	fmt.Println("content1:", string(data))
	fmt.Println("content2:", string(ReadFile2("example.txt").Q()))
	if data := ReadFile2("example.txt").Q(); len(data) > 0 {
		fmt.Println("content1:", string(data))
	}
	WriteFile2("example.txt", ReadFile2("example.txt").Q(), 0644).Q()
	OpenFile("example.txt").Q().Close().Q()
	func() Result[[]byte] {
		data := ReadFile2("example.txt").Q()
		return Ok(data)
	}().Q()
	Mkdir2("example_dir", 0755).Q()
	return Rename2("example.txt", "example_dir/example_renamed.txt")
}

func RunOptionQuestion() Option[[]byte] {
	WriteFile2("example.txt", []byte("Hello, World!"), 0644).Ok().Q()
	data := ReadFile2("example.txt").Ok().Q()
	fmt.Println("content1:", string(data))
	return Some(data)
}
