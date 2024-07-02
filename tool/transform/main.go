package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/arcane-craft/sugar/tool/transform/lib"
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
	names := strings.Split(os.Getenv("SUGAR_AVAILABLE_SYNTAX"), ",")
	if len(names) == 1 && names[0] == "" {
		names = DefaultSyntax()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := lib.SetTmpBuildTags(ctx, rootDir); err != nil {
		fmt.Println("change build tags failed:", err)
		return
	}

	for idx, p := range Programs(names) {
		err := p.Run(ctx, rootDir, idx == 0)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	if err := lib.SetProdBuildTags(ctx, rootDir); err != nil {
		fmt.Println("change build tags failed:", err)
		return
	}
}
