//go:build sugar_production

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
	varTJNIL9FA3O := WriteFile("example.txt", []byte("Hello, World!"), 0644)
	if varTJNIL9FA3O.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varTJNIL9FA3O.UnwrapErr()))
	}

	varDNAUG2UPL8 := ReadFile("example.txt")
	if varDNAUG2UPL8.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varDNAUG2UPL8.UnwrapErr()))
	}
	data := varDNAUG2UPL8.Unwrap()

	fmt.Println("content1:", string(data))
	varMT88AS5RP0 := ReadFile("example.txt")
	if varMT88AS5RP0.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varMT88AS5RP0.UnwrapErr()))
	}
	fmt.Println("content2:", string(varMT88AS5RP0.Unwrap()))
	varRCICRHK43G := ReadFile("example.txt")
	if varRCICRHK43G.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varRCICRHK43G.UnwrapErr()))
	}
	if data := varRCICRHK43G.Unwrap(); len(data) > 0 {
		fmt.Println("content1:", string(data))
		varNNIK80T43G := WriteFile("example.txt", []byte("Hello, World!"), 0644)
		if varNNIK80T43G.IsErr() {
			return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varNNIK80T43G.UnwrapErr()))
		}

	}
	varATG362K9JC := ReadFile("example.txt")
	if varATG362K9JC.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varATG362K9JC.UnwrapErr()))
	}
	varNO0DC79NC4 := WriteFile("example.txt", varATG362K9JC.Unwrap(), 0644)
	if varNO0DC79NC4.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varNO0DC79NC4.UnwrapErr()))
	}

	varSSGNBDM3GO := OpenFile("example.txt")
	if varSSGNBDM3GO.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varSSGNBDM3GO.UnwrapErr()))
	}
	var2IMR04IIR4 := varSSGNBDM3GO.Unwrap().Close()
	if var2IMR04IIR4.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", var2IMR04IIR4.UnwrapErr()))
	}

	varGJM4OU2PSK := func() Result[[]byte] {
		varK3S1CF79HK := ReadFile("example.txt")
		if varK3S1CF79HK.IsErr() {
			return Err[[]byte](fmt.Errorf("func() github.com/arcane-craft/sugar/result.Result[[]byte]: %w", varK3S1CF79HK.UnwrapErr()))
		}
		data := varK3S1CF79HK.Unwrap()

		return Ok(data)
	}()
	if varGJM4OU2PSK.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varGJM4OU2PSK.UnwrapErr()))
	}

	var54T1SAI3VG := Mkdir("example_dir", 0755)
	if var54T1SAI3VG.IsErr() {
		return Err[Empty](fmt.Errorf("func RunResultQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", var54T1SAI3VG.UnwrapErr()))
	}

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
	varG65MLTKE20 := DeocdeJSON([]byte(`{"hello":"world"}`))
	if varG65MLTKE20.IsNone() {
		return None[string]()
	}
	obj := varG65MLTKE20.Unwrap()

	varC4E6FQFSVS := Get[string]("hello", obj)
	if varC4E6FQFSVS.IsNone() {
		return None[string]()
	}
	val := varC4E6FQFSVS.Unwrap()

	fmt.Println(val)
	return Some(val)
}
