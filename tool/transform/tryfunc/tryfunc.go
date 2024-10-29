package tryfunc

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/arcane-craft/sugar/tool/transform/lib"
	"golang.org/x/tools/go/packages"
)

const (
	errorTypeName  = "error"
	tryFuncPkgPath = "github.com/arcane-craft/sugar/syntax/tryfunc"
	tryFuncName    = "Try"
	stdFmtPkgPath  = "fmt"
)

type TryStmt struct {
	*lib.Extent
	RetVars     []*lib.Extent
	AssignToken string
	CallExpr    *lib.Extent
	OuterStmt   *lib.Extent
}

func (t *TryStmt) String() string {
	return fmt.Sprintf("%s, RetVars: %s, AssignToken: %s, CallExpr: %s, OuterStmt: %s",
		t.Extent, lib.JoinStringers(t.RetVars, ","), t.AssignToken, t.CallExpr, t.OuterStmt)
}

type FuncResult struct {
	*lib.Extent
	Name string
	Type *lib.Extent
}

func (f *FuncResult) String() string {
	return fmt.Sprintf("%s, Name: %s, Type: %s", f.Extent, f.Name, f.Type)
}

type TrySyntax struct {
	*lib.Extent
	OuterFunc string
	Stmts     []*TryStmt
	Results   []*FuncResult
}

func (t *TrySyntax) String() string {
	return fmt.Sprintf("%s, Stmts: %s, Results: %s",
		t.Extent, lib.JoinStringers(t.Stmts, ","), lib.JoinStringers(t.Results, ","))
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
		new(ast.FuncDecl),
		new(ast.FuncLit),
	}
}

func (i *SyntaxInspector) findFuncResults(typ *ast.FuncType) (results []*FuncResult) {
	resultTypes := typ.Results
	if resultTypes.NumFields() > 0 {
		ident, ok := resultTypes.List[len(resultTypes.List)-1].Type.(*ast.Ident)
		if ok && ident.Name == errorTypeName {
			for idx, elem := range resultTypes.List {
				if len(elem.Names) > 0 {
					for _, name := range elem.Names {
						results = append(results, &FuncResult{
							Extent: &lib.Extent{
								Start: i.pkg.Fset.Position(elem.Pos()),
								End:   i.pkg.Fset.Position(elem.End()),
							},
							Name: name.Name,
							Type: &lib.Extent{
								Start: i.pkg.Fset.Position(elem.Type.Pos()),
								End:   i.pkg.Fset.Position(elem.Type.End()),
							},
						})
					}
				} else {
					resultName := "_"
					if idx == len(resultTypes.List)-1 {
						resultName = lib.GenVarName("err", i.pkg.Fset.Position(elem.Type.Pos()).String())
					}
					results = append(results, &FuncResult{
						Extent: &lib.Extent{
							Start: i.pkg.Fset.Position(elem.Pos()),
							End:   i.pkg.Fset.Position(elem.End()),
						},
						Name: resultName,
						Type: &lib.Extent{
							Start: i.pkg.Fset.Position(elem.Type.Pos()),
							End:   i.pkg.Fset.Position(elem.Type.End()),
						},
					})
				}
			}
		}
	}
	return
}

func (i *SyntaxInspector) getFuncIdent(expr ast.Expr) *ast.Ident {
	switch fun := expr.(type) {
	case *ast.Ident:
		return fun
	case *ast.SelectorExpr:
		return i.getFuncIdent(fun.Sel)
	case *ast.IndexExpr:
		return i.getFuncIdent(fun.X)
	}
	return nil
}

func (i *SyntaxInspector) inspectTryFuncCall(expr ast.Expr) (ast.Expr, int, bool) {
	if call, ok := expr.(*ast.CallExpr); ok {
		funcIdent := i.getFuncIdent(call.Fun)
		if funcIdent != nil && strings.HasPrefix(funcIdent.Name, tryFuncName) {
			object := i.pkg.TypesInfo.ObjectOf(funcIdent)
			if object != nil &&
				object.Pkg() != nil &&
				object.Pkg().Path() == tryFuncPkgPath &&
				len(call.Args) == 1 {
				var retNum int
				switch strings.TrimPrefix(funcIdent.Name, tryFuncName) {
				case "_":
					retNum = 0
				case "":
					retNum = 1
				case "2":
					retNum = 2
				case "3":
					retNum = 3
				}
				return call.Args[0], retNum, true
			}
		}
	}
	return nil, 0, false
}

func (i *SyntaxInspector) findTryStmt(node ast.Stmt, outer ast.Stmt) (ret *TryStmt) {
	switch stmt := node.(type) {
	case *ast.ExprStmt:
		originCall, retNum, ok := i.inspectTryFuncCall(stmt.X)
		if ok {
			ret = &TryStmt{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(stmt.Pos()),
					End:   i.pkg.Fset.Position(stmt.End()),
				},
				CallExpr: &lib.Extent{
					Start: i.pkg.Fset.Position(originCall.Pos()),
					End:   i.pkg.Fset.Position(originCall.End()),
				},
				AssignToken: token.ASSIGN.String(),
			}
			for idx := 0; idx < retNum; idx++ {
				ret.RetVars = append(ret.RetVars, nil)
			}
		}
	case *ast.AssignStmt:
		if len(stmt.Rhs) == 1 {
			originCall, _, ok := i.inspectTryFuncCall(stmt.Rhs[0])
			if ok {
				ret = &TryStmt{
					Extent: &lib.Extent{
						Start: i.pkg.Fset.Position(stmt.Pos()),
						End:   i.pkg.Fset.Position(stmt.End()),
					},
					CallExpr: &lib.Extent{
						Start: i.pkg.Fset.Position(originCall.Pos()),
						End:   i.pkg.Fset.Position(originCall.End()),
					},
					AssignToken: stmt.Tok.String(),
				}
				for _, l := range stmt.Lhs {
					ret.RetVars = append(ret.RetVars, &lib.Extent{
						Start: i.pkg.Fset.Position(l.Pos()),
						End:   i.pkg.Fset.Position(l.End()),
					})
				}
			}
		}
	}

	if ret != nil && outer != nil {
		var hasOuterStmt bool
		switch parent := outer.(type) {
		case *ast.IfStmt:
			hasOuterStmt = parent.Init == node
		case *ast.SwitchStmt:
			hasOuterStmt = parent.Init == node
		case *ast.TypeSwitchStmt:
			hasOuterStmt = parent.Init == node
		case *ast.ForStmt:
			hasOuterStmt = parent.Init == node
		}
		if hasOuterStmt {
			ret.OuterStmt = &lib.Extent{
				Start: i.pkg.Fset.Position(outer.Pos()),
				End:   i.pkg.Fset.Position(outer.End()),
			}
		}
	}
	return
}

func (i *SyntaxInspector) inspectFunc(typ *ast.FuncType, body *ast.BlockStmt) *TrySyntax {
	funcResults := i.findFuncResults(typ)
	if len(funcResults) > 0 {
		var outerStmt ast.Stmt
		var calls []*TryStmt
		ast.Inspect(body, func(child ast.Node) bool {
			if _, ok := child.(*ast.FuncLit); ok {
				return false
			}
			stmt, ok := child.(ast.Stmt)
			if !ok {
				return true
			}
			funcStmt := i.findTryStmt(stmt, outerStmt)
			if funcStmt != nil {
				calls = append(calls, funcStmt)
				outerStmt = nil
				return false
			}
			outerStmt = stmt
			return true
		})
		if len(calls) > 0 {
			return &TrySyntax{
				Extent: &lib.Extent{
					Start: i.pkg.Fset.Position(body.Pos() + 1),
					End:   i.pkg.Fset.Position(body.End() - 1),
				},
				Stmts:   calls,
				Results: funcResults,
			}
		}
	}
	return nil
}

func (i *SyntaxInspector) Inspect(node ast.Node, _ []ast.Node) *TrySyntax {
	var syntax *TrySyntax
	var outerFunc string
	switch fun := node.(type) {
	case *ast.FuncDecl:
		syntax = i.inspectFunc(fun.Type, fun.Body)
		outerFunc = i.pkg.TypesInfo.TypeOf(fun.Name).String()
		outerFunc = strings.Replace(outerFunc, "func", fmt.Sprintf("func %s", fun.Name), 1)
	case *ast.FuncLit:
		syntax = i.inspectFunc(fun.Type, fun.Body)
		outerFunc = i.pkg.TypesInfo.TypeOf(fun).String()
	}
	if syntax != nil {
		syntax.OuterFunc = outerFunc
	}
	return syntax
}

type Translator struct{}

func (*Translator) InpectTypes(p *packages.Package) []*lib.Extent {
	return nil
}

func (*Translator) InspectSyntax(p *packages.Package, _ []*lib.Extent) lib.SyntaxInspector[*TrySyntax] {
	return NewSyntaxInspector(p)
}

func genFuncResultTypeElem(name, typ string) string {
	return fmt.Sprintf("%s %s", name, typ)
}

func genFuncResultType(elems []string) string {
	return strings.Join(elems, ", ")
}

func genAssignStmt(lhs []string, token, callExpr string) string {
	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), token, callExpr)
}

func genErrHander(errVar string, retErrVar string) string {
	var stmts []string
	if errVar != retErrVar {
		stmts = append(stmts, fmt.Sprintf("%s = %s", retErrVar, errVar))
	}
	stmts = append(stmts, "return")
	return fmt.Sprintf("if %s != nil {\n %s \n}", errVar, strings.Join(stmts, "\n"))
}

func genErrWraper(retErrVar, fmtPkg, outerFunc string) string {
	if len(fmtPkg) > 0 {
		fmtPkg += "."
	}
	return fmt.Sprintf("defer func() {\n if %s != nil {\n %s = %sErrorf(\"%s: %%w\", %s) \n} \n}()", retErrVar, retErrVar, fmtPkg, outerFunc, retErrVar)
}

func (*Translator) Generate(info *lib.FileInfo[*TrySyntax], writer io.Writer) error {
	return lib.GenerateSyntax(info, writer, func(file *os.File, addImports map[string]string) ([]*lib.ReplaceBlock, error) {
		var blocks []*lib.ReplaceBlock
		for _, syntax := range info.Syntax {
			var resultTypeElems []string
			var resultStart, resultEnd token.Position
			var retErrVar string
			for idx, ret := range syntax.Results {
				elemType, err := lib.ReadExtent(file, ret.Type)
				if err != nil {
					return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
				}
				if resultStart.Offset <= 0 || ret.Start.Offset < resultStart.Offset {
					resultStart = ret.Start
				}
				if ret.End.Offset > resultEnd.Offset {
					resultEnd = ret.End
				}
				resultTypeElems = append(resultTypeElems, genFuncResultTypeElem(ret.Name, elemType))
				if idx == len(syntax.Results)-1 {
					retErrVar = ret.Name
				}
			}
			blocks = append(blocks, &lib.ReplaceBlock{
				Old: lib.Extent{
					Start: resultStart,
					End:   resultEnd,
				},
				New: genFuncResultType(resultTypeElems),
			})

			for _, stmt := range syntax.Stmts {
				var lhs []string
				for _, v := range stmt.RetVars {
					if v == nil {
						lhs = append(lhs, "_")
					} else {
						ret, err := lib.ReadExtent(file, v)
						if err != nil {
							return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
						}
						lhs = append(lhs, ret)
					}
				}
				var errVar string
				if stmt.AssignToken == token.DEFINE.String() {
					errVar = lib.GenVarName("err", stmt.CallExpr.Start.String())
				} else {
					errVar = retErrVar
				}
				lhs = append(lhs, errVar)
				callExpr, err := lib.ReadExtent(file, stmt.CallExpr)
				if err != nil {
					return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
				}
				assigneStmt := genAssignStmt(lhs, stmt.AssignToken, callExpr)
				errHandler := genErrHander(errVar, retErrVar)

				if stmt.OuterStmt != nil {
					blocks = append(blocks, &lib.ReplaceBlock{
						Old: lib.Extent{
							Start: stmt.OuterStmt.Start,
							End:   stmt.OuterStmt.Start,
						},
						New: assigneStmt + "\n" + errHandler + "\n",
					}, &lib.ReplaceBlock{
						Old: *stmt.Extent,
						New: strings.Join(lhs[:len(lhs)-1], ", "),
					})
				} else {
					blocks = append(blocks, &lib.ReplaceBlock{
						Old: *stmt.Extent,
						New: assigneStmt + "\n" + errHandler,
					})
				}
			}
		}
		for path := range addImports {
			if path == tryFuncPkgPath {
				delete(addImports, path)
			}
		}
		return blocks, nil
	})
}

func (t *Translator) Run(ctx context.Context, rootDir string, firstRun bool) error {
	err := lib.TranslateSyntax(ctx, rootDir, firstRun, t)
	if err != nil {
		return fmt.Errorf("translate tryfunc syntax failed: %w", err)
	}
	return nil
}
