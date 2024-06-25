package exception

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/arcane-craft/sugar/tool/transform/lib"
	"golang.org/x/tools/go/packages"
)

const (
	errorTypeName        = "error"
	exceptionPkgPath     = "github.com/arcane-craft/sugar/syntax/exception"
	tryFunName           = "Try"
	catchFunName         = "Catch"
	catchTargetErrorName = "Error"
	catchTargetTypeName  = "Type"
	finallyFunName       = "Finally"
	throwFunName         = "Throw"
	returnFunName        = "Return"
)

type ExceptionSyntaxInspector struct {
	pkg *packages.Package
}

func NewExceptionSyntaxInspector(pkg *packages.Package) *ExceptionSyntaxInspector {
	return &ExceptionSyntaxInspector{
		pkg: pkg,
	}
}

func (i *ExceptionSyntaxInspector) Nodes() []ast.Node {
	return []ast.Node{
		&ast.ExprStmt{},
	}
}

type CallStmt interface {
	fmt.Stringer
	Start() token.Position
	End() token.Position
}

type Func struct {
	lib.Extent
	Expr        *lib.Extent
	AssignToken string
	Vars        []*lib.Extent
	OuterStmt   *lib.Extent
}

func (m *Func) Start() token.Position {
	return m.Extent.Start
}

func (m *Func) End() token.Position {
	return m.Extent.End
}

func (m *Func) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Expr: %s, Vars: [%s], OuterStmt: %s",
		m.Extent, m.Expr, lib.JoinStringers(m.Vars, ";"), m.OuterStmt)
}

type Return struct {
	lib.Extent
	Args []*lib.Extent
}

func (m *Return) Start() token.Position {
	return m.Extent.Start
}

func (m *Return) End() token.Position {
	return m.Extent.End
}

func (m *Return) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Args: [%s]", m.Extent, lib.JoinStringers(m.Args, ";"))
}

type Throw struct {
	lib.Extent
	Err *lib.Extent
}

func (m *Throw) Start() token.Position {
	return m.Extent.Start
}

func (m *Throw) End() token.Position {
	return m.Extent.End
}

func (m *Throw) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Err: %s", m.Extent, m.Err)
}

type SyntaxBlock interface {
	Start() token.Position
	End() token.Position
	BodyStart() token.Position
	BodyEnd() token.Position
	CallStmts() []CallStmt
	fmt.Stringer
}

type Try struct {
	lib.Extent
	Body  *lib.Extent
	Calls []CallStmt
}

func (m *Try) Start() token.Position {
	return m.Extent.Start
}

func (m *Try) End() token.Position {
	return m.Extent.End
}

func (m *Try) BodyStart() token.Position {
	return m.Body.Start
}

func (m *Try) BodyEnd() token.Position {
	return m.Body.End
}

func (m *Try) CallStmts() []CallStmt {
	return m.Calls
}

func (m *Try) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Calls: [%s]", m.Extent, lib.JoinStringers(m.Calls, ";"))
}

const (
	CatchError = 0
	CatchType  = 1
)

type Catch struct {
	lib.Extent
	Body    *lib.Extent
	Type    int
	Targets []*lib.Extent
	Err     string
	Calls   []CallStmt
}

func (m *Catch) Start() token.Position {
	return m.Extent.Start
}

func (m *Catch) End() token.Position {
	return m.Extent.End
}

func (m *Catch) BodyStart() token.Position {
	return m.Body.Start
}

func (m *Catch) BodyEnd() token.Position {
	return m.Body.End
}

func (m *Catch) CallStmts() []CallStmt {
	return m.Calls
}

func (m *Catch) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Type: %d, Targets: [%s], Err: %s, Calls: [%s]",
		m.Extent, m.Type, lib.JoinStringers(m.Targets, ";"), m.Err, lib.JoinStringers(m.Calls, ";"))
}

type Finally struct {
	lib.Extent
	Body  *lib.Extent
	Calls []CallStmt
}

func (m *Finally) Start() token.Position {
	return m.Extent.Start
}

func (m *Finally) End() token.Position {
	return m.Extent.End
}

func (m *Finally) BodyStart() token.Position {
	return m.Body.Start
}

func (m *Finally) BodyEnd() token.Position {
	return m.Body.End
}

func (m *Finally) CallStmts() []CallStmt {
	return m.Calls
}

func (m *Finally) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s, Calls: [%s]", m.Extent, lib.JoinStringers(m.Calls, ";"))
}

type ExceptionSyntax struct {
	lib.Extent
	RetTypes   []*lib.Extent
	Blocks     []SyntaxBlock
	HasCatch   bool
	HasFinally bool
}

func (m *ExceptionSyntax) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RetTypes: [%s], Blocks: [%s]", lib.JoinStringers(m.RetTypes, ","), lib.JoinStringers(m.Blocks, ";"))
}

func (i *ExceptionSyntaxInspector) inspectSyntaxBody(body *ast.BlockStmt) (calls []CallStmt) {
	var outerStmt ast.Stmt
	ast.Inspect(body, func(child ast.Node) bool {
		if _, ok := child.(*ast.FuncLit); ok {
			return false
		}
		stmt, ok := child.(ast.Stmt)
		if !ok {
			return true
		}
		throwStmt := i.findThrowStmt(stmt)
		if throwStmt != nil {
			calls = append(calls, throwStmt)
			outerStmt = nil
			return false
		}
		returnStmt := i.findReturnStmt(stmt)
		if returnStmt != nil {
			calls = append(calls, returnStmt)
			outerStmt = nil
			return false
		}
		funcStmt := i.findCallStmt(stmt, outerStmt)
		if funcStmt != nil {
			calls = append(calls, funcStmt)
			return false
		}
		outerStmt = stmt
		return true
	})
	return
}

func (i *ExceptionSyntaxInspector) inspectTryBlock(node ast.Expr) (ret *Try) {
	callExpr, ok := node.(*ast.CallExpr)
	if ok && len(callExpr.Args) == 1 {
		var object types.Object
		switch fun := callExpr.Fun.(type) {
		case *ast.SelectorExpr:
			object = i.pkg.TypesInfo.ObjectOf(fun.Sel)
		case *ast.Ident:
			object = i.pkg.TypesInfo.ObjectOf(fun)
		}
		if object != nil && object.Pkg() != nil &&
			object.Pkg().Path() == exceptionPkgPath &&
			object.Name() == tryFunName {

			funcLit, ok := callExpr.Args[0].(*ast.FuncLit)
			if ok && funcLit.Body != nil {
				calls := i.inspectSyntaxBody(funcLit.Body)
				if len(calls) > 0 {
					ret = &Try{
						Extent: lib.Extent{
							Start: i.pkg.Fset.Position(callExpr.Pos()),
							End:   i.pkg.Fset.Position(callExpr.End()),
						},
						Body: &lib.Extent{
							Start: i.pkg.Fset.Position(funcLit.Body.Pos() + 1),
							End:   i.pkg.Fset.Position(funcLit.Body.End() - 1),
						},
						Calls: calls,
					}
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) queryCatchTarget(node ast.Expr) (typ int, targets []*lib.Extent, found bool) {
	callExpr, ok := node.(*ast.CallExpr)
	if ok {
		fun := callExpr.Fun
		indexExpr, ok := callExpr.Fun.(*ast.IndexExpr)
		if ok {
			fun = indexExpr.X
		}
		indexListExpr, ok := callExpr.Fun.(*ast.IndexListExpr)
		if ok {
			fun = indexListExpr.X
		}
		var targetType string
		switch f := fun.(type) {
		case *ast.SelectorExpr:
			targetType = f.Sel.Name
		case *ast.Ident:
			targetType = f.Name
		}
		if len(targetType) > 0 {
			if targetType == catchTargetErrorName {
				typ = CatchError
				for _, arg := range callExpr.Args {
					targets = append(targets, &lib.Extent{
						Start: i.pkg.Fset.Position(arg.Pos()),
						End:   i.pkg.Fset.Position(arg.End()),
					})
				}
				found = true
			} else if strings.HasPrefix(targetType, catchTargetTypeName) {
				typ = CatchType
				if indexExpr != nil {
					targets = append(targets, &lib.Extent{
						Start: i.pkg.Fset.Position(indexExpr.Index.Pos()),
						End:   i.pkg.Fset.Position(indexExpr.Index.End()),
					})
				} else if indexListExpr != nil {
					for _, elem := range indexListExpr.Indices {
						targets = append(targets, &lib.Extent{
							Start: i.pkg.Fset.Position(elem.Pos()),
							End:   i.pkg.Fset.Position(elem.End()),
						})
					}
				}
				found = true
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) inspectCatchBlock(node ast.Expr) (ret *Catch) {
	callExpr, ok := node.(*ast.CallExpr)
	if ok && len(callExpr.Args) == 2 {
		fun, ok := callExpr.Fun.(*ast.SelectorExpr)
		if ok {
			object := i.pkg.TypesInfo.ObjectOf(fun.Sel)
			if object != nil && object.Pkg() != nil &&
				object.Pkg().Path() == exceptionPkgPath &&
				object.Name() == catchFunName {
				catchType, targets, ok := i.queryCatchTarget(callExpr.Args[0])
				if ok {
					funcLit, ok := callExpr.Args[1].(*ast.FuncLit)
					if ok && funcLit.Type.Params != nil && funcLit.Type.Params.NumFields() == 1 {
						var errVar string
						if len(funcLit.Type.Params.List[0].Names) > 0 {
							errVar = funcLit.Type.Params.List[0].Names[0].Name
						}
						if funcLit.Body != nil {
							calls := i.inspectSyntaxBody(funcLit.Body)
							if len(calls) > 0 {
								ret = &Catch{
									Extent: lib.Extent{
										Start: i.pkg.Fset.Position(fun.Sel.Pos() - 1),
										End:   i.pkg.Fset.Position(callExpr.End()),
									},
									Body: &lib.Extent{
										Start: i.pkg.Fset.Position(funcLit.Body.Pos() + 1),
										End:   i.pkg.Fset.Position(funcLit.Body.End() - 1),
									},
									Type:    catchType,
									Targets: targets,
									Err:     errVar,
									Calls:   calls,
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) inspectFinallyBlock(node ast.Expr) (ret *Finally) {
	callExpr, ok := node.(*ast.CallExpr)
	if ok && len(callExpr.Args) == 1 {
		fun, ok := callExpr.Fun.(*ast.SelectorExpr)
		if ok {
			object := i.pkg.TypesInfo.ObjectOf(fun.Sel)
			if object != nil && object.Pkg() != nil &&
				object.Pkg().Path() == exceptionPkgPath &&
				object.Name() == finallyFunName {
				funcLit, ok := callExpr.Args[0].(*ast.FuncLit)
				if ok && funcLit.Body != nil {
					calls := i.inspectSyntaxBody(funcLit.Body)
					if len(calls) > 0 {
						ret = &Finally{
							Extent: lib.Extent{
								Start: i.pkg.Fset.Position(fun.Sel.Pos() - 1),
								End:   i.pkg.Fset.Position(callExpr.End()),
							},
							Body: &lib.Extent{
								Start: i.pkg.Fset.Position(funcLit.Body.Pos() + 1),
								End:   i.pkg.Fset.Position(funcLit.Body.End() - 1),
							},
							Calls: calls,
						}
					}
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) findFuncRetType(node ast.Node) (ret []*lib.Extent, finish bool) {
	var resultTypes *ast.FieldList
	switch fun := node.(type) {
	case *ast.FuncDecl:
		resultTypes = fun.Type.Results
	case *ast.FuncLit:
		resultTypes = fun.Type.Results
	default:
		return
	}
	if resultTypes != nil {
		finish = true
		if resultTypes.NumFields() > 0 {
			ident, ok := resultTypes.List[len(resultTypes.List)-1].Type.(*ast.Ident)
			if ok && ident.Name == errorTypeName {
				for _, elem := range resultTypes.List {
					ret = append(ret, &lib.Extent{
						Start: i.pkg.Fset.Position(elem.Type.Pos()),
						End:   i.pkg.Fset.Position(elem.Type.End()),
					})
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) queryCallStmtRetType(node ast.Expr) (num *int) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	switch typ := i.pkg.TypesInfo.TypeOf(call).(type) {
	case *types.Tuple:
		if typ.Len() > 0 &&
			typ.At(typ.Len()-1).Type().String() == errorTypeName {
			num = new(int)
			*num = typ.Len()
		}
	case *types.Named:
		if typ.Obj() != nil && typ.Obj().Name() == errorTypeName {
			num = new(int)
			*num = 1
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) findCallStmt(node ast.Stmt, outer ast.Stmt) (ret *Func) {
	switch stmt := node.(type) {
	case *ast.ExprStmt:
		retNum := i.queryCallStmtRetType(stmt.X)
		if retNum != nil {
			ret = &Func{
				Extent: lib.Extent{
					Start: i.pkg.Fset.Position(stmt.Pos()),
					End:   i.pkg.Fset.Position(stmt.End()),
				},
				Expr: &lib.Extent{
					Start: i.pkg.Fset.Position(stmt.X.Pos()),
					End:   i.pkg.Fset.Position(stmt.X.End()),
				},
			}
			for idx := 0; idx < *retNum; idx++ {
				ret.Vars = append(ret.Vars, &lib.Extent{
					Start: i.pkg.Fset.Position(stmt.X.Pos() + token.Pos(idx)),
					End:   i.pkg.Fset.Position(stmt.X.End()),
				})
			}
		}
	case *ast.AssignStmt:
		if len(stmt.Rhs) == 1 {
			retNum := i.queryCallStmtRetType(stmt.Rhs[0])
			if retNum != nil {
				if ident, ok := stmt.Lhs[len(stmt.Lhs)-1].(*ast.Ident); ok && ident.Name == "_" {
					callExpr := stmt.Rhs[0]
					ret = &Func{
						Extent: lib.Extent{
							Start: i.pkg.Fset.Position(stmt.Pos()),
							End:   i.pkg.Fset.Position(stmt.End()),
						},
						Expr: &lib.Extent{
							Start: i.pkg.Fset.Position(callExpr.Pos()),
							End:   i.pkg.Fset.Position(callExpr.End()),
						},
					}
					if *retNum > 1 {
						ret.AssignToken = stmt.Tok.String()
					}
					for idx := 0; idx < len(stmt.Lhs); idx++ {
						ret.Vars = append(ret.Vars, &lib.Extent{
							Start: i.pkg.Fset.Position(stmt.Lhs[idx].Pos()),
							End:   i.pkg.Fset.Position(stmt.Lhs[idx].End()),
						})
					}
				}
			}
		}
	}

	if outer != nil {
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

func (i *ExceptionSyntaxInspector) findThrowStmt(node ast.Stmt) (ret *Throw) {
	stmt, ok := node.(*ast.ExprStmt)
	if ok {
		call, ok := stmt.X.(*ast.CallExpr)
		if ok && len(call.Args) == 1 {
			var object types.Object
			switch fun := call.Fun.(type) {
			case *ast.SelectorExpr:
				object = i.pkg.TypesInfo.ObjectOf(fun.Sel)
			case *ast.Ident:
				object = i.pkg.TypesInfo.ObjectOf(fun)
			}
			if object != nil && object.Pkg() != nil &&
				object.Pkg().Path() == exceptionPkgPath &&
				object.Name() == throwFunName {
				ret = &Throw{
					Extent: lib.Extent{
						Start: i.pkg.Fset.Position(stmt.Pos()),
						End:   i.pkg.Fset.Position(stmt.End()),
					},
					Err: &lib.Extent{
						Start: i.pkg.Fset.Position(call.Args[0].Pos()),
						End:   i.pkg.Fset.Position(call.Args[0].End()),
					},
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) findReturnStmt(node ast.Stmt) (ret *Return) {
	stmt, ok := node.(*ast.ExprStmt)
	if ok {
		call, ok := stmt.X.(*ast.CallExpr)
		if ok {
			var object types.Object
			switch fun := call.Fun.(type) {
			case *ast.SelectorExpr:
				object = i.pkg.TypesInfo.ObjectOf(fun.Sel)
			case *ast.Ident:
				object = i.pkg.TypesInfo.ObjectOf(fun)
			}
			if object != nil && object.Pkg() != nil &&
				object.Pkg().Path() == exceptionPkgPath &&
				object.Name() == returnFunName {
				ret = &Return{
					Extent: lib.Extent{
						Start: i.pkg.Fset.Position(stmt.Pos()),
						End:   i.pkg.Fset.Position(stmt.End()),
					},
				}
				for _, arg := range call.Args {
					ret.Args = append(ret.Args, &lib.Extent{
						Start: i.pkg.Fset.Position(arg.Pos()),
						End:   i.pkg.Fset.Position(arg.End()),
					})
				}
			}
		}
	}
	return
}

func (i *ExceptionSyntaxInspector) Inspect(node ast.Node, stack []ast.Node) (syntax *ExceptionSyntax) {
	exprStmt := node.(*ast.ExprStmt)
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return
	}

	var blocks []SyntaxBlock
	var hasCatch, hasFinally bool
	finallyBlock := i.inspectFinallyBlock(callExpr)
	if finallyBlock != nil {
		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		callExpr, ok = sel.X.(*ast.CallExpr)
		if !ok {
			return
		}
		blocks = append(blocks, finallyBlock)
		hasFinally = true
	}
	for catch := i.inspectCatchBlock(callExpr); catch != nil; catch = i.inspectCatchBlock(callExpr) {
		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		callExpr, ok = sel.X.(*ast.CallExpr)
		if !ok {
			return
		}
		blocks = append(blocks, catch)
		hasCatch = true
	}
	tryBlock := i.inspectTryBlock(callExpr)
	if tryBlock == nil {
		return
	}
	blocks = append(blocks, tryBlock)

	if len(stack) < 2 {
		return
	}
	var outerFuncRetTypes []*lib.Extent
	for idx := len(stack) - 2; idx >= 0; idx-- {
		node := stack[idx]
		retTypes, finish := i.findFuncRetType(node)
		if !finish {
			continue
		}
		if len(retTypes) <= 0 {
			return
		}
		outerFuncRetTypes = retTypes
		break
	}

	slices.Reverse(blocks)
	syntax = &ExceptionSyntax{
		Extent: lib.Extent{
			Start: i.pkg.Fset.Position(exprStmt.Pos()),
			End:   i.pkg.Fset.Position(exprStmt.End()),
		},
		RetTypes:   outerFuncRetTypes,
		Blocks:     blocks,
		HasCatch:   hasCatch,
		HasFinally: hasFinally,
	}
	return
}

func GenBlock(stmts []string) string {
	return fmt.Sprintf("{\n%s\n}", strings.Join(stmts, "\n"))
}

func GenVarDecl(name string, typ string) string {
	return fmt.Sprintf("var %s %s", name, typ)
}

func GenAssignCall(lhs []string, token string, rhs string) string {
	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), token, rhs)
}

func GenErrHandler(errVar string, stmts []string) string {
	return fmt.Sprintf("if %s != nil {\n%s\n}", errVar, strings.Join(stmts, "\n"))
}

func GenErrIsHandler(errVar string, errVals []string, stmts []string) string {
	var checks []string
	for _, e := range errVals {
		checks = append(checks, fmt.Sprintf("errors.Is(%s, %s)", errVar, e))
	}
	return fmt.Sprintf("if %s {\n%s\n}", strings.Join(checks, " || "), strings.Join(stmts, "\n"))
}

func GenErrAsHandler(errVar string, errTypes []string, stmts []string) string {
	var checks []string
	for _, t := range errTypes {
		checks = append(checks, fmt.Sprintf("errors.As(%s, %s)", errVar, t))
	}
	return fmt.Sprintf("if %s {\n%s\n}", strings.Join(checks, " || "), strings.Join(stmts, "\n"))
}

func GenAssigneStmt(lhs []string, token string, rhs []string) string {
	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), token, strings.Join(rhs, ", "))
}

func GenGotoStmt(label string) string {
	return fmt.Sprintf("goto %s", label)
}

func GenLabelDecl(name string) string {
	return name + ":"
}

func GenReturns(exprs []string) string {
	return fmt.Sprintf("return %s", strings.Join(exprs, ", "))
}

func GenOrExpr(predicates ...string) string {
	return strings.Join(predicates, " || ")
}

func GenCompareExpr(left string, op string, right string) string {
	return fmt.Sprintf("%s %s %s", left, op, right)
}

func GenIfStmt(cond string, stmt []string) string {
	return fmt.Sprintf("if %s {\n %s \n}", cond, strings.Join(stmt, "\n"))
}

type Traslator struct{}

func NewTraslator() *Traslator {
	return &Traslator{}
}

func (*Traslator) InpectTypes(p *packages.Package) []*lib.Extent {
	return nil
}

func (*Traslator) InspectSyntax(p *packages.Package, _ []*lib.Extent) lib.SyntaxInspector[*ExceptionSyntax] {
	return NewExceptionSyntaxInspector(p)
}

func genFuncStmt(file *os.File, call *Func, catchErrVar, gotoLabel string) ([]string, error) {
	var stmts []string

	expr, err := lib.ReadExtent(file, call.Expr)
	if err != nil {
		return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
	}
	var vars []string
	if len(call.AssignToken) > 0 {
		vars, err = lib.ReadExtentList(file, call.Vars[:len(call.Vars)-1])
		if err != nil {
			return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
		}
	} else {
		for range call.Vars[:len(call.Vars)-1] {
			vars = append(vars, "_")
		}
	}
	errVar := lib.GenVarName("err", call.Vars[len(call.Vars)-1].String())
	if call.AssignToken == ":=" {
		stmts = append(stmts, GenAssignCall(append(vars, errVar), call.AssignToken, expr))
	} else if call.AssignToken == "=" {
		stmts = append(stmts, GenVarDecl(errVar, errorTypeName))
		stmts = append(stmts, GenAssignCall(append(vars, errVar), call.AssignToken, expr))
	} else {
		stmts = append(stmts, GenAssignCall(append(vars, errVar), ":=", expr))
	}
	stmts = append(stmts, GenErrHandler(errVar, []string{
		GenAssigneStmt([]string{catchErrVar}, "=", []string{errVar}),
		GenGotoStmt(gotoLabel),
	}))

	return stmts, nil
}

func genThrowStmt(file *os.File, call *Throw, catchErrVar, gotoLabel string) ([]string, error) {
	var stmts []string

	errExpr, err := lib.ReadExtent(file, call.Err)
	if err != nil {
		return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
	}
	stmts = append(stmts, GenAssigneStmt([]string{catchErrVar}, "=", []string{errExpr}))
	stmts = append(stmts, GenGotoStmt(gotoLabel))

	return stmts, nil
}

func genReturnStmt(file *os.File, call *Return, resultVars []string, hasReturnVar, catchErrVar, gotoLabel string) ([]string, error) {
	var stmts []string

	args, err := lib.ReadExtentList(file, call.Args)
	if err != nil {
		return nil, fmt.Errorf("lib.ReadExtentList() failed: %w", err)
	}
	stmts = append(stmts, GenAssigneStmt(resultVars, "=", args))
	if len(hasReturnVar) > 0 {
		stmts = append(stmts, GenAssigneStmt([]string{hasReturnVar}, "=", []string{"true"}))
		stmts = append(stmts, GenGotoStmt(gotoLabel))

	} else {
		stmts = append(stmts, GenReturns(append(resultVars, catchErrVar)))
	}

	return stmts, nil
}

func genBlockCalls(file *os.File, block SyntaxBlock, proc func(c CallStmt) ([]string, error)) ([]string, error) {
	var blockStmts []string
	gapStart := block.BodyStart()
	for _, c := range block.CallStmts() {
		gapExtent := &lib.Extent{
			Start: gapStart,
			End:   c.Start(),
		}
		if call, ok := c.(*Func); ok && call.OuterStmt != nil {
			gapExtent.End = call.OuterStmt.Start
		}
		gap, err := lib.ReadExtent(file, gapExtent)
		if err != nil {
			return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
		}
		blockStmts = append(blockStmts, gap)
		gapStart = c.End()

		stmts, err := proc(c)
		if err != nil {
			return nil, err
		}
		blockStmts = append(blockStmts, stmts...)

		if call, ok := c.(*Func); ok && call.OuterStmt != nil {
			gapExtent.Start = call.OuterStmt.Start
			gapExtent.End = call.Start()
			gap, err := lib.ReadExtent(file, gapExtent)
			if err != nil {
				return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
			}
			blockStmts = append(blockStmts, gap)
		}
	}
	gapExtent := &lib.Extent{
		Start: gapStart,
		End:   block.BodyEnd(),
	}
	gap, err := lib.ReadExtent(file, gapExtent)
	if err != nil {
		return nil, fmt.Errorf("lib.ReadExtent() failed: %w", err)
	}
	blockStmts = append(blockStmts, gap)
	return blockStmts, nil
}

func genErrHandlers(file *os.File, block *Catch, catchErrVar string, handlerStmts []string) ([]string, error) {
	targets, err := lib.ReadExtentList(file, block.Targets)
	if err != nil {
		return nil, fmt.Errorf("lib.ReadExtentList() failed: %w", err)
	}
	var blockStmts []string
	if block.Type == CatchError {
		if len(targets) > 0 {
			blockStmts = append(blockStmts,
				GenErrIsHandler(catchErrVar, targets,
					append([]string{
						GenAssigneStmt([]string{block.Err}, ":=", []string{catchErrVar}),
						GenAssigneStmt([]string{catchErrVar}, "=", []string{"nil"}),
					}, handlerStmts...),
				),
			)
		} else {
			blockStmts = append(blockStmts,
				GenErrHandler(catchErrVar,
					append([]string{
						GenAssigneStmt([]string{block.Err}, ":=", []string{catchErrVar}),
						GenAssigneStmt([]string{catchErrVar}, "=", []string{"nil"}),
					}, handlerStmts...),
				),
			)
		}
	} else if block.Type == CatchType {
		blockStmts = append(blockStmts,
			GenErrAsHandler(catchErrVar, targets,
				append([]string{
					GenAssigneStmt([]string{block.Err}, ":=", []string{catchErrVar}),
					GenAssigneStmt([]string{catchErrVar}, "=", []string{"nil"}),
				}, handlerStmts...),
			),
		)
	}
	return blockStmts, nil
}

func (*Traslator) Generate(info *lib.FileInfo[*ExceptionSyntax], writer io.Writer) error {
	return lib.GenerateSyntax(info, writer, func(file *os.File, addImports map[string]string) ([]*lib.ReplaceBlock, error) {
		var ret []*lib.ReplaceBlock
		for _, s := range info.Syntax {
			var stmts []string

			retTypes, err := lib.ReadExtentList(file, s.RetTypes)
			if err != nil {
				return nil, fmt.Errorf("lib.ReadExtentList() failed: %w", err)
			}
			var resultVars []string
			var catchErrVar string
			for idx, typ := range retTypes {
				if idx < len(retTypes)-1 {
					rVar := lib.GenVarName("result", s.RetTypes[idx].String())
					resultVars = append(resultVars, rVar)
					stmts = append(stmts, GenVarDecl(rVar, typ))
				} else {
					catchErrVar = lib.GenVarName("catchErr", s.RetTypes[idx].String())
					stmts = append(stmts, GenVarDecl(catchErrVar, typ))
				}
			}
			hasReturnVar := lib.GenVarName("hasRet", s.Start.String())
			stmts = append(stmts, GenVarDecl(hasReturnVar, "bool"))
			catchLabel := lib.GenVarName("Catch", s.Blocks[0].String())
			finallyLabel := lib.GenVarName("Finally", s.String())

			for blockIdx, b := range s.Blocks {
				var blockStmts []string
				switch block := b.(type) {
				case *Try:
					{
						handlerStmts, err := genBlockCalls(file, block, func(c CallStmt) (stmts []string, err error) {
							switch call := c.(type) {
							case *Func:
								stmts, err = genFuncStmt(file, call, catchErrVar, catchLabel)
								if err != nil {
									err = fmt.Errorf("genFuncStmt() failed: %w", err)
								}
							case *Throw:
								stmts, err = genThrowStmt(file, call, catchErrVar, catchLabel)
								if err != nil {
									err = fmt.Errorf("genThrowStmt() failed: %w", err)
								}
							case *Return:
								stmts, err = genReturnStmt(file, call, resultVars, hasReturnVar, "", finallyLabel)
								if err != nil {
									err = fmt.Errorf("genReturnStmt() failed: %w", err)
								}
							default:
								err = fmt.Errorf("unexpected call")
							}
							return
						})
						if err != nil {
							return nil, err
						}
						blockStmts = append(blockStmts, handlerStmts...)
						blockStmts = append(blockStmts, GenGotoStmt(finallyLabel))
					}
				case *Catch:
					{
						handlerStmts, err := genBlockCalls(file, block, func(c CallStmt) (stmts []string, err error) {
							switch call := c.(type) {
							case *Func:
								stmts, err = genFuncStmt(file, call, catchErrVar, finallyLabel)
								if err != nil {
									err = fmt.Errorf("genFuncStmt() failed: %w", err)
								}
							case *Throw:
								stmts, err = genThrowStmt(file, call, catchErrVar, finallyLabel)
								if err != nil {
									err = fmt.Errorf("genThrowStmt() failed: %w", err)
								}
							case *Return:
								stmts, err = genReturnStmt(file, call, resultVars, hasReturnVar, "", finallyLabel)
								if err != nil {
									err = fmt.Errorf("genReturnStmt() failed: %w", err)
								}
							default:
								err = fmt.Errorf("unexpected call")
							}
							return
						})
						if err != nil {
							return nil, err
						}
						handlerStmts = append(handlerStmts, GenGotoStmt(finallyLabel))
						handlerStmts, err = genErrHandlers(file, block, catchErrVar, handlerStmts)
						if err != nil {
							return nil, err
						}
						blockStmts = append(blockStmts, handlerStmts...)
					}
				case *Finally:
					{
						stmts = append(stmts, GenLabelDecl(finallyLabel))

						handlerStmts, err := genBlockCalls(file, block, func(c CallStmt) (stmts []string, err error) {
							call, ok := c.(*Return)
							if ok {
								stmts, err = genReturnStmt(file, call, resultVars, "", catchErrVar, "")
								if err != nil {
									err = fmt.Errorf("genReturnStmt() failed: %w", err)
								}
							} else {
								err = fmt.Errorf("unexpected call")
							}
							return
						})
						if err != nil {
							return nil, err
						}
						blockStmts = append(blockStmts, handlerStmts...)
						blockStmts = append(blockStmts, GenIfStmt(
							GenOrExpr(hasReturnVar, GenCompareExpr(catchErrVar, "!=", "nil")),
							[]string{
								GenReturns(append(resultVars, catchErrVar)),
							},
						))
					}
				}
				stmts = append(stmts, GenBlock(blockStmts))

				if blockIdx == 0 {
					stmts = append(stmts, GenLabelDecl(catchLabel))
					if !s.HasCatch {
						stmts = append(stmts, GenBlock([]string{
							GenGotoStmt(finallyLabel),
						}))
					}
				}
			}
			if !s.HasFinally {
				stmts = append(stmts, GenLabelDecl(finallyLabel))
				stmts = append(stmts, GenBlock([]string{
					GenIfStmt(
						GenOrExpr(hasReturnVar, GenCompareExpr(catchErrVar, "!=", "nil")),
						[]string{
							GenReturns(append(resultVars, catchErrVar)),
						},
					),
				}))
			}
			ret = append(ret, &lib.ReplaceBlock{
				Old: s.Extent,
				New: GenBlock(stmts),
			})
		}
		return ret, nil
	})
}
