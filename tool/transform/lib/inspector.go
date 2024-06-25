package lib

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
	tmpBuildTag  = "sugar_temp"
)

func BuildDirective(predicate bool) string {
	var op string
	if !predicate {
		op = "!"
	}
	return fmt.Sprintf("//go:build %s%s", op, prodBuildTag)
}

func TmpBuildDirective(predicate bool) string {
	var op string
	if !predicate {
		op = "!"
	}
	return fmt.Sprintf("//go:build %s%s", op, tmpBuildTag)
}

type SyntaxInspector[Syntax interface {
	fmt.Stringer
	comparable
}] interface {
	Nodes() []ast.Node
	Inspect(node ast.Node, stack []ast.Node) Syntax
}

type Extent struct {
	Start token.Position
	End   token.Position
}

func (m *Extent) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s-%s", m.Start, m.End)
}

func (m *Extent) IsEmpty() bool {
	if m == nil {
		return true
	}
	var zero Extent
	return *m == zero
}

type PackageInspector[Syntax interface {
	fmt.Stringer
	comparable
}] struct {
	pkg           *packages.Package
	imports       map[string]map[string]string
	importExtents map[string]*Extent
	buildTags     map[string]*Extent

	inspector SyntaxInspector[Syntax]
}

func NewPackageInspector[Syntax interface {
	fmt.Stringer
	comparable
}](pkg *packages.Package, inpector SyntaxInspector[Syntax]) *PackageInspector[Syntax] {
	imports := make(map[string]map[string]string)
	importExtents := make(map[string]*Extent)
	buildTags := make(map[string]*Extent)
	for _, file := range pkg.Syntax {
		for k, v := range FindFileBuildTags(pkg, file) {
			buildTags[k] = v
		}
		specs := make(map[string]string)
		pkgEndPos := pkg.Fset.Position(file.Name.End())
		importExtents[pkgEndPos.Filename] = &Extent{
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
		buildTags:     buildTags,
		inspector:     inpector,
	}
}

type FileInfo[Syntax interface {
	fmt.Stringer
	comparable
}] struct {
	Path         string
	BuildTag     *Extent
	PkgPath      string
	Imports      map[string]string
	ImportExtent *Extent
	Syntax       []Syntax
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
			syntax := i.inspector.Inspect(node, stack)
			var zero Syntax
			if syntax != zero {
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
		if extent, ok := i.buildTags[f.Path]; ok {
			f.BuildTag = extent
		}
		ret = append(ret, f)
	}
	return ret
}

func FindFileBuildTags(pkg *packages.Package, file *ast.File) map[string]*Extent {
	buildTags := make(map[string]*Extent)
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, TmpBuildDirective(false)) ||
				strings.Contains(c.Text, BuildDirective(false)) ||
				strings.Contains(c.Text, TmpBuildDirective(true)) ||
				strings.Contains(c.Text, BuildDirective(true)) {
				fileName := pkg.Fset.Position(c.Pos()).Filename
				buildTags[fileName] = &Extent{
					Start: pkg.Fset.Position(c.Pos()),
					End:   pkg.Fset.Position(c.End()),
				}
			}
		}
	}
	return buildTags
}

func FindPackageBuildTags(pkg *packages.Package) map[string]*Extent {
	buildTags := make(map[string]*Extent)
	if strings.HasPrefix(pkg.PkgPath, "github.com/arcane-craft/sugar/syntax") {
		return buildTags
	}
	for _, file := range pkg.Syntax {
		for k, v := range FindFileBuildTags(pkg, file) {
			buildTags[k] = v
		}
	}
	return buildTags
}
