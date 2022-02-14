package macro

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Marco Grammar:
//
// ```
// Macro := '#[' CallSequence ']'
// CallSequence := CallExpr (',' CallExpr)*
// CallExpr := CallWithoutArg | CallWithArgs
// CallWithoutArg := GoIdent
// CallWithArgs := GoIdent '(' Args? ')'
// Args := Arg (',' Arg)*
// Arg := GoIdent | GoBasicLit | True | False | GoSelectorExpr
// GoBasicLic := GoInt | GoFloat | GoString
// ```
//
// The permitted postions to put the macro are:
// - the last comment on top of the struct definition
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

func IsTyp(o interface{}, t interface{}) bool {
	return reflect.TypeOf(o).Elem() == reflect.TypeOf(t).Elem()
}

func ParseMacro(file, macro string, targe ast.Node, procCtx *ProcCtx) ([]*MacroCtx, error) {
	// construct a valid call expr in go, otherwise the returned `expr` will be approximate
	expr, err := parser.ParseExpr("_(" + macro + ")")
	if err != nil {
		return nil, err
	}
	ctxs := make([]*MacroCtx, 0, 1)
	for _, m := range expr.(*ast.CallExpr).Args {
		if IsTyp(m, (*ast.Ident)(nil)) {
			mc := &MacroCtx{procCtx, file, targe, m.(*ast.Ident).Name, nil}
			ctxs = append(ctxs, mc)
		} else if IsTyp(m, (*ast.CallExpr)(nil)) {
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
	case *ast.SelectorExpr:
		return nameOfSelectorExpr(v)
	}
	return nil, fmt.Errorf("deformed macro arg: %v", node)
}

func nameOfSelectorExpr(s *ast.SelectorExpr) (string, error) {
	ns := make([]string, 0, 1)
	for {
		ns = append([]string{s.Sel.Name}, ns...)
		if IsTyp(s.X, (*ast.SelectorExpr)(nil)) {
			s = s.X.(*ast.SelectorExpr)
		} else if IsTyp(s.X, (*ast.Ident)(nil)) {
			ns = append([]string{s.X.(*ast.Ident).Name}, ns...)
			break
		} else {
			return "", fmt.Errorf("deformed macro arg in SelectorExpr.X: %v", s.X)
		}
	}
	return strings.Join(ns, "."), nil
}

func IsStructDec(n ast.Node) (string, *ast.StructType, bool) {
	if !IsTyp(n, (*ast.GenDecl)(nil)) || n.(*ast.GenDecl).Tok != token.TYPE {
		return "", nil, false
	}
	d := n.(*ast.GenDecl)
	if len(d.Specs) != 1 || !IsTyp(d.Specs[0], (*ast.TypeSpec)(nil)) {
		return "", nil, false
	}
	t := d.Specs[0].(*ast.TypeSpec).Type
	if !IsTyp(t, (*ast.StructType)(nil)) {
		return "", nil, false
	}
	ts := d.Specs[0].(*ast.TypeSpec)
	return ts.Name.Name, ts.Type.(*ast.StructType), true
}

func isEnum(n ast.Node) ([]ast.Spec, bool) {
	if !IsTyp(n, (*ast.GenDecl)(nil)) || n.(*ast.GenDecl).Tok != token.CONST {
		return nil, false
	}
	return n.(*ast.GenDecl).Specs, true
}

func macroCtxsOfEnum(file *ast.File, filename string, specs []ast.Spec, procCtx *ProcCtx) ([]*MacroCtx, error) {
	ctxs := make([]*MacroCtx, 0, 1)
	for _, spec := range specs {
		if v, ok := spec.(*ast.ValueSpec); ok {
			if v.Comment == nil {
				continue
			}
			cmt := v.Comment.List[0].Text
			if c, ok := HasMacroLike(cmt); ok {
				m, err := ParseMacro(filename, c, spec, procCtx)
				if err != nil {
					return nil, err
				}
				ctxs = append(ctxs, m...)
			}
		}
	}
	return ctxs, nil
}

func macroCtxsInsideStruct(filename string, s *ast.StructType, procCtx *ProcCtx) ([]*MacroCtx, error) {
	ctxs := make([]*MacroCtx, 0, 1)
	for _, f := range s.Fields.List {
		if f.Comment == nil {
			continue
		}
		cmt := f.Comment.List[0].Text
		if c, ok := HasMacroLike(cmt); ok {
			m, err := ParseMacro(filename, c, f, procCtx)
			if err != nil {
				return nil, err
			}
			ctxs = append(ctxs, m...)
		}
	}
	return ctxs, nil
}

// the comments on the top of the struct definition are attached to the GenDecl which is the parent Node of the
// struct definition, so the comment needs to be passed by the callee
func macroCtxsOfStruct(filename string, comment string, s *ast.GenDecl, procCtx *ProcCtx) ([]*MacroCtx, error) {
	cmt, ok := HasMacroLike(comment)
	if !ok {
		return nil, nil
	}
	return ParseMacro(filename, cmt, s, procCtx)
}

func MacroCtxsOfFile(file *ast.File, filename string, procCtx *ProcCtx) ([]*MacroCtx, error) {
	ctxs := make([]*MacroCtx, 0, 1)
	for _, dec := range file.Decls {
		if n, ok := isEnum(dec); ok {
			cs, err := macroCtxsOfEnum(file, filename, n, procCtx)
			if err != nil {
				return nil, err
			}
			ctxs = append(ctxs, cs...)
		} else if _, n, ok := IsStructDec(dec); ok {
			// process the struct itself
			if dec.(*ast.GenDecl).Doc != nil {
				cmts := dec.(*ast.GenDecl).Doc.List
				cmt := cmts[len(cmts)-1].Text
				cs, err := macroCtxsOfStruct(filename, cmt, dec.(*ast.GenDecl), procCtx)
				if err != nil {
					return nil, err
				}
				ctxs = append(ctxs, cs...)
			}
			// process the fields of the struct
			cs, err := macroCtxsInsideStruct(filename, n, procCtx)
			if err != nil {
				return nil, err
			}
			ctxs = append(ctxs, cs...)
		}
	}
	return ctxs, nil
}

type WalkPkgHandle = func(*ast.File, string, *ProcCtx) error

func WalkPkgs(pkgs map[string]*ast.Package, handle WalkPkgHandle, pc *ProcCtx) error {
	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			err := handle(file, filename, pc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MacroCtxsOfWorkingDir(wd string, defDir string) ([]*MacroCtx, *ProcCtx, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, defDir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	procCtx := &ProcCtx{Path: wd, Fset: fset, Pkgs: pkgs}
	ctxs := make([]*MacroCtx, 0, 1)

	err = WalkPkgs(pkgs, func(f *ast.File, s string, pc *ProcCtx) error {
		var cs []*MacroCtx
		cs, err := MacroCtxsOfFile(f, s, procCtx)
		if err != nil {
			return err
		}
		ctxs = append(ctxs, cs...)
		return nil
	}, procCtx)

	if err != nil {
		return nil, nil, err
	}
	return ctxs, procCtx, nil
}
