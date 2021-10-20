package parser

import (
	"errors"
	"fmt"
)

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
	parser.lexer.ver = opts.Version
	parser.symtab = NewSymTab(opts.Externals)
	parser.ver = opts.Version
	parser.srcTyp = opts.SourceType
	return parser
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
	pg.loc.rng.end = p.lexer.src.Pos()
	return pg, nil
}

var errEof = errors.New("eof")

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-Statement
func (p *Parser) stmt() (node Node, err error) {
	tok := p.lexer.PeekStmtBegin()

	if tok.value > T_KEYWORD_BEGIN && tok.value < T_KEYWORD_END {
		switch tok.value {
		case T_VAR:
			node, err = p.varDecStmt(false, false)
		case T_FUNC:
			node, err = p.fnDec(false, false)
		case T_IF:
			node, err = p.ifStmt()
		case T_FOR:
			node, err = p.forStmt()
		case T_RETURN:
			node, err = p.retStmt()
		case T_WHILE:
			node, err = p.whileStmt()
		case T_CLASS:
			node, err = p.classDec(false)
		case T_BREAK:
			node, err = p.brkStmt()
		case T_CONTINUE:
			node, err = p.contStmt()
		case T_THROW:
			node, err = p.throwStmt()
		case T_TRY:
			node, err = p.tryStmt()
		case T_DO:
			node, err = p.doWhileStmt()
		case T_SWITCH:
			node, err = p.switchStmt()
		case T_DEBUGGER:
			node, err = p.debugStmt()
		case T_IMPORT:
			node, err = p.importDec()
		case T_EXPORT:
			node, err = p.exportDec()
		case T_WITH:
			node, err = p.withStmt()
		case T_NEW, T_THIS, T_SUPER:
			node, err = p.exprStmt()
		}
	} else if tok.value == T_BRACE_L {
		node, err = p.blockStmt(true)
	} else if p.aheadIsVarDec(tok) {
		node, err = p.varDecStmt(false, false)
	} else if p.aheadIsAsync(tok) {
		node, err = p.asyncFnDecStmt()
	} else if p.aheadIsLabel(tok) {
		node, err = p.labelStmt()
	} else if tok.value == T_SEMI {
		node, err = p.emptyStmt()
	} else if tok.value == T_EOF {
		node, err = nil, errEof
	} else {
		node, err = p.exprStmt()
	}

	if err != nil {
		return nil, err
	}

	p.lexer.beginStmt = false
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ExportDeclaration
func (p *Parser) exportDec() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	var err error
	node := &ExportDec{N_STMT_EXPORT, nil, false, false, nil, nil, nil}
	specs := make([]Node, 0)
	tok := p.lexer.Peek()
	if tok.value == T_MUL || tok.value == T_BRACE_L {
		ss, all, src, err := p.exportFrom()
		node.src = src
		node.all = all
		specs = append(specs, ss...)
		if err != nil {
			return nil, err
		}
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
	} else if p.aheadIsVarDec(tok) {
		node.dec, err = p.varDecStmt(false, false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_FUNC {
		node.dec, err = p.fnDec(false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsAsync(tok) {
		node.dec, err = p.fnDec(false, true)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_CLASS {
		node.dec, err = p.classDec(false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_DEFAULT {
		p.lexer.Next()
		tok := p.lexer.Peek()
		node.def = true
		if tok.value == T_FUNC {
			node.dec, err = p.fnDec(false, false)
		} else if p.aheadIsAsync(tok) {
			node.dec, err = p.fnDec(false, true)
		} else if tok.value == T_CLASS {
			node.dec, err = p.classDec(false)
		} else {
			node.dec, err = p.assignExpr(false)
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, p.errorTok(tok)
	}

	node.loc = p.finLoc(loc)
	node.specs = specs
	return node, nil
}

func (p *Parser) exportFrom() ([]Node, bool, Node, error) {
	tok := p.lexer.Next()
	var specs []Node
	var err error

	ns := false
	if tok.value == T_MUL {
		ns = true
		_, err = p.nextMustName("as")
		if err != nil {
			return nil, false, nil, err
		}

		id, err := p.ident()
		if err != nil {
			return nil, false, nil, err
		}
		specs = make([]Node, 1)
		specs[0] = &ExportSpec{N_EXPORT_SPEC, p.finLoc(p.locFromTok(tok)), true, id, nil}
	} else {
		specs, err = p.exportNamed()
		if err != nil {
			return nil, false, nil, err
		}
	}

	tok = p.lexer.Peek()
	if ns && !IsName(tok, "from") {
		return nil, false, nil, p.errorTok(tok)
	}

	var src Node
	if IsName(tok, "from") {
		p.lexer.Next()
		str, err := p.nextMustTok(T_STRING)
		if err != nil {
			return nil, false, nil, err
		}
		src = &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), str.HasLegacyOctalEscapeSeq()}
	}
	return specs, ns, src, nil
}

func (p *Parser) exportNamed() ([]Node, error) {
	specs := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		spec, err := p.exportSpec()
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
	}

	_, err := p.nextMustTok(T_BRACE_R)
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func (p *Parser) exportSpec() (Node, error) {
	loc := p.loc()
	local, err := p.ident()
	if err != nil {
		return nil, err
	}

	id := local
	if p.aheadIsName("as") {
		p.lexer.Next()
		id, err = p.ident()
		if err != nil {
			return nil, err
		}
	}

	return &ExportSpec{N_EXPORT_SPEC, p.finLoc(loc), false, local, id}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ImportDeclaration
func (p *Parser) importDec() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	specs := make([]Node, 0)
	tok := p.lexer.Peek()
	if tok.value != T_STRING {
		if tok.value == T_NAME {
			id, err := p.ident()
			if err != nil {
				return nil, err
			}
			spec := &ImportSpec{N_IMPORT_SPEC, p.finLoc(p.locFromTok(tok)), true, false, id, id}
			specs = append(specs, spec)
		} else {
			ss, err := p.importNamedOrNS()
			if err != nil {
				return nil, err
			}
			specs = append(specs, ss...)
		}

		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
			ss, err := p.importNamedOrNS()
			if err != nil {
				return nil, err
			}
			specs = append(specs, ss...)
		}

		_, err := p.nextMustName("from")
		if err != nil {
			return nil, err
		}
	}

	str, err := p.nextMustTok(T_STRING)
	if err != nil {
		return nil, err
	}
	src := &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), str.HasLegacyOctalEscapeSeq()}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	return &ImportDec{N_STMT_IMPORT, p.finLoc(loc), specs, src}, nil
}

func (p *Parser) importNamedOrNS() ([]Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_BRACE_L {
		return p.importNamed()
	} else if tok.value == T_MUL {
		return p.importNS()
	} else {
		return nil, p.errorTok(tok)
	}
}

func (p *Parser) importSpec() (Node, error) {
	loc := p.loc()
	binding, err := p.ident()
	if err != nil {
		return nil, err
	}

	id := binding
	if p.aheadIsName("as") {
		p.lexer.Next()
		binding, err = p.ident()
		if err != nil {
			return nil, err
		}
	}

	return &ImportSpec{N_IMPORT_SPEC, p.finLoc(loc), false, false, binding, id}, nil
}

func (p *Parser) importNamed() ([]Node, error) {
	p.lexer.Next()

	specs := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		spec, err := p.importSpec()
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
	}

	_, err := p.nextMustTok(T_BRACE_R)
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func (p *Parser) importNS() ([]Node, error) {
	loc := p.loc()
	p.lexer.Next()
	_, err := p.nextMustName("as")
	if err != nil {
		return nil, err
	}

	id, err := p.ident()
	if err != nil {
		return nil, err
	}

	specs := make([]Node, 1)
	specs[0] = &ImportSpec{N_IMPORT_SPEC, p.finLoc(loc), false, true, id, nil}
	return specs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec(expr bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	scope := p.symtab.EnterScope(false)
	scope.AddKind(SPK_STRICT)
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

	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}

	elems := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		if tok.value == T_SEMI {
			p.lexer.Next()
		}
		elem, err := p.classElem()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}

	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}

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
		name := tok.text
		if name == "constructor" {
			return p.method(static, nil, tok, false, false)
		} else if name == "get" || name == "set" {
			ahead := p.lexer.PeekGrow()
			isField := ahead.value == T_ASSIGN || ahead.value == T_SEMI || ahead.afterLineTerminator
			if !isField {
				p.lexer.Next()
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
	kindStr := ""
	if kind != nil {
		kindStr = kind.Text()
		loc = p.locFromTok(kind)
	}

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
	return &Method{N_METHOD, p.finLoc(loc), key, static, key.Type() != N_NAME, kindStr, value}, nil
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
	p.advanceIfSemi(false)
	return &Field{N_FIELD, p.finLoc(loc), key, static, key.Type() != N_NAME, value}, nil
}

func (p *Parser) classElemName() (Node, error) {
	return p.propName(true)
}

func (p *Parser) staticBlock() (Node, error) {
	block, err := p.blockStmt(true)
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

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	expr, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	tok := p.lexer.PeekStmtBegin()
	body, err := p.stmt()
	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	return &WithStmt{N_STMT_WITH, p.finLoc(loc), expr, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DebuggerStatement
func (p *Parser) debugStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	return &DebugStmt{N_STMT_DEBUG, p.finLoc(loc)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-TryStatement
func (p *Parser) tryStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	try, err := p.blockStmt(true)
	if err != nil {
		return nil, err
	}

	tok := p.lexer.Peek()
	if tok.value != T_CATCH && tok.value != T_FINALLY {
		return nil, p.errorTok(tok)
	}

	var catch Node
	if tok.value == T_CATCH {
		loc := p.loc()
		p.lexer.Next()
		if _, err := p.nextMustTok(T_PAREN_L); err != nil {
			return nil, err
		}
		param, err := p.bindingPattern()
		if err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_PAREN_R); err != nil {
			return nil, err
		}

		body, err := p.blockStmt(true)
		if err != nil {
			return nil, err
		}
		catch = &Catch{N_CATCH, p.finLoc(loc), param, body}
	}

	var fin Node
	if p.lexer.Peek().value == T_FINALLY {
		p.lexer.Next()
		fin, err = p.blockStmt(true)
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
	if tok.value != T_ILLEGAL && tok.value != T_EOF && !tok.afterLineTerminator {
		arg, err = p.expr(false)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, p.errorTok(tok)
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
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
	} else if tok.value != T_ILLEGAL &&
		tok.value != T_BRACE_R &&
		tok.value != T_PAREN_R &&
		tok.value != T_BRACKET_R &&
		// tok.value != T_COMMENT &&
		tok.value != T_EOF && !tok.afterLineTerminator {
		arg, err = p.expr(false)
		if err != nil {
			return nil, err
		}
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_FUNC) && !scope.IsKind(SPK_FUNC_INDIRECT) {
		return nil, p.errorAtLoc(loc, "Illegal return")
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

	scope := p.scope()
	labelName := label.Text()
	if scope.HasLabel(labelName) {
		return nil, p.errorAtLoc(loc, fmt.Sprintf("Label `%s` already declared", labelName))
	}

	node := &LabelStmt{N_STMT_LABEL, nil, label, nil}
	scope.Labels[labelName] = node

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	node.loc = p.finLoc(loc)
	node.body = body
	scope.Labels = make(map[string]Node)
	return node, nil
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

		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}

		if !p.scope().HasLabel(label.Text()) {
			return nil, p.errorAtLoc(label.loc, fmt.Sprintf("Undefined label `%s`", label.Text()))
		}

		return &BrkStmt{N_STMT_BRK, p.finLoc(loc), label}, nil
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_LOOP_DIRECT) &&
		!scope.IsKind(SPK_LOOP_INDIRECT) &&
		!scope.IsKind(SPK_SWITCH) {
		return nil, p.errorAtLoc(loc, "Illegal break")
	}
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

		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}

		if !p.scope().HasLabel(label.Text()) {
			return nil, p.errorAtLoc(label.loc, fmt.Sprintf("Undefined label `%s`", label.Text()))
		}

		return &ContStmt{N_STMT_CONT, p.finLoc(loc), label}, nil
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_LOOP_DIRECT) && !scope.IsKind(SPK_LOOP_INDIRECT) {
		return nil, p.errorAtLoc(loc, "Illegal continue")
	}
	return &ContStmt{N_STMT_CONT, p.finLoc(loc), nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-SwitchStatement
func (p *Parser) switchStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	cases := make([]*SwitchCase, 0)
	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}
	metDefault := false

	scope := p.symtab.EnterScope(false)
	scope.AddKind(SPK_SWITCH)

	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value != T_CASE && tok.value != T_DEFAULT {
			return nil, p.errorTok(tok)
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		if tok.value == T_DEFAULT && metDefault {
			return nil, p.errorAt(tok.value, &tok.begin, "Multiple default clauses")
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
	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}

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
	if _, err := p.nextMustTok(T_COLON); err != nil {
		return nil, err
	}

	cons := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_CASE || tok.value == T_DEFAULT || tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
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

func (p *Parser) forbidVarDec(stmt Node) error {
	if stmt.Type() == N_STMT_VAR_DEC {
		dec := stmt.(*VarDecStmt)
		if dec.kind != T_VAR {
			return p.errorAtLoc(stmt.Loc(), "Illegal lexical declaration")
		}
	}
	return nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-IfStatement
func (p *Parser) ifStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	tok := p.lexer.PeekStmtBegin()
	cons, err := p.stmt()
	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}
	if err = p.forbidVarDec(cons); err != nil {
		return nil, err
	}

	var alt Node
	if p.lexer.Peek().value == T_ELSE {
		p.lexer.Next()
		tok := p.lexer.PeekStmtBegin()
		alt, err = p.stmt()
		if err != nil {
			if err == errEof {
				return nil, p.errorTok(tok)
			}
			return nil, err
		}
	}
	return &IfStmt{N_STMT_IF, p.finLoc(loc), test, cons, alt}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DoWhileStatement
func (p *Parser) doWhileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.PeekStmtBegin()
	body, err := p.stmt()
	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	if _, err := p.nextMustTok(T_WHILE); err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}
	return &DoWhileStmt{N_STMT_DO_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WhileStatement
func (p *Parser) whileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr(false)
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	scope := p.symtab.EnterScope(false)
	scope.AddKind(SPK_LOOP_DIRECT)

	tok := p.lexer.PeekStmtBegin()
	body, err := p.stmt()
	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
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

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}

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
	isIn := IsName(tok, "in")
	if isIn || IsName(tok, "of") {
		if init == nil {
			return nil, p.errorTok(tok)
		}

		if init.Type() != N_STMT_VAR_DEC && !p.isSimpleLVal(init) {
			return nil, p.errorAtLoc(init.Loc(), "Assigning to rvalue")
		} else if init.Type() == N_STMT_VAR_DEC {
			varDec := init.(*VarDecStmt)
			if len(varDec.decList) > 1 {
				return nil, p.errorAtLoc(varDec.decList[1].Loc(), "Must have a single binding")
			}
		}

		p.lexer.Next()
		right, err := p.expr(false)
		if err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_PAREN_R); err != nil {
			return nil, err
		}
		tok := p.lexer.PeekStmtBegin()
		body, err := p.stmt()
		if err != nil {
			if err == errEof {
				return nil, p.errorTok(tok)
			}
			return nil, err
		}
		return &ForInOfStmt{N_STMT_FOR_IN_OF, p.finLoc(loc), isIn, await, init, right, body}, nil
	}

	if _, err := p.nextMustTok(T_SEMI); err != nil {
		return nil, err
	}
	var test Node
	if p.lexer.Peek().value == T_SEMI {
		p.lexer.Next()
	} else {
		test, err = p.expr(false)
		if err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_SEMI); err != nil {
			return nil, err
		}
	}

	var update Node
	if p.lexer.Peek().value != T_PAREN_R {
		update, err = p.expr(false)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}
	tok = p.lexer.PeekStmtBegin()
	body, err := p.stmt()
	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	return &ForStmt{N_STMT_FOR, p.finLoc(loc), init, test, update, body}, nil
}

func (p *Parser) aheadIsAsync(tok *Token) bool {
	if IsName(tok, "async") {
		ahead := p.lexer.PeekGrow()
		return (ahead.value == T_FUNC || ahead.value == T_PAREN_L) && !ahead.afterLineTerminator
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
	tok := p.lexer.Peek()
	if tok.value == T_FUNC {
		p.lexer.Next()
	}

	p.symtab.EnterScope(true)
	generator := p.lexer.Peek().value == T_MUL
	if generator {
		p.scope().AddKind(SPK_GENERATOR)
		p.lexer.Next()
	}

	var id Node
	var err error
	tok = p.lexer.Peek()
	if tok.value != T_PAREN_L {
		id, err = p.ident()
		if err != nil {
			return nil, err
		}
	}
	if !expr && id == nil {
		return nil, p.errorTok(tok)
	}

	params, err := p.formalParams()
	if err != nil {
		return nil, err
	}

	if generator {
		p.lexer.extMode(LM_GENERATOR, true)
	}

	tok = p.lexer.Peek()
	arrow := false
	if tok.value == T_ARROW {
		if expr {
			p.lexer.Next()
			arrow = true
		} else {
			return nil, p.errorTok(tok)
		}
	}

	body, err := p.fnBody()
	if err != nil {
		return nil, err
	}
	if generator {
		p.lexer.popMode()
	}

	// p.symtab.LeaveScope()
	// TODO: check formal params if in strict mode

	if arrow {
		return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), async, params, body}, nil
	}

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
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			p.lexer.Next()
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}

		param, err := p.bindingElem()
		if err != nil {
			return nil, err
		}

		if p.lexer.Peek().value == T_COMMA {
			tok := p.lexer.Next()
			if param.Type() == N_PATTERN_REST {
				msg := "Unexpected trailing comma after rest element"
				if p.lexer.Peek().value != T_PAREN_R {
					msg = "Rest element must be last element"
				}
				return nil, p.errorAt(tok.value, &tok.begin, msg)
			}
		}
		params = append(params, param)
	}

	return params, nil
}

func (p *Parser) fnBody() (Node, error) {
	return p.blockStmt(false)
}

func (p *Parser) scope() *Scope {
	return p.symtab.Cur
}

func (p *Parser) isStrictDirective(exprStmt Node) bool {
	expr := exprStmt.(*ExprStmt).expr
	if expr.Type() == N_LIT_STR {
		str := expr.(*StrLit).Raw()
		if str == "\"use strict\"" || str == "'use strict'" {
			return true
		}
	}
	return false
}

func (p *Parser) blockStmt(newScope bool) (*BlockStmt, error) {
	tok, err := p.nextMustTok(T_BRACE_L)
	if err != nil {
		return nil, err
	}
	if newScope {
		p.symtab.EnterScope(false)
	}
	loc := p.locFromTok(tok)

	stmts := make([]Node, 0)
	prologue := 0
	isPrologueClosed := false

	scope := p.scope()
	fnBody := scope.IsKind(SPK_FUNC)

	for {
		tok := p.lexer.PeekStmtBegin()
		if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			if fnBody && stmt.Type() == N_STMT_EXPR {
				if p.isStrictDirective(stmt) {
					scope.AddKind(SPK_STRICT)
					p.lexer.addMode(LM_STRICT)
					isPrologueClosed = true
				}
			}
			if !isPrologueClosed {
				prologue += 1
			}
			stmts = append(stmts, stmt)
		}
	}

	if isPrologueClosed && prologue > 0 {
		for i := 0; i < prologue; i++ {
			expr := stmts[i].(*ExprStmt).expr
			if expr.Type() == N_LIT_STR && expr.(*StrLit).legacyOctalEscapeSeq {
				return nil, p.errorAtLoc(expr.Loc(), "Octal escape sequences are not allowed in strict mode")
			}
		}
	}

	if newScope {
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

	node := &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]*VarDec, 0, 1)}
	kind := p.lexer.Next()

	if IsName(kind, "let") {
		node.kind = T_LET
	} else if IsName(kind, "const") {
		node.kind = T_CONST
	} else {
		node.kind = T_VAR
	}

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
	if !asExpr {
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
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
	tok, err := p.nextMustTok(T_NAME)
	if err != nil {
		return nil, err
	}
	ident := &Ident{N_NAME, nil, "", false}
	ident.loc = p.finLoc(p.locFromTok(tok))
	ident.val = tok.Text()
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
			return nil, p.errorTok(tok)
		}
		props = append(props, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.errorTok(tok)
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
	if key.Type() == N_PROP {
		return key, nil
	}

	tok := p.lexer.Peek()
	if tok.value != T_COLON {
		if key.Type() == N_NAME {
			return p.patternAssign(key)
		}
		return nil, p.errorTok(tok)
	}

	p.lexer.Next()
	value, err := p.bindingElem()
	if err != nil {
		return nil, err
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key), ""}, nil
}

func (p *Parser) propName(allowNamePVT bool) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	name, ok := tok.CanBePropKey()

	if ok && (name == "get" || name == "set") {
		ahead := p.lexer.Peek()
		isField := ahead.value == T_COLON ||
			ahead.value == T_ASSIGN ||
			ahead.value == T_SEMI ||
			ahead.afterLineTerminator
		if !isField {
			return p.propMethod(nil, name, false, false, false)
		}
	}

	if ok || (allowNamePVT && tok.value == T_NAME_PVT) {
		if tok.value == T_NUM {
			return &NumLit{N_LIT_NUM, p.finLoc(loc)}, nil
		}
		return &Ident{N_NAME, p.finLoc(loc), tok.Text(), tok.value == T_NAME_PVT}, nil
	}
	if tok.value == T_STRING {
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), tok.HasLegacyOctalEscapeSeq()}, nil
	}
	if tok.value == T_NUM {
		return &NumLit{N_LIT_NUM, p.finLoc(loc)}, nil
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
	return nil, p.errorTok(tok)
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
			return nil, p.errorTok(tok)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.errorTok(tok)
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

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	stmt.loc = p.finLoc(loc)

	// adjust col to include the open-close backquotes
	if expr.Type() == N_EXPR_TPL {
		if stmt.loc.begin.col > 0 {
			stmt.loc.begin.col -= 1
		}
		if stmt.loc.rng.start > 0 {
			stmt.loc.rng.start -= 1
		}
		stmt.loc.end.col += 1
		stmt.loc.rng.end += 1
	}
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
				return nil, p.errorAt(p.lexer.prtVal, &p.lexer.prtBegin, "")
			}
			exprs = append(exprs, expr)
		} else {
			break
		}
	}
	return &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), exprs}, nil
}

func (p *Parser) aheadIsYield() bool {
	if !p.scope().IsKind(SPK_GENERATOR) {
		return false
	}
	return IsName(p.lexer.Peek(), "yield")
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-YieldExpression
func (p *Parser) yieldExpr(notIn bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	if tok.afterLineTerminator {
		return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), false, nil}, nil
	}

	delegate := false
	if p.lexer.Peek().value == T_MUL {
		p.lexer.Next()
		delegate = true
	}

	arg, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}
	return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), delegate, arg}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-AssignmentExpression
func (p *Parser) assignExpr(notIn bool) (Node, error) {
	loc := p.loc()
	if p.aheadIsYield() {
		return p.yieldExpr(notIn)
	}

	lhs, err := p.condExpr(notIn)
	if err != nil {
		return nil, err
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == T_ILLEGAL {
		return lhs, nil
	}

	rhs, err := p.assignExpr(notIn)
	if err != nil {
		return nil, err
	}

	if !p.isSimpleLVal(lhs) {
		return nil, p.errorAtLoc(lhs.Loc(), "Assigning to rvalue")
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finLoc(loc), assign, lhs, p.unParen(rhs)}
	return node, nil
}

// https://tc39.es/ecma262/multipage/syntax-directed-operations.html#sec-static-semantics-assignmenttargettype
func (p *Parser) isSimpleLVal(expr Node) bool {
	switch expr.Type() {
	case N_NAME:
		node := expr.(*Ident)
		if p.scope().IsKind(SPK_ASYNC) && node.Text() == "await" {
			return false
		}
		return true
	case N_PATTERN_OBJ, N_PATTERN_ARRAY, N_PATTERN_ASSIGN, N_PATTERN_REST, N_EXPR_MEMBER:
		return true
	case N_EXPR_PAREN:
		node := expr.(*ParenExpr)
		return p.isSimpleLVal(node.expr)
	}
	return false
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
	op := tok.value
	if tok.IsUnary() || tok.value == T_ADD || tok.value == T_SUB {
		p.lexer.Next()
		arg, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		scope := p.scope()
		if scope.IsKind(SPK_STRICT) && tok.value == T_DELETE && arg.Type() == N_NAME {
			return nil, p.errorAtLoc(arg.Loc(), "Deleting local variable in strict mode")
		}

		return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), op, arg}, nil
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
		if !p.isSimpleLVal(arg) {
			return nil, p.errorAtLoc(arg.Loc(), "Assigning to rvalue")
		}
		return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, true, arg}, nil
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

	if !p.isSimpleLVal(arg) {
		return nil, p.errorAtLoc(arg.Loc(), "Assigning to rvalue")
	}

	p.lexer.Next()
	return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, false, arg}, nil
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
		} else if tok.value == T_TPL_HEAD {
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
	meta := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), tok.Text(), false}

	if p.lexer.Peek().value == T_DOT {
		p.lexer.Next()
		prop, err := p.ident()
		if err != nil {
			return nil, err
		}
		if prop.Text() != "meta" {
			return nil, p.errorAt(p.lexer.prtVal, &p.lexer.prtBegin, "")
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
	loc := p.loc()
	if tag != nil {
		tl := tag.Loc()
		loc.begin = tl.end.Clone()
		loc.rng.start = tl.rng.end
	} else {
		// move back one position to take the place of the beginning backquote
		if loc.begin.col > 0 {
			loc.begin.col -= 1
		}
		if loc.rng.start > 0 {
			loc.rng.start -= 1
		}
	}

	elems := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_TPL_TAIL {
			loc := p.loc()
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), false}
			elems = append(elems, str)
			break
		} else if tok.value == T_TPL_SPAN || tok.value == T_TPL_HEAD {
			loc := p.loc()
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), false}
			elems = append(elems, str)

			if tok.IsPlainTpl() {
				if tag == nil {
					return str, nil
				}
				break
			}

			expr, err := p.expr(false)
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		} else {
			return nil, p.errorTok(tok)
		}
	}

	loc = p.finLoc(loc)
	loc.end.col += 1
	loc.rng.end += 1
	return &TplExpr{N_EXPR_TPL, loc, tag, elems}, nil
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
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		arg, err := p.arg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}
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
		op := ahead.IsBin(notIn)
		if op == T_ILLEGAL || pcd < minPcd {
			break
		}
		p.lexer.Next()

		rhs, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		kind = ahead.Kind()
		for ahead.IsBin(notIn) != T_ILLEGAL && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			rhs, err = p.binExpr(rhs, pcd, notIn)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			kind = ahead.Kind()
		}
		pcd = kind.Pcd

		bin := &BinExpr{N_EXPR_BIN, nil, T_ILLEGAL, nil, nil}
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
		prop = &Ident{N_NAME, p.finLoc(loc), tok.Text(), tok.value == T_NAME_PVT}
	} else {
		return nil, p.errorTok(tok)
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
		return &NumLit{N_LIT_NUM, p.finLoc(loc)}, nil
	case T_STRING:
		p.lexer.Next()
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), tok.HasLegacyOctalEscapeSeq()}, nil
	case T_NULL:
		p.lexer.Next()
		return &NullLit{N_LIT_NULL, p.finLoc(loc)}, nil
	case T_TRUE, T_FALSE:
		p.lexer.Next()
		return &BoolLit{N_LIT_BOOL, p.finLoc(loc), tok.Text() == "true"}, nil
	case T_NAME:
		if p.aheadIsAsync(tok) {
			return p.fnDec(true, true)
		}
		p.lexer.Next()
		return &Ident{N_NAME, p.finLoc(loc), tok.Text(), false}, nil
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
		return &RegexpLit{N_LIT_REGEXP, p.finLoc(loc), tok.Text(), ext.Pattern(), ext.Flags()}, nil
	case T_CLASS:
		return p.classDec(true)
	case T_SUPER:
		p.lexer.Next()
		return &Super{N_SUPER, p.finLoc(loc)}, nil
	case T_IMPORT:
		return p.importCall()
	case T_TPL_HEAD:
		return p.tplExpr(nil)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) parenExpr() (Node, error) {
	loc := p.loc()
	params, err := p.argList()
	if err != nil {
		return nil, err
	}
	if p.lexer.Peek().value == T_ARROW {
		p.lexer.Next()
		p.symtab.EnterScope(true)

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

		return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), false, params, body}, nil
	}
	if len(params) == 0 {
		return nil, p.errorAt(p.lexer.prtVal, &p.lexer.prtBegin, "")
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
			return nil, p.errorTok(tok)
		}
		elems = append(elems, node)

		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.errorTok(tok)
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
	hasProto := false
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

		if node.Type() == N_PROP {
			prop := node.(*Prop)
			if !prop.computed {
				var propName string

				switch prop.key.Type() {
				case N_NAME:
					propName = prop.key.(*Ident).Text()
				case N_LIT_STR:
					propName = prop.key.(*StrLit).Text()
				case N_LIT_NUM:
					propName = prop.key.(*NumLit).Text()
				case N_LIT_BOOL:
					propName = prop.key.(*BoolLit).Text()
				case N_LIT_NULL:
					propName = prop.key.(*NullLit).Text()
				}

				if propName == "__proto__" {
					if hasProto {
						return nil, p.errorAtLoc(node.Loc(), "Redefinition of property")
					}
					hasProto = true
				}
			}
		}
		props = append(props, node)

		tok = p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
		} else if tok.value == T_BRACE_R {
			p.lexer.Next()
			break
		} else {
			return nil, p.errorTok(tok)
		}
	}
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props}, nil
}

func (p *Parser) objProp() (Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_DOT_TRI {
		return p.spread()
	}

	if tok.value == T_MUL {
		return p.propMethod(nil, "", true, false, false)
	} else if p.aheadIsAsync(tok) {
		return p.propMethod(nil, "", false, true, false)
	}
	return p.propField()
}

func (p *Parser) propField() (Node, error) {
	loc := p.loc()
	key, err := p.propName(false)
	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
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
		return p.propMethod(key, "", false, false, false)
	}

	return &Prop{N_PROP, p.finLoc(loc), key, value, !IsLitPropName(key), ""}, nil
}

func (p *Parser) propMethod(key Node, kind string, gen bool, async bool, allowNamePVT bool) (Node, error) {
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

	p.symtab.EnterScope(true)
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

func (p *Parser) advanceIfSemi(raise bool) error {
	tok := p.lexer.Peek()
	if tok.value == T_SEMI {
		p.lexer.Next()
	}
	if raise && tok.value != T_SEMI && tok.value != T_BRACE_R && !tok.afterLineTerminator && tok.value != T_EOF {
		return p.errorAt(tok.value, &tok.begin, "Unexpected token")
	}
	return nil
}

func (p *Parser) nextMustTok(val TokenValue) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != val {
		return nil, p.errorTok(tok)
	}
	return tok, nil
}

func (p *Parser) nextMustName(name string) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != T_NAME || tok.Text() != name {
		return nil, p.errorTok(tok)
	}
	return tok, nil
}

func (p *Parser) aheadIsName(name string) bool {
	tok := p.lexer.Peek()
	return tok.value == T_NAME && tok.Text() == name
}

func (p *Parser) advanceIfTok(val TokenValue) *Token {
	tok := p.lexer.Peek()
	if tok.value != val {
		return nil
	}
	return p.lexer.Next()
}

func (p *Parser) advanceIfTokIn(begin, end TokenValue) TokenValue {
	tok := p.lexer.Peek()
	if tok.value <= begin || tok.value >= end {
		return T_ILLEGAL
	}
	return p.lexer.Next().value
}

func (p *Parser) loc() *Loc {
	return p.lexer.Loc()
}

func (p *Parser) locFromTok(tok *Token) *Loc {
	return &Loc{
		src:   tok.raw.src,
		begin: tok.begin.Clone(),
		end:   &Pos{},
		rng:   &Range{tok.raw.lo, 0},
	}
}

func (p *Parser) finLoc(loc *Loc) *Loc {
	return p.lexer.FinLoc(loc)
}

func (p *Parser) errorTok(tok *Token) *ParserError {
	if tok.value != T_ILLEGAL {
		return NewParserError(fmt.Sprintf("Unexpected token `%s`", TokenKinds[tok.value].Name),
			p.lexer.src.path, tok.begin.line, tok.begin.col)
	}
	return NewParserError(tok.ErrMsg(), p.lexer.src.path, tok.begin.line, tok.begin.col)
}

func (p *Parser) errorAt(tok TokenValue, pos *Pos, errMsg string) *ParserError {
	if tok != T_ILLEGAL && errMsg == "" {
		return NewParserError(fmt.Sprintf("Unexpected token `%s`", TokenKinds[tok].Name),
			p.lexer.src.path, pos.line, pos.col)
	}
	return NewParserError(errMsg, p.lexer.src.path, pos.line, pos.col)
}

func (p *Parser) errorAtLoc(loc *Loc, errMsg string) *ParserError {
	return NewParserError(errMsg, p.lexer.src.path, loc.begin.line, loc.begin.col)
}

func IsLitPropName(node Node) bool {
	typ := node.Type()
	return typ == N_NAME || typ == N_LIT_STR || typ == N_LIT_NUM
}
