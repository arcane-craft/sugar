package main

import (
	"fmt"
	"os"
)

func WriteFile(name string, data []byte, perm os.FileMode) error {
	fmt.Println("WriteFile()")
	return nil
}

func ReadFile(name string) ([]byte, error) {
	fmt.Println("ReadFile()")
	return []byte("Hello, World!"), nil
}

func Mkdir(name string, perm os.FileMode) error {
	fmt.Println("Mkdir()")
	return nil
}

func Rename(oldpath, newpath string) error {
	fmt.Println("Rename()")
	return nil
}

func RunClassic() error {
	err := WriteFile("example.txt", []byte("Hello, World!"), 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile() failed: %w", err)
	}
	data, err := ReadFile("example.txt")
	if err != nil {
		return fmt.Errorf("os.ReadFile() failed: %w", err)
	}
	fmt.Println("content:", string(data))
	err = Mkdir("example_dir", 0755)
	if err != nil {
		return fmt.Errorf("os.Mkdir() failed: %w", err)
	}
	err = Rename("example.txt", "example_dir/example_renamed.txt")
	if err != nil {
		return fmt.Errorf("os.Rename() failed: %w", err)
	}
	return nil
}
