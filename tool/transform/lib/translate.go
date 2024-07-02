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

type Program interface {
	Run(ctx context.Context, rootDir string, firstRun bool) error
}

type SyntaxTranslator[Type, Syntax interface {
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
	translator SyntaxTranslator[Type, Syntax],
) error {

	var finished bool
	var buildTags []string

	for !finished {
		finished = true

		pkgs, err := LoadPackages(ctx, rootDir, buildTags...)
		if err != nil {
			return fmt.Errorf("load source packages failed: %w", err)
		}
		var instTypes []Type
		for _, p := range pkgs {
			instTypes = append(instTypes, translator.InpectTypes(p)...)
			for path, dep := range p.Imports {
				if p.PkgPath != path {
					instTypes = append(instTypes, translator.InpectTypes(dep)...)
				}
			}
		}

		var totalFiles []*FileInfo[Syntax]
		for _, p := range pkgs {
			files := NewPackageInspector(p, translator.InspectSyntax(p, instTypes)).Inspect()
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

					err = translator.Generate(info, file)
					if err != nil {
						fmt.Println("generate code of", info.Path, "failed:", err)
						return
					}

					if err = FormatCode(ctx, newFile, tmpBuildTag); err != nil {
						fmt.Println("format code of", newFile, "failed:", err)
					}

					if len(buildTags) > 0 {
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
						if info.BuildTag == nil {
							buf.Write([]byte(GenTmpBuildTags(false)))
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
		if len(buildTags) <= 0 || len(totalFiles) > 0 {
			finished = false
		}
		buildTags = []string{"-tags=" + tmpBuildTag}
	}
	return nil
}

func changeBuildTags(ctx context.Context, rootDir string, old string, new string) error {
	var finished bool
	var buildTags []string

	for !finished {
		finished = true

		pkgs, err := LoadPackages(ctx, rootDir, buildTags...)
		if err != nil {
			return fmt.Errorf("load source packages failed: %w", err)
		}
		for _, p := range pkgs {
			if err := replaceBuildTags(FindPackageBuildTags(p), old, new); err != nil {
				return fmt.Errorf("MakePackageTemp() failed: %w", err)
			}
		}

		if len(buildTags) <= 0 {
			buildTags = []string{"-tags=" + old}
			finished = false
		}
	}

	return nil
}

func SetTmpBuildTags(ctx context.Context, rootDir string) error {
	return changeBuildTags(ctx, rootDir, prodBuildTag, tmpBuildTag)
}

func SetProdBuildTags(ctx context.Context, rootDir string) error {
	return changeBuildTags(ctx, rootDir, tmpBuildTag, prodBuildTag)
}
