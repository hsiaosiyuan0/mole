package parser

import "errors"

type Parser struct {
	lexer  *Lexer
	symtab *SymTab
	ver    ESVersion
	srcTyp SourceType
}

type ParserOpts struct {
	Externals  []string
	Version    ESVersion
	SourceType SourceType
}

func NewParserOpts() *ParserOpts {
	return &ParserOpts{
		Externals:  make([]string, 0),
		Version:    ES12,
		SourceType: ST_MODULE,
	}
}

func NewParser(src *Source, opts *ParserOpts) *Parser {
	parser := &Parser{}
	parser.lexer = NewLexer(src)
	parser.symtab = NewSymTab(opts.Externals)
	parser.ver = opts.Version
	parser.srcTyp = opts.SourceType
	return parser
}

func (p *Parser) isStrict() bool {
	if p.srcTyp == ST_MODULE {
		return true
	}
	return p.symtab.Cur.Strict
}

func (p *Parser) Prog() (Node, error) {
	loc := p.loc()
	pg := NewProg()
	for {
		stmt, err := p.stmt()
		if err != nil {
			if err == errEof {
				break
			}
			return nil, err
		}
		if stmt != nil {
			pg.stmts = append(pg.stmts, stmt)
		}
	}
	pg.loc = p.finLoc(loc)
	pg.loc.end.line = p.lexer.src.line
	pg.loc.end.col = p.lexer.src.col
	return pg, nil
}

var errEof = errors.New("eof")

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-Statement
func (p *Parser) stmt() (Node, error) {
	p.lexer.beginStmt = true
	tok := p.lexer.Peek()
	switch tok.value {
	case T_FUNC:
		return p.fnDec(false, false)
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
		return p.classDec(false)
	case T_THROW:
		return p.throwStmt()
	case T_TRY:
		return p.tryStmt()
	case T_BRACE_L:
		return p.blockStmt(false)
	case T_DO:
		return p.doWhileStmt()
	case T_SWITCH:
		return p.switchStmt()
	case T_DEBUGGER:
		return p.debugStmt()
	case T_SEMI:
		return p.emptyStmt()
	case T_COMMENT:
		p.lexer.Next()
		return nil, nil
	case T_WITH:
		return p.withStmt()
	case T_EOF:
		return nil, errEof
	}
	if p.aheadIsVarDec(tok) {
		return p.varDecStmt(false, false)
	} else if p.aheadIsAsync(tok) {
		return p.asyncFnDecStmt()
	} else if p.aheadIsLabel(tok) {
		return p.labelStmt()
	}
	return p.exprStmt()
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec(expr bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	scope := p.symtab.EnterScope()
	scope.Strict = true
	p.lexer.pushMode(LM_STRICT)

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
		p.lexer.Next()
		super, err = p.lhs()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.classBody()
	if err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()
	p.lexer.popMode()

	typ := N_STMT_CLASS
	if expr {
		typ = N_EXPR_CLASS
	}
	return &ClassDec{typ, p.finLoc(loc), id, super, body}, nil
}

func (p *Parser) classBody() (Node, error) {
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
		value, err = p.assignExpr(false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.method(static, key, nil, false, false)
	}
	return &Field{N_FIELD, p.finLoc(loc), key, static, key.Type() != N_NAME, value}, nil
}

func (p *Parser) classElemName() (Node, error) {
	return p.propName(true)
}

func (p *Parser) staticBlock() (Node, error) {
	block, err := p.blockStmt(false)
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

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WithStatement
func (p *Parser) withStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	expr, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	return &WithStmt{N_STMT_WITH, p.finLoc(loc), expr, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DebuggerStatement
func (p *Parser) debugStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()
	p.advanceIfSemi(true)
	return &DebugStmt{N_STMT_DEBUG, p.finLoc(loc)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-TryStatement
func (p *Parser) tryStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	try, err := p.blockStmt(false)
	if err != nil {
		return nil, err
	}

	tok := p.lexer.Peek()
	if tok.value != T_CATCH && tok.value != T_FINALLY {
		return nil, p.error(tok.begin)
	}

	var catch Node
	if tok.value == T_CATCH {
		loc := p.loc()
		p.lexer.Next()
		p.nextMustTok(T_PAREN_L)
		param, err := p.bindingPattern()
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_PAREN_R)

		body, err := p.blockStmt(false)
		if err != nil {
			return nil, err
		}
		catch = &Catch{N_CATCH, p.finLoc(loc), param, body}
	}

	var fin Node
	if p.lexer.Peek().value == T_FINALLY {
		p.lexer.Next()
		fin, err = p.blockStmt(false)
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
		arg, err = p.expr(false)
		if err != nil {
			return nil, err
		}
	}
	p.advanceIfSemi(false)
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
	} else if tok.value != T_ILLEGAL &&
		tok.value != T_BRACE_R &&
		tok.value != T_PAREN_R &&
		tok.value != T_BRACKET_R &&
		tok.value != T_COMMENT &&
		tok.value != T_EOF && !tok.afterLineTerminator {
		arg, err = p.expr(false)
		if err != nil {
			return nil, err
		}
	}
	p.advanceIfSemi(false)
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

	tok := p.lexer.Peek()
	if tok.value == T_NAME && !tok.afterLineTerminator {
		label, err := p.ident()
		if err != nil {
			return nil, err
		}
		p.advanceIfSemi(false)
		return &BrkStmt{N_STMT_BRK, p.finLoc(loc), label}, nil
	}

	p.advanceIfSemi(false)
	return &BrkStmt{N_STMT_BRK, p.finLoc(loc), nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ContinueStatement
func (p *Parser) contStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	if tok.value == T_NAME && !tok.afterLineTerminator {
		label, err := p.ident()
		if err != nil {
			return nil, err
		}
		p.advanceIfSemi(false)
		return &ContStmt{N_STMT_CONT, p.finLoc(loc), label}, nil
	}

	p.advanceIfSemi(false)
	return &ContStmt{N_STMT_CONT, p.finLoc(loc), nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-SwitchStatement
func (p *Parser) switchStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr(false)
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
			return nil, p.error(tok.begin)
		}
		if tok.value == T_DEFAULT && metDefault {
			return nil, p.error(tok.begin)
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
		test, err = p.expr(false)
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
		if stmt != nil {
			cons = append(cons, stmt)
		}
	}
	return &SwitchCase{N_SWITCH_CASE, p.finLoc(loc), test, cons}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-IfStatement
func (p *Parser) ifStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr(false)
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
	test, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	p.nextMustTok(T_PAREN_R)

	p.advanceIfSemi(true)
	return &DoWhileStmt{N_STMT_DO_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WhileStatement
func (p *Parser) whileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	p.nextMustTok(T_PAREN_L)
	test, err := p.expr(false)
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
		init, err = p.varDecStmt(true, true)
		if err != nil {
			return nil, err
		}
	} else if tok.value != T_SEMI {
		init, err = p.expr(true)
		if err != nil {
			return nil, err
		}
	}

	tok = p.lexer.Peek()
	if IsName(tok, "of") || IsName(tok, "in") {
		if init == nil {
			return nil, p.error(tok.begin)
		}

		p.lexer.Next()
		right, err := p.expr(false)
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_PAREN_R)
		body, err := p.stmt()
		if err != nil {
			return nil, err
		}
		return &ForInOfStmt{N_STMT_FOR_IN_OF, p.finLoc(loc), IsName(tok, "in"), await, init, right, body}, nil
	}

	p.nextMustTok(T_SEMI)
	var test Node
	if p.lexer.Peek().value == T_SEMI {
		p.lexer.Next()
	} else {
		test, err = p.expr(false)
		if err != nil {
			return nil, err
		}
		p.nextMustTok(T_SEMI)
	}

	var update Node
	if p.lexer.Peek().value != T_PAREN_R {
		update, err = p.expr(false)
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
	return p.fnDec(false, true)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-HoistableDeclaration
func (p *Parser) fnDec(expr bool, async bool) (Node, error) {
	loc := p.loc()
	if async {
		p.lexer.Next()
	}
	p.lexer.Next()

	p.symtab.EnterScope()

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
		return nil, p.error(tok.begin)
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

	p.symtab.LeaveScope()
	// TODO: check formal params if in strict mode

	typ := N_STMT_FN
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
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
		params = append(params, param)
	}
	return params, nil
}

func (p *Parser) fnBody() (Node, error) {
	return p.blockStmt(true)
}

func (p *Parser) blockStmt(fnBody bool) (*BlockStmt, error) {
	tok, err := p.nextMustTok(T_BRACE_L)
	if err != nil {
		return nil, err
	}
	if !fnBody {
		p.symtab.EnterScope()
	}
	loc := p.locFromTok(tok)

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
		if stmt != nil {
			if fnBody && stmt.Type() == N_STMT_EXPR {
				expr := stmt.(*ExprStmt).expr
				if expr.Type() == N_LIT_STR && expr.(*StrLit).Text() == "use strict" {
					p.symtab.Cur.Strict = true
					p.lexer.addMode(LM_STRICT)
				}
			}
			stmts = append(stmts, stmt)
		}
	}
	if !fnBody {
		p.symtab.LeaveScope()
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
func (p *Parser) varDecStmt(notIn bool, asExpr bool) (Node, error) {
	loc := p.loc()
	kind := p.lexer.Next()

	node := &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]*VarDec, 0, 1)}
	for {
		dec, err := p.varDec(notIn)
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
	if !asExpr {
		p.advanceIfSemi(true)
	}
	node.loc = p.finLoc(loc)
	return node, nil
}

func (p *Parser) varDec(notIn bool) (*VarDec, error) {
	binding, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}
	loc := binding.Loc().Clone()

	var init Node
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr(notIn)
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
	ident := &Ident{N_NAME, &Loc{}, nil, false}
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
			return nil, p.error(node.Loc().begin)
		}
		props = append(props, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(tok.begin)
		}
	}
	return &ObjPattern{N_PATTERN_OBJ, p.finLoc(loc), props}, nil
}

func (p *Parser) patternProp() (Node, error) {
	loc := p.loc()

	key, err := p.propName(false)
	if err != nil {
		return nil, err
	}

	if p.lexer.Peek().value != T_COLON {
		if key.Type() == N_NAME {
			return p.patternAssign(key)
		}
		return nil, p.error(key.Loc().begin)
	}

	p.lexer.Next()
	value, err := p.bindingElem()
	if err != nil {
		return nil, err
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key), nil}, nil
}

func (p *Parser) propName(allowNamePVT bool) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	_, ok := tok.CanBePropKey()
	if ok || (allowNamePVT && tok.value == T_NAME_PVT) {
		if tok.value == T_NUM {
			return &NumLit{N_LIT_NUM, p.finLoc(loc), tok}, nil
		}
		return &Ident{N_NAME, p.finLoc(loc), tok, tok.value == T_NAME_PVT}, nil
	}
	if tok.value == T_STRING {
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok}, nil
	}
	if tok.value == T_NUM {
		return &NumLit{N_LIT_NUM, p.finLoc(loc), tok}, nil
	}
	if tok.value == T_BRACKET_L {
		name, err := p.assignExpr(false)
		if err != nil {
			return nil, err
		}
		_, err = p.nextMustTok(T_BRACKET_R)
		if err != nil {
			return nil, err
		}
		return name, nil
	}
	return nil, p.error(loc.begin)
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
			return nil, p.error(node.Loc().begin)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(tok.begin)
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
		init, err = p.assignExpr(false)
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
	loc := p.loc()
	stmt := &ExprStmt{N_STMT_EXPR, &Loc{}, nil}
	expr, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	stmt.expr = p.unParen(expr)
	p.advanceIfSemi(false)
	stmt.loc = p.finLoc(loc)
	return stmt, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-Expression
func (p *Parser) expr(notIn bool) (Node, error) {
	return p.seqExpr(notIn)
}

func (p *Parser) seqExpr(notIn bool) (Node, error) {
	loc := p.loc()
	expr, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}
	if p.lexer.Peek().value != T_COMMA {
		return expr, nil
	}

	exprs := make([]Node, 0)
	exprs = append(exprs, expr)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
			expr, err = p.assignExpr(notIn)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.error(p.lexer.Loc().begin)
			}
			exprs = append(exprs, expr)
		} else {
			break
		}
	}
	return &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), exprs}, nil
}

func (p *Parser) assignExpr(notIn bool) (Node, error) {
	loc := p.loc()
	lhs, err := p.condExpr(notIn)
	if err != nil {
		return nil, err
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == nil {
		return lhs, nil
	}

	rhs, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finLoc(loc), assign, lhs, p.unParen(rhs)}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ConditionalExpression
func (p *Parser) condExpr(notIn bool) (Node, error) {
	loc := p.loc()
	test, err := p.binExpr(nil, 0, notIn)
	if err != nil {
		return nil, err
	}

	if hook := p.advanceIfTok(T_HOOK); hook == nil {
		return test, nil
	}

	cons, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}

	_, err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_COND, p.finLoc(loc), test, cons, alt}
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

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-NewExpression
func (p *Parser) newExpr() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var expr Node
	var err error
	expr, err = p.memberExpr(nil, false)
	if err != nil {
		return nil, err
	}

	var args []Node
	if p.lexer.Peek().value == T_PAREN_L {
		args, err = p.argList()
		if err != nil {
			return nil, err
		}
	}

	var ret Node
	ret = &NewExpr{N_EXPR_NEW, p.finLoc(loc), expr, args}
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_L {
			if ret, err = p.callExpr(ret); err != nil {
				return nil, err
			}
		} else if tok.value == T_BRACKET_L || tok.value == T_DOT {
			if ret, err = p.memberExpr(ret, true); err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return ret, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-CallExpression
func (p *Parser) callExpr(callee Node) (Node, error) {
	var loc *Loc
	var err error
	if callee == nil {
		loc = p.loc()
		callee, err = p.primaryExpr()
		if err != nil {
			return nil, err
		}
	} else {
		loc = callee.Loc().Clone()
	}

	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_L {
			args, err := p.argList()
			if err != nil {
				return nil, err
			}
			callee = &CallExpr{N_EXPR_CALL, p.finLoc(loc), p.unParen(callee), args}
		} else if tok.value == T_BRACKET_L || tok.value == T_DOT {
			callee, err = p.memberExpr(callee, true)
			if err != nil {
				return nil, err
			}
			return callee, nil
		} else if tok.value == T_TPL_SPAN {
			callee, err = p.tplExpr(callee)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return callee, nil
}

// https://262.ecma-international.org/12.0/#prod-ImportCall
func (p *Parser) importCall() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	meta := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), tok, false}

	if p.lexer.Peek().value == T_DOT {
		p.lexer.Next()
		prop, err := p.ident()
		if err != nil {
			return nil, err
		}
		if prop.Text() != "meta" {
			return nil, p.error(prop.loc.begin)
		}
		return &MetaProp{N_META_PROP, p.finLoc(loc), meta, prop}, nil
	}

	_, err := p.nextMustTok(T_PAREN_L)
	if err != nil {
		return nil, err
	}
	src, err := p.assignExpr(false)
	if err != nil {
		return nil, err
	}
	_, err = p.nextMustTok(T_PAREN_R)
	if err != nil {
		return nil, err
	}
	return &ImportCall{N_IMPORT_CALL, p.finLoc(loc), src}, nil
}

func (p *Parser) tplExpr(tag Node) (Node, error) {
	loc := p.locFromNode(tag)

	elems := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_TPL_TAIL {
			loc := p.loc()
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, p.finLoc(loc), tok}
			elems = append(elems, str)
			break
		} else if tok.value == T_TPL_SPAN {
			loc := p.loc()
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, p.finLoc(loc), tok}
			elems = append(elems, str)

			expr, err := p.expr(false)
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		} else {
			return nil, p.error(tok.begin)
		}
	}

	return &TplExpr{N_EXPR_TPL, p.finLoc(loc), tag, elems}, nil
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
	return p.assignExpr(false)
}

func (p *Parser) binExpr(lhs Node, minPcd int, notIn bool) (Node, error) {
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
		if !ahead.IsBin(notIn) || pcd < minPcd {
			break
		}
		op := p.lexer.Next()
		rhs, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		kind = ahead.Kind()
		for ahead.IsBin(notIn) && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			rhs, err = p.binExpr(rhs, pcd, notIn)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			kind = ahead.Kind()
		}
		pcd = kind.Pcd

		bin := &BinExpr{N_EXPR_BIN, nil, nil, nil, nil}
		bin.loc = p.finLoc(lhs.Loc().Clone())
		bin.op = op
		bin.lhs = p.unParen(lhs)
		bin.rhs = p.unParen(rhs)
		lhs = bin
	}
	return lhs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-MemberExpression
func (p *Parser) memberExpr(obj Node, call bool) (Node, error) {
	var err error
	if obj == nil {
		if p.lexer.Peek().value == T_NEW {
			if obj, err = p.newExpr(); err != nil {
				return nil, err
			}
		} else if obj, err = p.memberExprObj(); err != nil {
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
	if call && p.lexer.Peek().value == T_PAREN_L {
		return p.callExpr(obj)
	}
	return obj, nil
}

func (p *Parser) memberExprObj() (Node, error) {
	return p.primaryExpr()
}

func (p *Parser) memberExprPropSubscript(obj Node) (Node, error) {
	p.lexer.Next()
	prop, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), p.unParen(obj), prop, true, false}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node) (Node, error) {
	p.lexer.Next()

	loc := p.loc()
	tok := p.lexer.Next()
	_, ok := tok.CanBePropKey()

	var prop Node
	if (ok && tok.value != T_NUM) || tok.value == T_NAME_PVT {
		prop = &Ident{N_NAME, p.finLoc(loc), tok, tok.value == T_NAME_PVT}
	} else {
		return nil, p.error(tok.begin)
	}

	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), p.unParen(obj), prop, false, false}
	return node, nil
}

func (p *Parser) primaryExpr() (Node, error) {
	tok := p.lexer.Peek()
	loc := p.locFromTok(tok)

	switch tok.value {
	case T_NUM:
		p.lexer.Next()
		return &NumLit{N_LIT_NUM, p.finLoc(loc), tok}, nil
	case T_STRING:
		p.lexer.Next()
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok}, nil
	case T_NULL:
		p.lexer.Next()
		return &NullLit{N_LIT_NULL, p.finLoc(loc)}, nil
	case T_TRUE, T_FALSE:
		p.lexer.Next()
		return &BoolLit{N_LIT_BOOL, p.finLoc(loc), tok}, nil
	case T_NAME:
		if p.aheadIsAsync(tok) {
			return p.fnDec(true, true)
		}
		p.lexer.Next()
		return &Ident{N_NAME, p.finLoc(loc), tok, false}, nil
	case T_THIS:
		p.lexer.Next()
		return &ThisExpr{N_EXPR_THIS, p.finLoc(loc)}, nil
	case T_PAREN_L:
		return p.parenExpr()
	case T_BRACKET_L:
		return p.arrLit()
	case T_BRACE_L:
		return p.objLit()
	case T_FUNC:
		return p.fnDec(true, false)
	case T_REGEXP:
		p.lexer.Next()
		ext := tok.ext.(*TokExtRegexp)
		return &RegexpLit{N_LIT_REGEXP, p.finLoc(loc), tok, ext.Pattern(), ext.Flags()}, nil
	case T_CLASS:
		return p.classDec(true)
	case T_SUPER:
		p.lexer.Next()
		return &Super{N_SUPER, p.finLoc(loc)}, nil
	case T_IMPORT:
		return p.importCall()
	}
	return nil, p.error(tok.begin)
}

func (p *Parser) parenExpr() (Node, error) {
	loc := p.loc()
	params, err := p.argList()
	if err != nil {
		return nil, err
	}
	if p.lexer.Peek().value == T_ARROW {
		p.lexer.Next()
		p.symtab.EnterScope()

		var body Node
		var err error
		if p.lexer.Peek().value == T_BRACE_L {
			body, err = p.fnBody()
		} else {
			body, err = p.expr(false)
		}
		if err != nil {
			return nil, err
		}

		p.symtab.LeaveScope()
		// TODO: check params

		return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), false, false, params, body}, nil
	}
	if len(params) == 0 {
		return nil, p.error(p.loc().begin)
	}
	if len(params) == 1 {
		return &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), params[0]}, nil
	}
	// TODO: check spread
	return &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), params}, nil
}

func (p *Parser) unParen(expr Node) Node {
	if expr.Type() == N_EXPR_PAREN {
		return expr.(*ParenExpr).Expr()
	}
	return expr
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
			return nil, p.error(node.Loc().begin)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(tok.begin)
		}
	}
	return &ArrLit{N_LIT_ARR, p.finLoc(loc), elems}, nil
}

func (p *Parser) arrElem() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(false)
}

func (p *Parser) spread() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	node, err := p.assignExpr(false)
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
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		}

		node, err := p.objProp()
		if err != nil {
			return nil, err
		}
		props = append(props, node)

		tok = p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.error(tok.begin)
		}
	}
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props}, nil
}

func (p *Parser) objProp() (Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_DOT_TRI {
		return p.spread()
	}

	if name, ok := tok.CanBePropKey(); ok {
		if name == "get" || name == "set" {
			ahead := p.lexer.PeekGrow()
			isField := ahead.value == T_COLON ||
				ahead.value == T_ASSIGN ||
				ahead.value == T_SEMI ||
				ahead.afterLineTerminator
			if !isField {
				p.lexer.Next()
				return p.propMethod(nil, tok, false, false, false)
			}
		}
	} else if tok.value == T_MUL {
		return p.propMethod(nil, nil, true, false, false)
	} else if p.aheadIsAsync(tok) {
		return p.propMethod(nil, nil, false, true, false)
	}

	return p.propField()
}

func (p *Parser) propField() (Node, error) {
	loc := p.loc()
	key, err := p.propName(false)
	if err != nil {
		return nil, err
	}

	var value Node
	tok := p.lexer.Peek()
	if tok.value == T_COLON {
		p.lexer.Next()
		value, err = p.assignExpr(false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.propMethod(key, nil, false, false, false)
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key), nil}, nil
}

func (p *Parser) propMethod(key Node, kind *Token, gen bool, async bool, allowNamePVT bool) (Node, error) {
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
		key, err = p.propName(allowNamePVT)
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
	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key), kind}, nil
}

func (p *Parser) advanceIfSemi(cmt bool) {
	if cmt && p.lexer.Peek().value == T_COMMENT {
		p.lexer.Next()
	}
	if p.lexer.Peek().value == T_SEMI {
		p.lexer.Next()
	}
}

func (p *Parser) nextMustTok(val TokenValue) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != val {
		return nil, p.error(tok.begin)
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
	return p.lexer.Loc()
}

func (p *Parser) locFromNode(node Node) *Loc {
	if node != nil {
		return node.Loc().Clone()
	}
	return p.loc()
}

func (p *Parser) locFromTok(tok *Token) *Loc {
	return &Loc{
		src:   tok.raw.src,
		begin: tok.begin.Clone(),
		end:   &Pos{},
		rng:   &Range{},
	}
}

func (p *Parser) finLoc(loc *Loc) *Loc {
	return p.lexer.FinLoc(loc)
}

func (p *Parser) error(loc *Pos) *ParserError {
	return NewParserError("unexpected token at",
		p.lexer.src.path, loc.line, loc.col)
}

func IsLitPropName(node Node) bool {
	typ := node.Type()
	return typ == N_NAME || typ == N_LIT_STR || typ == N_LIT_NUM
}
