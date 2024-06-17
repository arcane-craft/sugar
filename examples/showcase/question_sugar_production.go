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

func RunQuestion() Result[Empty] {
	varLG4CPGKQIO := WriteFile("example.txt", []byte("Hello, World!"), 0644)
	if varLG4CPGKQIO.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varLG4CPGKQIO.UnwrapErr()))
	}

	varA6DGJQJ4HG := ReadFile("example.txt")
	if varA6DGJQJ4HG.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varA6DGJQJ4HG.UnwrapErr()))
	}
	data := varA6DGJQJ4HG.Unwrap()

	fmt.Println("content1:", string(data))
	varR84QQ3VU68 := ReadFile("example.txt")
	if varR84QQ3VU68.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varR84QQ3VU68.UnwrapErr()))
	}
	fmt.Println("content2:", string(varR84QQ3VU68.Unwrap()))
	varBOBQPU848C := ReadFile("example.txt")
	if varBOBQPU848C.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varBOBQPU848C.UnwrapErr()))
	}
	if data := varBOBQPU848C.Unwrap(); len(data) > 0 {
		fmt.Println("content1:", string(data))
	}
	varI84GK8RJ6O := ReadFile("example.txt")
	if varI84GK8RJ6O.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varI84GK8RJ6O.UnwrapErr()))
	}
	varFISU6MJM64 := WriteFile("example.txt", varI84GK8RJ6O.Unwrap(), 0644)
	if varFISU6MJM64.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varFISU6MJM64.UnwrapErr()))
	}

	varJVJQ0ONNGG := OpenFile("example.txt")
	if varJVJQ0ONNGG.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varJVJQ0ONNGG.UnwrapErr()))
	}
	varSDJABSGJ7G := varJVJQ0ONNGG.Unwrap().Close()
	if varSDJABSGJ7G.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varSDJABSGJ7G.UnwrapErr()))
	}

	varE999TVMTCC := func() Result[[]byte] {
		var6RC26CH800 := ReadFile("example.txt")
		if var6RC26CH800.IsErr() {
			return Err[[]byte](fmt.Errorf("func() github.com/arcane-craft/sugar/result.Result[[]byte]: %w", var6RC26CH800.UnwrapErr()))
		}
		data := var6RC26CH800.Unwrap()

		return Ok(data)
	}()
	if varE999TVMTCC.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varE999TVMTCC.UnwrapErr()))
	}

	varP1VHL72LSS := Mkdir("example_dir", 0755)
	if varP1VHL72LSS.IsErr() {
		return Err[Empty](fmt.Errorf("func RunQuestion() github.com/arcane-craft/sugar/result.Result[github.com/arcane-craft/sugar.Empty]: %w", varP1VHL72LSS.UnwrapErr()))
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
	varDT85DBP2V0 := DeocdeJSON([]byte(`{"hello":"world"}`))
	if varDT85DBP2V0.IsNone() {
		return None[string]()
	}
	obj := varDT85DBP2V0.Unwrap()

	varMBV60NL5P8 := Get[string]("hello", obj)
	if varMBV60NL5P8.IsNone() {
		return None[string]()
	}
	val := varMBV60NL5P8.Unwrap()

	fmt.Println(val)
	return Some(val)
}
