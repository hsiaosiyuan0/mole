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
		if p.lexer.Peek().value == T_EOF {
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

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-Statement
func (p *Parser) stmt() (Node, error) {
	tok := p.lexer.Peek()
	switch tok.value {
	case T_FUNC:
		return p.fnDecStmt(false, false)
	case T_IF:
		return p.ifStmt()
	case T_BREAK:
		return p.brkStmt()
	case T_CONTINUE:
		return p.contStmt()
	case T_FOR:
		return p.forStmt()
	case T_RETURN:
		return p.retStmt()
	case T_WHILE:
		return p.whileStmt()
	case T_CLASS:
		return p.classDec()
	case T_THROW:
		return p.throwStmt()
	case T_TRY:
		return p.tryStmt()
	case T_BRACE_L:
		return p.blockStmt()
	case T_DO:
		return p.doWhileStmt()
	case T_SWITCH:
		return p.switchStmt()
	case T_DEBUGGER:
		return p.debugStmt()
	case T_SEMI:
		return p.emptyStmt()
	}
	if p.aheadIsVarDec(tok) {
		return p.varDecStmt()
	} else if p.aheadIsAsync(tok) {
		return p.asyncFnDecStmt()
	} else if p.aheadIsLabel(tok) {
		return p.labelStmt()
	}
	return p.exprStmt()
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var id Node
	var err error
	if p.lexer.Peek().value != T_BRACE_L {
		id, err = p.ident()
		if err != nil {
			return nil, err
		}
	}

	var super Node
	if p.lexer.Peek().value == T_EXTENDS {
		super, err = p.lhs()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.classBody()
	if err != nil {
		return nil, err
	}
	return &ClassStmt{N_STMT_CLASS, p.finLoc(loc), id, super, body}, nil
}

func (p *Parser) classBody() (*ClassBody, error) {
	loc := p.loc()
	p.nextMustTok(T_BRACE_L)
	elems := make([]Node, 0)
	for {
		if p.lexer.Peek().value == T_BRACE_R {
			break
		}
		if p.lexer.Peek().value == T_SEMI {
			p.lexer.Next()
		}
		elem, err := p.classElem()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}
	p.nextMustTok(T_BRACE_R)
	return &ClassBody{N_ClASS_BODY, p.finLoc(loc), elems}, nil
}

func (p *Parser) classElem() (Node, error) {
	static := p.lexer.Peek().value == T_STATIC
	if static {
		p.lexer.Next()
		if p.lexer.Peek().value == T_BRACE_L {
			return p.staticBlock()
		}
	}

	tok := p.lexer.Peek()
	if tok.value == T_NAME {
		name := tok.Text()
		if name == "constructor" {
			return p.method(static, nil, tok, false, false)
		} else if name == "get" || name == "set" {
			ahead := p.lexer.PeekGrow()
			isField := ahead.value == T_ASSIGN || ahead.value == T_SEMI || ahead.afterLineTerminator
			if !isField {
				return p.method(static, nil, tok, false, false)
			}
		}
	} else if tok.value == T_MUL {
		return p.method(static, nil, tok, true, false)
	} else if p.aheadIsAsync(tok) {
		return p.method(static, nil, tok, false, true)
	}

	return p.field(static)
}

func (p *Parser) method(static bool, key Node, kind *Token, gen bool, async bool) (Node, error) {
	loc := p.loc()

	if async {
		p.lexer.Next()
		gen = p.lexer.Peek().value == T_MUL
	}
	if gen {
		p.lexer.Next()
	}

	var err error
	if key == nil {
		key, err = p.classElemName()
		if err != nil {
			return nil, err
		}
	}

	fnLoc := p.loc()
	params, err := p.formalParams()
	if err != nil {
		return nil, err
	}

	if gen {
		p.lexer.extMode(LM_GENERATOR, true)
	}
	body, err := p.fnBody()
	if gen {
		p.lexer.popMode()
	}
	if err != nil {
		return nil, err
	}

	value := &FnDec{N_EXPR_FN, p.finLoc(fnLoc), nil, gen, async, params, body}
	return &Method{N_METHOD, p.finLoc(loc), key, static, key.Type() != N_NAME, kind, value}, nil
}

func (p *Parser) field(static bool) (Node, error) {
	loc := p.loc()
	key, err := p.classElemName()
	if err != nil {
		return nil, err
	}
	var value Node
	tok := p.lexer.Peek()
	if tok.value == T_ASSIGN {
		value, err = p.assignExpr()
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.method(static, key, nil, false, false)
	}
	return &Field{N_FIELD, p.finLoc(loc), key, static, key.Type() != N_NAME, value}, nil
}

func (p *Parser) classElemName() (Node, error) {
	return p.propName()
}

func (p *Parser) staticBlock() (Node, error) {
	block, err := p.blockStmt()
	if err != nil {
		return nil, err
	}
	return &StaticBlock{N_STATIC_BLOCK, block.loc, block.body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-EmptyStatement
func (p *Parser) emptyStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()
	return &EmptyStmt{N_STMT_EMPTY, p.finLoc(loc)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DebuggerStatement
func (p *Parser) debugStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()
	p.advanceIfSemi()
	return &DebugStmt{N_STMT_DEBUG, p.finLoc(loc)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-TryStatement
func (p *Parser) tryStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	try, err := p.blockStmt()
	if err != nil {
		return nil, err
	}

	tok := p.lexer.Peek()
	if tok.value != T_CATCH && tok.value != T_FINALLY {
		return nil, p.error(&tok.loc)
	}

	var catch *Catch
	if tok.value == T_CATCH {
		loc := p.loc()
		p.lexer.Next()
		p.nextMustTok(T_PAREN_L)
		param, err := p.bindingPattern()
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_PAREN_R)

		body, err := p.blockStmt()
		if err != nil {
			return nil, err
		}
		catch = &Catch{N_CATCH, p.finLoc(loc), param, body}
	}

	var fin *BlockStmt
	if p.lexer.Peek().value == T_FINALLY {
		p.lexer.Next()
		fin, err = p.blockStmt()
		if err != nil {
			return nil, err
		}
	}

	return &TryStmt{N_STMT_TRY, p.finLoc(loc), try, catch, fin}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ThrowStatement
func (p *Parser) throwStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	var arg Node
	var err error
	if tok.value == T_SEMI {
		p.lexer.Next()
	} else if tok.value != T_ILLEGAL && tok.value != T_EOF && !tok.afterLineTerminator {
		arg, err = p.expr()
		if err != nil {
			return nil, err
		}
	}
	return &ThrowStmt{N_STMT_THROW, p.finLoc(loc), arg}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ReturnStatement
func (p *Parser) retStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	var arg Node
	var err error
	if tok.value == T_SEMI {
		p.lexer.Next()
	} else if tok.value != T_ILLEGAL && tok.value != T_EOF && !tok.afterLineTerminator {
		arg, err = p.expr()
		if err != nil {
			return nil, err
		}
	}
	return &RetStmt{N_STMT_RET, p.finLoc(loc), arg}, nil
}

func (p *Parser) aheadIsLabel(tok *Token) bool {
	if tok.value == T_NAME {
		ahead := p.lexer.PeekGrow()
		return ahead.value == T_COLON
	}
	return false
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-LabelledStatement
func (p *Parser) labelStmt() (Node, error) {
	loc := p.loc()
	label, err := p.ident()
	if err != nil {
		return nil, err
	}

	// advance `:`
	p.lexer.Next()

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}
	return &LabelStmt{N_STMT_LABEL, p.finLoc(loc), label, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-BreakStatement
func (p *Parser) brkStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var label *Ident
	var err error
	tok := p.lexer.Peek()
	if tok.value == T_SEMI {
		p.lexer.Next()
	} else if tok.value == T_NAME {
		label, err = p.ident()
		if err != nil {
			return nil, err
		}
		p.advanceIfSemi()
	}

	return &BrkStmt{N_STMT_BRK, p.finLoc(loc), label}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ContinueStatement
func (p *Parser) contStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var label *Ident
	var err error
	tok := p.lexer.Peek()
	if tok.value == T_SEMI {
		p.lexer.Next()
	} else if tok.value == T_NAME {
		label, err = p.ident()
		if err != nil {
			return nil, err
		}
		p.advanceIfSemi()
	}

	return &ContStmt{N_STMT_CONT, p.finLoc(loc), label}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-SwitchStatement
func (p *Parser) switchStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	cases := make([]*SwitchCase, 0)
	p.nextMustTok(T_BRACE_L)
	metDefault := false
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value != T_CASE && tok.value != T_DEFAULT {
			return nil, p.error(&tok.loc)
		}
		if tok.value == T_DEFAULT && metDefault {
			return nil, p.error(&tok.loc)
		}

		caseClause, err := p.switchCase(tok)
		if err != nil {
			return nil, err
		}
		if caseClause != nil {
			metDefault = caseClause.test == nil
			cases = append(cases, caseClause)
		} else {
			break
		}
	}
	p.nextMustTok(T_BRACE_R)

	return &SwitchStmt{N_STMT_SWITCH, p.finLoc(loc), test, cases}, nil
}

func (p *Parser) switchCase(tok *Token) (*SwitchCase, error) {
	loc := p.loc()
	p.lexer.nextTok()

	var test Node
	var err error
	if tok.value == T_CASE {
		test, err = p.expr()
		if err != nil {
			return nil, err
		}
	}
	p.nextMustTok(T_COLON)

	cons := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_CASE || tok.value == T_DEFAULT || tok.value == T_BRACE_R {
			break
		}
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}
		cons = append(cons, stmt)
	}
	return &SwitchCase{N_SWITCH_CASE, p.finLoc(loc), test, cons}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-IfStatement
func (p *Parser) ifStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	cons, err := p.stmt()
	if err != nil {
		return nil, err
	}

	var alt Node
	if p.lexer.Peek().value == T_ELSE {
		p.lexer.Next()
		alt, err = p.stmt()
		if err != nil {
			return nil, err
		}
	}
	return &IfStmt{N_STMT_IF, p.finLoc(loc), test, cons, alt}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DoWhileStatement
func (p *Parser) doWhileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	p.nextMustTok(T_WHILE)
	p.nextMustTok(T_PAREN_L)
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	return &DoWhileStmt{N_STMT_DO_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WhileStatement
func (p *Parser) whileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{N_STMT_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ForStatement
func (p *Parser) forStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	await := false
	if IsName(p.lexer.Peek(), "await") {
		await = true
		p.lexer.Next()
	}

	p.nextMustTok(T_PAREN_L)

	tok := p.lexer.Peek()

	var init Node
	var err error
	if tok.value == T_LET || tok.value == T_CONST || tok.value == T_VAR {
		init, err = p.varDecStmt()
		if err != nil {
			return nil, err
		}
	} else if tok.value != T_SEMI {
		init, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	if init != nil && init.Type() == N_EXPR_BIN && init.(*BinExpr).op.value == T_IN {
		p.nextMustTok(T_PAREN_R)
		body, err := p.stmt()
		if err != nil {
			return nil, err
		}
		expr := init.(*BinExpr)
		return &ForInOfStmt{N_STMT_FOR_IN_OF, p.finLoc(loc), true, await, expr.lhs, expr.rhs, body}, nil
	}

	tok = p.lexer.Peek()
	if IsName(tok, "of") {
		if init == nil {
			return nil, p.error(&tok.loc)
		}

		p.lexer.Next()
		right, err := p.expr()
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_PAREN_R)
		body, err := p.stmt()
		if err != nil {
			return nil, err
		}
		return &ForInOfStmt{N_STMT_FOR_IN_OF, p.finLoc(loc), false, await, init, right, body}, nil
	}

	p.nextMustTok(T_SEMI)
	var test Node
	if p.lexer.Peek().value == T_SEMI {
		p.lexer.Next()
	} else {
		test, err = p.expr()
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_SEMI)
	}

	var update Node
	if p.lexer.Peek().value != T_PAREN_R {
		update, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	p.nextMustTok(T_PAREN_R)
	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	return &ForStmt{N_STMT_FOR, p.finLoc(loc), init, test, update, body}, nil
}

func (p *Parser) aheadIsAsync(tok *Token) bool {
	if IsName(tok, "async") {
		ahead := p.lexer.PeekGrow()
		return ahead.value == T_FUNC && !ahead.afterLineTerminator
	}
	return false
}

func (p *Parser) asyncFnDecStmt() (Node, error) {
	return p.fnDecStmt(false, true)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-HoistableDeclaration
func (p *Parser) fnDecStmt(expr bool, async bool) (Node, error) {
	loc := p.loc()
	if async {
		p.lexer.Next()
	}
	p.lexer.Next()

	generator := p.lexer.Peek().value == T_MUL
	if generator {
		p.lexer.Next()
	}

	var id Node
	var err error
	tok := p.lexer.Peek()
	if tok.value != T_PAREN_L {
		id, err = p.ident()
		if err != nil {
			return nil, err
		}
	}
	if !expr && id == nil {
		return nil, p.error(&tok.loc)
	}

	params, err := p.formalParams()
	if err != nil {
		return nil, err
	}

	if generator {
		p.lexer.extMode(LM_GENERATOR, true)
	}
	body, err := p.fnBody()
	if generator {
		p.lexer.popMode()
	}

	if err != nil {
		return nil, err
	}

	typ := N_STMT_FN_DEC
	if expr {
		typ = N_EXPR_FN
	}
	return &FnDec{typ, p.finLoc(loc), id, generator, async, params, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-FormalParameters
func (p *Parser) formalParams() ([]Node, error) {
	p.lexer.Next()
	params := make([]Node, 0)
	for {
		if p.lexer.Peek().value == T_PAREN_R {
			p.lexer.Next()
			break
		}
		param, err := p.bindingElem()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

func (p *Parser) fnBody() (Node, error) {
	return p.blockStmt()
}

func (p *Parser) blockStmt() (*BlockStmt, error) {
	loc := p.loc()
	p.lexer.Next()

	stmts := make([]Node, 0)
	for {
		if p.lexer.Peek().value == T_BRACE_R {
			p.lexer.Next()
			break
		}
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return &BlockStmt{N_STMT_BLOCK, p.finLoc(loc), stmts}, nil
}

func (p *Parser) aheadIsVarDec(tok *Token) bool {
	if tok.value == T_VAR {
		return true
	}
	return IsName(tok, "let") || IsName(tok, "const")
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-VariableStatement
func (p *Parser) varDecStmt() (Node, error) {
	loc := p.loc()
	kind := p.lexer.Next()

	node := NewVarDecStmt()
	for {
		dec, err := p.varDec()
		if err != nil {
			return nil, err
		}
		node.decList = append(node.decList, dec)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		} else {
			break
		}
	}
	if IsName(kind, "let") {
		node.kind = T_LET
	} else if IsName(kind, "const") {
		node.kind = T_CONST
	} else {
		node.kind = T_VAR
	}
	node.loc = p.finLoc(loc)
	return node, nil
}

func (p *Parser) varDec() (*VarDec, error) {
	loc := p.loc()

	binding, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}

	var init Node
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr()
		if err != nil {
			return nil, err
		}
	}

	dec := &VarDec{N_VAR_DEC, p.finLoc(loc), binding, init}
	return dec, nil
}

func (p *Parser) ident() (*Ident, error) {
	loc := p.loc()
	tok, err := p.nextMustTok(T_NAME)
	if err != nil {
		return nil, err
	}
	ident := NewIdent()
	ident.loc = p.finLoc(loc)
	ident.val = tok
	return ident, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#sec-destructuring-binding-patterns
func (p *Parser) bindingPattern() (Node, error) {
	tok := p.lexer.Peek()
	var binding Node
	var err error
	if tok.value == T_BRACE_L {
		binding, err = p.patternObj()
	} else if tok.value == T_BRACKET_L {
		binding, err = p.patternArr()
	} else {
		binding, err = p.ident()
	}
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func (p *Parser) patternObj() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	props := make([]Node, 0, 1)
	for {
		node, err := p.patternProp()
		if err != nil {
			return nil, err
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PATTERN_REST && tok.value != T_BRACE_R {
			return nil, p.error(&node.Loc().begin)
		}
		props = append(props, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(&tok.loc)
		}
	}
	return &ObjPattern{N_PATTERN_OBJ, p.finLoc(loc), props}, nil
}

func (p *Parser) patternProp() (Node, error) {
	loc := p.loc()

	key, err := p.propName()
	if err != nil {
		return nil, err
	}

	if p.lexer.Peek().value != T_COLON {
		if key.Type() == N_NAME {
			return p.patternAssign(key)
		}
		return nil, p.error(&key.Loc().begin)
	}

	p.lexer.Next()
	value, err := p.bindingElem()
	if err != nil {
		return nil, err
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key)}, nil
}

func (p *Parser) propName() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	if tok.value == T_NAME || tok.value == T_NAME_PVT {
		return &Ident{N_NAME, p.finLoc(loc), tok, tok.value == T_NAME_PVT}, nil
	}
	if tok.value == T_STRING {
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok}, nil
	}
	if tok.value == T_NUM {
		return &NumLit{N_LIT_NUM, p.finLoc(loc), tok}, nil
	}
	if tok.value == T_BRACKET_L {
		name, err := p.assignExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.nextMustTok(T_BRACKET_R)
		if err != nil {
			return nil, err
		}
		return name, nil
	}
	return nil, p.error(&loc.begin)
}

func (p *Parser) patternArr() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	elems := make([]Node, 0, 1)
	for {
		elems = append(elems, p.elision()...)
		if p.lexer.Peek().value == T_BRACKET_R {
			p.lexer.Next()
			break
		}

		node, err := p.bindingElem()
		if err != nil {
			return nil, err
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PATTERN_REST && tok.value != T_BRACKET_R {
			return nil, p.error(&node.Loc().begin)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(&tok.loc)
		}
	}
	return &ArrayPattern{N_PATTERN_ARRAY, p.finLoc(loc), elems}, nil
}

func (p *Parser) elision() []Node {
	ret := make([]Node, 0, 1)
	for {
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
			ret = append(ret, nil)
		} else {
			break
		}
	}
	return ret
}

func (p *Parser) bindingElem() (Node, error) {
	tok := p.lexer.Peek()
	var binding Node
	var err error
	if tok.value == T_BRACE_L {
		binding, err = p.patternObj()
	} else if tok.value == T_BRACKET_L {
		binding, err = p.patternArr()
	} else if tok.value == T_DOT_TRI {
		binding, err = p.patternRest()
	} else {
		binding, err = p.ident()
	}
	if err != nil {
		return nil, err
	}
	return p.patternAssign(binding)
}

func (p *Parser) patternAssign(ident Node) (Node, error) {
	var init Node
	var err error
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr()
		if err != nil {
			return nil, err
		}
	}

	if init == nil {
		return ident, nil
	}
	return &AssignPattern{N_PATTERN_ASSIGN, p.finLoc(ident.Loc().Clone()), ident, init}, nil
}

func (p *Parser) patternRest() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var arg Node
	var err error
	tok := p.lexer.Peek()
	if tok.value == T_BRACE_L {
		arg, err = p.patternObj()
	} else if tok.value == T_BRACKET_L {
		arg, err = p.patternArr()
	} else {
		arg, err = p.ident()
	}

	if err != nil {
		return nil, err
	}
	return &RestPattern{N_PATTERN_REST, p.finLoc(loc), arg}, nil
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

	p.advanceIfSemi()
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

	_, err = p.nextMustTok(T_COLON)
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

	var expr Node
	var err error
	if p.lexer.Peek().value == T_NEW {
		expr, err = p.newExpr()
		if err != nil {
			return nil, err
		}
	} else {
		expr, err = p.memberExpr(nil)
		if err != nil {
			return nil, err
		}
	}

	node := NewNewExpr()
	node.expr = expr
	node.loc = p.finLoc(loc)
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-CallExpression
func (p *Parser) callExpr(callee Node) (Node, error) {
	// TODO: SuperCall ImportCall
	loc := p.loc()

	var err error
	if callee == nil {
		callee, err = p.memberExpr(nil)
		if err != nil {
			return nil, err
		}
	}

	tok := p.lexer.Peek()
	if tok.value == T_BRACKET_L {
		callee, err = p.memberExpr(callee)
		if err != nil {
			return nil, err
		}
	}

	if tok.value != T_PAREN_L {
		return callee, nil
	}

	args, err := p.argList()
	if err != nil {
		return nil, err
	}
	node := &CallExpr{N_EXPR_CALL, p.finLoc(loc), callee, args}

	if p.lexer.Peek().value == T_PAREN_L {
		return p.callExpr(node)
	}
	return node, nil
}

func (p *Parser) argList() ([]Node, error) {
	p.lexer.Next()
	args := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_PAREN_R {
			break
		}
		arg, err := p.arg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	p.nextMustTok(T_PAREN_R)
	return args, nil
}

func (p *Parser) arg() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr()
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
	if p.lexer.Peek().value == T_PAREN_L {
		return p.callExpr(obj)
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
	if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, false}
	return node, nil
}

func (p *Parser) aheadIsPvt() bool {
	return p.lexer.src.AheadIsCh('#')
}

func (p *Parser) namePvt() (Node, error) {
	loc := p.loc()
	id := p.lexer.Next()
	if id.value != T_NAME && id.value != T_NAME_PVT {
		return nil, p.error(&id.loc)
	}
	return &Ident{N_NAME, p.finLoc(loc), id, id.value == T_NAME_PVT}, nil
}

func (p *Parser) memberExprPropDot(obj Node) (Node, error) {
	p.lexer.Next()
	prop, err := p.namePvt()
	if err != nil {
		return nil, err
	}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, false}
	return node, nil
}

func (p *Parser) primaryExpr() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Peek()

	switch tok.value {
	case T_NUM:
		p.lexer.Next()
		node := NewNumLit()
		node.loc = p.finLoc(loc)
		node.val = tok
		return node, nil
	case T_NAME:
		p.lexer.Next()
		node := NewIdent()
		node.loc = p.finLoc(loc)
		node.val = tok
		return node, nil
	case T_BRACKET_L:
		return p.arrLit()
	case T_BRACE_L:
		return p.objLit()
	case T_FUNC:
		return p.fnDecStmt(true, false)
	}
	if IsName(tok, "async") {
		return p.fnDecStmt(true, true)
	}
	return nil, p.error(&tok.loc)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ArrayLiteral
func (p *Parser) arrLit() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	elems := make([]Node, 0, 1)
	for {
		elems = append(elems, p.elision()...)
		if p.lexer.Peek().value == T_BRACKET_R {
			p.lexer.Next()
			break
		}

		node, err := p.arrElem()
		if err != nil {
			return nil, err
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PATTERN_REST && tok.value != T_BRACKET_R {
			return nil, p.error(&node.Loc().begin)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(&tok.loc)
		}
	}
	return &ArrLit{N_LIT_ARR, p.finLoc(loc), elems}, nil
}

func (p *Parser) arrElem() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr()
}

func (p *Parser) spread() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	node, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	return &Spread{N_SPREAD, p.finLoc(loc), node}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ObjectLiteral
func (p *Parser) objLit() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	props := make([]Node, 0, 1)
	for {
		node, err := p.objProp()
		if err != nil {
			return nil, err
		}

		props = append(props, node)

		tok := p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(&tok.loc)
		}
	}
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props}, nil
}

func (p *Parser) objProp() (Node, error) {
	loc := p.loc()

	tok := p.lexer.Peek()
	if tok.value == T_DOT_TRI {
		return p.spread()
	}

	key, err := p.propName()
	if err != nil {
		return nil, err
	}

	if p.lexer.Peek().value != T_COLON {
		if key.Type() == N_NAME {
			return key, nil
		}
		return nil, p.error(&key.Loc().begin)
	}

	p.lexer.Next()
	value, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key)}, nil

	// TODO: MethodDefinition
}

func (p *Parser) advanceIfSemi() {
	if p.lexer.Peek().value == T_SEMI {
		p.lexer.Next()
	}
}

func (p *Parser) nextMustTok(val TokenValue) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != val {
		return nil, p.error(&tok.loc)
	}
	return tok, nil
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

func IsLitPropName(node Node) bool {
	typ := node.Type()
	return typ == N_NAME || typ == N_LIT_STR || typ == N_LIT_NUM
}
