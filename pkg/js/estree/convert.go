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

func expressions(exprs []parser.Node) []Expression {
	s := make([]Expression, len(exprs))
	for i, expr := range exprs {
		s[i] = convert(expr)
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

func cases(cs []*parser.SwitchCase) []*SwitchCase {
	s := make([]*SwitchCase, len(cs))
	for i, c := range cs {
		s[i] = &SwitchCase{
			Type:       "SwitchCase",
			Loc:        loc(c.Loc()),
			Test:       convert(c.Test()),
			Consequent: statements(c.Cons()),
		}
	}
	return s
}

func declarations(decList []*parser.VarDec) []*VariableDeclarator {
	s := make([]*VariableDeclarator, len(decList))
	for i, d := range decList {
		s[i] = &VariableDeclarator{
			Type: "VariableDeclarator",
			Loc:  loc(d.Loc()),
			Id:   convert(d.Id()),
			Init: convert(d.Init()),
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
					Type: "TemplateElement",
					Loc:  lc,
					Tail: last,
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
				Type: "TemplateElement",
				Loc:  loc(elem.Loc()),
				Tail: last,
				Value: &TemplateElementValue{
					Cooked: str.Text(),
					Raw:    str.Raw(),
				},
			})
		}
	}
	return &TemplateLiteral{
		Type:        "TemplateLiteral",
		Loc:         loc(tplLoc),
		Quasis:      quasis,
		Expressions: exprs,
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
		new := node.(*parser.NewExpr)
		return &NewExpression{
			Type:      "NewExpression",
			Loc:       loc(node.Loc()),
			Callee:    convert(new.Callee()),
			Arguments: expressions(new.Args()),
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
		opv := bin.Op().Value()

		if opv == parser.T_AND || opv == parser.T_OR {
			return &LogicalExpression{
				Type:     "LogicalExpression",
				Loc:      loc(node.Loc()),
				Operator: op,
				Left:     lhs,
				Right:    rhs,
			}
		}
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
			Kind:     prop.Kind(),
			Computed: prop.Computed(),
			Method:   prop.Method(),
		}
	case parser.N_STMT_BLOCK:
		return blockStmt(node.(*parser.BlockStmt))
	case parser.N_EXPR_FN:
		fn := node.(*parser.FnDec)
		return &FunctionExpression{
			Type:       "FunctionExpression",
			Loc:        loc(fn.Loc()),
			Id:         convert(fn.Id()),
			Params:     fnParams(fn.Params()),
			Body:       convert(fn.Body()),
			Generator:  fn.Generator(),
			Async:      fn.Async(),
			Expression: false,
		}
	case parser.N_STMT_FN:
		fn := node.(*parser.FnDec)
		return &FunctionDeclaration{
			Type:      "FunctionDeclaration",
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
	case parser.N_SPREAD:
		n := node.(*parser.Spread)
		return &SpreadElement{
			Type:     "SpreadElement",
			Loc:      loc(n.Loc()),
			Argument: convert(n.Arg()),
		}
	case parser.N_PATTERN_REST:
		n := node.(*parser.RestPattern)
		return &RestElement{
			Type:     "RestElement",
			Loc:      loc(n.Loc()),
			Argument: convert(n.Arg()),
		}
	case parser.N_STMT_IF:
		ifStmt := node.(*parser.IfStmt)
		return &IfStatement{
			Type:       "IfStatement",
			Loc:        loc(ifStmt.Loc()),
			Test:       convert(ifStmt.Test()),
			Consequent: convert(ifStmt.Cons()),
			Alternate:  convert(ifStmt.Alt()),
		}
	case parser.N_EXPR_CALL:
		call := node.(*parser.CallExpr)
		return &CallExpression{
			Type:      "CallExpression",
			Loc:       loc(call.Loc()),
			Callee:    convert(call.Callee()),
			Arguments: expressions(call.Args()),
		}
	case parser.N_STMT_SWITCH:
		swc := node.(*parser.SwitchStmt)
		return &SwitchStatement{
			Type:         "SwitchStatement",
			Loc:          loc(swc.Loc()),
			Discriminant: convert(swc.Test()),
			Cases:        cases(swc.Cases()),
		}
	case parser.N_STMT_VAR_DEC:
		varDec := node.(*parser.VarDecStmt)
		return &VariableDeclaration{
			Type:         "VariableDeclaration",
			Loc:          loc(varDec.Loc()),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList()),
		}
	case parser.N_EXPR_MEMBER:
		mem := node.(*parser.MemberExpr)
		return &MemberExpression{
			Type:     "MemberExpression",
			Loc:      loc(mem.Loc()),
			Object:   convert(mem.Obj()),
			Property: convert(mem.Prop()),
			Computed: mem.Compute(),
			Optional: mem.Optional(),
		}
	case parser.N_EXPR_SEQ:
		seq := node.(*parser.SeqExpr)
		return &SequenceExpression{
			Type:        "SequenceExpression",
			Loc:         loc(seq.Loc()),
			Expressions: expressions(seq.Elems()),
		}
	case parser.N_EXPR_UPDATE:
		up := node.(*parser.UpdateExpr)
		return &UpdateExpression{
			Type:     "UpdateExpression",
			Loc:      loc(up.Loc()),
			Operator: up.Op().Text(),
			Argument: convert(up.Arg()),
			Prefix:   up.Prefix(),
		}
	case parser.N_EXPR_UNARY:
		un := node.(*parser.UnaryExpr)
		return &UnaryExpression{
			Type:     "UnaryExpression",
			Loc:      loc(un.Loc()),
			Operator: un.Op().Text(),
			Prefix:   true,
			Argument: convert(un.Arg()),
		}
	case parser.N_EXPR_COND:
		cond := node.(*parser.CondExpr)
		return &ConditionalExpression{
			Type:       "ConditionalExpression",
			Loc:        loc(cond.Loc()),
			Test:       convert(cond.Test()),
			Consequent: convert(cond.Cons()),
			Alternate:  convert(cond.Alt()),
		}
	case parser.N_STMT_EMPTY:
		return &EmptyStatement{
			Type: "EmptyStatement",
			Loc:  loc(node.Loc()),
		}
	case parser.N_STMT_DO_WHILE:
		stmt := node.(*parser.DoWhileStmt)
		return &DoWhileStatement{
			Type: "DoWhileStatement",
			Loc:  loc(node.Loc()),
			Test: convert(stmt.Test()),
			Body: convert(stmt.Body()),
		}
	case parser.N_LIT_BOOL:
		b := node.(*parser.BoolLit)
		return &Literal{
			Type:  "Literal",
			Loc:   loc(node.Loc()),
			Value: b.Value(),
		}
	case parser.N_STMT_WHILE:
		stmt := node.(*parser.WhileStmt)
		return &WhileStatement{
			Type: "WhileStatement",
			Loc:  loc(node.Loc()),
			Test: convert(stmt.Test()),
			Body: convert(stmt.Body()),
		}
	case parser.N_STMT_FOR:
		stmt := node.(*parser.ForStmt)
		return &ForStatement{
			Type:   "ForStatement",
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
				Loc:   loc(node.Loc()),
				Left:  convert(stmt.Left()),
				Right: convert(stmt.Right()),
				Body:  convert(stmt.Body()),
			}
		}
		return &ForOfStatement{
			Type:  "ForOfStatement",
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
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
		}
	case parser.N_STMT_BRK:
		stmt := node.(*parser.BrkStmt)
		return &BreakStatement{
			Type:  "BreakStatement",
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
		}
	case parser.N_STMT_LABEL:
		stmt := node.(*parser.LabelStmt)
		return &LabeledStatement{
			Type:  "LabeledStatement",
			Loc:   loc(stmt.Loc()),
			Label: convert(stmt.Label()),
			Body:  convert(stmt.Body()),
		}
	case parser.N_STMT_THROW:
		stmt := node.(*parser.ThrowStmt)
		return &ThrowStatement{
			Type:     "ThrowStatement",
			Loc:      loc(stmt.Loc()),
			Argument: convert(stmt.Arg()),
		}
	case parser.N_STMT_TRY:
		stmt := node.(*parser.TryStmt)
		return &TryStatement{
			Type:      "TryStatement",
			Loc:       loc(stmt.Loc()),
			Block:     convert(stmt.Try()),
			Handler:   convert(stmt.Catch()),
			Finalizer: convert(stmt.Fin()),
		}
	case parser.N_CATCH:
		expr := node.(*parser.Catch)
		return &CatchClause{
			Type:  "CatchClause",
			Loc:   loc(expr.Loc()),
			Param: convert(expr.Param()),
			Body:  convert(expr.Body()),
		}
	case parser.N_STMT_DEBUG:
		return &DebuggerStatement{
			Type: "DebuggerStatement",
			Loc:  loc(node.Loc()),
		}
	case parser.N_STMT_WITH:
		stmt := node.(*parser.WithStmt)
		return &WithStatement{
			Type:   "WithStatement",
			Loc:    loc(node.Loc()),
			Object: convert(stmt.Expr()),
			Body:   convert(stmt.Body()),
		}
	case parser.N_IMPORT_CALL:
		stmt := node.(*parser.ImportCall)
		return &ImportExpression{
			Type:   "ImportExpression",
			Loc:    loc(node.Loc()),
			Source: convert(stmt.Src()),
		}
	case parser.N_META_PROP:
		stmt := node.(*parser.MetaProp)
		return &MetaProperty{
			Type:     "MetaProperty",
			Loc:      loc(node.Loc()),
			Meta:     convert(stmt.Meta()),
			Property: convert(stmt.Prop()),
		}
	case parser.N_STMT_CLASS:
		stmt := node.(*parser.ClassDec)
		return &ClassDeclaration{
			Type:       "ClassDeclaration",
			Loc:        loc(stmt.Loc()),
			Id:         convert(stmt.Id()),
			SuperClass: convert(stmt.Super()),
			Body:       convert(stmt.Body()),
		}
	case parser.N_ClASS_BODY:
		stmt := node.(*parser.ClassBody)
		return &ClassBody{
			Type: "ClassBody",
			Loc:  loc(stmt.Loc()),
			Body: expressions(stmt.Elems()),
		}
	case parser.N_METHOD:
		n := node.(*parser.Method)
		return &MethodDefinition{
			Type:     "MethodDefinition",
			Loc:      loc(n.Loc()),
			Key:      convert(n.Key()),
			Value:    convert(n.Value()),
			Kind:     n.Kind(),
			Computed: n.Computed(),
			Static:   n.Static(),
		}
	case parser.N_SUPER:
		n := node.(*parser.Super)
		return &Super{
			Type: "Super",
			Loc:  loc(n.Loc()),
		}
	case parser.N_EXPR_TPL:
		tpl := node.(*parser.TplExpr)
		if tpl.Tag() == nil {
			return tplLiteral(tpl.Loc(), tpl.Elems())
		}
		return &TaggedTemplateExpression{
			Type:  "TaggedTemplateExpression",
			Loc:   loc(tpl.LocWithTag()),
			Tag:   convert(tpl.Tag()),
			Quasi: tplLiteral(tpl.Loc(), tpl.Elems()),
		}
	}
	return nil
}
