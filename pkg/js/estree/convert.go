package estree

import (
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func pos(p *parser.Pos) *Position {
	return &Position{Line: p.Line(), Column: p.Column()}
}

func loc(s *parser.Loc) *SrcLoc {
	return &SrcLoc{
		Source: s.Source(),
		Start:  pos(s.Begin()),
		End:    pos(s.End()),
		Range:  rng(s),
	}
}

func rng(s *parser.Loc) *SrcRange {
	rng := s.Range()
	return &SrcRange{
		Start: rng.Start(),
		End:   rng.End(),
	}
}

func program(n *parser.Prog) *Program {
	stmts := n.Body()
	body := make([]Node, len(stmts))
	for i, s := range stmts {
		body[i] = convert(s)
	}
	return &Program{
		Type: "Program",
		Loc:  loc(n.Loc()),
		Body: body,
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
		Loc:      loc(n.Loc()),
		Elements: elems,
	}
}

func obj(n *parser.ObjLit) *ObjectExpression {
	ps := n.Props()
	props := make([]*Property, len(ps))
	for i, p := range ps {
		props[i] = convert(p).(*Property)
	}
	return &ObjectExpression{
		Type:       "ObjectExpression",
		Loc:        loc(n.Loc()),
		Properties: props,
	}
}

func fnParams(params []parser.Node) []Node {
	ps := make([]Node, len(params))
	for i, p := range params {
		ps[i] = convert(p)
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

func blockStmt(block *parser.BlockStmt) *BlockStatement {
	return &BlockStatement{
		Type: "BlockStatement",
		Loc:  loc(block.Loc()),
		Body: statements(block.Body()),
	}
}

func convert(node parser.Node) Node {
	if node == nil {
		return nil
	}
	switch node.Type() {
	case parser.N_STMT_EXPR:
		expr := convert(node.(*parser.ExprStmt).Expr())
		return &ExpressionStatement{
			Type:       "ExpressionStatement",
			Loc:        loc(node.Loc()),
			Expression: expr,
		}
	case parser.N_EXPR_NEW:
		expr := convert(node.(*parser.NewExpr).Expr())
		return &NewExpression{
			Type:   "NewExpression",
			Loc:    loc(node.Loc()),
			Callee: expr,
		}
	case parser.N_NAME:
		id := node.(*parser.Ident)
		name := id.Text()
		return &Identifier{
			Type: "Identifier",
			Loc:  loc(id.Loc()),
			Name: name,
		}
	case parser.N_EXPR_THIS:
		return &ThisExpression{
			Type: "ThisExpression",
			Loc:  loc(node.Loc()),
		}
	case parser.N_LIT_NULL:
		return &Literal{
			Type: "Literal",
			Loc:  loc(node.Loc()),
		}
	case parser.N_LIT_NUM:
		return &Literal{
			Type:  "Literal",
			Loc:   loc(node.Loc()),
			Value: node.(*parser.NumLit).ToFloat(),
		}
	case parser.N_LIT_STR:
		return &Literal{
			Type:  "Literal",
			Loc:   loc(node.Loc()),
			Value: node.(*parser.StrLit).Text(),
		}
	case parser.N_LIT_REGEXP:
		regexp := node.(*parser.RegexpLit)
		return &RegExpLiteral{
			Type:   "Literal",
			Loc:    loc(node.Loc()),
			Regexp: &Regexp{regexp.Pattern(), regexp.Flags()},
		}
	case parser.N_EXPR_BIN:
		bin := node.(*parser.BinExpr)
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
		op := bin.Op().Text()
		return &BinaryExpression{
			Type:     "BinaryExpression",
			Loc:      loc(node.Loc()),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	case parser.N_EXPR_ASSIGN:
		bin := node.(*parser.AssignExpr)
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
		op := bin.Op().Text()
		return &AssignmentExpression{
			Type:     "AssignmentExpression",
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
			Type:     "Property",
			Loc:      loc(prop.Loc()),
			Key:      convert(prop.Key()),
			Value:    convert(prop.Value()),
			Kind:     "init",
			Computed: prop.Computed(),
		}
	case parser.N_METHOD:
		method := node.(*parser.Method)
		return &Property{
			Type:     "Property",
			Loc:      loc(method.Loc()),
			Key:      convert(method.Key()),
			Value:    convert(method.Value()),
			Kind:     method.Kind(),
			Computed: method.Computed(),
			Method:   true,
		}
	case parser.N_STMT_BLOCK:
		return blockStmt(node.(*parser.BlockStmt))
	case parser.N_EXPR_FN:
		fn := node.(*parser.FnDec)
		return &FunctionExpression{
			Type:      "FunctionExpression",
			Loc:       loc(fn.Loc()),
			Id:        convert(fn.Id()),
			Params:    fnParams(fn.Params()),
			Body:      convert(fn.Body()),
			Generator: fn.Generator(),
			Async:     fn.Async(),
		}
	case parser.N_STMT_RET:
		ret := node.(*parser.RetStmt)
		return &ReturnStatement{
			Type:     "ReturnStatement",
			Loc:      loc(ret.Loc()),
			Argument: convert(ret.Arg()),
		}
	}
	return nil
}
