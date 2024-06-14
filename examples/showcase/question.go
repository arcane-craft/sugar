package main

import (
	"fmt"
	"os"

	. "github.com/arcane-craft/sugar"
	. "github.com/arcane-craft/sugar/result"
)

func QReadFile(name string) Result[[]byte] {
	return FromR1Err(ReadFile(name))
}

func QWriteFile(name string, data []byte, perm os.FileMode) Result[Empty] {
	return FromErr(WriteFile(name, data, perm))
}

func QMkdir(name string, perm os.FileMode) Result[Empty] {
	return FromErr(Mkdir(name, perm))
}

func QRename(oldpath, newpath string) Result[Empty] {
	return FromErr(Rename(oldpath, newpath))
}

func RunQuestion() Result[Empty] {
	QWriteFile("example.txt", []byte("Hello, World!"), 0644).Q()
	data := QReadFile("example.txt").Q()
	fmt.Println("content:", string(data))
	QMkdir("example_dir", 0755).Q()
	return QRename("example.txt", "example_dir/example_renamed.txt")
}
