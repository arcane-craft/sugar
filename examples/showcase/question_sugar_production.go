//go:build sugar_production

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
	varST6P35MNAS := WriteFile2("example.txt", []byte("Hello, World!"), 0644)
	if varST6P35MNAS.IsErr() {
		return Err[Empty](varST6P35MNAS.UnwrapErr())
	}

	varFTT4P543LO := ReadFile2("example.txt")
	if varFTT4P543LO.IsErr() {
		return Err[Empty](varFTT4P543LO.UnwrapErr())
	}
	data := varFTT4P543LO.Unwrap()

	fmt.Println("content1:", string(data))
	varGJAEVSDNS4 := ReadFile2("example.txt")
	if varGJAEVSDNS4.IsErr() {
		return Err[Empty](varGJAEVSDNS4.UnwrapErr())
	}
	fmt.Println("content2:", string(varGJAEVSDNS4.Unwrap()))
	varUEQK03HDPC := ReadFile2("example.txt")
	if varUEQK03HDPC.IsErr() {
		return Err[Empty](varUEQK03HDPC.UnwrapErr())
	}
	if data := varUEQK03HDPC.Unwrap(); len(data) > 0 {
		fmt.Println("content1:", string(data))
	}
	varU2UTTF90EG := ReadFile2("example.txt")
	if varU2UTTF90EG.IsErr() {
		return Err[Empty](varU2UTTF90EG.UnwrapErr())
	}
	var7OFN9MET64 := WriteFile2("example.txt", varU2UTTF90EG.Unwrap(), 0644)
	if var7OFN9MET64.IsErr() {
		return Err[Empty](var7OFN9MET64.UnwrapErr())
	}

	varRMF8PAD6HS := OpenFile("example.txt")
	if varRMF8PAD6HS.IsErr() {
		return Err[Empty](varRMF8PAD6HS.UnwrapErr())
	}
	var94QLG8ESBC := varRMF8PAD6HS.Unwrap().Close()
	if var94QLG8ESBC.IsErr() {
		return Err[Empty](var94QLG8ESBC.UnwrapErr())
	}

	varBO2431Q6L8 := func() Result[[]byte] {
		varIU2GIB8D0K := ReadFile2("example.txt")
		if varIU2GIB8D0K.IsErr() {
			return Err[[]byte](varIU2GIB8D0K.UnwrapErr())
		}
		data := varIU2GIB8D0K.Unwrap()

		return Ok(data)
	}()
	if varBO2431Q6L8.IsErr() {
		return Err[Empty](varBO2431Q6L8.UnwrapErr())
	}

	var7J0BIB89PK := Mkdir2("example_dir", 0755)
	if var7J0BIB89PK.IsErr() {
		return Err[Empty](var7J0BIB89PK.UnwrapErr())
	}

	return Rename2("example.txt", "example_dir/example_renamed.txt")
}

func RunOptionQuestion() Option[[]byte] {
	varGV1J5K7F6S := WriteFile2("example.txt", []byte("Hello, World!"), 0644).Ok()
	if varGV1J5K7F6S.IsNone() {
		return None[[]byte]()
	}

	var0UR61RRT4S := ReadFile2("example.txt").Ok()
	if var0UR61RRT4S.IsNone() {
		return None[[]byte]()
	}
	data := var0UR61RRT4S.Unwrap()

	fmt.Println("content1:", string(data))
	return Some(data)
}
