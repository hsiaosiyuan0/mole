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

func ConvertProg(n *parser.Prog) *Program {
	stmts := n.Body()
	body := make([]Node, len(stmts))
	for i, s := range stmts {
		body[i] = convert(s)
	}
	return &Program{
		Type:  "Program",
		Start: start(n.Loc()),
		End:   end(n.Loc()),
		Loc:   loc(n.Loc()),
		Body:  body,
	}
}

func arrExpr(n *parser.ArrLit) *ArrayExpression {
	exprs := n.Elems()
	elems := make([]Expression, len(exprs))
	for i, e := range exprs {
		elems[i] = convert(e)
	}
	return &ArrayExpression{
		Type:     "ArrayExpression",
		Start:    start(n.Loc()),
		End:      end(n.Loc()),
		Loc:      loc(n.Loc()),
		Elements: elems,
	}
}

func obj(n *parser.ObjLit) *ObjectExpression {
	ps := n.Props()
	props := make([]Node, len(ps))
	for i, p := range ps {
		props[i] = convert(p)
	}
	return &ObjectExpression{
		Type:       "ObjectExpression",
		Start:      start(n.Loc()),
		End:        end(n.Loc()),
		Loc:        loc(n.Loc()),
		Properties: props,
	}
}

func fnParams(params []parser.Node) []Node {
	ps := make([]Node, len(params))
	for i, p := range params {
		fp := tsParamProp(p)
		if fp == nil {
			fp = convert(p)
		}
		ps[i] = fp
	}
	return ps
}

func statements(stmts []parser.Node) []Statement {
	s := make([]Statement, len(stmts))
	for i, stmt := range stmts {
		s[i] = convert(stmt)
	}
	return s
}

func expressions(exprs []parser.Node) []Expression {
	s := make([]Expression, len(exprs))
	for i, expr := range exprs {
		s[i] = convert(expr)
	}
	return s
}

func blockStmt(n *parser.BlockStmt) *BlockStatement {
	return &BlockStatement{
		Type:  "BlockStatement",
		Start: start(n.Loc()),
		End:   end(n.Loc()),
		Loc:   loc(n.Loc()),
		Body:  statements(n.Body()),
	}
}

func cases(cs []parser.Node) []*SwitchCase {
	s := make([]*SwitchCase, len(cs))
	for i, c := range cs {
		sc := c.(*parser.SwitchCase)
		s[i] = &SwitchCase{
			Type:       "SwitchCase",
			Start:      start(sc.Loc()),
			End:        end(sc.Loc()),
			Loc:        loc(sc.Loc()),
			Test:       convert(sc.Test()),
			Consequent: statements(sc.Cons()),
		}
	}
	return s
}

func declarations(decList []parser.Node) []*VariableDeclarator {
	s := make([]*VariableDeclarator, len(decList))
	for i, d := range decList {
		dc := d.(*parser.VarDec)
		s[i] = &VariableDeclarator{
			Type:  "VariableDeclarator",
			Start: start(dc.Loc()),
			End:   end(dc.Loc()),
			Loc:   loc(dc.Loc()),
			Id:    convert(dc.Id()),
			Init:  convert(dc.Init()),
		}
	}
	return s
}

func tplLiteral(tplLoc *parser.Loc, elems []parser.Node) *TemplateLiteral {
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
			exprs = append(exprs, convert(elem))
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

func importSpec(spec *parser.ImportSpec) Node {
	if spec.Default() {
		return &ImportDefaultSpecifier{
			Type:  "ImportDefaultSpecifier",
			Start: start(spec.Loc()),
			End:   end(spec.Loc()),
			Loc:   loc(spec.Loc()),
			Local: convert(spec.Local()),
		}
	} else if spec.NameSpace() {
		return &ImportNamespaceSpecifier{
			Type:  "ImportNamespaceSpecifier",
			Start: start(spec.Loc()),
			End:   end(spec.Loc()),
			Loc:   loc(spec.Loc()),
			Local: convert(spec.Local()),
		}
	}
	return &ImportSpecifier{
		Type:     "ImportSpecifier",
		Start:    start(spec.Loc()),
		End:      end(spec.Loc()),
		Loc:      loc(spec.Loc()),
		Local:    convert(spec.Local()),
		Imported: convert(spec.Id()),
	}
}

func importSpecs(specs []parser.Node) []Node {
	ret := make([]Node, len(specs))
	for i, spec := range specs {
		ret[i] = importSpec(spec.(*parser.ImportSpec))
	}
	return ret
}

func exportAll(node *parser.ExportDec) Node {
	var spec parser.Node
	if len(node.Specs()) == 1 {
		spec = node.Specs()[0].(*parser.ExportSpec).Local()
	}
	return &ExportAllDeclaration{
		Type:     "ExportAllDeclaration",
		Start:    start(node.Loc()),
		End:      end(node.Loc()),
		Loc:      loc(node.Loc()),
		Source:   convert(node.Src()),
		Exported: convert(spec),
	}
}

func exportDefault(node *parser.ExportDec) Node {
	return &ExportDefaultDeclaration{
		Type:        "ExportDefaultDeclaration",
		Start:       start(node.Loc()),
		End:         end(node.Loc()),
		Loc:         loc(node.Loc()),
		Declaration: convert(node.Dec()),
	}
}

func exportSpecs(specs []parser.Node) []Node {
	ret := make([]Node, len(specs))
	for i, spec := range specs {
		s := spec.(*parser.ExportSpec)
		ret[i] = &ExportSpecifier{
			Type:     "ExportSpecifier",
			Start:    start(s.Loc()),
			End:      end(s.Loc()),
			Loc:      loc(s.Loc()),
			Local:    convert(s.Local()),
			Exported: convert(s.Id()),
		}
	}
	return ret
}

func exportNamed(node *parser.ExportDec) Node {
	return &ExportNamedDeclaration{
		Type:        "ExportNamedDeclaration",
		Start:       start(node.Loc()),
		End:         end(node.Loc()),
		Loc:         loc(node.Loc()),
		Declaration: convert(node.Dec()),
		Specifiers:  exportSpecs(node.Specs()),
		Source:      convert(node.Src()),
	}
}

func elems(nodes []parser.Node) []Node {
	ret := make([]Node, len(nodes))
	for i, node := range nodes {
		ret[i] = convert(node)
	}
	return ret
}

func ident(node parser.Node) Node {
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
			TypeAnnotation: typAnnot(ti),
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

func convert(node parser.Node) Node {
	if node == nil {
		return nil
	}
	switch node.Type() {
	case parser.N_EXPR_PAREN:
		return convert(node.(*parser.ParenExpr).Expr())
	case parser.N_STMT_EXPR:
		exprStmt := node.(*parser.ExprStmt)
		expr := convert(exprStmt.Expr())
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
			Callee:    convert(new.Callee()),
			Arguments: expressions(new.Args()),
		}
	case parser.N_NAME:
		return ident(node)
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
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
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
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
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
		return arrExpr(node.(*parser.ArrLit))
	case parser.N_LIT_OBJ:
		return obj(node.(*parser.ObjLit))
	case parser.N_PROP:
		prop := node.(*parser.Prop)
		return &Property{
			Type:      "Property",
			Start:     start(prop.Loc()),
			End:       end(prop.Loc()),
			Loc:       loc(prop.Loc()),
			Key:       convert(prop.Key()),
			Value:     convert(prop.Value()),
			Kind:      prop.Kind(),
			Computed:  prop.Computed(),
			Shorthand: prop.Shorthand(),
			Method:    prop.Method(),
		}
	case parser.N_STMT_BLOCK:
		return blockStmt(node.(*parser.BlockStmt))
	case parser.N_EXPR_FN:
		n := node.(*parser.FnDec)
		if n.TypInfo() != nil {
			ti := n.TypInfo()
			return &TSFunctionExpression{
				Type:           "FunctionExpression",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Id:             convert(n.Id()),
				Params:         fnParams(n.Params()),
				Body:           convert(n.Body()),
				Generator:      n.Generator(),
				Async:          n.Async(),
				Expression:     false,
				TypeParameters: typParams(ti),
				ReturnType:     typAnnot(ti),
			}
		}
		return &FunctionExpression{
			Type:       "FunctionExpression",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Id:         convert(n.Id()),
			Params:     fnParams(n.Params()),
			Body:       convert(n.Body()),
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
				Params:         fnParams(n.Params()),
				Body:           convert(n.Body()),
				Generator:      false,
				Async:          n.Async(),
				Expression:     n.Expr(),
				TypeParameters: typParams(ti),
				ReturnType:     typAnnot(ti),
			}
		}
		return &ArrowFunctionExpression{
			Type:       "ArrowFunctionExpression",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Id:         nil,
			Params:     fnParams(n.Params()),
			Body:       convert(n.Body()),
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
					Id:             convert(n.Id()),
					Params:         fnParams(n.Params()),
					Body:           convert(n.Body()),
					Generator:      false,
					Async:          n.Async(),
					TypeParameters: typParams(ti),
					ReturnType:     typAnnot(ti),
				}
			}
			return &TSFunctionDeclaration{
				Type:           "FunctionDeclaration",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Id:             convert(n.Id()),
				Params:         fnParams(n.Params()),
				Body:           convert(n.Body()),
				Generator:      n.Generator(),
				Async:          n.Async(),
				TypeParameters: typParams(ti),
				ReturnType:     typAnnot(ti),
			}
		}
		return &FunctionDeclaration{
			Type:      "FunctionDeclaration",
			Start:     start(n.Loc()),
			End:       end(n.Loc()),
			Loc:       loc(n.Loc()),
			Id:        convert(n.Id()),
			Params:    fnParams(n.Params()),
			Body:      convert(n.Body()),
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
			Argument: convert(node.Arg()),
		}
	case parser.N_STMT_RET:
		ret := node.(*parser.RetStmt)
		return &ReturnStatement{
			Type:     "ReturnStatement",
			Start:    start(ret.Loc()),
			End:      end(ret.Loc()),
			Loc:      loc(ret.Loc()),
			Argument: convert(ret.Arg()),
		}
	case parser.N_SPREAD:
		n := node.(*parser.Spread)
		return &SpreadElement{
			Type:     "SpreadElement",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Argument: convert(n.Arg()),
		}
	case parser.N_PAT_ARRAY:
		n := node.(*parser.ArrPat)
		ti := n.TypInfo()
		if ti != nil {
			lc := parser.LocWithTypeInfo(n, false)
			return &TSArrayPattern{
				Type:     "ArrayPattern",
				Start:    start(lc),
				End:      end(lc),
				Loc:      loc(lc),
				Elements: elems(n.Elems()),
				Optional: ti.Optional(),
			}
		}
		return &ArrayPattern{
			Type:     "ArrayPattern",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Elements: elems(n.Elems()),
		}
	case parser.N_PAT_ASSIGN:
		n := node.(*parser.AssignPat)
		return &AssignmentPattern{
			Type:  "AssignmentPattern",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Left:  convert(n.Left()),
			Right: convert(n.Right()),
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
				Properties:     elems(n.Props()),
				TypeAnnotation: typAnnot(ti),
			}
		}
		return &ObjectPattern{
			Type:       "ObjectPattern",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Properties: elems(n.Props()),
		}
	case parser.N_PAT_REST:
		n := node.(*parser.RestPat)
		if n.TypInfo() != nil {
			return &TSRestElement{
				Type:           "RestElement",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Argument:       convert(n.Arg()),
				Optional:       n.Optional(),
				TypeAnnotation: typAnnot(n.TypInfo()),
			}
		}
		return &RestElement{
			Type:     "RestElement",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Argument: convert(n.Arg()),
		}
	case parser.N_STMT_IF:
		ifStmt := node.(*parser.IfStmt)
		return &IfStatement{
			Type:       "IfStatement",
			Start:      start(ifStmt.Loc()),
			End:        end(ifStmt.Loc()),
			Loc:        loc(ifStmt.Loc()),
			Test:       convert(ifStmt.Test()),
			Consequent: convert(ifStmt.Cons()),
			Alternate:  convert(ifStmt.Alt()),
		}
	case parser.N_EXPR_CALL:
		n := node.(*parser.CallExpr)
		if n.TypInfo() != nil {
			return &TSCallExpression{
				Type:           "CallExpression",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Callee:         convert(n.Callee()),
				Arguments:      expressions(n.Args()),
				Optional:       n.Optional(),
				TypeParameters: typArgs(n.TypInfo()),
			}
		}
		return &CallExpression{
			Type:      "CallExpression",
			Start:     start(n.Loc()),
			End:       end(n.Loc()),
			Loc:       loc(n.Loc()),
			Callee:    convert(n.Callee()),
			Arguments: expressions(n.Args()),
			Optional:  n.Optional(),
		}
	case parser.N_STMT_SWITCH:
		swc := node.(*parser.SwitchStmt)
		return &SwitchStatement{
			Type:         "SwitchStatement",
			Start:        start(swc.Loc()),
			End:          end(swc.Loc()),
			Loc:          loc(swc.Loc()),
			Discriminant: convert(swc.Test()),
			Cases:        cases(swc.Cases()),
		}
	case parser.N_STMT_VAR_DEC:
		varDec := node.(*parser.VarDecStmt)
		return &VariableDeclaration{
			Type:         "VariableDeclaration",
			Start:        start(varDec.Loc()),
			End:          end(varDec.Loc()),
			Loc:          loc(varDec.Loc()),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList()),
		}
	case parser.N_EXPR_MEMBER:
		mem := node.(*parser.MemberExpr)
		return &MemberExpression{
			Type:     "MemberExpression",
			Loc:      loc(mem.Loc()),
			Start:    start(mem.Loc()),
			End:      end(mem.Loc()),
			Object:   convert(mem.Obj()),
			Property: convert(mem.Prop()),
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
			Expressions: expressions(seq.Elems()),
		}
	case parser.N_EXPR_UPDATE:
		up := node.(*parser.UpdateExpr)
		return &UpdateExpression{
			Type:     "UpdateExpression",
			Start:    start(up.Loc()),
			End:      end(up.Loc()),
			Loc:      loc(up.Loc()),
			Operator: up.OpText(),
			Argument: convert(up.Arg()),
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
			Argument: convert(un.Arg()),
		}
	case parser.N_EXPR_COND:
		cond := node.(*parser.CondExpr)
		return &ConditionalExpression{
			Type:       "ConditionalExpression",
			Start:      start(cond.Loc()),
			End:        end(cond.Loc()),
			Loc:        loc(cond.Loc()),
			Test:       convert(cond.Test()),
			Consequent: convert(cond.Cons()),
			Alternate:  convert(cond.Alt()),
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
			Test:  convert(stmt.Test()),
			Body:  convert(stmt.Body()),
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
			Test:  convert(stmt.Test()),
			Body:  convert(stmt.Body()),
		}
	case parser.N_STMT_FOR:
		stmt := node.(*parser.ForStmt)
		return &ForStatement{
			Type:   "ForStatement",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Init:   convert(stmt.Init()),
			Test:   convert(stmt.Test()),
			Update: convert(stmt.Update()),
			Body:   convert(stmt.Body()),
		}
	case parser.N_STMT_FOR_IN_OF:
		stmt := node.(*parser.ForInOfStmt)
		if stmt.In() {
			return &ForInStatement{
				Type:  "ForInStatement",
				Start: start(stmt.Loc()),
				End:   end(stmt.Loc()),
				Loc:   loc(stmt.Loc()),
				Left:  convert(stmt.Left()),
				Right: convert(stmt.Right()),
				Body:  convert(stmt.Body()),
			}
		}
		return &ForOfStatement{
			Type:  "ForOfStatement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Left:  convert(stmt.Left()),
			Right: convert(stmt.Right()),
			Body:  convert(stmt.Body()),
			Await: stmt.Await(),
		}
	case parser.N_STMT_CONT:
		stmt := node.(*parser.ContStmt)
		return &ContinueStatement{
			Type:  "ContinueStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
		}
	case parser.N_STMT_BRK:
		stmt := node.(*parser.BrkStmt)
		return &BreakStatement{
			Type:  "BreakStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
		}
	case parser.N_STMT_LABEL:
		stmt := node.(*parser.LabelStmt)
		return &LabeledStatement{
			Type:  "LabeledStatement",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
			Body:  convert(stmt.Body()),
		}
	case parser.N_STMT_THROW:
		stmt := node.(*parser.ThrowStmt)
		return &ThrowStatement{
			Type:     "ThrowStatement",
			Start:    start(stmt.Loc()),
			End:      end(stmt.Loc()),
			Loc:      loc(stmt.Loc()),
			Argument: convert(stmt.Arg()),
		}
	case parser.N_STMT_TRY:
		stmt := node.(*parser.TryStmt)
		return &TryStatement{
			Type:      "TryStatement",
			Start:     start(stmt.Loc()),
			End:       end(stmt.Loc()),
			Loc:       loc(stmt.Loc()),
			Block:     convert(stmt.Try()),
			Handler:   convert(stmt.Catch()),
			Finalizer: convert(stmt.Fin()),
		}
	case parser.N_CATCH:
		expr := node.(*parser.Catch)
		return &CatchClause{
			Type:  "CatchClause",
			Start: start(expr.Loc()),
			End:   end(expr.Loc()),
			Loc:   loc(expr.Loc()),
			Param: convert(expr.Param()),
			Body:  convert(expr.Body()),
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
			Object: convert(stmt.Expr()),
			Body:   convert(stmt.Body()),
		}
	case parser.N_IMPORT_CALL:
		stmt := node.(*parser.ImportCall)
		return &ImportExpression{
			Type:   "ImportExpression",
			Start:  start(stmt.Loc()),
			End:    end(stmt.Loc()),
			Loc:    loc(stmt.Loc()),
			Source: convert(stmt.Src()),
		}
	case parser.N_META_PROP:
		stmt := node.(*parser.MetaProp)
		return &MetaProperty{
			Type:     "MetaProperty",
			Start:    start(stmt.Loc()),
			End:      end(stmt.Loc()),
			Loc:      loc(stmt.Loc()),
			Meta:     convert(stmt.Meta()),
			Property: convert(stmt.Prop()),
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
				Id:                  convert(stmt.Id()),
				TypeParameters:      convertTsTyp(typParams),
				SuperClass:          convert(stmt.Super()),
				SuperTypeParameters: convertTsTyp(superTypArgs),
				Implements:          elems(implements),
				Body:                convert(stmt.Body()),
				Abstract:            stmt.Abstract(),
			}
		}
		return &ClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Id:         convert(stmt.Id()),
			SuperClass: convert(stmt.Super()),
			Body:       convert(stmt.Body()),
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
				Id:                  convert(stmt.Id()),
				TypeParameters:      convertTsTyp(typParams),
				SuperClass:          convert(stmt.Super()),
				SuperTypeParameters: convertTsTyp(superTypArgs),
				Implements:          elems(implements),
				Body:                convert(stmt.Body()),
				Abstract:            stmt.Abstract(),
			}
		}
		return &ClassExpression{
			Type:       "ClassExpression",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Id:         convert(stmt.Id()),
			SuperClass: convert(stmt.Super()),
			Body:       convert(stmt.Body()),
			Abstract:   stmt.Abstract(),
		}
	case parser.N_ClASS_BODY:
		stmt := node.(*parser.ClassBody)
		return &ClassBody{
			Type:  "ClassBody",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Body:  expressions(stmt.Elems()),
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
				Key:           convert(n.Key()),
				Value:         convert(n.Val()),
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
			Key:      convert(n.Key()),
			Value:    convert(n.Val()),
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
					Parameters:     elems([]parser.Node{n.Key()}),
					TypeAnnotation: typAnnot(ti),
				}
			}

			return &TSPropertyDefinition{
				Type:           "PropertyDefinition",
				Start:          start(n.Loc()),
				End:            end(n.Loc()),
				Loc:            loc(n.Loc()),
				Key:            convert(n.Key()),
				Value:          convert(n.Value()),
				Computed:       n.Computed(),
				Static:         n.Static(),
				Abstract:       ti.Abstract(),
				Optional:       ti.Optional(),
				Definite:       ti.Definite(),
				Readonly:       ti.Readonly(),
				Override:       ti.Override(),
				Declare:        ti.Declare(),
				Accessibility:  ti.AccMod().String(),
				TypeAnnotation: typAnnot(ti),
			}
		}
		return &PropertyDefinition{
			Type:     "PropertyDefinition",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			Key:      convert(n.Key()),
			Value:    convert(n.Value()),
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
			return tplLiteral(tpl.Loc(), tpl.Elems())
		}
		return &TaggedTemplateExpression{
			Type:  "TaggedTemplateExpression",
			Start: start(tpl.LocWithTag()),
			End:   end(tpl.LocWithTag()),
			Loc:   loc(tpl.LocWithTag()),
			Tag:   convert(tpl.Tag()),
			Quasi: tplLiteral(tpl.Loc(), tpl.Elems()),
		}
	case parser.N_STMT_IMPORT:
		stmt := node.(*parser.ImportDec)
		return &ImportDeclaration{
			Type:       "ImportDeclaration",
			Start:      start(stmt.Loc()),
			End:        end(stmt.Loc()),
			Loc:        loc(stmt.Loc()),
			Specifiers: importSpecs(stmt.Specs()),
			Source:     convert(stmt.Src()),
		}
	case parser.N_STMT_EXPORT:
		stmt := node.(*parser.ExportDec)
		if stmt.All() {
			return exportAll(stmt)
		}
		if stmt.Default() {
			return exportDefault(stmt)
		}
		if stmt.Dec() != nil && stmt.Dec().Type() == parser.N_TS_NAMESPACE {
			n := stmt.Dec().(*parser.TsNS)
			if n.Alias() {
				return &TSNamespaceExportDeclaration{
					Type:  "TSNamespaceExportDeclaration",
					Start: start(node.Loc()),
					End:   end(node.Loc()),
					Loc:   loc(node.Loc()),
					Id:    convert(n.Id()),
				}
			}
		}
		return exportNamed(stmt)
	case parser.N_STATIC_BLOCK:
		stmt := node.(*parser.StaticBlock)
		return &StaticBlock{
			Type:  "StaticBlock",
			Start: start(stmt.Loc()),
			End:   end(stmt.Loc()),
			Loc:   loc(stmt.Loc()),
			Body:  statements(stmt.Body()),
		}
	case parser.N_EXPR_CHAIN:
		node := node.(*parser.ChainExpr)
		return &ChainExpression{
			Type:       "ChainExpression",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: convert(node.Expr()),
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
			Object:   convert(node.Obj()),
			Property: convert(node.Prop()),
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
				Children: elems(node.Children()),
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
			OpeningElement: convert(node.Open()),
			Children:       elems(node.Children()),
			ClosingElement: convert(node.Close()),
		}
	case parser.N_JSX_OPEN:
		node := node.(*parser.JsxOpen)
		return &JSXOpeningElement{
			Type:        "JSXOpeningElement",
			Start:       start(node.Loc()),
			End:         end(node.Loc()),
			Loc:         loc(node.Loc()),
			Name:        convert(node.Name()),
			Attributes:  elems(node.Attrs()),
			SelfClosing: node.Closed(),
		}
	case parser.N_JSX_CLOSE:
		node := node.(*parser.JsxClose)
		return &JSXClosingElement{
			Type:  "JSXClosingElement",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Name:  convert(node.Name()),
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
			Expression: convert(node.Expr()),
		}
	case parser.N_JSX_ATTR:
		node := node.(*parser.JsxAttr)
		return &JSXAttribute{
			Type:  "JSXAttribute",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Name:  convert(node.Name()),
			Value: convert(node.Value()),
		}
	case parser.N_JSX_ATTR_SPREAD:
		node := node.(*parser.JsxSpreadAttr)
		return &JSXSpreadAttribute{
			Type:     "JSXSpreadAttribute",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			Argument: convert(node.Arg()),
		}
	case parser.N_JSX_CHILD_SPREAD:
		node := node.(*parser.JsxSpreadChild)
		return &JSXSpreadChild{
			Type:       "JSXSpreadChild",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: convert(node.Expr()),
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
	return convertTsTyp(node)
}
