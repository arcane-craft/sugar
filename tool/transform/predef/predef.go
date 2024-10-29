package predef

import (
	"context"
	"fmt"
	"go/ast"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/arcane-craft/sugar/tool/transform/lib"
	"golang.org/x/tools/go/packages"
)

const (
	predefPkgPath = "github.com/arcane-craft/sugar/syntax/predef"

	identFunction       = "Function__"
	identPrettyFunction = "PrettyFunction__"
	identPackage        = "Package__"
	identFile           = "File__"
	identLine           = "Line__"
)

type Identifier struct {
	*lib.Extent
	Name  string
	Value string
}

func (i *Identifier) String() string {
	return fmt.Sprintf("%s, Name: %s, Value: %s", i.Extent, i.Name, i.Value)
}

type PredefSyntax struct {
	Ident *Identifier
}

func (s *PredefSyntax) String() string {
	return fmt.Sprintf("Ident: %s", s.Ident)
}

type SyntaxInspector struct {
	pkg *packages.Package
}

func NewSyntaxInspector(pkg *packages.Package) *SyntaxInspector {
	return &SyntaxInspector{
		pkg: pkg,
	}
}

func (i *SyntaxInspector) Nodes() []ast.Node {
	return []ast.Node{
		new(ast.Ident),
	}
}

func (i *SyntaxInspector) findNearestFunc(stack []ast.Node) ast.Node {
	if len(stack) > 1 {
		for i := len(stack) - 2; i >= 0; i-- {
			switch fun := stack[i].(type) {
			case *ast.FuncDecl:
				return fun
			case *ast.FuncLit:
				return fun
			}
		}
	}
	return nil
}

func (i *SyntaxInspector) Inspect(node ast.Node, stack []ast.Node) *PredefSyntax {
	ident := node.(*ast.Ident)
	object := i.pkg.TypesInfo.ObjectOf(ident)
	if object != nil && object.Pkg() != nil && object.Pkg().Path() == predefPkgPath {
		identityNode := ast.Expr(ident)
		if len(stack) > 1 {
			if _, ok := stack[len(stack)-2].(*ast.ValueSpec); ok {
				return nil
			}
			selector, ok := stack[len(stack)-2].(*ast.SelectorExpr)
			if ok {
				identityNode = selector
			}
		}
		var identifier *Identifier
		switch ident.Name {
		case identFunction:
			var funcName string
			if expr := i.findNearestFunc(stack); expr != nil {
				funcDecl, ok := expr.(*ast.FuncDecl)
				if ok {
					funcName = funcDecl.Name.Name
				}
			}
			identifier = &Identifier{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(identityNode.Pos()),
					End:   i.pkg.Fset.Position(identityNode.End()),
				},
				Name:  ident.Name,
				Value: funcName,
			}
		case identPrettyFunction:
			var funcType string
			if expr := i.findNearestFunc(stack); expr != nil {
				switch fun := expr.(type) {
				case *ast.FuncDecl:
					funcType = i.pkg.TypesInfo.TypeOf(fun.Name).String()
					funcType = strings.Replace(funcType, "func", "func "+fun.Name.Name, 1)
				case *ast.FuncLit:
					funcType = i.pkg.TypesInfo.TypeOf(fun).String()
				}
			}
			identifier = &Identifier{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(identityNode.Pos()),
					End:   i.pkg.Fset.Position(identityNode.End()),
				},
				Name:  ident.Name,
				Value: funcType,
			}
		case identPackage:
			identifier = &Identifier{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(identityNode.Pos()),
					End:   i.pkg.Fset.Position(identityNode.End()),
				},
				Name:  ident.Name,
				Value: i.pkg.PkgPath,
			}
		case identFile:
			identifier = &Identifier{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(identityNode.Pos()),
					End:   i.pkg.Fset.Position(identityNode.End()),
				},
				Name:  ident.Name,
				Value: i.pkg.Fset.Position(identityNode.Pos()).Filename,
			}
		case identLine:
			identifier = &Identifier{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(identityNode.Pos()),
					End:   i.pkg.Fset.Position(identityNode.End()),
				},
				Name:  ident.Name,
				Value: strconv.Itoa(i.pkg.Fset.Position(identityNode.Pos()).Line),
			}
		}
		if identifier != nil {
			return &PredefSyntax{
				Ident: identifier,
			}
		}
	}

	return nil
}

type Translator struct{}

func (*Translator) InpectTypes(p *packages.Package) []*lib.Extent {
	return nil
}

func (*Translator) InspectSyntax(p *packages.Package, _ []*lib.Extent) lib.SyntaxInspector[*PredefSyntax] {
	return NewSyntaxInspector(p)
}

func genStringLit(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func (*Translator) Generate(info *lib.FileInfo[*PredefSyntax], writer io.Writer) error {
	return lib.GenerateSyntax(info, writer, func(file *os.File, addImports map[string]string) ([]*lib.ReplaceBlock, error) {
		var blocks []*lib.ReplaceBlock
		for _, syntax := range info.Syntax {
			blocks = append(blocks, &lib.ReplaceBlock{
				Old: *syntax.Ident.Extent,
				New: genStringLit(syntax.Ident.Value),
			})
		}
		for path := range addImports {
			if path == predefPkgPath {
				delete(addImports, path)
			}
		}
		return blocks, nil
	})
}

func (t *Translator) Run(ctx context.Context, rootDir string, firstRun bool) error {
	err := lib.TranslateSyntax(ctx, rootDir, firstRun, t)
	if err != nil {
		return fmt.Errorf("translate predefine syntax failed: %w", err)
	}
	return nil
}
