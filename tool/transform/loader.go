package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func LoadPackages(ctx context.Context, path string, buildFlags ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedModule,
		Context:    ctx,
		Logf:       log.Printf,
		Dir:        path,
		Env:        os.Environ(),
		BuildFlags: buildFlags,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("packages.Load() failed: %w", err)
	}
	return pkgs, nil
}
