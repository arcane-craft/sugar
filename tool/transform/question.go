package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

const (
	sugarPkgPath  = "github.com/arcane-craft/sugar"
	questionIface = "Question"
	questionImpl  = "QuestionImpl"
	questionFun   = "Q"
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

func (m QuestionInstanceType) String() string {
	return fmt.Sprintf("Name: %s", m.Name)
}

func (i *QuestionTypeInspector) inspectQuestionIfaceType(ident *ast.Ident) *QuestionInstanceType {
	if ifaceType := i.pkg.TypesInfo.TypeOf(ident); ifaceType != nil {
		return &QuestionInstanceType{
			Name: GetPkgPathFromType(ifaceType) + "." + ident.Name,
		}
	}
	return nil
}

func (i *QuestionTypeInspector) inspectQuestionEmbed(methodType ast.Expr, target *ast.Ident) *QuestionInstanceType {
	indexExpr, ok := methodType.(*ast.IndexExpr)
	if ok {
		switch node := indexExpr.X.(type) {
		case *ast.SelectorExpr:
			if ok && node.Sel.Name == questionIface {
				if obj := i.pkg.TypesInfo.ObjectOf(node.Sel); obj != nil {
					if objPkg := obj.Pkg(); objPkg != nil &&
						objPkg.Path() == sugarPkgPath {
						return i.inspectQuestionIfaceType(target)
					}
				}
			}
		case *ast.Ident:
			if ok && node.Name == questionIface {
				if obj := i.pkg.TypesInfo.ObjectOf(node); obj != nil {
					if objPkg := obj.Pkg(); objPkg != nil &&
						objPkg.Path() == sugarPkgPath {
						return i.inspectQuestionIfaceType(target)
					}
				}
			}
		}
	}
	return nil
}

func (i *QuestionTypeInspector) InspectQuestionTypes() []*QuestionInstanceType {
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
							instType := i.inspectQuestionEmbed(method.Type, typeSpec.Name)
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
	Extent
	Expr        *Extent
	AssignVar   *Extent
	AssignToken string
	ExprType    string
	OuterStmt   *Extent
}

type QImplType struct {
	MainType  string
	InnerType string
}

type QuestionSyntax struct {
	Call    *QuestionCall
	RetType *QImplType
}

func (m QuestionSyntax) String() string {
	return fmt.Sprintf("Call: %+v, RetType: %+v", m.Call, m.RetType)
}

func (i *QuestionSyntaxInspector) Nodes() []ast.Node {
	return []ast.Node{
		&ast.CallExpr{},
	}
}

func (i *QuestionSyntaxInspector) isQuestionMethod(sel *ast.SelectorExpr) bool {
	if sel.Sel.Name == questionFun {
		instanceType := i.pkg.TypesInfo.TypeOf(sel.X).String()
		if _, ok := i.instanceTypes[GetNameFromTypeStr(instanceType)]; ok {
			return true
		}
	}
	return false
}

func (i *QuestionSyntaxInspector) inspectQuestionFnCall(expr ast.Expr) (*Extent, string) {
	call, ok := expr.(*ast.CallExpr)
	if ok && len(call.Args) <= 0 {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok && i.isQuestionMethod(sel) {
			return &Extent{
				Start: i.pkg.Fset.Position(sel.X.Pos()),
				End:   i.pkg.Fset.Position(sel.X.End()),
			}, i.pkg.TypesInfo.TypeOf(sel.X).String()
		}
	}
	return nil, ""
}

func (i *QuestionSyntaxInspector) InspectFuncRetType(t ast.Expr) *QImplType {
	expr, ok := t.(*ast.IndexExpr)
	if ok {
		return &QImplType{
			MainType:  GetNameFromType(i.pkg.TypesInfo.TypeOf(expr)),
			InnerType: i.pkg.TypesInfo.TypeOf(expr.Index).String(),
		}
	}
	return nil
}

func (i *QuestionSyntaxInspector) InspectFuncType(typ *ast.FuncType) *QImplType {
	if typ.Results != nil && typ.Results.NumFields() == 1 {
		retType := i.InspectFuncRetType(typ.Results.List[0].Type)
		if retType != nil {
			if _, ok := i.instanceTypes[retType.MainType]; ok {
				return retType
			}
		}
	}
	return nil
}

func (i *QuestionSyntaxInspector) InspectSyntax(n ast.Node, stack []ast.Node) (syntax *QuestionSyntax) {
	callExpr := n.(*ast.CallExpr)
	exprExt, exprType := i.inspectQuestionFnCall(callExpr)
	call := &QuestionCall{
		Expr:     exprExt,
		ExprType: exprType,
	}
	if exprExt != nil && len(stack) > 1 {
		var retType *QImplType
	FindOuterFunc:
		for idx := len(stack) - 2; idx >= 0; idx-- {
			switch fn := stack[idx].(type) {
			case *ast.FuncLit:
				retType = i.InspectFuncType(fn.Type)
				break FindOuterFunc
			case *ast.FuncDecl:
				retType = i.InspectFuncType(fn.Type)
				break FindOuterFunc
			}
		}
		if retType != nil {
			switch parent := stack[len(stack)-2].(type) {
			case *ast.AssignStmt:
				for idx, rhs := range parent.Rhs {
					rhsCall, ok := rhs.(*ast.CallExpr)
					if ok && rhsCall == callExpr {
						call.Extent = Extent{
							Start: i.pkg.Fset.Position(parent.Pos()),
							End:   i.pkg.Fset.Position(parent.End()),
						}
						call.AssignVar = &Extent{
							Start: i.pkg.Fset.Position(parent.Lhs[idx].Pos()),
							End:   i.pkg.Fset.Position(parent.Lhs[idx].End()),
						}
						call.AssignToken = parent.Tok.String()
						break
					}
				}
			case *ast.ExprStmt:
				call.Extent = Extent{
					Start: i.pkg.Fset.Position(parent.Pos()),
					End:   i.pkg.Fset.Position(parent.End()),
				}
			default:
				if len(stack) > 2 {
				FindOuterStmt:
					for idx := len(stack) - 3; idx >= 0; idx-- {
						switch stmt := stack[idx].(type) {
						case ast.Stmt:
							call.OuterStmt = &Extent{
								Start: i.pkg.Fset.Position(stmt.Pos()),
								End:   i.pkg.Fset.Position(stmt.End()),
							}
							call.Extent = Extent{
								Start: i.pkg.Fset.Position(callExpr.Pos()),
								End:   i.pkg.Fset.Position(callExpr.End()),
							}
							break FindOuterStmt
						}
					}
				}
			}
			syntax = &QuestionSyntax{
				Call:    call,
				RetType: retType,
			}
		}
	}
	return
}

func GenAssginStmt(assignVar, assignToken, call string) string {
	return fmt.Sprintf("%s %s %s\n", assignVar, assignToken, call)
}

func GenErrorHandler(resultVar, retType string) string {
	return fmt.Sprintf("if %s.IsErr() {\nreturn Err[%s](%s.UnwrapErr())\n}\n", resultVar, retType, resultVar)
}

func GenUnwrapStmt(assignVar, assignToken, resultVar string) string {
	return fmt.Sprintf("%s %s %s.Unwrap()\n", assignVar, assignToken, resultVar)
}

func GenUnwrapExpr(resultVar string) string {
	return fmt.Sprintf("%s.Unwrap()", resultVar)
}

func GenerateQuestionSyntax(info *FileInfo[QuestionSyntax], writer io.Writer) error {
	return GenerateSyntax(info, writer, func(file *os.File, addImports map[string]string) ([]*ReplaceBlock, error) {
		var ret []*ReplaceBlock
		for _, syntax := range info.Syntax {
			if syntax.Call.AssignVar != nil {
				assignVar, err := ReadExtent(file, *syntax.Call.AssignVar)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				callExpr, err := ReadExtent(file, *syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				resultVar := GetRandVarName(syntax.Call.AssignVar.Start.String())
				retType, adds := ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(resultVar, token.DEFINE.String(), callExpr)
				result += GenErrorHandler(resultVar, retType)
				result += GenUnwrapStmt(assignVar, syntax.Call.AssignToken, resultVar)
				ret = append(ret, &ReplaceBlock{
					Old: syntax.Call.Extent,
					New: result,
				})
			} else if syntax.Call.OuterStmt != nil {
				callExpr, err := ReadExtent(file, *syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				resultVar := GetRandVarName(syntax.Call.Expr.Start.String())
				retType, adds := ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(resultVar, token.DEFINE.String(), callExpr)
				result += GenErrorHandler(resultVar, retType)
				ret = append(ret,
					&ReplaceBlock{
						Old: Extent{
							Start: syntax.Call.OuterStmt.Start,
							End:   syntax.Call.OuterStmt.Start,
						},
						New: result,
					},
					&ReplaceBlock{
						Old: syntax.Call.Extent,
						New: GenUnwrapExpr(resultVar),
					})
			} else {
				callExpr, err := ReadExtent(file, *syntax.Call.Expr)
				if err != nil {
					return nil, fmt.Errorf("ReadExtent() failed: %w", err)
				}
				resultVar := GetRandVarName(syntax.Call.Expr.Start.String())
				retType, adds := ResetTypeStrPkgName(syntax.RetType.InnerType, info.Imports, info.PkgPath)
				if len(adds) > 0 {
					for k, v := range adds {
						addImports[k] = v
					}
				}
				result := GenAssginStmt(resultVar, token.DEFINE.String(), callExpr)
				result += GenErrorHandler(resultVar, retType)
				ret = append(ret, &ReplaceBlock{
					Old: syntax.Call.Extent,
					New: result,
				})
			}
		}
		return ret, nil
	})
}
