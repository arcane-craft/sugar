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
	var file *os.File
	{
		var resultUPH82AR1SS []byte
		var catchErr4K9AHBGELS error
		var hasRetG3VL22R4DG bool
		{

			var errMSCVA62MQ0 error
			file, errMSCVA62MQ0 = os.Open("hello.txt")
			if errMSCVA62MQ0 != nil {
				catchErr4K9AHBGELS = errMSCVA62MQ0
				goto CatchEOMC3UC1OO
			}

			content, errL0OS45QF7O := io.ReadAll(file)
			if errL0OS45QF7O != nil {
				catchErr4K9AHBGELS = errL0OS45QF7O
				goto CatchEOMC3UC1OO
			}

			resultUPH82AR1SS = content
			hasRetG3VL22R4DG = true
			goto FinallyB5JUD32NLO

			goto FinallyB5JUD32NLO
		}
	CatchEOMC3UC1OO:
		{
			if errors.As(catchErr4K9AHBGELS, *os.PathError) {
				err := catchErr4K9AHBGELS
				catchErr4K9AHBGELS = nil

				_, errRMSCQGU0GS := fmt.Println("error occured:", err)
				if errRMSCQGU0GS != nil {
					catchErr4K9AHBGELS = errRMSCQGU0GS
					goto FinallyB5JUD32NLO
				}

				catchErr4K9AHBGELS = err
				goto FinallyB5JUD32NLO

				goto FinallyB5JUD32NLO
			}
		}
	FinallyB5JUD32NLO:
		{

			if file != nil {

				file.Close()

			}

			if hasRetG3VL22R4DG || catchErr4K9AHBGELS != nil {
				return resultUPH82AR1SS, catchErr4K9AHBGELS
			}
		}
	}
	return nil, nil
}
