//go:build sugar_production

package main

import (
	"errors"
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
	{
		var resultH9VD59RS04 []byte
		var catchErrEACBAU03SS error
		var hasRetJ4RA0PNIJG bool
		{

			errND678M6OIK := WriteFile("example.txt", []byte("Hello, World!"), 0644)
			if errND678M6OIK != nil {
				catchErrEACBAU03SS = errND678M6OIK
				goto CatchUBBOEM0QOG
			}

			data, errQSR81CLNH8 := ReadFile("example.txt")
			if errQSR81CLNH8 != nil {
				catchErrEACBAU03SS = errQSR81CLNH8
				goto CatchUBBOEM0QOG
			}

			_, errRFTFE4BD5C := fmt.Println("content1:", string(data))
			if errRFTFE4BD5C != nil {
				catchErrEACBAU03SS = errRFTFE4BD5C
				goto CatchUBBOEM0QOG
			}

			data, err0CLG22HJ6G := ReadFile("example.txt")
			if err0CLG22HJ6G != nil {
				catchErrEACBAU03SS = err0CLG22HJ6G
				goto CatchUBBOEM0QOG
			}
			if len(data) > 0 {

				_, errE2GVJ4MNFG := fmt.Println("content1:", string(data))
				if errE2GVJ4MNFG != nil {
					catchErrEACBAU03SS = errE2GVJ4MNFG
					goto CatchUBBOEM0QOG
				}

			}
			f, err := OpenFile("example.txt")
			if err != nil {

				catchErrEACBAU03SS = err
				goto CatchUBBOEM0QOG

			}

			errDVABHN7ES0 := f.Close()
			if errDVABHN7ES0 != nil {
				catchErrEACBAU03SS = errDVABHN7ES0
				goto CatchUBBOEM0QOG
			}

			err6VL84C2IM8 := Mkdir("example_dir", 0755)
			if err6VL84C2IM8 != nil {
				catchErrEACBAU03SS = err6VL84C2IM8
				goto CatchUBBOEM0QOG
			}

			resultH9VD59RS04 = data
			hasRetJ4RA0PNIJG = true
			goto Finally3GPFC07HG4

			goto Finally3GPFC07HG4
		}
	CatchUBBOEM0QOG:
		{
			if errors.Is(catchErrEACBAU03SS, os.ErrPermission) {
				err := catchErrEACBAU03SS
				catchErrEACBAU03SS = nil

				_, err60VD505L8O := fmt.Println("catch error:", err)
				if err60VD505L8O != nil {
					catchErrEACBAU03SS = err60VD505L8O
					goto Finally3GPFC07HG4
				}

				resultH9VD59RS04 = nil
				hasRetJ4RA0PNIJG = true
				goto Finally3GPFC07HG4

				goto Finally3GPFC07HG4
			}
		}
		{
			if errors.As(catchErrEACBAU03SS, *os.PathError) {
				err := catchErrEACBAU03SS
				catchErrEACBAU03SS = nil

				_, errPDTQ416PIC := fmt.Println("catch error type:", err)
				if errPDTQ416PIC != nil {
					catchErrEACBAU03SS = errPDTQ416PIC
					goto Finally3GPFC07HG4
				}

				catchErrEACBAU03SS = err
				goto Finally3GPFC07HG4

				goto Finally3GPFC07HG4
			}
		}
	Finally3GPFC07HG4:
		{

			resultH9VD59RS04 = []byte{}
			return resultH9VD59RS04, catchErrEACBAU03SS

			if hasRetJ4RA0PNIJG || catchErrEACBAU03SS != nil {
				return resultH9VD59RS04, catchErrEACBAU03SS
			}
		}
	}
	return nil, nil
}

func Run2() ([]byte, error) {

	{
		var result4MOLFL1TVS []byte
		var catchErrHOR7F7GJHO error
		var hasRetOIL5MBR1V4 bool
		{

			errOP2SLMBJ1G := WriteFile("example.txt", []byte("Hello, World!"), 0644)
			if errOP2SLMBJ1G != nil {
				catchErrHOR7F7GJHO = errOP2SLMBJ1G
				goto CatchLVOKD82L6G
			}

			data, errANME00MMT0 := ReadFile("example.txt")
			if errANME00MMT0 != nil {
				catchErrHOR7F7GJHO = errANME00MMT0
				goto CatchLVOKD82L6G
			}

			_, errPRAHFCLLCC := fmt.Println("content1:", string(data))
			if errPRAHFCLLCC != nil {
				catchErrHOR7F7GJHO = errPRAHFCLLCC
				goto CatchLVOKD82L6G
			}

			result4MOLFL1TVS = data
			hasRetOIL5MBR1V4 = true
			goto Finally67MKL14HSO

			goto Finally67MKL14HSO
		}
	CatchLVOKD82L6G:
		{
			if errors.As(catchErrHOR7F7GJHO, *os.PathError) {
				err := catchErrHOR7F7GJHO
				catchErrHOR7F7GJHO = nil

				_, errURAC23OTHS := fmt.Println("catch error:", err)
				if errURAC23OTHS != nil {
					catchErrHOR7F7GJHO = errURAC23OTHS
					goto Finally67MKL14HSO
				}

				goto Finally67MKL14HSO
			}
		}
	Finally67MKL14HSO:
		{
			if hasRetOIL5MBR1V4 || catchErrHOR7F7GJHO != nil {
				return result4MOLFL1TVS, catchErrHOR7F7GJHO
			}
		}
	}

	return nil, nil
}
