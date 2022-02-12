package macro

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
)

// the marco grammar:
//
// ```
// Macro := '#[' CallSequence ']'
// CallSequence := CallExpr (',' CallExpr)*
// CallExpr := CallWithoutArg | CallWithArgs
// CallWithoutArg := GoIdent
// CallWithArgs := GoIdent '(' Args? ')'
// Args := Arg (',' Arg)*
// Arg := GoIdent | GoString | GoInt | GoFloat | True | False
// ```
//
// macro can be placed:
// - as last comment on top of the struct definition
// - immediately follow the field of struct and/or enum
//
// ```go
// // #[visitor]
// type BinExpr struct {
//   lhs Node // #[visitor]
//   rhs Node // #[visitor], some other comments
// }
//
// const (
//   N_PROG // #[visitor]
// )
// ```

type ProcCtx struct {
	Path string                  // the path to start the entire process
	Fset *token.FileSet          // files been processed by `go/parser/ParseDir`
	Pkgs map[string]*ast.Package // processed results, produced by `go/parser/ParseDir`
}

type MacroCtx struct {
	ProcCtx *ProcCtx
	File    string        // file being walked through to find out the macros
	Node    ast.Node      // the node being attached with macro
	Name    string        // the name of this macro
	Args    []interface{} // the args defined in the macro expr
}

type MacroImpl = func(MacroCtx)

func HasMacroLike(cmt string) (string, bool) {
	reg := regexp.MustCompile(`^//\s*#\[([^]]*?)\]`)
	matched := reg.FindStringSubmatch(cmt)
	if len(matched) != 2 {
		return "", false
	}
	return matched[1], true
}

func typ(t interface{}) reflect.Type {
	return reflect.TypeOf(t).Elem()
}

func ParseMacro(file, macro string, targe ast.Node, procCtx *ProcCtx) ([]*MacroCtx, error) {
	// construct a valid call expr in go, otherwise the returned `expr` will be approximate
	expr, err := parser.ParseExpr("_(" + macro + ")")
	if err != nil {
		return nil, err
	}
	ctxs := make([]*MacroCtx, 0, 1)
	for _, m := range expr.(*ast.CallExpr).Args {
		if typ(m) == typ((*ast.Ident)(nil)) {
			mc := &MacroCtx{procCtx, file, targe, m.(*ast.Ident).Name, nil}
			ctxs = append(ctxs, mc)
		} else if typ(m) == typ((*ast.CallExpr)(nil)) {
			c := m.(*ast.CallExpr)
			name, err := nameOfCallExpr(c)
			if err != nil {
				return nil, err
			}
			args, err := argsOfCallExpr(c)
			if err != nil {
				return nil, err
			}
			mc := &MacroCtx{procCtx, file, targe, name, args}
			ctxs = append(ctxs, mc)
		} else {
			return nil, fmt.Errorf("deformed macro: %v", m)
		}
	}
	return ctxs, nil
}

func nameOfCallExpr(c *ast.CallExpr) (string, error) {
	id := c.Fun
	if typ(id) != typ((*ast.Ident)(nil)) {
		return "", fmt.Errorf("deformed macro name: %v", id)
	}
	return id.(*ast.Ident).Name, nil
}

func argsOfCallExpr(c *ast.CallExpr) ([]interface{}, error) {
	args := make([]interface{}, 0, 1)
	for _, arg := range c.Args {
		a, err := parseArg(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	return args, nil
}

func parseArg(node ast.Node) (interface{}, error) {
	switch v := node.(type) {
	case *ast.Ident:
		if v.Name == "true" {
			return true, nil
		}
		if v.Name == "false" {
			return false, nil
		}
		return v.Name, nil
	case *ast.BasicLit:
		if v.Kind == token.INT || v.Kind == token.FLOAT {
			return v.Value, nil
		}
		if v.Kind == token.STRING {
			s, err := strconv.Unquote(v.Value)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, fmt.Errorf("deformed macro arg: %v", node)
}
