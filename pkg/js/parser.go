package js

type Parser struct {
	lexer  *Lexer
	symtab *SymTab
}

func NewParser(src *Source, externals []string) *Parser {
	parser := &Parser{}
	parser.lexer = NewLexer(src)
	parser.symtab = NewSymTab(externals)
	return parser
}

func (p *Parser) Prog() (Node, error) {
	pg := NewProg()
	for {
		if p.lexer.src.AheadIsEof() {
			break
		}
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}
		pg.stmts = append(pg.stmts, stmt)
	}
	return pg, nil
}

func (p *Parser) stmt() (Node, error) {
	return p.exprStmt()
}

func (p *Parser) exprStmt() (Node, error) {
	stmt := NewExprStmt()
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	stmt.expr = expr
	return stmt, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-Expression
func (p *Parser) expr() (Node, error) {
	return p.assignExpr()
}

func (p *Parser) assignExpr() (Node, error) {
	loc := p.loc()
	lhs, err := p.condExpr()
	if err != nil {
		return nil, err
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == nil {
		return lhs, nil
	}

	rhs, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finLoc(loc), assign, lhs, rhs}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ConditionalExpression
func (p *Parser) condExpr() (Node, error) {
	loc := p.loc()
	test, err := p.binExpr(nil, 0)
	if err != nil {
		return nil, err
	}

	if hook := p.advanceIfTok(T_HOOK); hook == nil {
		return test, nil
	}

	cons, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_BIN, p.finLoc(loc), test, cons, alt}
	return node, nil
}

func (p *Parser) unaryExpr() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Peek()
	if tok.IsUnary() || tok.value == T_ADD || tok.value == T_SUB {
		p.lexer.Next()
		arg, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), tok, arg}, nil
	}
	return p.updateExpr()
}

func (p *Parser) updateExpr() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Peek()
	if tok.value == T_INC || tok.value == T_DEC {
		p.lexer.Next()
		arg, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}
		return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok, true, arg}, nil
	}

	arg, err := p.lhs()
	if err != nil {
		return nil, err
	}

	tok = p.lexer.Peek()
	postfix := !tok.afterLineTerminator && (tok.value == T_INC || tok.value == T_DEC)
	if !postfix {
		return arg, nil
	}

	p.lexer.Next()
	return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok, false, arg}, nil
}

func (p *Parser) lhs() (Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_NEW {
		return p.newExpr()
	}
	return p.callExpr(nil)
}

func (p *Parser) newExpr() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	expr, err := p.memberExpr(nil)
	if err != nil {
		return nil, err
	}

	node := NewNewExpr()
	node.expr = expr
	node.loc = p.finLoc(loc)
	return node, nil
}

func (p *Parser) callExpr(callee Node) (Node, error) {
	// TODO: SuperCall ImportCall
	return p.memberExpr(nil)
}

func (p *Parser) binExpr(lhs Node, minPcd int) (Node, error) {
	var err error
	if lhs == nil {
		if lhs, err = p.unaryExpr(); err != nil {
			return nil, err
		}
	}

	ahead := p.lexer.Peek()
	kind := ahead.Kind()
	pcd := kind.Pcd
	for {
		if !ahead.IsBin() || pcd < minPcd {
			break
		}
		op := p.lexer.Next()
		rhs, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		kind = ahead.Kind()
		for ahead.IsBin() && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			rhs, err = p.binExpr(rhs, pcd)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			kind = ahead.Kind()
		}
		pcd = kind.Pcd

		bin := NewBinExpr()
		bin.loc = p.finLoc(lhs.Loc())
		bin.op = op
		bin.lhs = lhs
		bin.rhs = rhs
		lhs = bin
	}
	return lhs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-MemberExpression
func (p *Parser) memberExpr(obj Node) (Node, error) {
	var err error
	if obj == nil {
		if obj, err = p.memberExprObj(); err != nil {
			return nil, err
		}
	}
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACKET_L {
			if obj, err = p.memberExprPropSubscript(obj); err != nil {
				return nil, err
			}
		} else if tok.value == T_DOT {
			if obj, err = p.memberExprPropDot(obj); err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return obj, nil
}

func (p *Parser) memberExprObj() (Node, error) {
	return p.primaryExpr()
}

func (p *Parser) memberExprPropSubscript(obj Node) (Node, error) {
	p.lexer.Next()
	prop, err := p.expr()
	if err != nil {
		return nil, err
	}
	if err := p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, false}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node) (Node, error) {
	p.lexer.Next()
	loc := p.loc()
	var private bool
	if private = p.lexer.src.AheadIsCh('#'); private {
		p.lexer.src.Read()
	}
	id := p.lexer.Next()
	if id.value != T_NAME {
		return nil, p.error(&id.loc)
	}
	prop := &Ident{N_NAME, p.finLoc(loc), id, private}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, false}
	return node, nil
}

func (p *Parser) primaryExpr() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	switch tok.value {
	case T_NUM:
		node := NewNumLit()
		node.loc = p.finLoc(loc)
		node.val = tok
		return node, nil
	case T_NAME:
		node := NewIdent()
		node.loc = p.finLoc(loc)
		node.val = tok
		return node, nil
	}
	return nil, p.error(&tok.loc)
}

func (p *Parser) nextMustTok(val TokenValue) error {
	tok := p.lexer.Next()
	if tok.value != val {
		return p.error(&tok.loc)
	}
	return nil
}

func (p *Parser) advanceIfTok(val TokenValue) *Token {
	tok := p.lexer.Peek()
	if tok.value != val {
		return nil
	}
	return p.lexer.Next()
}

func (p *Parser) advanceIfTokIn(begin, end TokenValue) *Token {
	tok := p.lexer.Peek()
	if tok.value <= begin || tok.value >= end {
		return nil
	}
	return p.lexer.Next()
}

func (p *Parser) loc() *Loc {
	loc := &Loc{}
	loc.src = p.lexer.src
	loc.begin.line = p.lexer.src.line
	loc.begin.col = p.lexer.src.col
	return loc
}

func (p *Parser) finLoc(loc *Loc) *Loc {
	loc.end.line = p.lexer.src.line
	loc.end.col = p.lexer.src.col
	return loc
}

func (p *Parser) error(loc *Position) *ParserError {
	return NewParserError("unexpected token at",
		p.lexer.src.path, loc.line, loc.col)
}
