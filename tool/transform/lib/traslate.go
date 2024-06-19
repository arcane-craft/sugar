package lib

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

type SyntaxTraslator[Type, Syntax interface {
	fmt.Stringer
	comparable
}] interface {
	InpectTypes(p *packages.Package) []Type
	InspectSyntax(p *packages.Package, instTypes []Type) SyntaxInspector[Syntax]
	Generate(info *FileInfo[Syntax], writer io.Writer) error
}

func TranslateSyntax[Type, Syntax interface {
	fmt.Stringer
	comparable
}](
	ctx context.Context, rootDir string, firstRun bool,
	traslator SyntaxTraslator[Type, Syntax],
) error {

	var finished bool
	var buildFlags []string
	modifiedFiles := map[string]struct{}{}

	for !finished {
		finished = true

		pkgs, err := LoadPackages(ctx, rootDir, buildFlags...)
		if err != nil {
			return fmt.Errorf("load source packages failed: %w", err)
		}
		var instTypes []Type
		for _, p := range pkgs {
			instTypes = append(instTypes, traslator.InpectTypes(p)...)
			for path, dep := range p.Imports {
				if p.PkgPath != path {
					instTypes = append(instTypes, traslator.InpectTypes(dep)...)
				}
			}
		}

		var totalFiles []*FileInfo[Syntax]
		for _, p := range pkgs {
			files := NewPackageInspector(p, traslator.InspectSyntax(p, instTypes)).Inspect()
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

					err = traslator.Generate(info, file)
					if err != nil {
						fmt.Println("generate code of", info.Path, "failed:", err)
						return
					}

					if err = FormatCode(ctx, newFile, tmpBuildTag); err != nil {
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
						modifiedFiles[info.Path] = struct{}{}
					} else {
						buf := bytes.NewBuffer(nil)
						if info.BuildFlag == nil {
							buf.Write([]byte(GenTmpBuildFlags(false)))
						}
						bs, err := os.ReadFile(info.Path)
						if err != nil {
							fmt.Println("read file", info.Path, "failed:", err)
							return
						}
						buf.Write(bytes.Replace(bs, []byte(prodBuildTag), []byte(tmpBuildTag), 1))
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
						modifiedFiles[info.Path] = struct{}{}
					}
				}()
			}
			totalFiles = append(totalFiles, files...)
		}
		if len(buildFlags) <= 0 || len(totalFiles) > 0 {
			finished = false
		}
		buildFlags = []string{"-tags=" + tmpBuildTag}
	}
	for f := range modifiedFiles {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("os.ReadFile(): %w", err)
		}
		content = bytes.Replace(content, []byte(tmpBuildTag), []byte(prodBuildTag), 1)
		if err := os.WriteFile(f, content, 0644); err != nil {
			return fmt.Errorf("os.WriteFile() %w", err)
		}
	}
	return nil
}
