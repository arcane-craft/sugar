package question

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"
	"strings"

	"github.com/arcane-craft/sugar/tool/transform/lib"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

const (
	questionPkgPath = "github.com/arcane-craft/sugar/syntax/question"
	questionIface   = "Question"
	questionFun     = "Q"
	resultPkgPath   = "github.com/arcane-craft/sugar/result"
	optionPkgPath   = "github.com/arcane-craft/sugar/option"
	stdFmtPkgPath   = "fmt"
)

type QuestionTypeInspector struct {
	pkg *packages.Package
}

func NewQuestionTypeInspector(pkg *packages.Package) *QuestionTypeInspector {
	return &QuestionTypeInspector{pkg}
}

type QuestionInstanceType struct {
	Name string
}

func (m *QuestionInstanceType) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Name: %s", m.Name)
}

func (i *QuestionTypeInspector) findQuestionIfaceType(ident *ast.Ident) *QuestionInstanceType {
	if ifaceType := i.pkg.TypesInfo.TypeOf(ident); ifaceType != nil {
		return &QuestionInstanceType{
			Name: lib.GetPkgPathFromType(ifaceType) + "." + ident.Name,
		}
	}
	return nil
}

func (i *QuestionTypeInspector) findQuestionEmbed(methodType ast.Expr, target *ast.Ident) *QuestionInstanceType {
	indexExpr, ok := methodType.(*ast.IndexExpr)
	if ok {
		switch node := indexExpr.X.(type) {
		case *ast.SelectorExpr:
			if ok && node.Sel.Name == questionIface {
				if obj := i.pkg.TypesInfo.ObjectOf(node.Sel); obj != nil {
					if objPkg := obj.Pkg(); objPkg != nil &&
						objPkg.Path() == questionPkgPath {
						return i.findQuestionIfaceType(target)
					}
				}
			}
		case *ast.Ident:
			if ok && node.Name == questionIface {
				if obj := i.pkg.TypesInfo.ObjectOf(node); obj != nil {
					if objPkg := obj.Pkg(); objPkg != nil &&
						objPkg.Path() == questionPkgPath {
						return i.findQuestionIfaceType(target)
					}
				}
			}
		}
	}
	return nil
}

func (i *QuestionTypeInspector) Inspect() []*QuestionInstanceType {
	var ret []*QuestionInstanceType
	ins := inspector.New(i.pkg.Syntax)
	ins.Nodes([]ast.Node{
		&ast.GenDecl{},
	}, func(n ast.Node, _ bool) bool {
		decl := n.(*ast.GenDecl)
		if decl.Tok.String() == ast.Typ.String() {
			for _, spec := range decl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if ok && ifaceType.Methods != nil {
						for _, method := range ifaceType.Methods.List {
							instType := i.findQuestionEmbed(method.Type, typeSpec.Name)
							if instType != nil {
								ret = append(ret, instType)
								break
							}
						}
					}
				}
			}
		}
		return false
	})
	return ret
}

type QuestionSyntaxInspector struct {
	pkg           *packages.Package
	instanceTypes map[string]*QuestionInstanceType
}

func NewQuestionSyntaxInspector(pkg *packages.Package, instances []*QuestionInstanceType) *QuestionSyntaxInspector {
	instanceTypes := make(map[string]*QuestionInstanceType)
	for _, inst := range instances {
		instanceTypes[inst.Name] = inst
	}
	return &QuestionSyntaxInspector{
		pkg:           pkg,
		instanceTypes: instanceTypes,
	}
}

type QuestionCall struct {
	lib.Extent
	Expr        *lib.Extent
	AssignVar   *lib.Extent
	AssignToken string
	ExprType    string
	OuterStmt   *lib.Extent
}

type QImplType struct {
	MainType  string
	InnerType string
}

type QuestionSyntax struct {
	Call    *QuestionCall
	OuterFn string
	RetType *QImplType
}

func (m *QuestionSyntax) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Call: %+v, OuterFn: %s, RetType: %+v", m.Call, m.OuterFn, m.RetType)
}

func (i *QuestionSyntaxInspector) Nodes() []ast.Node {
	return []ast.Node{
		&ast.CallExpr{},
	}
}

func (i *QuestionSyntaxInspector) isQuestionMethod(sel *ast.SelectorExpr) bool {
	if sel.Sel.Name == questionFun {
		instanceType := i.pkg.TypesInfo.TypeOf(sel.X).String()
		if _, ok := i.instanceTypes[lib.GetNameFromTypeStr(instanceType)]; ok {
			return true
		}
	}
	return false
}

func (i *QuestionSyntaxInspector) findQuestionFnCall(expr ast.Expr) (*lib.Extent, string) {
	call, ok := expr.(*ast.CallExpr)
	if ok && len(call.Args) <= 0 {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok && i.isQuestionMethod(sel) {
			return &lib.Extent{
				Start: i.pkg.Fset.Position(sel.X.Pos()),
				End:   i.pkg.Fset.Position(sel.X.End()),
			}, i.pkg.TypesInfo.TypeOf(sel.X).String()
		}
	}
	return nil, ""
}

func (i *QuestionSyntaxInspector) queryFuncRetType(t ast.Expr) *QImplType {
	expr, ok := t.(*ast.IndexExpr)
	if ok {
		return &QImplType{
			MainType:  lib.GetNameFromType(i.pkg.TypesInfo.TypeOf(expr)),
			InnerType: i.pkg.TypesInfo.TypeOf(expr.Index).String(),
		}
	}
	return nil
}

func (i *QuestionSyntaxInspector) queryFuncType(typ *ast.FuncType) *QImplType {
	if typ.Results != nil && typ.Results.NumFields() == 1 {
		retType := i.queryFuncRetType(typ.Results.List[0].Type)
		if retType != nil {
			if _, ok := i.instanceTypes[retType.MainType]; ok {
				return retType
			}
		}
	}
	return nil
}

func (i *QuestionSyntaxInspector) Inspect(n ast.Node, stack []ast.Node) (syntax *QuestionSyntax) {
	callExpr := n.(*ast.CallExpr)
	exprExt, exprType := i.findQuestionFnCall(callExpr)
	call := &QuestionCall{
		Expr:     exprExt,
		ExprType: exprType,
	}
	if exprExt != nil && len(stack) > 1 {
		var retType *QImplType
		var outerFn string
	FindOuterFunc:
		for idx := len(stack) - 2; idx >= 0; idx-- {
			switch fn := stack[idx].(type) {
			case *ast.FuncLit:
				retType = i.queryFuncType(fn.Type)
				outerFn = i.pkg.TypesInfo.TypeOf(fn).String()
				break FindOuterFunc
			case *ast.FuncDecl:
				retType = i.queryFuncType(fn.Type)
				outerFn = strings.Replace(i.pkg.TypesInfo.TypeOf(fn.Name).String(), "func", "func "+fn.Name.Name+"", 1)
				break FindOuterFunc
			}
		}
		if retType != nil && retType.MainType == lib.GetNameFromTypeStr(exprType) {
			if len(stack) > 2 {
				for idx := len(stack) - 2; idx >= 0; idx-- {
					stmt, ok := stack[idx].(ast.Stmt)
					if ok {
						if len(stack) > 3 {
							switch parent := stack[idx-1].(type) {
							case *ast.IfStmt:
								if parent.Init == stmt {
									call.OuterStmt = &lib.Extent{
										Start: i.pkg.Fset.Position(parent.Pos()),
										End:   i.pkg.Fset.Position(parent.End()),
									}
								}
							case *ast.SwitchStmt:
								if parent.Init == stmt {
									call.OuterStmt = &lib.Extent{
										Start: i.pkg.Fset.Position(parent.Pos()),
										End:   i.pkg.Fset.Position(parent.End()),
									}
								}
							case *ast.TypeSwitchStmt:
								if parent.Init == stmt {
									call.OuterStmt = &lib.Extent{
										Start: i.pkg.Fset.Position(parent.Pos()),
										End:   i.pkg.Fset.Position(parent.End()),
									}
								}
							case *ast.ForStmt:
								if parent.Init == stmt {
									call.OuterStmt = &lib.Extent{
										Start: i.pkg.Fset.Position(parent.Pos()),
										End:   i.pkg.Fset.Position(parent.End()),
									}
								}
							default:
							}
							if call.OuterStmt != nil {
								call.Extent = lib.Extent{
									Start: i.pkg.Fset.Position(callExpr.Pos()),
									End:   i.pkg.Fset.Position(callExpr.End()),
								}
								break
							}
						}
						switch current := stmt.(type) {
						case *ast.AssignStmt:
							for idx, rhs := range current.Rhs {
								rhsCall, ok := rhs.(*ast.CallExpr)
								if ok && rhsCall == callExpr {
									call.Extent = lib.Extent{
										Start: i.pkg.Fset.Position(current.Pos()),
										End:   i.pkg.Fset.Position(current.End()),
									}
									call.AssignVar = &lib.Extent{
										Start: i.pkg.Fset.Position(current.Lhs[idx].Pos()),
										End:   i.pkg.Fset.Position(current.Lhs[idx].End()),
									}
									call.AssignToken = current.Tok.String()
									break
								}
							}
						case *ast.ExprStmt:
							if current.X == ast.Expr(callExpr) {
								call.Extent = lib.Extent{
									Start: i.pkg.Fset.Position(current.Pos()),
									End:   i.pkg.Fset.Position(current.End()),
								}
							}
						default:
						}
						if call.Extent.IsEmpty() {
							call.Extent = lib.Extent{
								Start: i.pkg.Fset.Position(callExpr.Pos()),
								End:   i.pkg.Fset.Position(callExpr.End()),
							}
							call.OuterStmt = &lib.Extent{
								Start: i.pkg.Fset.Position(stmt.Pos()),
								End:   i.pkg.Fset.Position(stmt.End()),
							}
						}
						break
					}
				}
			}
			syntax = &QuestionSyntax{
				Call:    call,
				OuterFn: outerFn,
				RetType: retType,
			}
		}
	}
	return
}

func GenAssginStmt(assignVar, assignToken, call string) string {
	return fmt.Sprintf("%s %s %s\n", assignVar, assignToken, call)
}

func GenErrorHandler(resultPkg, fmtPkg, resultVar, retType, outerFn string) string {
	if len(resultPkg) > 0 {
		resultPkg += "."
	}
	if len(fmtPkg) > 0 {
		fmtPkg += "."
	}
	return fmt.Sprintf("if %s.IsErr() {\nreturn %sErr[%s](%sErrorf(\"%s: %%w\", %s.UnwrapErr()))\n}\n",
		resultVar, resultPkg, retType, fmtPkg, outerFn, resultVar)
}

func GenUnwrapStmt(assignVar, assignToken, receiverVar string) string {
	return fmt.Sprintf("%s %s %s.Unwrap()\n", assignVar, assignToken, receiverVar)
}

func GenUnwrapExpr(receiverVar string) string {
	return fmt.Sprintf("%s.Unwrap()", receiverVar)
}

func GenNoneHandler(optionPkg, optionVar, retType string) string {
	if len(optionPkg) > 0 {
		optionPkg += "."
	}
	return fmt.Sprintf("if %s.IsNone() {\nreturn %sNone[%s]()\n}\n", optionVar, optionPkg, retType)
}

func GenerateQuestionSyntax(info *lib.FileInfo[*QuestionSyntax], writer io.Writer) error {
	return lib.GenerateSyntax(info, writer, func(file *os.File, addImports map[string]string) ([]*lib.ReplaceBlock, error) {
		var ret []*lib.ReplaceBlock
		for _, syntax := range info.Syntax {
			if syntax.Call.AssignVar != nil {
				assignVar, err := lib.ReadExtent(file, syntax.Call.AssignVar)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				callExpr, err := lib.ReadExtent(file, syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				receiverVar := lib.GenVarName("var", syntax.Call.Expr.String())
				retType, adds := lib.ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(receiverVar, token.DEFINE.String(), callExpr)
				if strings.HasSuffix(lib.GetNameFromTypeStr(syntax.Call.ExprType), "Result") {
					resultPkg, ok := info.Imports[resultPkgPath]
					if !ok {
						resultPkg = lib.GenPkgName(path.Base(resultPkgPath), resultPkgPath)
						addImports[resultPkgPath] = resultPkg
					}
					fmtPkg, ok := info.Imports[stdFmtPkgPath]
					if !ok {
						fmtPkg = lib.GenPkgName(path.Base(stdFmtPkgPath), stdFmtPkgPath)
						addImports[stdFmtPkgPath] = fmtPkg
					}
					result += GenErrorHandler(resultPkg, fmtPkg, receiverVar, retType, syntax.OuterFn)
				} else {
					optionPkg, ok := info.Imports[optionPkgPath]
					if !ok {
						optionPkg = lib.GenPkgName(path.Base(optionPkgPath), optionPkgPath)
						addImports[optionPkgPath] = optionPkg
					}
					result += GenNoneHandler(optionPkg, receiverVar, retType)
				}
				result += GenUnwrapStmt(assignVar, syntax.Call.AssignToken, receiverVar)
				ret = append(ret, &lib.ReplaceBlock{
					Old: syntax.Call.Extent,
					New: result,
				})
			} else if syntax.Call.OuterStmt != nil {
				callExpr, err := lib.ReadExtent(file, syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				receiverVar := lib.GenVarName("var", syntax.Call.Expr.String())
				retType, adds := lib.ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(receiverVar, token.DEFINE.String(), callExpr)
				if strings.HasSuffix(lib.GetNameFromTypeStr(syntax.Call.ExprType), "Result") {
					resultPkg, ok := info.Imports[resultPkgPath]
					if !ok {
						resultPkg = lib.GenPkgName(path.Base(resultPkgPath), resultPkgPath)
						addImports[resultPkgPath] = resultPkg
					}
					fmtPkg, ok := info.Imports[stdFmtPkgPath]
					if !ok {
						fmtPkg = lib.GenPkgName(path.Base(stdFmtPkgPath), stdFmtPkgPath)
						addImports[stdFmtPkgPath] = fmtPkg
					}
					result += GenErrorHandler(resultPkg, fmtPkg, receiverVar, retType, syntax.OuterFn)
				} else {
					optionPkg, ok := info.Imports[optionPkgPath]
					if !ok {
						optionPkg = lib.GenPkgName(path.Base(optionPkgPath), optionPkgPath)
						addImports[optionPkgPath] = optionPkg
					}
					result += GenNoneHandler(optionPkg, receiverVar, retType)
				}
				ret = append(ret,
					&lib.ReplaceBlock{
						Old: lib.Extent{
							Start: syntax.Call.OuterStmt.Start,
							End:   syntax.Call.OuterStmt.Start,
						},
						New: result,
					},
					&lib.ReplaceBlock{
						Old: syntax.Call.Extent,
						New: GenUnwrapExpr(receiverVar),
					})
			} else {
				callExpr, err := lib.ReadExtent(file, syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				receiverVar := lib.GenVarName("var", syntax.Call.Expr.String())
				retType, adds := lib.ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(receiverVar, token.DEFINE.String(), callExpr)
				if strings.HasSuffix(lib.GetNameFromTypeStr(syntax.Call.ExprType), "Result") {
					resultPkg, ok := info.Imports[resultPkgPath]
					if !ok {
						resultPkg = lib.GenPkgName(path.Base(resultPkgPath), resultPkgPath)
						addImports[resultPkgPath] = resultPkg
					}
					fmtPkg, ok := info.Imports[stdFmtPkgPath]
					if !ok {
						fmtPkg = lib.GenPkgName(path.Base(stdFmtPkgPath), stdFmtPkgPath)
						addImports[stdFmtPkgPath] = fmtPkg
					}
					result += GenErrorHandler(resultPkg, fmtPkg, receiverVar, retType, syntax.OuterFn)
				} else {
					optionPkg, ok := info.Imports[optionPkgPath]
					if !ok {
						optionPkg = lib.GenPkgName(path.Base(optionPkgPath), optionPkgPath)
						addImports[optionPkgPath] = optionPkg
					}
					result += GenNoneHandler(optionPkg, receiverVar, retType)
				}
				ret = append(ret, &lib.ReplaceBlock{
					Old: syntax.Call.Extent,
					New: result,
				})
			}
		}
		return ret, nil
	})
}

type Translator struct{}

func (*Translator) InpectTypes(p *packages.Package) []*QuestionInstanceType {
	return NewQuestionTypeInspector(p).Inspect()
}

func (*Translator) InspectSyntax(p *packages.Package, instTypes []*QuestionInstanceType) lib.SyntaxInspector[*QuestionSyntax] {
	return NewQuestionSyntaxInspector(p, instTypes)
}

func (*Translator) Generate(info *lib.FileInfo[*QuestionSyntax], writer io.Writer) error {
	return GenerateQuestionSyntax(info, writer)
}

func (t *Translator) Run(ctx context.Context, rootDir string, firstRun bool) error {
	err := lib.TranslateSyntax(ctx, rootDir, firstRun, t)
	if err != nil {
		return fmt.Errorf("translate question syntax failed: %w", err)
	}
	return nil
}
