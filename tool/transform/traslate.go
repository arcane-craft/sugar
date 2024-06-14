package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func TranslateSyntax[Type, Syntax fmt.Stringer](
	ctx context.Context, rootDir string, firstRun bool,
	inpectTypes func(p *packages.Package) []*Type,
	inspectSyntax func(p *packages.Package, instTypes []*Type) SyntaxInspector[Syntax],
	generate func(info *FileInfo[Syntax], writer io.Writer) error,
) error {

	var finished bool
	var buildFlags []string

	for !finished {
		finished = true

		pkgs, err := LoadPackages(ctx, rootDir, buildFlags...)
		if err != nil {
			return fmt.Errorf("load source packages failed: %w", err)
		}
		var instTypes []*Type
		for _, p := range pkgs {
			instTypes = append(instTypes, inpectTypes(p)...)
			for path, dep := range p.Imports {
				if p.PkgPath != path {
					instTypes = append(instTypes, inpectTypes(dep)...)
				}
			}
		}

		for _, t := range instTypes {
			fmt.Println(t)
		}

		var totalFiles []*FileInfo[Syntax]
		for _, p := range pkgs {
			files := NewPackageInspector(p, inspectSyntax(p, instTypes)).Inspect()
			for _, info := range files {
				func() {
					ext := filepath.Ext(info.Path)
					newFile := strings.TrimSuffix(info.Path, ext) + "_" + prodBuildTag + ext
					if _, err := os.Stat(newFile); err == nil {
						if !firstRun {
							return
						}
					}
					file, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						fmt.Println("open file", newFile, "failed:", err)
						return
					}
					defer file.Close()

					err = generate(info, file)
					if err != nil {
						fmt.Println("generate code of", info.Path, "failed:", err)
						return
					}

					if err = FormatCode(ctx, newFile, prodBuildTag); err != nil {
						fmt.Println("format code of", newFile, "failed:", err)
					}

					if len(buildFlags) > 0 {
						if err := os.Remove(info.Path); err != nil {
							fmt.Println("remove file", info.Path, "failed:", err)
							return
						}
						if err := os.Rename(newFile, info.Path); err != nil {
							fmt.Println("rename file", newFile, "failed:", err)
							return
						}
					} else {
						buf := bytes.NewBuffer(nil)
						if info.BuildFlag == nil {
							buf.Write([]byte(GenBuildFlags(false)))
						}
						bs, err := os.ReadFile(info.Path)
						if err != nil {
							fmt.Println("read file", info.Path, "failed:", err)
							return
						}
						buf.Write(bs)
						tmpOldFileNam := info.Path + ".old"
						if err := os.WriteFile(tmpOldFileNam, buf.Bytes(), 0644); err != nil {
							fmt.Println("write file", tmpOldFileNam, "failed:", err)
							return
						}
						if err := os.Remove(info.Path); err != nil {
							fmt.Println("remove file", info.Path, "failed:", err)
							return
						}
						if err := os.Rename(tmpOldFileNam, info.Path); err != nil {
							fmt.Println("rename file", tmpOldFileNam, "failed:", err)
							return
						}
					}
				}()
			}
			totalFiles = append(totalFiles, files...)
		}
		for _, f := range totalFiles {
			fmt.Printf("%s:\n", f.Path)
			for _, s := range f.Syntax {
				fmt.Printf("%s", s)
			}
		}
		if len(buildFlags) <= 0 || len(totalFiles) > 0 {
			finished = false
		}
		buildFlags = []string{"-tags=" + prodBuildTag}
	}
	return nil
}
