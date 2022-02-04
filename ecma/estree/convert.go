package estree

import (
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func pos(p *parser.Pos) *Position {
	return &Position{Line: p.Line(), Column: p.Column()}
}

func loc(s *parser.Loc) *SrcLoc {
	return &SrcLoc{
		Source: s.Source(),
		Start:  pos(s.Begin()),
		End:    pos(s.End()),
	}
}

func start(s *parser.Loc) int {
	return s.Range().Start()
}

func end(s *parser.Loc) int {
	return s.Range().End()
}

func ConvertProg(n *parser.Prog, ctx *ConvertCtx) *Program {
	stmts := n.Body()
	body := make([]Node, len(stmts))
	for i, s := range stmts {
		body[i] = Convert(s, ctx)
	}
	return &Program{
		Type:  "Program",
		Start: start(n.Loc()),
		End:   end(n.Loc()),
		Loc:   loc(n.Loc()),
		Body:  body,
	}
}

func arrExpr(n *parser.ArrLit, ctx *ConvertCtx) *ArrayExpression {
	exprs := n.Elems()
	elems := make([]Expression, len(exprs))
	for i, e := range exprs {
		elems[i] = Convert(e, ctx)
	}
	return &ArrayExpression{
		Type:     "ArrayExpression",
		Start:    start(n.Loc()),
		End:      end(n.Loc()),
		Loc:      loc(n.Loc()),
		Elements: elems,
	}
}

func obj(n *parser.ObjLit, ctx *ConvertCtx) *ObjectExpression {
	ps := n.Props()
	props := make([]Node, len(ps))
	for i, p := range ps {
		props[i] = Convert(p, ctx)
	}
	return &ObjectExpression{
		Type:       "ObjectExpression",
		Start:      start(n.Loc()),
		End:        end(n.Loc()),
		Loc:        loc(n.Loc()),
		Properties: props,
	}
}

func fnParams(params []parser.Node, ctx *ConvertCtx) []Node {
	ps := make([]Node, len(params))
	for i, p := range params {
		fp := tsParamProp(p, ctx)
		if fp == nil {
			fp = Convert(p, ctx)
		}
		ps[i] = fp
	}
	return ps
}

func statements(stmts []parser.Node, ctx *ConvertCtx) []Statement {
	s := make([]Statement, len(stmts))
	for i, stmt := range stmts {
		s[i] = Convert(stmt, ctx)
	}
	return s
}

func expressions(exprs []parser.Node, ctx *ConvertCtx) []Expression {
	s := make([]Expression, len(exprs))
	for i, expr := range exprs {
		s[i] = Convert(expr, ctx)
	}
	return s
}

func blockStmt(n *parser.BlockStmt, ctx *ConvertCtx) *BlockStatement {
	return &BlockStatement{
		Type:  "BlockStatement",
		Start: start(n.Loc()),
		End:   end(n.Loc()),
		Loc:   loc(n.Loc()),
		Body:  statements(n.Body(), ctx),
	}
}

func cases(cs []parser.Node, ctx *ConvertCtx) []*SwitchCase {
	s := make([]*SwitchCase, len(cs))
	for i, c := range cs {
		sc := c.(*parser.SwitchCase)
		s[i] = &SwitchCase{
			Type:       "SwitchCase",
			Start:      start(sc.Loc()),
			End:        end(sc.Loc()),
			Loc:        loc(sc.Loc()),
			Test:       Convert(sc.Test(), ctx),
			Consequent: statements(sc.Cons(), ctx),
		}
	}
	return s
}

func declarations(decList []parser.Node, ctx *ConvertCtx) []*VariableDeclarator {
	s := make([]*VariableDeclarator, len(decList))
	for i, d := range decList {
		dc := d.(*parser.VarDec)
		s[i] = &VariableDeclarator{
			Type:  "VariableDeclarator",
			Start: start(dc.Loc()),
			End:   end(dc.Loc()),
			Loc:   loc(dc.Loc()),
			Id:    Convert(dc.Id(), ctx),
			Init:  Convert(dc.Init(), ctx),
		}
	}
	return s
}

func tplLiteral(tplLoc *parser.Loc, elems []parser.Node, ctx *ConvertCtx) *TemplateLiteral {
	quasis := make([]Expression, 0)
	exprs := make([]Expression, 0)
	cnt := len(elems)
	for i, elem := range elems {
		first := i == 0
		last := i == cnt-1
		if elem.Type() != parser.N_LIT_STR {
			if first || last {
				lc := loc(tplLoc)
				if first {
					lc.End.Column = lc.Start.Column
				} else {
					lc.Start.Column = lc.End.Column
				}
				quasis = append(quasis, &TemplateElement{
					Type:  "TemplateElement",
					Start: start(tplLoc),
					End:   end(tplLoc),
					Loc:   lc,
					Tail:  last,
					Value: &TemplateElementValue{
						Cooked: "",
						Raw:    "",
					},
				})
			}
			exprs = append(exprs, Convert(elem, ctx))
		} else {
			str := elem.(*parser.StrLit)
			quasis = append(quasis, &TemplateElement{
				Type:  "TemplateElement",
				Start: start(elem.Loc()),
				End:   end(elem.Loc()),
				Loc:   loc(elem.Loc()),
				Tail:  last,
				Value: &TemplateElementValue{
					Cooked: str.Text(),
					Raw:    str.Raw(),
				},
			})
		}
	}
	return &TemplateLiteral{
		Type:        "TemplateLiteral",
		Start:       start(tplLoc),
		End:         end(tplLoc),
		Loc:         loc(tplLoc),
		Quasis:      quasis,
		Expressions: exprs,
	}
}

func importSpec(spec *parser.ImportSpec, ctx *ConvertCtx) Node {
	if spec.Default() {
		return &ImportDefaultSpecifier{
			Type:  "ImportDefaultSpecifier",
			Start: start(spec.Loc()),
			End:   end(spec.Loc()),
			Loc:   loc(spec.Loc()),
			Local: Convert(spec.Local(), ctx),
		}
	} else if spec.NameSpace() {
		return &ImportNamespaceSpecifier{
			Type:  "ImportNamespaceSpecifier",
			Start: start(spec.Loc()),
			End:   end(spec.Loc()),
			Loc:   loc(spec.Loc()),
			Local: Convert(spec.Local(), ctx),
		}
	}
	return &ImportSpecifier{
		Type:     "ImportSpecifier",
		Start:    start(spec.Loc()),
		End:      end(spec.Loc()),
		Loc:      loc(spec.Loc()),
		Local:    Convert(spec.Local(), ctx),
		Imported: Convert(spec.Id(), ctx),
	}
}

func importSpecs(specs []parser.Node, ctx *ConvertCtx) []Node {
	ret := make([]Node, len(specs))
	for i, spec := range specs {
		ret[i] = importSpec(spec.(*parser.ImportSpec), ctx)
	}
	return ret
}

func exportAll(node *parser.ExportDec, ctx *ConvertCtx) Node {
	var spec parser.Node
	if len(node.Specs()) == 1 {
		spec = node.Specs()[0].(*parser.ExportSpec).Local()
	}
	return &ExportAllDeclaration{
		Type:     "ExportAllDeclaration",
		Start:    start(node.Loc()),
		End:      end(node.Loc()),
		Loc:      loc(node.Loc()),
		Source:   Convert(node.Src(), ctx),
		Exported: Convert(spec, ctx),
	}
}

func exportDefault(node *parser.ExportDec, ctx *ConvertCtx) Node {
	return &ExportDefaultDeclaration{
		Type:        "ExportDefaultDeclaration",
		Start:       start(node.Loc()),
		End:         end(node.Loc()),
		Loc:         loc(node.Loc()),
		Declaration: Convert(node.Dec(), ctx),
	}
}

func exportSpecs(specs []parser.Node, ctx *ConvertCtx) []Node {
	ret := make([]Node, len(specs))
	for i, spec := range specs {
		s := spec.(*parser.ExportSpec)
		ret[i] = &ExportSpecifier{
			Type:     "ExportSpecifier",
			Start:    start(s.Loc()),
			End:      end(s.Loc()),
			Loc:      loc(s.Loc()),
			Local:    Convert(s.Local(), ctx),
			Exported: Convert(s.Id(), ctx),
		}
	}
	return ret
}

func exportNamed(node *parser.ExportDec, ctx *ConvertCtx) Node {
	return &ExportNamedDeclaration{
		Type:        "ExportNamedDeclaration",
		Start:       start(node.Loc()),
		End:         end(node.Loc()),
		Loc:         loc(node.Loc()),
		Declaration: Convert(node.Dec(), ctx),
		Specifiers:  exportSpecs(node.Specs(), ctx),
		Source:      Convert(node.Src(), ctx),
	}
}

func elems(nodes []parser.Node, ctx *ConvertCtx) []Node {
	ret := make([]Node, len(nodes))
	for i, node := range nodes {
		ret[i] = Convert(node, ctx)
	}
	return ret
}

func ident(node parser.Node, ctx *ConvertCtx) Node {
	n := node.(*parser.Ident)
	name := n.Text()
	if n.IsPrivate() {
		return &PrivateIdentifier{
			Type:  "PrivateIdentifier",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Name:  name,
		}
	}
	if n.TypInfo() != nil {
		ti := n.TypInfo()
		lc := parser.LocWithTypeInfo(n, false)
		return &TSIdentifier{
			Type:           "Identifier",
			Start:          start(lc),
			End:            end(lc),
			Loc:            loc(lc),
			Name:           name,
			Optional:       optional(ti),
			TypeAnnotation: typAnnot(ti, ctx),
		}
	}
	return &Identifier{
		Type:  "Identifier",
		Start: start(n.Loc()),
		End:   end(n.Loc()),
		Loc:   loc(n.Loc()),
		Name:  name,
	}
}

type ConvertCtx struct {
	Scope *ConvertScope
}

func NewConvertCtx() *ConvertCtx {
	return &ConvertCtx{
		Scope: &ConvertScope{},
	}
}

func (c *ConvertCtx) enter() *ConvertScope {
	scope := &ConvertScope{}
	scope.Up = c.Scope
	c.Scope = scope
	return scope
}

func (c *ConvertCtx) leave() {
	scope := c.Scope
	c.Scope = scope.Up
}

type ConvertScopeFlag uint64

const (
	CSF_NONE      ConvertScopeFlag = 0
	CSF_INTERFACE ConvertScopeFlag = 1 << iota
)

func (f ConvertScopeFlag) On(flag ConvertScopeFlag) ConvertScopeFlag {
	return f | flag
}

func (f ConvertScopeFlag) Off(flag ConvertScopeFlag) ConvertScopeFlag {
	return f & ^flag
}

type ConvertScope struct {
	Up   *ConvertScope
	Flag ConvertScopeFlag
}

func Convert(node parser.Node, ctx *ConvertCtx) Node {
	if node == nil {
		return nil
	}
	switch node.Type() {
	case parser.N_EXPR_PAREN:
		return Convert(node.(*parser.ParenExpr).Expr(), ctx)
	case parser.N_STMT_EXPR:
		exprStmt := node.(*parser.ExprStmt)
		expr := Convert(exprStmt.Expr(), ctx)
		if exprStmt.Dir() {
			return &Directive{
				Type:       "ExpressionStatement",
				Start:      start(node.Loc()),
				End:        end(node.Loc()),
				Loc:        loc(node.Loc()),
				Expression: expr,
				Directive:  exprStmt.DirStr(),
			}
		}
		return &ExpressionStatement{
			Type:       "ExpressionStatement",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: expr,
		}
	case parser.N_EXPR_NEW:
		new := node.(*parser.NewExpr)
		return &NewExpression{
			Type:      "NewExpression",
			Start:     start(node.Loc()),
			End:       end(node.Loc()),
			Loc:       loc(node.Loc()),
			Callee:    Convert(new.Callee(), ctx),
			Arguments: expressions(new.Args(), ctx),
		}
	case parser.N_NAME:
		return ident(node, ctx)
	case parser.N_EXPR_THIS:
		return &ThisExpression{
			Type:  "ThisExpression",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_LIT_NULL:
		return &Literal{
			Type:  "Literal",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_LIT_NUM:
		num := node.(*parser.NumLit)
		if num.IsBigint() {
			return &BigIntLiteral{
				Type:   "Literal",
				Start:  start(node.Loc()),
				End:    end(node.Loc()),
				Loc:    loc(node.Loc()),
				Value:  num.Float(),
				Raw:    num.Text(),
				Bigint: num.ToBigint().String(),
			}
		}
		return &Literal{
			Type:  "Literal",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Value: num.Float(),
			Raw:   num.Text(),
		}
	case parser.N_LIT_STR:
		str := node.(*parser.StrLit)
		return &Literal{
			Type:  "Literal",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Value: str.Text(),
			Raw:   str.Raw(),
		}
	case parser.N_LIT_REGEXP:
		regexp := node.(*parser.RegLit)
		return &RegExpLiteral{
			Type:   "Literal",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Regexp: &Regexp{regexp.Pattern(), regexp.Flags()},
		}
	case parser.N_EXPR_BIN:
		bin := node.(*parser.BinExpr)
		lhs := Convert(bin.Lhs(), ctx)
		rhs := Convert(bin.Rhs(), ctx)
		op := bin.OpText()
		opv := bin.Op()

		if opv == parser.T_AND || opv == parser.T_OR || opv == parser.T_NULLISH {
			return &LogicalExpression{
				Type:     "LogicalExpression",
				Start:    start(node.Loc()),
				End:      end(node.Loc()),
				Loc:      loc(node.Loc()),
				Operator: op,
				Left:     lhs,
				Right:    rhs,
			}
		}
		if opv == parser.T_TS_AS {
			return &TSAsExpression{
				Type:           "TSAsExpression",
				Start:          start(node.Loc()),
				End:            end(node.Loc()),
				Loc:            loc(node.Loc()),
				Expression:     lhs,
				TypeAnnotation: rhs,
			}
		}
		return &BinaryExpression{
			Type:     "BinaryExpression",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	case parser.N_EXPR_ASSIGN:
		bin := node.(*parser.AssignExpr)
		lhs := Convert(bin.Lhs(), ctx)
		rhs := Convert(bin.Rhs(), ctx)
		op := bin.OpText()
		return &AssignmentExpression{
			Type:     "AssignmentExpression",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	case parser.N_LIT_ARR:
		return arrExpr(node.(*parser.ArrLit), ctx)
	case parser.N_LIT_OBJ:
		return obj(node.(*parser.ObjLit), ctx)
	case parser.N_PROP:
		prop := node.(*parser.Prop)
		return &Property{
			Type:      "Property",
			Start:     start(prop.Loc()),
			End:       end(prop.Loc()),
			Loc:       loc(prop.Loc()),
			Key:       Convert(prop.Key(), ctx),
			Value:     Convert(prop.Value(), ctx),
			Kind:      prop.Kind(),
			Computed:  prop.Computed(),
			Shorthand: prop.Shorthand(),
			Method:    prop.Method(),
		}
	case parser.N_STMT_BLOCK:
		return blockStmt(node.(*parser.BlockStmt), ctx)
	case parser.N_EXPR_FN:
		n := node.(*parser.FnDec)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			return &TSFunctionExpression{
				Type:           "FunctionExpression",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Id:             Convert(n.Id(), ctx),
				Params:         fnParams(n.Params(), ctx),
				Body:           Convert(n.Body(), ctx),
				Generator:      n.Generator(),
				Async:          n.Async(),
				Expression:     false,
				TypeParameters: typParams(ti, ctx),
				ReturnType:     typAnnot(ti, ctx),
			}
		}
		return &FunctionExpression{
			Type:       "FunctionExpression",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Id:         Convert(n.Id(), ctx),
			Params:     fnParams(n.Params(), ctx),
			Body:       Convert(n.Body(), ctx),
			Generator:  n.Generator(),
			Async:      n.Async(),
			Expression: false,
		}
	case parser.N_EXPR_ARROW:
		n := node.(*parser.ArrowFn)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			lc := parser.LocWithTypeInfo(n, false)
			return &TSArrowFunctionExpression{
				Type:           "ArrowFunctionExpression",
				Start:          start(lc),
				End:            end(lc),
				Loc:            loc(lc),
				Id:             nil,
				Params:         fnParams(n.Params(), ctx),
				Body:           Convert(n.Body(), ctx),
				Generator:      false,
				Async:          n.Async(),
				Expression:     n.Expr(),
				TypeParameters: typParams(ti, ctx),
				ReturnType:     typAnnot(ti, ctx),
			}
		}
		return &ArrowFunctionExpression{
			Type:       "ArrowFunctionExpression",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Id:         nil,
			Params:     fnParams(n.Params(), ctx),
			Body:       Convert(n.Body(), ctx),
			Generator:  false,
			Async:      n.Async(),
			Expression: n.Expr(),
		}
	case parser.N_STMT_FN:
		n := node.(*parser.FnDec)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			if n.Body() == nil {
				return &TSDeclareFunction{
					Type:           "TSDeclareFunction",
					Start:          start(n.Loc()),
					End:            end(n.Loc()),
					Loc:            loc(n.Loc()),
					Id:             Convert(n.Id(), ctx),
					Params:         fnParams(n.Params(), ctx),
					Body:           Convert(n.Body(), ctx),
					Generator:      false,
					Async:          n.Async(),
					TypeParameters: typParams(ti, ctx),
					ReturnType:     typAnnot(ti, ctx),
				}
			}
			return &TSFunctionDeclaration{
				Type:           "FunctionDeclaration",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Id:             Convert(n.Id(), ctx),
				Params:         fnParams(n.Params(), ctx),
				Body:           Convert(n.Body(), ctx),
				Generator:      n.Generator(),
				Async:          n.Async(),
				TypeParameters: typParams(ti, ctx),
				ReturnType:     typAnnot(ti, ctx),
			}
		}
		return &FunctionDeclaration{
			Type:      "FunctionDeclaration",
			Start:     start(n.Loc()),
			End:       end(n.Loc()),
			Loc:       loc(n.Loc()),
			Id:        Convert(n.Id(), ctx),
			Params:    fnParams(n.Params(), ctx),
			Body:      Convert(n.Body(), ctx),
			Generator: n.Generator(),
			Async:     n.Async(),
		}
	case parser.N_EXPR_YIELD:
		node := node.(*parser.YieldExpr)
		return &YieldExpression{
			Type:     "YieldExpression",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Delegate: node.Delegate(),
			Argument: Convert(node.Arg(), ctx),
		}
	case parser.N_STMT_RET:
		ret := node.(*parser.RetStmt)
		return &ReturnStatement{
			Type:     "ReturnStatement",
			Start:    start(ret.Loc()),
			End:      end(ret.Loc()),
			Loc:      loc(ret.Loc()),
			Argument: Convert(ret.Arg(), ctx),
		}
	case parser.N_SPREAD:
		n := node.(*parser.Spread)
		return &SpreadElement{
			Type:     "SpreadElement",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Argument: Convert(n.Arg(), ctx),
		}
	case parser.N_PAT_ARRAY:
		n := node.(*parser.ArrPat)
		ti := n.TypInfo()
		if ti != nil {
			lc := parser.LocWithTypeInfo(n, false)
			return &TSArrayPattern{
				Type:           "ArrayPattern",
				Start:          start(lc),
				End:            end(lc),
				Loc:            loc(lc),
				Elements:       elems(n.Elems(), ctx),
				Optional:       ti.Optional(),
				TypeAnnotation: typAnnot(ti, ctx),
			}
		}
		return &ArrayPattern{
			Type:     "ArrayPattern",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Elements: elems(n.Elems(), ctx),
		}
	case parser.N_PAT_ASSIGN:
		n := node.(*parser.AssignPat)
		return &AssignmentPattern{
			Type:  "AssignmentPattern",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Left:  Convert(n.Left(), ctx),
			Right: Convert(n.Right(), ctx),
		}
	case parser.N_PAT_OBJ:
		n := node.(*parser.ObjPat)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			lc := parser.LocWithTypeInfo(n, false)
			return &TSObjectPattern{
				Type:           "ObjectPattern",
				Start:          start(lc),
				End:            end(lc),
				Loc:            loc(lc),
				Properties:     elems(n.Props(), ctx),
				TypeAnnotation: typAnnot(ti, ctx),
			}
		}
		return &ObjectPattern{
			Type:       "ObjectPattern",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Properties: elems(n.Props(), ctx),
		}
	case parser.N_PAT_REST:
		n := node.(*parser.RestPat)
		if n.TypInfo() != nil {
			return &TSRestElement{
				Type:           "RestElement",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Argument:       Convert(n.Arg(), ctx),
				Optional:       n.Optional(),
				TypeAnnotation: typAnnot(n.TypInfo(), ctx),
			}
		}
		return &RestElement{
			Type:     "RestElement",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Argument: Convert(n.Arg(), ctx),
		}
	case parser.N_STMT_IF:
		ifStmt := node.(*parser.IfStmt)
		return &IfStatement{
			Type:       "IfStatement",
			Start:      start(ifStmt.Loc()),
			End:        end(ifStmt.Loc()),
			Loc:        loc(ifStmt.Loc()),
			Test:       Convert(ifStmt.Test(), ctx),
			Consequent: Convert(ifStmt.Cons(), ctx),
			Alternate:  Convert(ifStmt.Alt(), ctx),
		}
	case parser.N_EXPR_CALL:
		n := node.(*parser.CallExpr)
		if n.TypInfo() != nil {
			return &TSCallExpression{
				Type:           "CallExpression",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Callee:         Convert(n.Callee(), ctx),
				Arguments:      expressions(n.Args(), ctx),
				Optional:       n.Optional(),
				TypeParameters: typArgs(n.TypInfo(), ctx),
			}
		}
		return &CallExpression{
			Type:      "CallExpression",
			Start:     start(n.Loc()),
			End:       end(n.Loc()),
			Loc:       loc(n.Loc()),
			Callee:    Convert(n.Callee(), ctx),
			Arguments: expressions(n.Args(), ctx),
			Optional:  n.Optional(),
		}
	case parser.N_STMT_SWITCH:
		swc := node.(*parser.SwitchStmt)
		return &SwitchStatement{
			Type:         "SwitchStatement",
			Start:        start(swc.Loc()),
			End:          end(swc.Loc()),
			Loc:          loc(swc.Loc()),
			Discriminant: Convert(swc.Test(), ctx),
			Cases:        cases(swc.Cases(), ctx),
		}
	case parser.N_STMT_VAR_DEC:
		varDec := node.(*parser.VarDecStmt)
		return &VariableDeclaration{
			Type:         "VariableDeclaration",
			Start:        start(varDec.Loc()),
			End:          end(varDec.Loc()),
			Loc:          loc(varDec.Loc()),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList(), ctx),
		}
	case parser.N_EXPR_MEMBER:
		mem := node.(*parser.MemberExpr)
		return &MemberExpression{
			Type:     "MemberExpression",
			Loc:      loc(mem.Loc()),
			Start:    start(mem.Loc()),
			End:      end(mem.Loc()),
			Object:   Convert(mem.Obj(), ctx),
			Property: Convert(mem.Prop(), ctx),
			Computed: mem.Compute(),
			Optional: mem.Optional(),
		}
	case parser.N_EXPR_SEQ:
		seq := node.(*parser.SeqExpr)
		return &SequenceExpression{
			Type:        "SequenceExpression",
			Start:       start(seq.Loc()),
			End:         end(seq.Loc()),
			Loc:         loc(seq.Loc()),
			Expressions: expressions(seq.Elems(), ctx),
		}
	case parser.N_EXPR_UPDATE:
		up := node.(*parser.UpdateExpr)
		return &UpdateExpression{
			Type:     "UpdateExpression",
			Start:    start(up.Loc()),
			End:      end(up.Loc()),
			Loc:      loc(up.Loc()),
			Operator: up.OpText(),
			Argument: Convert(up.Arg(), ctx),
			Prefix:   up.Prefix(),
		}
	case parser.N_EXPR_UNARY:
		un := node.(*parser.UnaryExpr)
		typ := "UnaryExpression"
		if un.Op() == parser.T_AWAIT {
			typ = "AwaitExpression"
		}
		return &UnaryExpression{
			Type:     typ,
			Start:    start(un.Loc()),
			End:      end(un.Loc()),
			Loc:      loc(un.Loc()),
			Operator: un.OpText(),
			Prefix:   true,
			Argument: Convert(un.Arg(), ctx),
		}
	case parser.N_EXPR_COND:
		cond := node.(*parser.CondExpr)
		return &ConditionalExpression{
			Type:       "ConditionalExpression",
			Start:      start(cond.Loc()),
			End:        end(cond.Loc()),
			Loc:        loc(cond.Loc()),
			Test:       Convert(cond.Test(), ctx),
			Consequent: Convert(cond.Cons(), ctx),
			Alternate:  Convert(cond.Alt(), ctx),
		}
	case parser.N_STMT_EMPTY:
		return &EmptyStatement{
			Type:  "EmptyStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_STMT_DO_WHILE:
		stmt := node.(*parser.DoWhileStmt)
		return &DoWhileStatement{
			Type:  "DoWhileStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Test:  Convert(stmt.Test(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_LIT_BOOL:
		b := node.(*parser.BoolLit)
		return &Literal{
			Type:  "Literal",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Value: b.Value(),
		}
	case parser.N_STMT_WHILE:
		stmt := node.(*parser.WhileStmt)
		return &WhileStatement{
			Type:  "WhileStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Test:  Convert(stmt.Test(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_STMT_FOR:
		stmt := node.(*parser.ForStmt)
		return &ForStatement{
			Type:   "ForStatement",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Init:   Convert(stmt.Init(), ctx),
			Test:   Convert(stmt.Test(), ctx),
			Update: Convert(stmt.Update(), ctx),
			Body:   Convert(stmt.Body(), ctx),
		}
	case parser.N_STMT_FOR_IN_OF:
		stmt := node.(*parser.ForInOfStmt)
		if stmt.In() {
			return &ForInStatement{
				Type:  "ForInStatement",
				Start: start(stmt.Loc()),
				End:   end(stmt.Loc()),
				Loc:   loc(stmt.Loc()),
				Left:  Convert(stmt.Left(), ctx),
				Right: Convert(stmt.Right(), ctx),
				Body:  Convert(stmt.Body(), ctx),
			}
		}
		return &ForOfStatement{
			Type:  "ForOfStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Left:  Convert(stmt.Left(), ctx),
			Right: Convert(stmt.Right(), ctx),
			Body:  Convert(stmt.Body(), ctx),
			Await: stmt.Await(),
		}
	case parser.N_STMT_CONT:
		stmt := node.(*parser.ContStmt)
		return &ContinueStatement{
			Type:  "ContinueStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: Convert(stmt.Label(), ctx),
		}
	case parser.N_STMT_BRK:
		stmt := node.(*parser.BrkStmt)
		return &BreakStatement{
			Type:  "BreakStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: Convert(stmt.Label(), ctx),
		}
	case parser.N_STMT_LABEL:
		stmt := node.(*parser.LabelStmt)
		return &LabeledStatement{
			Type:  "LabeledStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: Convert(stmt.Label(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_STMT_THROW:
		stmt := node.(*parser.ThrowStmt)
		return &ThrowStatement{
			Type:     "ThrowStatement",
			Start:    start(stmt.Loc()),
			End:      end(stmt.Loc()),
			Loc:      loc(stmt.Loc()),
			Argument: Convert(stmt.Arg(), ctx),
		}
	case parser.N_STMT_TRY:
		stmt := node.(*parser.TryStmt)
		return &TryStatement{
			Type:      "TryStatement",
			Start:     start(stmt.Loc()),
			End:       end(stmt.Loc()),
			Loc:       loc(stmt.Loc()),
			Block:     Convert(stmt.Try(), ctx),
			Handler:   Convert(stmt.Catch(), ctx),
			Finalizer: Convert(stmt.Fin(), ctx),
		}
	case parser.N_CATCH:
		expr := node.(*parser.Catch)
		return &CatchClause{
			Type:  "CatchClause",
			Start: start(expr.Loc()),
			End:   end(expr.Loc()),
			Loc:   loc(expr.Loc()),
			Param: Convert(expr.Param(), ctx),
			Body:  Convert(expr.Body(), ctx),
		}
	case parser.N_STMT_DEBUG:
		return &DebuggerStatement{
			Type:  "DebuggerStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_STMT_WITH:
		stmt := node.(*parser.WithStmt)
		return &WithStatement{
			Type:   "WithStatement",
			Start:  start(stmt.Loc()),
			End:    end(stmt.Loc()),
			Loc:    loc(node.Loc()),
			Object: Convert(stmt.Expr(), ctx),
			Body:   Convert(stmt.Body(), ctx),
		}
	case parser.N_IMPORT_CALL:
		stmt := node.(*parser.ImportCall)
		return &ImportExpression{
			Type:   "ImportExpression",
			Start:  start(stmt.Loc()),
			End:    end(stmt.Loc()),
			Loc:    loc(stmt.Loc()),
			Source: Convert(stmt.Src(), ctx),
		}
	case parser.N_META_PROP:
		stmt := node.(*parser.MetaProp)
		return &MetaProperty{
			Type:     "MetaProperty",
			Start:    start(stmt.Loc()),
			End:      end(stmt.Loc()),
			Loc:      loc(stmt.Loc()),
			Meta:     Convert(stmt.Meta(), ctx),
			Property: Convert(stmt.Prop(), ctx),
		}
	case parser.N_STMT_CLASS:
		stmt := node.(*parser.ClassDec)
		superTypArgs := stmt.SuperTypArgs()
		typParams := stmt.TypParams()
		implements := stmt.Implements()
		if superTypArgs != nil || typParams != nil || implements != nil {
			return &TSClassDeclaration{
				Type:                "ClassDeclaration",
				Start:               start(stmt.Loc()),
				End:                 end(stmt.Loc()),
				Loc:                 loc(stmt.Loc()),
				Id:                  Convert(stmt.Id(), ctx),
				TypeParameters:      ConvertTsTyp(typParams, ctx),
				SuperClass:          Convert(stmt.Super(), ctx),
				SuperTypeParameters: ConvertTsTyp(superTypArgs, ctx),
				Implements:          elems(implements, ctx),
				Body:                Convert(stmt.Body(), ctx),
				Abstract:            stmt.Abstract(),
			}
		}
		return &ClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Id:         Convert(stmt.Id(), ctx),
			SuperClass: Convert(stmt.Super(), ctx),
			Body:       Convert(stmt.Body(), ctx),
			Abstract:   stmt.Abstract(),
		}
	case parser.N_EXPR_CLASS:
		stmt := node.(*parser.ClassDec)
		superTypArgs := stmt.SuperTypArgs()
		typParams := stmt.TypParams()
		implements := stmt.Implements()
		if superTypArgs != nil || typParams != nil || implements != nil {
			return &TSClassExpression{
				Type:                "ClassExpression",
				Start:               start(stmt.Loc()),
				End:                 end(stmt.Loc()),
				Loc:                 loc(stmt.Loc()),
				Id:                  Convert(stmt.Id(), ctx),
				TypeParameters:      ConvertTsTyp(typParams, ctx),
				SuperClass:          Convert(stmt.Super(), ctx),
				SuperTypeParameters: ConvertTsTyp(superTypArgs, ctx),
				Implements:          elems(implements, ctx),
				Body:                Convert(stmt.Body(), ctx),
				Abstract:            stmt.Abstract(),
			}
		}
		return &ClassExpression{
			Type:       "ClassExpression",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Id:         Convert(stmt.Id(), ctx),
			SuperClass: Convert(stmt.Super(), ctx),
			Body:       Convert(stmt.Body(), ctx),
			Abstract:   stmt.Abstract(),
		}
	case parser.N_ClASS_BODY:
		stmt := node.(*parser.ClassBody)
		return &ClassBody{
			Type:  "ClassBody",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Body:  expressions(stmt.Elems(), ctx),
		}
	case parser.N_METHOD:
		n := node.(*parser.Method)
		f := n.Val().(*parser.FnDec)
		if f.TypInfo() != nil {
			ti := f.TypInfo()
			return &TSMethodDefinition{
				Type:          "MethodDefinition",
				Start:         start(n.Loc()),
				End:           end(n.Loc()),
				Loc:           loc(n.Loc()),
				Key:           Convert(n.Key(), ctx),
				Value:         Convert(n.Val(), ctx),
				Kind:          n.Kind(),
				Computed:      n.Computed(),
				Static:        n.Static(),
				Optional:      ti.Optional(),
				Definite:      ti.Definite(),
				Abstract:      ti.Abstract(),
				Override:      ti.Override(),
				Readonly:      ti.Readonly(),
				Accessibility: ti.AccMod().String(),
			}
		}
		return &MethodDefinition{
			Type:     "MethodDefinition",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Key:      Convert(n.Key(), ctx),
			Value:    Convert(n.Val(), ctx),
			Kind:     n.Kind(),
			Computed: n.Computed(),
			Static:   n.Static(),
		}
	case parser.N_FIELD:
		n := node.(*parser.Field)
		if n.TypInfo() != nil {
			ti := n.TypInfo()

			if n.IsTsSig() {
				return &TSIndexSignature{
					Type:           "TSIndexSignature",
					Start:          start(n.Loc()),
					End:            end(n.Loc()),
					Loc:            loc(n.Loc()),
					Static:         n.Static(),
					Abstract:       ti.Abstract(),
					Optional:       ti.Optional(),
					Declare:        ti.Declare(),
					Readonly:       ti.Readonly(),
					Accessibility:  ti.AccMod().String(),
					Parameters:     elems([]parser.Node{n.Key()}, ctx),
					TypeAnnotation: typAnnot(ti, ctx),
				}
			}

			return &TSPropertyDefinition{
				Type:           "PropertyDefinition",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Key:            Convert(n.Key(), ctx),
				Value:          Convert(n.Value(), ctx),
				Computed:       n.Computed(),
				Static:         n.Static(),
				Abstract:       ti.Abstract(),
				Optional:       ti.Optional(),
				Definite:       ti.Definite(),
				Readonly:       ti.Readonly(),
				Override:       ti.Override(),
				Declare:        ti.Declare(),
				Accessibility:  ti.AccMod().String(),
				TypeAnnotation: typAnnot(ti, ctx),
			}
		}
		return &PropertyDefinition{
			Type:     "PropertyDefinition",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Key:      Convert(n.Key(), ctx),
			Value:    Convert(n.Value(), ctx),
			Computed: n.Computed(),
			Static:   n.Static(),
		}
	case parser.N_SUPER:
		n := node.(*parser.Super)
		return &Super{
			Type:  "Super",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
		}
	case parser.N_EXPR_TPL:
		tpl := node.(*parser.TplExpr)
		if tpl.Tag() == nil {
			return tplLiteral(tpl.Loc(), tpl.Elems(), ctx)
		}
		return &TaggedTemplateExpression{
			Type:  "TaggedTemplateExpression",
			Start: start(tpl.LocWithTag()),
			End:   end(tpl.LocWithTag()),
			Loc:   loc(tpl.LocWithTag()),
			Tag:   Convert(tpl.Tag(), ctx),
			Quasi: tplLiteral(tpl.Loc(), tpl.Elems(), ctx),
		}
	case parser.N_STMT_IMPORT:
		stmt := node.(*parser.ImportDec)
		return &ImportDeclaration{
			Type:       "ImportDeclaration",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Specifiers: importSpecs(stmt.Specs(), ctx),
			Source:     Convert(stmt.Src(), ctx),
		}
	case parser.N_STMT_EXPORT:
		stmt := node.(*parser.ExportDec)
		if stmt.All() {
			return exportAll(stmt, ctx)
		}
		if stmt.Default() {
			return exportDefault(stmt, ctx)
		}
		if stmt.Dec() != nil {
			if stmt.Dec().Type() == parser.N_TS_NAMESPACE {
				n := stmt.Dec().(*parser.TsNS)
				if n.Alias() {
					return &TSNamespaceExportDeclaration{
						Type:  "TSNamespaceExportDeclaration",
						Start: start(node.Loc()),
						End:   end(node.Loc()),
						Loc:   loc(node.Loc()),
						Id:    Convert(n.Id(), ctx),
					}
				}
			}
			if stmt.Dec().Type() == parser.N_TS_IMPORT_REQUIRE {
				n := ConvertTsTyp(stmt.Dec(), ctx).(*TSImportEqualsDeclaration)
				n.IsExport = true
				return n
			}
		}
		return exportNamed(stmt, ctx)
	case parser.N_STATIC_BLOCK:
		stmt := node.(*parser.StaticBlock)
		return &StaticBlock{
			Type:  "StaticBlock",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Body:  statements(stmt.Body(), ctx),
		}
	case parser.N_EXPR_CHAIN:
		node := node.(*parser.ChainExpr)
		return &ChainExpression{
			Type:       "ChainExpression",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_ID:
		node := node.(*parser.JsxIdent)
		return &JSXIdentifier{
			Type:  "JSXIdentifier",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Name:  node.Text(),
		}
	case parser.N_JSX_NS:
		node := node.(*parser.JsxNsName)
		return &JSXNamespacedName{
			Type:      "JSXNamespacedName",
			Start:     start(node.Loc()),
			End:       end(node.Loc()),
			Loc:       loc(node.Loc()),
			Namespace: node.NS(),
			Name:      node.Name(),
		}
	case parser.N_JSX_MEMBER:
		node := node.(*parser.JsxMemberExpr)
		return &JSXMemberExpression{
			Type:     "JSXMemberExpression",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Object:   Convert(node.Obj(), ctx),
			Property: Convert(node.Prop(), ctx),
		}
	case parser.N_JSX_ELEM:
		node := node.(*parser.JsxElem)
		if node.IsFragment() {
			open := node.Open().(*parser.JsxOpen)
			close := node.Close().(*parser.JsxClose)
			return &JSXFragment{
				Type:  "JSXFragment",
				Start: start(node.Loc()),
				End:   end(node.Loc()),
				Loc:   loc(node.Loc()),
				OpeningFragment: &JSXOpeningFragment{
					Type:  "JSXOpeningFragment",
					Start: start(open.Loc()),
					End:   end(open.Loc()),
					Loc:   loc(open.Loc()),
				},
				Children: elems(node.Children(), ctx),
				ClosingFragment: &JSXClosingFragment{
					Type:  "JSXClosingFragment",
					Start: start(close.Loc()),
					End:   end(close.Loc()),
					Loc:   loc(close.Loc()),
				},
			}
		}
		return &JSXElement{
			Type:           "JSXElement",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			OpeningElement: Convert(node.Open(), ctx),
			Children:       elems(node.Children(), ctx),
			ClosingElement: Convert(node.Close(), ctx),
		}
	case parser.N_JSX_OPEN:
		node := node.(*parser.JsxOpen)
		return &JSXOpeningElement{
			Type:        "JSXOpeningElement",
			Start:       start(node.Loc()),
			End:         end(node.Loc()),
			Loc:         loc(node.Loc()),
			Name:        Convert(node.Name(), ctx),
			Attributes:  elems(node.Attrs(), ctx),
			SelfClosing: node.Closed(),
		}
	case parser.N_JSX_CLOSE:
		node := node.(*parser.JsxClose)
		return &JSXClosingElement{
			Type:  "JSXClosingElement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Name:  Convert(node.Name(), ctx),
		}
	case parser.N_JSX_TXT:
		node := node.(*parser.JsxText)
		return &JSXText{
			Type:  "JSXText",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Value: node.Value(),
			Raw:   node.Raw(),
		}
	case parser.N_JSX_EXPR_SPAN:
		node := node.(*parser.JsxExprSpan)
		return &JSXExpressionContainer{
			Type:       "JSXExpressionContainer",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_ATTR:
		node := node.(*parser.JsxAttr)
		return &JSXAttribute{
			Type:  "JSXAttribute",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Name:  Convert(node.Name(), ctx),
			Value: Convert(node.Value(), ctx),
		}
	case parser.N_JSX_ATTR_SPREAD:
		node := node.(*parser.JsxSpreadAttr)
		return &JSXSpreadAttribute{
			Type:     "JSXSpreadAttribute",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Argument: Convert(node.Arg(), ctx),
		}
	case parser.N_JSX_CHILD_SPREAD:
		node := node.(*parser.JsxSpreadChild)
		return &JSXSpreadChild{
			Type:       "JSXSpreadChild",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_EMPTY:
		node := node.(*parser.JsxEmpty)
		return &JSXEmptyExpression{
			Type:  "JSXEmptyExpression",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	}

	// bypass the ts related node like `Ambient`
	return ConvertTsTyp(node, ctx)
}
