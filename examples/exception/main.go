//go:build !sugar_production

package main

import (
	"fmt"
	"os"

	. "github.com/arcane-craft/sugar/syntax/exception"
)

func main() {
	Run()
}

type File struct{}

func (File) Close() error {
	fmt.Println("Close()")
	return nil
}

func OpenFile(name string) (*File, error) {
	fmt.Println("OpenFile()")
	return &File{}, nil
}

func ReadFile(name string) ([]byte, error) {
	fmt.Println("ReadFile()")
	return []byte("Hello, World!"), nil
}

func WriteFile(name string, data []byte, perm os.FileMode) error {
	fmt.Println("WriteFile()")
	return nil
}

func Mkdir(name string, perm os.FileMode) error {
	fmt.Println("Mkdir()")
	return nil
}

func Rename(oldpath, newpath string) error {
	fmt.Println("Rename()")
	return nil
}

func Run() ([]byte, error) {
	Try(func() {
		_ = WriteFile("example.txt", []byte("Hello, World!"), 0644)
		data, _ := ReadFile("example.txt")
		fmt.Println("content1:", string(data))
		if data, _ := ReadFile("example.txt"); len(data) > 0 {
			fmt.Println("content1:", string(data))
		}
		f, err := OpenFile("example.txt")
		if err != nil {
			Throw(err)
		}
		f.Close()
		Mkdir("example_dir", 0755)
		Return(data)
	}).Catch(Error(os.ErrPermission), func(err error) {
		fmt.Println("catch error:", err)
		Return(nil)
	}).Catch(Type[*os.PathError](), func(err error) {
		fmt.Println("catch error type:", err)
		Throw(err)
	}).Finally(func() {
		Return([]byte{})
	})
	return nil, nil
}

func Run2() ([]byte, error) {

	Try(func() {
		_ = WriteFile("example.txt", []byte("Hello, World!"), 0644)
		data, _ := ReadFile("example.txt")
		fmt.Println("content1:", string(data))
		Return(data)
	}).Catch(Type[*os.PathError](), func(err error) {
		fmt.Println("catch error:", err)
	})

	return nil, nil
}
