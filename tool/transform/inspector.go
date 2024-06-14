package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

const (
	prodBuildTag = "sugar_production"
)

func BuildDirective(predicate bool) string {
	var op string
	if !predicate {
		op = "!"
	}
	return fmt.Sprintf("//go:build %s%s", op, prodBuildTag)
}

type SyntaxInspector[Syntax fmt.Stringer] interface {
	Nodes() []ast.Node
	InspectSyntax(node ast.Node, stack []ast.Node) *Syntax
}

type Extent struct {
	Start token.Position
	End   token.Position
}

type PackageInspector[Syntax fmt.Stringer] struct {
	pkg           *packages.Package
	imports       map[string]map[string]string
	importExtents map[string]Extent
	buildFlags    map[string]Extent

	inspector SyntaxInspector[Syntax]
}

func NewPackageInspector[Syntax fmt.Stringer](pkg *packages.Package, inpector SyntaxInspector[Syntax]) *PackageInspector[Syntax] {
	imports := make(map[string]map[string]string)
	importExtents := make(map[string]Extent)
	buildFlags := make(map[string]Extent)
	for _, file := range pkg.Syntax {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				if strings.Contains(c.Text, BuildDirective(false)) || strings.Contains(c.Text, BuildDirective(true)) {
					fileName := pkg.Fset.Position(c.Pos()).Filename
					buildFlags[fileName] = Extent{
						Start: pkg.Fset.Position(c.Pos()),
						End:   pkg.Fset.Position(c.End()),
					}
				}
			}
		}
		specs := make(map[string]string)
		pkgEndPos := pkg.Fset.Position(file.Name.End())
		importExtents[pkgEndPos.Filename] = Extent{
			Start: pkgEndPos,
			End:   pkgEndPos,
		}
		for _, spec := range file.Imports {
			importPath := strings.Trim(spec.Path.Value, "\"")
			importName := path.Base(importPath)
			if spec.Name != nil {
				importName = spec.Name.Name
				if importName == "." {
					importName = ""
				}
			}
			specs[importPath] = importName
		}
		fileName := pkg.Fset.Position(file.Pos()).Filename
		imports[fileName] = specs
	}
	return &PackageInspector[Syntax]{
		pkg:           pkg,
		imports:       imports,
		importExtents: importExtents,
		buildFlags:    buildFlags,
		inspector:     inpector,
	}
}

type FileInfo[Syntax fmt.Stringer] struct {
	Path         string
	BuildFlag    *Extent
	PkgPath      string
	Imports      map[string]string
	ImportExtent Extent
	Syntax       []*Syntax
}

func (i *PackageInspector[Syntax]) Inspect() []*FileInfo[Syntax] {
	ins := inspector.New(i.pkg.Syntax)
	ins.Nodes([]ast.Node{
		&ast.GenDecl{},
		&ast.CallExpr{},
	}, func(n ast.Node, _ bool) bool {
		node, ok := n.(*ast.GenDecl)
		if ok && node.Tok == token.IMPORT {
			end := i.pkg.Fset.Position(node.End())
			extent := i.importExtents[end.Filename]
			if extent.End.Offset < end.Offset {
				extent.End = end
				i.importExtents[end.Filename] = extent
			}
			return false
		}
		return true
	})
	fileMap := make(map[string]*FileInfo[Syntax])
	ins.WithStack(i.inspector.Nodes(),
		func(node ast.Node, _ bool, stack []ast.Node) bool {
			syntax := i.inspector.InspectSyntax(node, stack)
			if syntax != nil {
				fileName := i.pkg.Fset.Position(node.Pos()).Filename
				file := fileMap[fileName]
				if file == nil {
					file = &FileInfo[Syntax]{
						Path:         fileName,
						PkgPath:      i.pkg.PkgPath,
						Imports:      i.imports[fileName],
						ImportExtent: i.importExtents[fileName],
					}
					fileMap[fileName] = file
				}
				file.Syntax = append(file.Syntax, syntax)
				return false
			}
			return true
		})
	var ret []*FileInfo[Syntax]
	for _, f := range fileMap {
		if extent, ok := i.buildFlags[f.Path]; ok {
			f.BuildFlag = &extent
		}
		ret = append(ret, f)
	}
	return ret
}
