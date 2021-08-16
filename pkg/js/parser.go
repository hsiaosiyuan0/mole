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
	return p.condExpr()
}

func (p *Parser) condExpr() (Node, error) {
	return p.binExpr(nil, 0)
}

func (p *Parser) unaryExpr() (Node, error) {
	return p.lhs()
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

	expr, err := p.memberExpr()
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
	return p.memberExpr()
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

func (p *Parser) memberExpr() (Node, error) {
	return p.primaryExpr()
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
	return nil, p.error(loc)
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

func (p *Parser) error(loc *Loc) *ParserError {
	return NewParserError("unexpected token at",
		p.lexer.src.path,
		loc.begin.line,
		loc.begin.col)
}
