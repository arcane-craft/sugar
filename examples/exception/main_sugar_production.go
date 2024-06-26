//go:build sugar_production

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	. "github.com/arcane-craft/sugar/syntax/exception"
)

func main() {
	Run()
}

func Run() ([]byte, error) {
	{
		var resultUPH82AR1SS []byte
		var catchErr4K9AHBGELS error
		var hasRet5PRN3G50I8 bool
		{

			file, errVJ948GMHAK := os.Open("hello.txt")
			if errVJ948GMHAK != nil {
				catchErr4K9AHBGELS = errVJ948GMHAK
				goto CatchTBC0MLSJ4S
			}

			defer file.Close()

			content, errL0OS45QF7O := io.ReadAll(file)
			if errL0OS45QF7O != nil {
				catchErr4K9AHBGELS = errL0OS45QF7O
				goto CatchTBC0MLSJ4S
			}

			resultUPH82AR1SS = content
			hasRet5PRN3G50I8 = true
			goto Finally8KDIGIEIBS

			goto Finally8KDIGIEIBS
		}
	CatchTBC0MLSJ4S:
		{
			if errors.As(catchErr4K9AHBGELS, *os.PathError) {
				err := catchErr4K9AHBGELS
				catchErr4K9AHBGELS = nil

				_, errRMSCQGU0GS := fmt.Println("error occured:", err)
				if errRMSCQGU0GS != nil {
					catchErr4K9AHBGELS = errRMSCQGU0GS
					goto Finally8KDIGIEIBS
				}

				catchErr4K9AHBGELS = err
				goto Finally8KDIGIEIBS

				goto Finally8KDIGIEIBS
			}
		}
	Finally8KDIGIEIBS:
		{
			if hasRet5PRN3G50I8 || catchErr4K9AHBGELS != nil {
				return resultUPH82AR1SS, catchErr4K9AHBGELS
			}
		}
	}
	return nil, nil
}
