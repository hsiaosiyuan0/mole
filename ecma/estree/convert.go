package estree

import (
	"math"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

func pos(p span.Pos) Position {
	return Position{Line: int(p.Line), Column: int(p.Col)}
}

func locOfRng(rng span.Range, s *span.Source, ctx *ConvertCtx) *SrcLoc {
	if ctx.LineCol {
		begin, end := s.LineCol(rng)
		return &SrcLoc{
			Start: pos(begin),
			End:   pos(end),
		}
	}
	return nil
}

func locOfNode(n parser.Node, s *span.Source, ctx *ConvertCtx) *SrcLoc {
	return locOfRng(n.Range(), s, ctx)
}

func ConvertProg(n *parser.Prog, ctx *ConvertCtx) *Program {
	stmts := n.Body()
	body := make([]Node, len(stmts))
	for i, s := range stmts {
		body[i] = Convert(s, ctx)
	}
	return &Program{
		Type:  "Program",
		Start: int(n.Range().Lo),
		End:   int(n.Range().Hi),
		Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
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
		Start:    int(n.Range().Lo),
		End:      int(n.Range().Hi),
		Loc:      locOfNode(n, ctx.Parser.Source(), ctx),
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
		Start:      int(n.Range().Lo),
		End:        int(n.Range().Hi),
		Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
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
		Start: int(n.Range().Lo),
		End:   int(n.Range().Hi),
		Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
		Body:  statements(n.Body(), ctx),
	}
}

func cases(cs []parser.Node, ctx *ConvertCtx) []*SwitchCase {
	s := make([]*SwitchCase, len(cs))
	for i, c := range cs {
		sc := c.(*parser.SwitchCase)
		s[i] = &SwitchCase{
			Type:       "SwitchCase",
			Start:      int(sc.Range().Lo),
			End:        int(sc.Range().Hi),
			Loc:        locOfNode(sc, ctx.Parser.Source(), ctx),
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
			Start: int(dc.Range().Lo),
			End:   int(dc.Range().Hi),
			Loc:   locOfNode(dc, ctx.Parser.Source(), ctx),
			Id:    Convert(dc.Id(), ctx),
			Init:  Convert(dc.Init(), ctx),
		}
	}
	return s
}

func tplLiteral(tplLoc span.Range, elems []parser.Node, ctx *ConvertCtx) *TemplateLiteral {
	quasis := make([]Expression, 0)
	exprs := make([]Expression, 0)
	cnt := len(elems)
	for i, elem := range elems {
		first := i == 0
		last := i == cnt-1
		if elem.Type() != parser.N_LIT_STR {
			if first || last {
				lc := locOfRng(tplLoc, ctx.Parser.Source(), ctx)
				if first {
					lc.End.Column = lc.Start.Column
				} else {
					lc.Start.Column = lc.End.Column
				}
				quasis = append(quasis, &TemplateElement{
					Type:  "TemplateElement",
					Start: int(tplLoc.Lo),
					End:   int(tplLoc.Hi),
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
				Start: int(elem.Range().Lo),
				End:   int(elem.Range().Hi),
				Loc:   locOfNode(elem, ctx.Parser.Source(), ctx),
				Tail:  last,
				Value: &TemplateElementValue{
					Cooked: ctx.Parser.NodeText(str),
					Raw:    ctx.Parser.RngText(str.Range()),
				},
			})
		}
	}
	return &TemplateLiteral{
		Type:        "TemplateLiteral",
		Start:       int(tplLoc.Lo),
		End:         int(tplLoc.Hi),
		Loc:         locOfRng(tplLoc, ctx.Parser.Source(), ctx),
		Quasis:      quasis,
		Expressions: exprs,
	}
}

func importSpec(spec *parser.ImportSpec, ctx *ConvertCtx) Node {
	if spec.Default() {
		return &ImportDefaultSpecifier{
			Type:  "ImportDefaultSpecifier",
			Start: int(spec.Range().Lo),
			End:   int(spec.Range().Hi),
			Loc:   locOfNode(spec, ctx.Parser.Source(), ctx),
			Local: Convert(spec.Local(), ctx),
		}
	} else if spec.NameSpace() {
		return &ImportNamespaceSpecifier{
			Type:  "ImportNamespaceSpecifier",
			Start: int(spec.Range().Lo),
			End:   int(spec.Range().Hi),
			Loc:   locOfNode(spec, ctx.Parser.Source(), ctx),
			Local: Convert(spec.Local(), ctx),
		}
	}
	return &ImportSpecifier{
		Type:       "ImportSpecifier",
		Start:      int(spec.Range().Lo),
		End:        int(spec.Range().Hi),
		Loc:        locOfNode(spec, ctx.Parser.Source(), ctx),
		Local:      Convert(spec.Local(), ctx),
		Imported:   Convert(spec.Id(), ctx),
		ImportKind: spec.Kind(),
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
		Start:    int(node.Range().Lo),
		End:      int(node.Range().Hi),
		Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
		Source:   Convert(node.Src(), ctx),
		Exported: Convert(spec, ctx),
	}
}

func exportDefault(node *parser.ExportDec, ctx *ConvertCtx) Node {
	return &ExportDefaultDeclaration{
		Type:        "ExportDefaultDeclaration",
		Start:       int(node.Range().Lo),
		End:         int(node.Range().Hi),
		Loc:         locOfNode(node, ctx.Parser.Source(), ctx),
		Declaration: Convert(node.Dec(), ctx),
	}
}

func exportSpecs(specs []parser.Node, ctx *ConvertCtx) []Node {
	ret := make([]Node, len(specs))
	for i, spec := range specs {
		s := spec.(*parser.ExportSpec)
		ret[i] = &ExportSpecifier{
			Type:       "ExportSpecifier",
			Start:      int(s.Range().Lo),
			End:        int(s.Range().Hi),
			Loc:        locOfNode(s, ctx.Parser.Source(), ctx),
			Local:      Convert(s.Local(), ctx),
			Exported:   Convert(s.Id(), ctx),
			ExportKind: s.Kind(),
		}
	}
	return ret
}

func exportNamed(node *parser.ExportDec, ctx *ConvertCtx) Node {
	return &ExportNamedDeclaration{
		Type:        "ExportNamedDeclaration",
		Start:       int(node.Range().Lo),
		End:         int(node.Range().Hi),
		Loc:         locOfNode(node, ctx.Parser.Source(), ctx),
		Declaration: Convert(node.Dec(), ctx),
		Specifiers:  exportSpecs(node.Specs(), ctx),
		Source:      Convert(node.Src(), ctx),
		ExportKind:  node.Kind(),
	}
}

func elems(nodes []parser.Node, ctx *ConvertCtx) []Node {
	ret := make([]Node, len(nodes))
	for i, node := range nodes {
		ret[i] = Convert(node, ctx)
	}
	return ret
}

func FirstLoc(s *span.Source, ctx *ConvertCtx, locs ...span.Range) (span.Range, *SrcLoc) {
	k := 0
	lo := uint32(math.MaxUint32)
	hi := uint32(math.MaxUint32)
	for i, loc := range locs {
		if loc.Empty() {
			continue
		}
		if loc.Lo < lo || (loc.Lo == lo && loc.Hi < hi) {
			lo = loc.Lo
			hi = loc.Hi
			k = i
		}
	}
	rng := locs[k]
	return rng, locOfRng(rng, s, ctx)
}

func LastLoc(s *span.Source, ctx *ConvertCtx, locs ...span.Range) (span.Range, *SrcLoc) {
	k := 0
	lo := uint32(0)
	hi := uint32(0)
	for i, loc := range locs {
		if loc.Empty() {
			continue
		}
		if loc.Hi > hi || (loc.Hi == hi && loc.Lo > lo) {
			lo = loc.Lo
			hi = loc.Hi
			k = i
		}
	}
	rng := locs[k]
	return rng, locOfRng(rng, s, ctx)
}

func CalcLoc(node parser.Node, s *span.Source, ctx *ConvertCtx) *SrcLoc {
	if util.IsNilPtr(node) {
		return nil
	}
	return locOfNode(node, s, ctx)
}

func locWithTypeInfo(node parser.Node, includeParamProp bool, s *span.Source, ctx *ConvertCtx) (rng span.Range, loc *SrcLoc) {
	nw, ok := node.(parser.NodeWithTypInfo)
	if !ok {
		return node.Range(), CalcLoc(node, s, ctx)
	}

	ti := nw.TypInfo()
	loc = CalcLoc(node, s, ctx)

	starLocList := []span.Range{}
	if ti.TypParams() != nil {
		starLocList = append(starLocList, ti.TypParams().Range())
	}
	starLocList = append(starLocList, node.Range())

	if includeParamProp {
		starLocList = append(starLocList, ti.BeginRng())
	}

	firstRng, firstLoc := FirstLoc(s, ctx, starLocList...)
	rng.Lo = firstRng.Lo
	loc.Start = firstLoc.Start

	endLocList := []span.Range{}
	if ti.TypAnnot() != nil {
		endLocList = append(endLocList, ti.TypAnnot().Range())
	}
	endLocList = append(endLocList, []span.Range{ti.Ques(), node.Range()}...)

	endRng, endLoc := LastLoc(s, ctx, endLocList...)
	rng.Hi = endRng.Hi
	loc.End = endLoc.End

	return
}

func TplLocWithTag(n *parser.TplExpr, s *span.Source, ctx *ConvertCtx) *SrcLoc {
	loc := CalcLoc(n, s, ctx)
	if n.Tag() != nil {
		tl := CalcLoc(n.Tag(), s, ctx)
		loc.Start = tl.Start
		loc.End = tl.End
	}
	return loc
}

func ident(node parser.Node, ctx *ConvertCtx) Node {
	n := node.(*parser.Ident)
	name := ctx.Parser.NodeText(node)
	if n.IsPrivate() {
		return &PrivateIdentifier{
			Type:  "PrivateIdentifier",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
			Name:  name,
		}
	}
	if n.TypInfo() != nil {
		ti := n.TypInfo()
		rng, loc := locWithTypeInfo(n, false, ctx.Parser.Source(), ctx)
		return &TSIdentifier{
			Type:           "Identifier",
			Start:          int(rng.Lo),
			End:            int(rng.Hi),
			Loc:            loc,
			Name:           name,
			Optional:       optional(ti),
			TypeAnnotation: typAnnot(ti, ctx),
			Decorators:     elems(parser.DecoratorsOf(n), ctx),
		}
	}
	return &Identifier{
		Type:  "Identifier",
		Start: int(n.Range().Lo),
		End:   int(n.Range().Hi),
		Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
		Name:  name,
	}
}

func locWithDecorator(node parser.Node, ds []parser.Node, s *span.Source, ctx *ConvertCtx) (rng span.Range, loc *SrcLoc) {
	rng = node.Range()
	loc = locOfNode(node, s, ctx)
	if len(ds) == 0 {
		return
	}
	d := ds[0]
	dLoc := locOfNode(d, s, ctx)
	loc.Start = dLoc.Start
	rng.Lo = d.Range().Lo
	return
}

type ConvertCtx struct {
	Parser  *parser.Parser
	Scope   *ConvertScope
	LineCol bool
}

func NewConvertCtx(p *parser.Parser) *ConvertCtx {
	return &ConvertCtx{
		Parser:  p,
		Scope:   &ConvertScope{},
		LineCol: true,
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
				Start:      int(node.Range().Lo),
				End:        int(node.Range().Hi),
				Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
				Expression: expr,
				Directive:  exprStmt.Expr().(*parser.StrLit).Val(),
			}
		}
		return &ExpressionStatement{
			Type:       "ExpressionStatement",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: expr,
		}
	case parser.N_EXPR_NEW:
		n := node.(*parser.NewExpr)
		if n.TypInfo() != nil {
			return &TSNewExpression{
				Type:           "NewExpression",
				Start:          int(n.Range().Lo),
				End:            int(n.Range().Hi),
				Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
				Callee:         Convert(n.Callee(), ctx),
				Arguments:      expressions(n.Args(), ctx),
				TypeParameters: typArgs(n.TypInfo(), ctx),
			}
		}
		return &NewExpression{
			Type:      "NewExpression",
			Start:     int(node.Range().Lo),
			End:       int(node.Range().Hi),
			Loc:       locOfNode(node, ctx.Parser.Source(), ctx),
			Callee:    Convert(n.Callee(), ctx),
			Arguments: expressions(n.Args(), ctx),
		}
	case parser.N_NAME:
		return ident(node, ctx)
	case parser.N_EXPR_THIS:
		return &ThisExpression{
			Type:  "ThisExpression",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_LIT_NULL:
		return &Literal{
			Type:  "Literal",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_LIT_NUM:
		num := node.(*parser.NumLit)
		t := ctx.Parser.NodeText(node)

		f := parser.NodeToFloat(node, ctx.Parser.Source())
		if math.IsInf(f, 0) {
			f = 0
		}

		if parser.NodeIsBigint(num, ctx.Parser.Source()) {
			return &BigIntLiteral{
				Type:   "Literal",
				Start:  int(node.Range().Lo),
				End:    int(node.Range().Hi),
				Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
				Value:  f,
				Raw:    t,
				Bigint: parser.ParseBigint(t).String(),
			}
		}
		return &Literal{
			Type:  "Literal",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Value: f,
			Raw:   t,
		}
	case parser.N_LIT_STR:
		str := node.(*parser.StrLit)
		return &Literal{
			Type:  "Literal",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Value: parser.NodeText(str, ctx.Parser.Source()),
			Raw:   ctx.Parser.RngText(str.Range()),
		}
	case parser.N_LIT_REGEXP:
		regexp := node.(*parser.RegLit)
		return &RegExpLiteral{
			Type:   "Literal",
			Start:  int(node.Range().Lo),
			End:    int(node.Range().Hi),
			Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
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
				Start:    int(node.Range().Lo),
				End:      int(node.Range().Hi),
				Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
				Operator: op,
				Left:     lhs,
				Right:    rhs,
			}
		}
		if opv == parser.T_TS_AS {
			return &TSAsExpression{
				Type:           "TSAsExpression",
				Start:          int(node.Range().Lo),
				End:            int(node.Range().Hi),
				Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
				Expression:     lhs,
				TypeAnnotation: rhs,
			}
		}
		return &BinaryExpression{
			Type:     "BinaryExpression",
			Start:    int(node.Range().Lo),
			End:      int(node.Range().Hi),
			Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	case parser.N_EXPR_ASSIGN:
		bin := node.(*parser.AssignExpr)
		lhs := Convert(bin.Lhs(), ctx)
		rhs := Convert(bin.Rhs(), ctx)
		op := bin.OpName()
		return &AssignmentExpression{
			Type:     "AssignmentExpression",
			Start:    int(node.Range().Lo),
			End:      int(node.Range().Hi),
			Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
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
			Start:     int(prop.Range().Lo),
			End:       int(prop.Range().Hi),
			Loc:       locOfNode(prop, ctx.Parser.Source(), ctx),
			Key:       Convert(prop.Key(), ctx),
			Value:     Convert(prop.Val(), ctx),
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
				Start:          int(n.Range().Lo),
				End:            int(n.Range().Hi),
				Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:      int(n.Range().Lo),
			End:        int(n.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
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
			rng, loc := locWithTypeInfo(n, false, ctx.Parser.Source(), ctx)
			return &TSArrowFunctionExpression{
				Type:           "ArrowFunctionExpression",
				Start:          int(rng.Lo),
				End:            int(rng.Hi),
				Loc:            loc,
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
			Start:      int(n.Range().Lo),
			End:        int(n.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
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
					Start:          int(n.Range().Lo),
					End:            int(n.Range().Hi),
					Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
				Start:          int(n.Range().Lo),
				End:            int(n.Range().Hi),
				Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:     int(n.Range().Lo),
			End:       int(n.Range().Hi),
			Loc:       locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:    int(node.Range().Lo),
			End:      int(node.Range().Hi),
			Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
			Delegate: node.Delegate(),
			Argument: Convert(node.Arg(), ctx),
		}
	case parser.N_STMT_RET:
		ret := node.(*parser.RetStmt)
		return &ReturnStatement{
			Type:     "ReturnStatement",
			Start:    int(ret.Range().Lo),
			End:      int(ret.Range().Hi),
			Loc:      locOfNode(ret, ctx.Parser.Source(), ctx),
			Argument: Convert(ret.Arg(), ctx),
		}
	case parser.N_SPREAD:
		n := node.(*parser.Spread)
		return &SpreadElement{
			Type:     "SpreadElement",
			Start:    int(n.Range().Lo),
			End:      int(n.Range().Hi),
			Loc:      locOfNode(n, ctx.Parser.Source(), ctx),
			Argument: Convert(n.Arg(), ctx),
		}
	case parser.N_PAT_ARRAY:
		n := node.(*parser.ArrPat)
		ti := n.TypInfo()
		if ti != nil {
			rng, loc := locWithTypeInfo(n, false, ctx.Parser.Source(), ctx)
			return &TSArrayPattern{
				Type:           "ArrayPattern",
				Start:          int(rng.Lo),
				End:            int(rng.Hi),
				Loc:            loc,
				Elements:       elems(n.Elems(), ctx),
				Optional:       ti.Optional(),
				TypeAnnotation: typAnnot(ti, ctx),
			}
		}
		return &ArrayPattern{
			Type:     "ArrayPattern",
			Start:    int(n.Range().Lo),
			End:      int(n.Range().Hi),
			Loc:      locOfNode(n, ctx.Parser.Source(), ctx),
			Elements: elems(n.Elems(), ctx),
		}
	case parser.N_PAT_ASSIGN:
		n := node.(*parser.AssignPat)
		return &AssignmentPattern{
			Type:  "AssignmentPattern",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
			Left:  Convert(n.Lhs(), ctx),
			Right: Convert(n.Rhs(), ctx),
		}
	case parser.N_PAT_OBJ:
		n := node.(*parser.ObjPat)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			rng, loc := locWithTypeInfo(n, false, ctx.Parser.Source(), ctx)
			return &TSObjectPattern{
				Type:           "ObjectPattern",
				Start:          int(rng.Lo),
				End:            int(rng.Hi),
				Loc:            loc,
				Properties:     elems(n.Props(), ctx),
				TypeAnnotation: typAnnot(ti, ctx),
			}
		}
		return &ObjectPattern{
			Type:       "ObjectPattern",
			Start:      int(n.Range().Lo),
			End:        int(n.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
			Properties: elems(n.Props(), ctx),
		}
	case parser.N_PAT_REST:
		n := node.(*parser.RestPat)
		if n.TypInfo() != nil {
			return &TSRestElement{
				Type:           "RestElement",
				Start:          int(n.Range().Lo),
				End:            int(n.Range().Hi),
				Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
				Argument:       Convert(n.Arg(), ctx),
				Optional:       n.Optional(),
				TypeAnnotation: typAnnot(n.TypInfo(), ctx),
			}
		}
		return &RestElement{
			Type:     "RestElement",
			Start:    int(n.Range().Lo),
			End:      int(n.Range().Hi),
			Loc:      locOfNode(n, ctx.Parser.Source(), ctx),
			Argument: Convert(n.Arg(), ctx),
		}
	case parser.N_STMT_IF:
		ifStmt := node.(*parser.IfStmt)
		return &IfStatement{
			Type:       "IfStatement",
			Start:      int(ifStmt.Range().Lo),
			End:        int(ifStmt.Range().Hi),
			Loc:        locOfNode(ifStmt, ctx.Parser.Source(), ctx),
			Test:       Convert(ifStmt.Test(), ctx),
			Consequent: Convert(ifStmt.Cons(), ctx),
			Alternate:  Convert(ifStmt.Alt(), ctx),
		}
	case parser.N_EXPR_CALL:
		n := node.(*parser.CallExpr)
		if n.TypInfo() != nil {
			return &TSCallExpression{
				Type:           "CallExpression",
				Start:          int(n.Range().Lo),
				End:            int(n.Range().Hi),
				Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
				Callee:         Convert(n.Callee(), ctx),
				Arguments:      expressions(n.Args(), ctx),
				Optional:       n.Optional(),
				TypeParameters: typArgs(n.TypInfo(), ctx),
			}
		}
		return &CallExpression{
			Type:      "CallExpression",
			Start:     int(n.Range().Lo),
			End:       int(n.Range().Hi),
			Loc:       locOfNode(n, ctx.Parser.Source(), ctx),
			Callee:    Convert(n.Callee(), ctx),
			Arguments: expressions(n.Args(), ctx),
			Optional:  n.Optional(),
		}
	case parser.N_STMT_SWITCH:
		swc := node.(*parser.SwitchStmt)
		return &SwitchStatement{
			Type:         "SwitchStatement",
			Start:        int(swc.Range().Lo),
			End:          int(swc.Range().Hi),
			Loc:          locOfNode(swc, ctx.Parser.Source(), ctx),
			Discriminant: Convert(swc.Test(), ctx),
			Cases:        cases(swc.Cases(), ctx),
		}
	case parser.N_STMT_VAR_DEC:
		varDec := node.(*parser.VarDecStmt)
		return &VariableDeclaration{
			Type:         "VariableDeclaration",
			Start:        int(varDec.Range().Lo),
			End:          int(varDec.Range().Hi),
			Loc:          locOfNode(varDec, ctx.Parser.Source(), ctx),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList(), ctx),
		}
	case parser.N_EXPR_MEMBER:
		mem := node.(*parser.MemberExpr)
		return &MemberExpression{
			Type:     "MemberExpression",
			Loc:      locOfNode(mem, ctx.Parser.Source(), ctx),
			Start:    int(mem.Range().Lo),
			End:      int(mem.Range().Hi),
			Object:   Convert(mem.Obj(), ctx),
			Property: Convert(mem.Prop(), ctx),
			Computed: mem.Compute(),
			Optional: mem.Optional(),
		}
	case parser.N_EXPR_SEQ:
		seq := node.(*parser.SeqExpr)
		return &SequenceExpression{
			Type:        "SequenceExpression",
			Start:       int(seq.Range().Lo),
			End:         int(seq.Range().Hi),
			Loc:         locOfNode(seq, ctx.Parser.Source(), ctx),
			Expressions: expressions(seq.Elems(), ctx),
		}
	case parser.N_EXPR_UPDATE:
		up := node.(*parser.UpdateExpr)
		return &UpdateExpression{
			Type:     "UpdateExpression",
			Start:    int(up.Range().Lo),
			End:      int(up.Range().Hi),
			Loc:      locOfNode(up, ctx.Parser.Source(), ctx),
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
			Start:    int(un.Range().Lo),
			End:      int(un.Range().Hi),
			Loc:      locOfNode(un, ctx.Parser.Source(), ctx),
			Operator: un.OpText(),
			Prefix:   true,
			Argument: Convert(un.Arg(), ctx),
		}
	case parser.N_EXPR_COND:
		cond := node.(*parser.CondExpr)
		return &ConditionalExpression{
			Type:       "ConditionalExpression",
			Start:      int(cond.Range().Lo),
			End:        int(cond.Range().Hi),
			Loc:        locOfNode(cond, ctx.Parser.Source(), ctx),
			Test:       Convert(cond.Test(), ctx),
			Consequent: Convert(cond.Cons(), ctx),
			Alternate:  Convert(cond.Alt(), ctx),
		}
	case parser.N_STMT_EMPTY:
		return &EmptyStatement{
			Type:  "EmptyStatement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_STMT_DO_WHILE:
		stmt := node.(*parser.DoWhileStmt)
		return &DoWhileStatement{
			Type:  "DoWhileStatement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Test:  Convert(stmt.Test(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_LIT_BOOL:
		b := node.(*parser.BoolLit)
		return &Literal{
			Type:  "Literal",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Value: b.Val(),
		}
	case parser.N_STMT_WHILE:
		stmt := node.(*parser.WhileStmt)
		return &WhileStatement{
			Type:  "WhileStatement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Test:  Convert(stmt.Test(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_STMT_FOR:
		stmt := node.(*parser.ForStmt)
		return &ForStatement{
			Type:   "ForStatement",
			Start:  int(node.Range().Lo),
			End:    int(node.Range().Hi),
			Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
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
				Start: int(stmt.Range().Lo),
				End:   int(stmt.Range().Hi),
				Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
				Left:  Convert(stmt.Left(), ctx),
				Right: Convert(stmt.Right(), ctx),
				Body:  Convert(stmt.Body(), ctx),
			}
		}
		return &ForOfStatement{
			Type:  "ForOfStatement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Left:  Convert(stmt.Left(), ctx),
			Right: Convert(stmt.Right(), ctx),
			Body:  Convert(stmt.Body(), ctx),
			Await: stmt.Await(),
		}
	case parser.N_STMT_CONT:
		stmt := node.(*parser.ContStmt)
		return &ContinueStatement{
			Type:  "ContinueStatement",
			Start: int(stmt.Range().Lo),
			End:   int(stmt.Range().Hi),
			Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
			Label: Convert(stmt.Label(), ctx),
		}
	case parser.N_STMT_BRK:
		stmt := node.(*parser.BrkStmt)
		return &BreakStatement{
			Type:  "BreakStatement",
			Start: int(stmt.Range().Lo),
			End:   int(stmt.Range().Hi),
			Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
			Label: Convert(stmt.Label(), ctx),
		}
	case parser.N_STMT_LABEL:
		stmt := node.(*parser.LabelStmt)
		return &LabeledStatement{
			Type:  "LabeledStatement",
			Start: int(stmt.Range().Lo),
			End:   int(stmt.Range().Hi),
			Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
			Label: Convert(stmt.Label(), ctx),
			Body:  Convert(stmt.Body(), ctx),
		}
	case parser.N_STMT_THROW:
		stmt := node.(*parser.ThrowStmt)
		return &ThrowStatement{
			Type:     "ThrowStatement",
			Start:    int(stmt.Range().Lo),
			End:      int(stmt.Range().Hi),
			Loc:      locOfNode(stmt, ctx.Parser.Source(), ctx),
			Argument: Convert(stmt.Arg(), ctx),
		}
	case parser.N_STMT_TRY:
		stmt := node.(*parser.TryStmt)
		return &TryStatement{
			Type:      "TryStatement",
			Start:     int(stmt.Range().Lo),
			End:       int(stmt.Range().Hi),
			Loc:       locOfNode(stmt, ctx.Parser.Source(), ctx),
			Block:     Convert(stmt.Try(), ctx),
			Handler:   Convert(stmt.Catch(), ctx),
			Finalizer: Convert(stmt.Fin(), ctx),
		}
	case parser.N_CATCH:
		expr := node.(*parser.Catch)
		return &CatchClause{
			Type:  "CatchClause",
			Start: int(expr.Range().Lo),
			End:   int(expr.Range().Hi),
			Loc:   locOfNode(expr, ctx.Parser.Source(), ctx),
			Param: Convert(expr.Param(), ctx),
			Body:  Convert(expr.Body(), ctx),
		}
	case parser.N_STMT_DEBUG:
		return &DebuggerStatement{
			Type:  "DebuggerStatement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_STMT_WITH:
		stmt := node.(*parser.WithStmt)
		return &WithStatement{
			Type:   "WithStatement",
			Start:  int(stmt.Range().Lo),
			End:    int(stmt.Range().Hi),
			Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
			Object: Convert(stmt.Expr(), ctx),
			Body:   Convert(stmt.Body(), ctx),
		}
	case parser.N_IMPORT_CALL:
		stmt := node.(*parser.ImportCall)
		return &ImportExpression{
			Type:   "ImportExpression",
			Start:  int(stmt.Range().Lo),
			End:    int(stmt.Range().Hi),
			Loc:    locOfNode(stmt, ctx.Parser.Source(), ctx),
			Source: Convert(stmt.Src(), ctx),
		}
	case parser.N_META_PROP:
		stmt := node.(*parser.MetaProp)
		return &MetaProperty{
			Type:     "MetaProperty",
			Start:    int(stmt.Range().Lo),
			End:      int(stmt.Range().Hi),
			Loc:      locOfNode(stmt, ctx.Parser.Source(), ctx),
			Meta:     Convert(stmt.Meta(), ctx),
			Property: Convert(stmt.Prop(), ctx),
		}
	case parser.N_STMT_CLASS:
		stmt := node.(*parser.ClassDec)
		superTypArgs := stmt.SuperTypArgs()
		typParams := stmt.TypParams()
		implements := stmt.Implements()
		rng, loc := locWithDecorator(stmt, parser.DecoratorsOf(stmt), ctx.Parser.Source(), ctx)
		if superTypArgs != nil || typParams != nil || implements != nil {
			return &TSClassDeclaration{
				Type:                "ClassDeclaration",
				Start:               int(rng.Lo),
				End:                 int(rng.Hi),
				Loc:                 loc,
				Id:                  Convert(stmt.Id(), ctx),
				TypeParameters:      ConvertTsTyp(typParams, ctx),
				SuperClass:          Convert(stmt.Super(), ctx),
				SuperTypeParameters: ConvertTsTyp(superTypArgs, ctx),
				Implements:          elems(implements, ctx),
				Body:                Convert(stmt.Body(), ctx),
				Abstract:            stmt.Abstract(),
				Declare:             stmt.Declare(),
				Decorators:          elems(parser.DecoratorsOf(stmt), ctx),
			}
		}
		return &ClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      int(rng.Lo),
			End:        int(rng.Hi),
			Loc:        loc,
			Id:         Convert(stmt.Id(), ctx),
			SuperClass: Convert(stmt.Super(), ctx),
			Body:       Convert(stmt.Body(), ctx),
			Abstract:   stmt.Abstract(),
			Declare:    stmt.Declare(),
			Decorators: elems(parser.DecoratorsOf(stmt), ctx),
		}
	case parser.N_EXPR_CLASS:
		stmt := node.(*parser.ClassDec)
		superTypArgs := stmt.SuperTypArgs()
		typParams := stmt.TypParams()
		implements := stmt.Implements()
		if superTypArgs != nil || typParams != nil || implements != nil {
			return &TSClassExpression{
				Type:                "ClassExpression",
				Start:               int(stmt.Range().Lo),
				End:                 int(stmt.Range().Hi),
				Loc:                 locOfNode(stmt, ctx.Parser.Source(), ctx),
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
			Start:      int(stmt.Range().Lo),
			End:        int(stmt.Range().Hi),
			Loc:        locOfNode(stmt, ctx.Parser.Source(), ctx),
			Id:         Convert(stmt.Id(), ctx),
			SuperClass: Convert(stmt.Super(), ctx),
			Body:       Convert(stmt.Body(), ctx),
			Abstract:   stmt.Abstract(),
		}
	case parser.N_CLASS_BODY:
		stmt := node.(*parser.ClassBody)
		return &ClassBody{
			Type:  "ClassBody",
			Start: int(stmt.Range().Lo),
			End:   int(stmt.Range().Hi),
			Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
			Body:  expressions(stmt.Elems(), ctx),
		}
	case parser.N_METHOD:
		n := node.(*parser.Method)
		f := n.Val().(*parser.FnDec)
		rng, loc := locWithDecorator(n, parser.DecoratorsOf(n), ctx.Parser.Source(), ctx)
		if f.TypInfo() != nil {
			ti := f.TypInfo()
			return &TSMethodDefinition{
				Type:          "MethodDefinition",
				Start:         int(rng.Lo),
				End:           int(rng.Hi),
				Loc:           loc,
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
				Decorators:    elems(parser.DecoratorsOf(n), ctx),
			}
		}
		return &MethodDefinition{
			Type:       "MethodDefinition",
			Start:      int(rng.Lo),
			End:        int(rng.Hi),
			Loc:        loc,
			Key:        Convert(n.Key(), ctx),
			Value:      Convert(n.Val(), ctx),
			Kind:       n.Kind(),
			Computed:   n.Computed(),
			Static:     n.Static(),
			Decorators: elems(parser.DecoratorsOf(n), ctx),
		}
	case parser.N_FIELD:
		n := node.(*parser.Field)
		rng, loc := locWithDecorator(n, parser.DecoratorsOf(n), ctx.Parser.Source(), ctx)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			if n.IsTsSig() {
				return &TSIndexSignature{
					Type:           "TSIndexSignature",
					Start:          int(rng.Lo),
					End:            int(rng.Hi),
					Loc:            loc,
					Static:         n.Static(),
					Abstract:       ti.Abstract(),
					Optional:       ti.Optional(),
					Declare:        ti.Declare(),
					Readonly:       ti.Readonly(),
					Accessibility:  ti.AccMod().String(),
					Parameters:     elems([]parser.Node{n.Key()}, ctx),
					TypeAnnotation: typAnnot(ti, ctx),
					Decorators:     elems(parser.DecoratorsOf(n), ctx),
				}
			}

			return &TSPropertyDefinition{
				Type:           "PropertyDefinition",
				Start:          int(rng.Lo),
				End:            int(rng.Hi),
				Loc:            loc,
				Key:            Convert(n.Key(), ctx),
				Value:          Convert(n.Val(), ctx),
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
				Decorators:     elems(parser.DecoratorsOf(n), ctx),
			}
		}
		return &PropertyDefinition{
			Type:       "PropertyDefinition",
			Start:      int(rng.Lo),
			End:        int(rng.Hi),
			Loc:        loc,
			Key:        Convert(n.Key(), ctx),
			Value:      Convert(n.Val(), ctx),
			Computed:   n.Computed(),
			Static:     n.Static(),
			Decorators: elems(parser.DecoratorsOf(n), ctx),
		}
	case parser.N_SUPER:
		n := node.(*parser.Super)
		return &Super{
			Type:  "Super",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
		}
	case parser.N_EXPR_TPL:
		tpl := node.(*parser.TplExpr)
		if tpl.Tag() == nil {
			return tplLiteral(tpl.Range(), tpl.Elems(), ctx)
		}
		locWithTag := tpl.Tag().Range()
		locWithTag.Hi = tpl.Range().Hi
		tag := tpl.Tag()
		if wt, ok := tag.(parser.NodeWithTypInfo); ok {
			ti := wt.TypInfo()
			if ti != nil {
				return &TSTaggedTemplateExpression{
					Type:           "TaggedTemplateExpression",
					Start:          int(locWithTag.Lo),
					End:            int(locWithTag.Hi),
					Loc:            locOfRng(locWithTag, ctx.Parser.Source(), ctx),
					Tag:            Convert(tpl.Tag(), ctx),
					Quasi:          tplLiteral(tpl.Range(), tpl.Elems(), ctx),
					TypeParameters: typArgs(ti, ctx),
				}
			}
		}
		return &TaggedTemplateExpression{
			Type:  "TaggedTemplateExpression",
			Start: int(locWithTag.Lo),
			End:   int(locWithTag.Hi),
			Loc:   locOfRng(locWithTag, ctx.Parser.Source(), ctx),
			Tag:   Convert(tpl.Tag(), ctx),
			Quasi: tplLiteral(tpl.Range(), tpl.Elems(), ctx),
		}
	case parser.N_STMT_IMPORT:
		stmt := node.(*parser.ImportDec)
		return &ImportDeclaration{
			Type:       "ImportDeclaration",
			Start:      int(stmt.Range().Lo),
			End:        int(stmt.Range().Hi),
			Loc:        locOfNode(stmt, ctx.Parser.Source(), ctx),
			Specifiers: importSpecs(stmt.Specs(), ctx),
			Source:     Convert(stmt.Src(), ctx),
			ImportKind: stmt.Kind(),
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
						Start: int(node.Range().Lo),
						End:   int(node.Range().Hi),
						Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
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
			Start: int(stmt.Range().Lo),
			End:   int(stmt.Range().Hi),
			Loc:   locOfNode(stmt, ctx.Parser.Source(), ctx),
			Body:  statements(stmt.Body(), ctx),
		}
	case parser.N_EXPR_CHAIN:
		node := node.(*parser.ChainExpr)
		return &ChainExpression{
			Type:       "ChainExpression",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_ID:
		node := node.(*parser.JsxIdent)
		return &JSXIdentifier{
			Type:  "JSXIdentifier",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Name:  node.Val(),
		}
	case parser.N_JSX_NS:
		node := node.(*parser.JsxNsName)
		return &JSXNamespacedName{
			Type:      "JSXNamespacedName",
			Start:     int(node.Range().Lo),
			End:       int(node.Range().Hi),
			Loc:       locOfNode(node, ctx.Parser.Source(), ctx),
			Namespace: node.NS(),
			Name:      node.Name(),
		}
	case parser.N_JSX_MEMBER:
		node := node.(*parser.JsxMember)
		return &JSXMemberExpression{
			Type:     "JSXMemberExpression",
			Start:    int(node.Range().Lo),
			End:      int(node.Range().Hi),
			Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
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
				Start: int(node.Range().Lo),
				End:   int(node.Range().Hi),
				Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
				OpeningFragment: &JSXOpeningFragment{
					Type:  "JSXOpeningFragment",
					Start: int(open.Range().Lo),
					End:   int(open.Range().Hi),
					Loc:   locOfNode(open, ctx.Parser.Source(), ctx),
				},
				Children: elems(node.Children(), ctx),
				ClosingFragment: &JSXClosingFragment{
					Type:  "JSXClosingFragment",
					Start: int(close.Range().Lo),
					End:   int(close.Range().Hi),
					Loc:   locOfNode(close, ctx.Parser.Source(), ctx),
				},
			}
		}
		return &JSXElement{
			Type:           "JSXElement",
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
			OpeningElement: Convert(node.Open(), ctx),
			Children:       elems(node.Children(), ctx),
			ClosingElement: Convert(node.Close(), ctx),
		}
	case parser.N_JSX_OPEN:
		node := node.(*parser.JsxOpen)
		if wt, ok := node.Name().(parser.NodeWithTypInfo); ok {
			ti := wt.TypInfo()
			if ti != nil {
				return &TSXOpeningElement{
					Type:           "JSXOpeningElement",
					Start:          int(node.Range().Lo),
					End:            int(node.Range().Hi),
					Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
					Name:           Convert(node.Name(), ctx),
					Attributes:     elems(node.Attrs(), ctx),
					SelfClosing:    node.Closed(),
					TypeParameters: typArgs(ti, ctx),
				}
			}
		}
		return &JSXOpeningElement{
			Type:        "JSXOpeningElement",
			Start:       int(node.Range().Lo),
			End:         int(node.Range().Hi),
			Loc:         locOfNode(node, ctx.Parser.Source(), ctx),
			Name:        Convert(node.Name(), ctx),
			Attributes:  elems(node.Attrs(), ctx),
			SelfClosing: node.Closed(),
		}
	case parser.N_JSX_CLOSE:
		node := node.(*parser.JsxClose)
		return &JSXClosingElement{
			Type:  "JSXClosingElement",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Name:  Convert(node.Name(), ctx),
		}
	case parser.N_JSX_TXT:
		node := node.(*parser.JsxText)
		return &JSXText{
			Type:  "JSXText",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Value: node.Val(),
			Raw:   ctx.Parser.RngText(node.Range()),
		}
	case parser.N_JSX_EXPR_SPAN:
		node := node.(*parser.JsxExprSpan)
		return &JSXExpressionContainer{
			Type:       "JSXExpressionContainer",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_ATTR:
		node := node.(*parser.JsxAttr)
		return &JSXAttribute{
			Type:  "JSXAttribute",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Name:  Convert(node.Name(), ctx),
			Value: Convert(node.Val(), ctx),
		}
	case parser.N_JSX_ATTR_SPREAD:
		node := node.(*parser.JsxSpreadAttr)
		return &JSXSpreadAttribute{
			Type:     "JSXSpreadAttribute",
			Start:    int(node.Range().Lo),
			End:      int(node.Range().Hi),
			Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
			Argument: Convert(node.Arg(), ctx),
		}
	case parser.N_JSX_CHILD_SPREAD:
		node := node.(*parser.JsxSpreadChild)
		return &JSXSpreadChild{
			Type:       "JSXSpreadChild",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: Convert(node.Expr(), ctx),
		}
	case parser.N_JSX_EMPTY:
		node := node.(*parser.JsxEmpty)
		return &JSXEmptyExpression{
			Type:  "JSXEmptyExpression",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_DECORATOR:
		node := node.(*parser.Decorator)
		return &Decorator{
			Type:       "Decorator",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: Convert(node.Expr(), ctx),
		}
	}

	// bypass the ts related node like `Ambient`
	return ConvertTsTyp(node, ctx)
}
