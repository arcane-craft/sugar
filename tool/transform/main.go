package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/arcane-craft/sugar/tool/transform/exception"
	"github.com/arcane-craft/sugar/tool/transform/lib"
	"github.com/arcane-craft/sugar/tool/transform/question"
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

	if err := lib.SetTmpBuildTags(ctx, rootDir); err != nil {
		fmt.Println("change build tags failed:", err)
		return
	}

	err := lib.TranslateSyntax(ctx, rootDir, true, question.NewTraslator())
	if err != nil {
		fmt.Println("translate question syntax failed:", err)
		return
	}

	err = lib.TranslateSyntax(ctx, rootDir, false, exception.NewTraslator())
	if err != nil {
		fmt.Println("translate exception syntax failed:", err)
		return
	}

	if err := lib.SetProdBuildTags(ctx, rootDir); err != nil {
		fmt.Println("change build tags failed:", err)
		return
	}
}
