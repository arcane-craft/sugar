package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/tools/go/packages"
)

func main() {
	var rootDir string
	if len(os.Args) < 2 {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("get workspace failed:", err)
			os.Exit(2)
			return
		}
		rootDir = pwd
	} else {
		rootDir = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := TranslateSyntax(ctx, rootDir, true,
		func(p *packages.Package) []*QuestionInstanceType {
			return NewQuestionTypeInspector(p).InspectQuestionTypes()
		},
		func(p *packages.Package, instTypes []*QuestionInstanceType) SyntaxInspector[QuestionSyntax] {
			return NewQuestionSyntaxInspector(p, instTypes)
		},
		func(info *FileInfo[QuestionSyntax], writer io.Writer) error {
			return GenerateQuestionSyntax(info, writer)
		},
	)
	if err != nil {
		fmt.Println("translate delay syntax failed:", err)
		return
	}
}
