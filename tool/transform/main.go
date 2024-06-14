package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/arcane-craft/sugar/tool/transform/lib"
	"github.com/arcane-craft/sugar/tool/transform/question"
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

	err := lib.TranslateSyntax(ctx, rootDir, true,
		func(p *packages.Package) []*question.QuestionInstanceType {
			return question.NewQuestionTypeInspector(p).InspectQuestionTypes()
		},
		func(p *packages.Package, instTypes []*question.QuestionInstanceType) lib.SyntaxInspector[question.QuestionSyntax] {
			return question.NewQuestionSyntaxInspector(p, instTypes)
		},
		func(info *lib.FileInfo[question.QuestionSyntax], writer io.Writer) error {
			return question.GenerateQuestionSyntax(info, writer)
		},
	)
	if err != nil {
		fmt.Println("translate delay syntax failed:", err)
		return
	}
}
