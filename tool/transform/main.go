package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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

	err := lib.TranslateSyntax(ctx, rootDir, true, question.NewTraslator())
	if err != nil {
		fmt.Println("translate delay syntax failed:", err)
		return
	}
}
