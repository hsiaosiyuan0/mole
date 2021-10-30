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

	scope := p.scope()
	scope.AddKind(SPK_GLOBAL)

	stmts, err := p.stmts(T_ILLEGAL)
	if err != nil {
		return nil, err
	}

	pg.stmts = stmts
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
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		node, err = p.fnDec(false, true)
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
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
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
		} else if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			node.dec, err = p.fnDec(false, true)
		} else if tok.value == T_CLASS {
			node.dec, err = p.classDec(false)
		} else {
			node.dec, err = p.assignExpr(false, true)
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
		_, err = p.nextMustName("as", false)
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
	if ns && !IsName(tok, "from", false) {
		return nil, false, nil, p.errorTok(tok)
	}

	var src Node
	if IsName(tok, "from", false) {
		p.lexer.Next()
		str, err := p.nextMustTok(T_STRING)
		if err != nil {
			return nil, false, nil, err
		}
		src = &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), str.HasLegacyOctalEscapeSeq(), nil}
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
		tok := p.lexer.Next()
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
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

		_, err := p.nextMustName("from", false)
		if err != nil {
			return nil, err
		}
	}

	str, err := p.nextMustTok(T_STRING)
	if err != nil {
		return nil, err
	}
	legacyOctalEscapeSeq := str.HasLegacyOctalEscapeSeq()
	if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
		return nil, p.errorAtLoc(p.locFromTok(str), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
	}
	src := &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), legacyOctalEscapeSeq, nil}

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
	_, err := p.nextMustName("as", false)
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
	if !expr && id == nil {
		return nil, p.errorAtLoc(p.loc(), ERR_CLASS_NAME_REQUIRED)
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

	loc := p.loc()
	tok := p.lexer.Peek()
	if tok.value == T_NAME {
		name := tok.text
		if name == "constructor" {
			return p.method(loc, nil, false, PK_CTOR, false, false, true, true, static)
		} else if name == "get" || name == "set" {
			ahead := p.lexer.PeekGrow()
			isField := ahead.value == T_ASSIGN || ahead.value == T_SEMI || ahead.afterLineTerminator
			if !isField {
				p.lexer.Next()

				k := PK_INIT
				if name == "get" {
					k = PK_GETTER
				} else {
					k = PK_SETTER
				}
				return p.method(loc, nil, false, k, false, false, true, true, static)
			}
		} else if p.aheadIsAsync(tok, true, true) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			return p.method(nil, nil, false, PK_METHOD, false, true, true, true, static)
		}
	} else if tok.value == T_MUL {
		return p.method(nil, nil, false, PK_METHOD, true, false, true, true, static)
	}

	return p.field(static)
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
		p.lexer.Next()
		value, err = p.assignExpr(false, true)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.method(nil, key, false, PK_METHOD, false, false, true, true, static)
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

	if p.scope().IsKind(SPK_STRICT) {
		return nil, p.errorAtLoc(p.finLoc(loc), ERR_WITH_STMT_IN_STRICT)
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

		scope := p.symtab.EnterScope(false)
		scope.AddKind(SPK_CATCH)

		names, _ := p.collectNames([]Node{param})
		for _, nameNode := range names {
			ref := NewRef()
			ref.Node = nameNode.(*Ident)
			ref.BindKind = BK_LET
			if err := p.addLocalBinding(nil, ref, true); err != nil {
				return nil, err
			}
		}

		body, err := p.blockStmt(false)
		p.symtab.LeaveScope()

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
		if tok.afterLineTerminator {
			return nil, p.errorAtLoc(loc, ERR_ILLEGAL_NEWLINE_AFTER_THROW)
		}
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
		return nil, p.errorAtLoc(loc, ERR_ILLEGAL_RETURN)
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
		return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_DUP_LABEL, labelName))
	}

	node := &LabelStmt{N_STMT_LABEL, nil, label, nil}
	scope.uniqueLabels[labelName] = 1
	scope.Labels = append(scope.Labels, node)

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	node.loc = p.finLoc(loc)
	node.body = body
	// reset to check next label chain
	scope.uniqueLabels = make(map[string]int)
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
			return nil, p.errorAtLoc(label.loc, fmt.Sprintf(ERR_UNDEF_LABEL, label.Text()))
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
		return nil, p.errorAtLoc(loc, ERR_ILLEGAL_BREAK)
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
			return nil, p.errorAtLoc(label.loc, fmt.Sprintf(ERR_UNDEF_LABEL, label.Text()))
		}

		return &ContStmt{N_STMT_CONT, p.finLoc(loc), label}, nil
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_LOOP_DIRECT) && !scope.IsKind(SPK_LOOP_INDIRECT) {
		return nil, p.errorAtLoc(loc, ERR_ILLEGAL_CONTINUE)
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
			return nil, p.errorAt(tok.value, &tok.begin, ERR_MULTI_DEFAULT)
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
			return p.errorAtLoc(stmt.Loc(), ERR_ILLEGAL_LEXICAL_DEC)
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
	if IsName(p.lexer.Peek(), "await", false) {
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
	isIn := IsName(tok, "in", false)
	if isIn || IsName(tok, "of", false) {
		if init == nil {
			return nil, p.errorTok(tok)
		}

		if init.Type() != N_STMT_VAR_DEC && !p.isSimpleLVal(init) {
			return nil, p.errorAtLoc(init.Loc(), ERR_ASSIGN_TO_RVALUE)
		} else if init.Type() == N_STMT_VAR_DEC {
			varDec := init.(*VarDecStmt)
			if len(varDec.decList) > 1 {
				return nil, p.errorAtLoc(varDec.decList[1].Loc(), ERR_DUP_BINDING)
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

func (p *Parser) aheadIsAsync(tok *Token, prop bool, pvt bool) bool {
	if IsName(tok, "async", true) {
		ahead := p.lexer.PeekGrow()
		return !ahead.afterLineTerminator &&
			(ahead.value == T_FUNC || ahead.value == T_PAREN_L || ahead.value == T_MUL ||
				(prop && (ahead.value == T_BRACKET_L || ahead.value == T_NAME) || pvt && ahead.value == T_NAME_PVT))
	}
	return false
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

	parentScope := p.scope()
	scope := p.symtab.EnterScope(true)
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

		// name of the function expression will not add a ref record
		if !expr {
			ref := NewRef()
			ref.Node = id.(*Ident)
			// TODO: from es6 the function declaration is block-level scope,
			// BK_VAR => BK_LET from es6
			ref.BindKind = BK_VAR
			if err := p.addLocalBinding(parentScope, ref, true); err != nil {
				return nil, err
			}
		}
	}
	if !expr && id == nil {
		return nil, p.errorTok(tok)
	}

	// the arg check is skipped here, its correctness is guaranteed by
	// below `argsToFormalParams`
	args, err := p.argList(false)
	if err != nil {
		return nil, err
	}

	params, err := p.argsToFormalParams(args)
	if err != nil {
		return nil, err
	}

	paramNames, firstComplicated := p.collectNames(params)
	for _, paramName := range paramNames {
		ref := NewRef()
		ref.Node = paramName.(*Ident)
		ref.BindKind = BK_PARAM
		// duplicate-checking is enable in strict mode so here skip doing checking,
		// checking is delegated to below `checkParams`
		p.addLocalBinding(nil, ref, false)
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

	isStrict := scope.IsKind(SPK_STRICT)
	if id != nil && isStrict && isProhibitedName(id.(*Ident).Text()) {
		return nil, p.errorAtLoc(id.Loc(), ERR_RESERVED_WORD_IN_STRICT_MODE)
	}

	if err := p.checkParams(paramNames, firstComplicated, isStrict); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()

	if arrow {
		return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), async, params, body, nil}, nil
	}

	typ := N_STMT_FN
	if expr {
		typ = N_EXPR_FN
	}
	return &FnDec{typ, p.finLoc(loc), id, generator, async, params, body, nil}, nil
}

func (p *Parser) collectNames(nodes []Node) (names []Node, firstComplicated *Loc) {
	names = make([]Node, 0)
	for _, param := range nodes {
		if firstComplicated == nil && param.Type() != N_NAME {
			firstComplicated = param.Loc()
		}
		names = append(names, p.namesInPattern(param)...)
	}
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#sec-parameter-lists-static-semantics-early-errors
// `isSimpleParamList` should be true if function body directly contains `use strict` directive
func (p *Parser) checkParams(names []Node, firstComplicated *Loc, isStrict bool) error {
	var dupLoc *Loc
	unique := make(map[string]bool)
	for _, id := range names {
		name := id.(*Ident).Text()
		if isStrict && isProhibitedName(name) {
			return p.errorAtLoc(id.Loc(), ERR_RESERVED_WORD_IN_STRICT_MODE)
		}

		if dupLoc == nil {
			if _, ok := unique[name]; ok {
				dupLoc = id.Loc()

			} else {
				unique[name] = true
			}
		}
	}

	if isStrict && firstComplicated != nil {
		return p.errorAtLoc(firstComplicated, ERR_STRICT_DIRECTIVE_AFTER_NOT_SIMPLE)
	}

	if dupLoc != nil {
		if isStrict {
			return p.errorAtLoc(dupLoc, ERR_DUP_PARAM_NAME)
		}
		if firstComplicated != nil {
			return p.errorAtLoc(dupLoc, ERR_DUP_PARAM_NAME)
		}
	}
	return nil
}

func (p *Parser) namesInPattern(node Node) []Node {
	out := make([]Node, 0)
	if node == nil {
		return out
	}
	switch node.Type() {
	case N_PAT_ARRAY:
		elems := node.(*ArrPat).elems
		for _, node := range elems {
			names := p.namesInPattern(node)
			out = append(out, names...)
		}
	case N_PAT_ASSIGN:
		names := p.namesInPattern(node.(*AssignPat).lhs)
		out = append(out, names...)
	case N_PAT_OBJ:
		props := node.(*ObjPat).props
		for _, node := range props {
			var names []Node
			if node.Type() == N_NAME {
				names = p.namesInPattern(node)
			} else if node.Type() == N_PROP {
				val := node.(*Prop).value
				names = p.namesInPattern(val)
			} else {
				names = p.namesInPattern(node)
			}
			out = append(out, names...)
		}
	case N_PAT_REST:
		id := node.(*RestPat).arg.(*Ident)
		out = append(out, id)
	case N_NAME:
		out = append(out, node)
	}
	return out
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

func (p *Parser) stmts(terminal TokenValue) ([]Node, error) {
	stmts := make([]Node, 0)
	// the index in above `stmts` contains the last
	// stmt in Directive Prologue
	prologue := 0

	scope := p.scope()
	for {
		tok := p.lexer.PeekStmtBegin()
		if terminal != T_ILLEGAL {
			if tok.value == terminal {
				p.lexer.Next()
				break
			} else if tok.value == T_EOF {
				return nil, p.errorTok(tok)
			}
		} else if tok.value == T_EOF {
			break
		}
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			// StrictDirective processing
			if (scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_GLOBAL)) && stmt.Type() == N_STMT_EXPR {
				if p.isStrictDirective(stmt) {
					scope.AddKind(SPK_STRICT)
					p.lexer.addMode(LM_STRICT)

					if prologue > 0 {
						for i := 0; i < prologue; i++ {
							expr := stmts[i].(*ExprStmt).expr
							if expr.Type() == N_LIT_STR && expr.(*StrLit).legacyOctalEscapeSeq {
								return nil, p.errorAtLoc(expr.Loc(), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
							}
						}
					}
				}
			}
			prologue += 1
			stmts = append(stmts, stmt)
		}
	}
	return stmts, nil
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

	stmts, err := p.stmts(T_BRACE_R)
	if err != nil {
		return nil, err
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
	if tok.value == T_LET || tok.value == T_CONST {
		return true
	}
	return IsName(tok, "let", false) || IsName(tok, "const", false)
}

func (p *Parser) addLocalBinding(s *Scope, ref *Ref, checkDup bool) error {
	if s == nil {
		s = p.scope()
	}
	ok := s.AddLocal(ref, checkDup)
	if ok {
		return nil
	}
	if !ok {
		name := ref.Node.Text()
		return p.errorAtLoc(ref.Node.loc, fmt.Sprintf(ERR_ID_DUP_DEF, name))
	}
	return nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-VariableStatement
func (p *Parser) varDecStmt(notIn bool, asExpr bool) (Node, error) {
	loc := p.loc()

	node := &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]*VarDec, 0, 1)}
	kind := p.lexer.Next()

	isConst := false
	node.kind = T_VAR
	bindKind := BK_VAR
	if IsName(kind, "let", false) {
		node.kind = T_LET
		bindKind = BK_LET
	} else if IsName(kind, "const", false) {
		isConst = true
		node.kind = T_CONST
		bindKind = BK_CONST
	}

	lvs := make([]Node, 0)
	for {
		dec, err := p.varDec(notIn)
		if err != nil {
			return nil, err
		}
		lvs = append(lvs, dec.id)

		if isConst && dec.init == nil {
			return nil, p.errorAtLoc(dec.loc, ERR_CONST_DEC_INIT_REQUIRED)
		}

		node.decList = append(node.decList, dec)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		} else {
			break
		}
	}

	names, _ := p.collectNames(lvs)
	for _, nameNode := range names {
		id := nameNode.(*Ident)
		ref := NewRef()
		ref.Node = id
		ref.BindKind = bindKind
		if err := p.addLocalBinding(nil, ref, true); err != nil {
			return nil, err
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
		init, err = p.assignExpr(notIn, true)
		if err != nil {
			return nil, err
		}
	}
	return &VarDec{N_VAR_DEC, p.finLoc(loc), binding, init}, nil
}

var prohibitedNames = map[string]bool{
	"arguments":  true,
	"eval":       true,
	"yield":      true,
	"await":      true,
	"implements": true,
	"interface":  true,
	"let":        true,
	"package":    true,
	"protected":  true,
	"public":     true,
	"static":     true,
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#sec-identifiers-static-semantics-early-errors
func isProhibitedName(name string) bool {
	_, ok := prohibitedNames[name]
	return ok
}

func (p *Parser) ident() (*Ident, error) {
	tok, err := p.nextMustTok(T_NAME)
	if err != nil {
		return nil, err
	}
	ident := &Ident{N_NAME, nil, "", false, tok.ContainsEscape(), nil}
	ident.loc = p.finLoc(p.locFromTok(tok))
	ident.val = tok.Text()

	if p.scope().IsKind(SPK_STRICT) && isProhibitedName(ident.val) {
		return nil, p.errorAtLoc(ident.loc, ERR_RESERVED_WORD_IN_STRICT_MODE)
	}

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
		if node.Type() == N_PAT_REST && tok.value != T_BRACE_R {
			if tok.value == T_COMMA {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_REST_TRAILING_COMMA)
			}
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
	return &ObjPat{N_PAT_OBJ, p.finLoc(loc), props, nil}, nil
}

func (p *Parser) patternProp() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		binding, err := p.patternRest()
		if err != nil {
			return nil, err
		}
		return binding, nil
	}

	key, err := p.propName(false)
	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
	}

	loc := key.Loc().Clone()
	tok := p.lexer.Peek()
	opLoc := p.locFromTok(tok)
	assign := tok.value == T_ASSIGN
	var value Node
	if tok.value == T_COLON {
		p.lexer.Next()
		value, err = p.bindingElem(true)
		if err != nil {
			return nil, err
		}
	} else if assign {
		if key.Type() == N_NAME {
			return p.patternAssign(key, true)
		}
		return nil, p.errorTok(tok)
	}

	shorthand := false
	if value == nil {
		value = key
		shorthand = true
	}

	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, !IsLitPropName(key), shorthand, assign, PK_INIT}, nil
}

func (p *Parser) propName(allowNamePVT bool) (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	name, ok := tok.CanBePropKey()

	if ok && (name == "get" || name == "set") {
		ahead := p.lexer.Peek()
		isField := ahead.value == T_COLON ||
			ahead.value == T_ASSIGN ||
			ahead.value == T_SEMI ||
			ahead.afterLineTerminator
		if !isField {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}

			k := PK_INIT
			if name == "get" {
				k = PK_GETTER
			} else {
				k = PK_SETTER
			}
			return p.method(nil, nil, false, k, false, false, false, false, false)
		}
	}

	if ok || (allowNamePVT && tok.value == T_NAME_PVT) {
		if tok.value == T_NUM {
			return &NumLit{N_LIT_NUM, p.finLoc(loc), nil}, nil
		}
		return &Ident{N_NAME, p.finLoc(loc), tok.Text(), tok.value == T_NAME_PVT, tok.ContainsEscape(), nil}, nil
	}
	if tok.value == T_STRING {
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, p.errorAtLoc(p.locFromTok(tok), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil}, nil
	}
	if tok.value == T_NUM {
		return &NumLit{N_LIT_NUM, p.finLoc(loc), nil}, nil
	}
	if tok.value == T_BRACKET_L {
		name, err := p.assignExpr(false, true)
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

		node, err := p.bindingElem(false)
		if err != nil {
			return nil, err
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PAT_REST && tok.value != T_BRACKET_R {
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
	return &ArrPat{N_PAT_ARRAY, p.finLoc(loc), elems, nil}, nil
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

func (p *Parser) bindingElem(asProp bool) (Node, error) {
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
	return p.patternAssign(binding, asProp)
}

func (p *Parser) patternAssign(ident Node, asProp bool) (Node, error) {
	var init Node
	var err error
	var opLoc *Loc
	if p.lexer.Peek().value == T_ASSIGN {
		tok := p.lexer.Next()
		opLoc = p.locFromTok(tok)
		init, err = p.assignExpr(false, true)
		if err != nil {
			return nil, err
		}
	}

	if init == nil {
		return ident, nil
	}

	loc := ident.Loc()
	val := &AssignPat{N_PAT_ASSIGN, p.finLoc(loc.Clone()), ident, init, nil}
	if !asProp {
		return val, nil
	}
	return &Prop{N_PROP, p.finLoc(loc.Clone()), val.lhs, opLoc, val, false, true, true, PK_INIT}, nil
}

func (p *Parser) patternRest() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	arg, err := p.ident()
	if err != nil {
		return nil, err
	}

	tok := p.lexer.Peek()
	if tok.value == T_COMMA {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_REST_TRAILING_COMMA)
	}

	return &RestPat{N_PAT_REST, p.finLoc(loc), arg, nil}, nil
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
	expr, err := p.assignExpr(notIn, true)
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
			expr, err = p.assignExpr(notIn, true)
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
	return &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), exprs, nil}, nil
}

func (p *Parser) aheadIsYield() bool {
	if !p.scope().IsKind(SPK_GENERATOR) {
		return false
	}
	return IsName(p.lexer.Peek(), "yield", false)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-YieldExpression
func (p *Parser) yieldExpr(notIn bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	if tok.afterLineTerminator {
		return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), false, nil, nil}, nil
	}

	delegate := false
	if p.lexer.Peek().value == T_MUL {
		p.lexer.Next()
		delegate = true
	}

	arg, err := p.assignExpr(notIn, true)
	if err != nil {
		return nil, err
	}
	return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), delegate, arg, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-AssignmentExpression
func (p *Parser) assignExpr(notIn bool, checkLhs bool) (Node, error) {
	loc := p.loc()
	if p.aheadIsYield() {
		return p.yieldExpr(notIn)
	}

	lhs, err := p.condExpr(notIn)
	if err != nil {
		return nil, err
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == nil {
		return lhs, nil
	}
	op := assign.value
	opLoc := p.locFromTok(assign)

	rhs, err := p.assignExpr(notIn, checkLhs)
	if err != nil {
		return nil, err
	}

	// set `depth` to 1 to permit expr like `i + 2 = 42`
	// and so just do the arg to param transform silently
	lhs, err = p.argToParam(lhs, 1, false, true)
	if err != nil {
		return nil, err
	}

	if checkLhs && !p.isSimpleLVal(lhs) {
		return nil, p.errorAtLoc(lhs.Loc(), ERR_ASSIGN_TO_RVALUE)
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finLoc(loc), op, opLoc, lhs, p.unParen(rhs), nil}
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
	case N_PAT_OBJ, N_PAT_ARRAY, N_PAT_ASSIGN, N_PAT_REST, N_EXPR_MEMBER:
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

	cons, err := p.assignExpr(notIn, true)
	if err != nil {
		return nil, err
	}

	_, err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr(notIn, true)
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_COND, p.finLoc(loc), test, cons, alt, nil}
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
			return nil, p.errorAtLoc(arg.Loc(), ERR_DELETE_LOCAL_IN_STRICT)
		}

		return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), op, arg, nil}, nil
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
			return nil, p.errorAtLoc(arg.Loc(), ERR_ASSIGN_TO_RVALUE)
		}
		return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, true, arg, nil}, nil
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
		return nil, p.errorAtLoc(arg.Loc(), ERR_ASSIGN_TO_RVALUE)
	}

	p.lexer.Next()
	return &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, false, arg, nil}, nil
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
		args, err = p.argList(true)
		if err != nil {
			return nil, err
		}
	}

	var ret Node
	ret = &NewExpr{N_EXPR_NEW, p.finLoc(loc), expr, args, nil}
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
			args, err := p.argList(true)
			if err != nil {
				return nil, err
			}
			callee = &CallExpr{N_EXPR_CALL, p.finLoc(loc), p.unParen(callee), args, nil}
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
	meta := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), tok.Text(), false, tok.ContainsEscape(), nil}

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
	src, err := p.assignExpr(false, true)
	if err != nil {
		return nil, err
	}
	_, err = p.nextMustTok(T_PAREN_R)
	if err != nil {
		return nil, err
	}
	return &ImportCall{N_IMPORT_CALL, p.finLoc(loc), src, nil}, nil
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
		if tok.value == T_TPL_TAIL || tok.value == T_TPL_SPAN || tok.value == T_TPL_HEAD {
			cooked := ""
			if ext := tok.ext.(*TokExtTplSpan); ext != nil {
				if ext.IllegalEscape != nil {
					// raise error for bad escape sequence if the template is not tagged
					if tag == nil {
						return nil, p.errorAt(tok.value, &tok.begin, ext.IllegalEscape.Err)
					}
				} else {
					cooked = tok.Text()
				}
			}

			loc := p.loc()
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, p.finLoc(loc), cooked, false, nil}
			elems = append(elems, str)

			if tok.value == T_TPL_TAIL || tok.IsPlainTpl() {
				break
			}

			expr, err := p.expr(false)
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		} else {
			return nil, p.errorTok(tok)
		}
	}

	loc = p.finLoc(loc)
	loc.end.col += 1
	loc.rng.end += 1
	return &TplExpr{N_EXPR_TPL, loc, tag, elems, nil}, nil
}

func (p *Parser) argsToFormalParams(args []Node) ([]Node, error) {
	params := make([]Node, len(args))
	var err error
	for i, arg := range args {
		params[i], err = p.argToParam(arg, 0, false, false)
		if err != nil {
			return nil, err
		}
	}
	return params, nil
}

func (p *Parser) argToParam(arg Node, depth int, prop bool, destruct bool) (Node, error) {
	switch arg.Type() {
	case N_LIT_ARR:
		arr := arg.(*ArrLit)
		pat := &ArrPat{
			typ:   N_PAT_ARRAY,
			loc:   arr.loc,
			elems: make([]Node, len(arr.elems)),
		}
		var err error
		for i, node := range arr.elems {
			pat.elems[i], err = p.argToParam(node, depth+1, false, destruct)
			if err != nil {
				return nil, err
			}
		}
		return pat, nil
	case N_LIT_OBJ:
		n := arg.(*ObjLit)
		pat := &ObjPat{
			typ:   N_PAT_OBJ,
			loc:   n.loc,
			props: make([]Node, len(n.props)),
		}
		isProp := true
		if depth > 0 {
			isProp = prop
		}
		for i, prop := range n.props {
			pp, err := p.argToParam(prop, depth+1, isProp, destruct)
			if err != nil {
				return nil, err
			}
			pat.props[i] = pp
		}
		return pat, nil
	case N_PROP:
		n := arg.(*Prop)
		var err error
		if n.value != nil {
			if n.value.Type() == N_EXPR_FN && depth == 1 {
				return nil, p.errorAtLoc(arg.Loc(), ERR_OBJ_PATTERN_CANNOT_FN)
			}

			n.value, err = p.argToParam(n.value, depth+1, prop, destruct)
			if err != nil {
				return nil, err
			}
		}
		if n.assign {
			if !prop && depth > 0 {
				return nil, p.errorAtLoc(n.opLoc, ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT)
			}

			n.value = &AssignPat{
				typ: N_PAT_ASSIGN,
				loc: n.loc,
				lhs: n.key,
				rhs: n.value,
			}
		}
		return n, nil
	case N_EXPR_ASSIGN:
		n := arg.(*AssignExpr)
		if n.op != T_ASSIGN {
			return nil, p.errorAtLoc(n.opLoc, ERR_UNEXPECTED_TOKEN)
		}

		lhs, err := p.argToParam(n.lhs, depth+1, false, destruct)
		if err != nil {
			return nil, err
		}

		err = p.checkArg(n.rhs, false)
		if err != nil {
			return nil, err
		}
		p := &AssignPat{
			typ: N_PAT_ASSIGN,
			loc: n.loc,
			lhs: lhs,
			rhs: n.rhs,
		}
		return p, nil
	case N_NAME, N_PAT_REST:
		return arg, nil
	case N_SPREAD:
		n := arg.(*Spread)
		if n.trailingCommaLoc != nil {
			return nil, p.errorAtLoc(n.trailingCommaLoc, ERR_REST_TRAILING_COMMA)
		}

		if n.arg.Type() == N_NAME {
			// `({...(obj)} = foo)` raises error`Parenthesized pattern` in acorn
			// however it's legal in babel-parser, chrome and firefox
			//
			// use `destruct` to require the caller to indicate the parsing state
			// is in destructing or not
			if !destruct {
				if extra, ok := n.arg.Extra().(*ExprExtra); ok {
					if extra != nil && extra.OuterParen != nil {
						return nil, p.errorAtLoc(extra.OuterParen, ERR_INVALID_PAREN_ASSIGN_PATTERN)
					}
				}
			}
		} else {
			if n.arg.Type() == N_EXPR_ASSIGN {
				return nil, p.errorAtLoc(n.arg.Loc(), ERR_REST_CANNOT_SET_DEFAULT)
			}
			return nil, p.errorAtLoc(n.arg.Loc(), ERR_REST_ARG_NOT_SIMPLE)
		}
		return &RestPat{
			typ: N_PAT_REST,
			loc: n.loc,
			arg: n.arg,
		}, nil
	}
	if depth == 0 {
		return nil, p.errorAtLoc(arg.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	return arg, nil
}

func (p *Parser) checkArg(arg Node, spread bool) error {
	switch arg.Type() {
	case N_LIT_OBJ:
		n := arg.(*ObjLit)
		for _, prop := range n.props {
			err := p.checkArg(prop, true)
			if err != nil {
				return err
			}
		}
	case N_PROP:
		n := arg.(*Prop)
		if n.assign {
			return p.errorAtLoc(n.opLoc, ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT)
		}
		var err error
		if n.value != nil {
			err = p.checkArg(n.value, true)
			if err != nil {
				return err
			}
		}
	case N_PAT_REST, N_SPREAD:
		if !spread {
			return p.errorAtLoc(arg.Loc(), ERR_UNEXPECTED_TOKEN)
		}
	}
	return nil
}

func (p *Parser) checkArgs(args []Node, spread bool) error {
	for _, arg := range args {
		err := p.checkArg(arg, spread)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) argList(check bool) ([]Node, error) {
	p.lexer.Next()
	args := make([]Node, 0)

	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		arg, err := p.arg()
		if err != nil {
			return nil, err
		}

		if p.lexer.Peek().value == T_COMMA {
			tok := p.lexer.Next()
			if arg.Type() == N_SPREAD {
				msg := ERR_REST_TRAILING_COMMA
				if p.lexer.Peek().value != T_PAREN_R {
					msg = ERR_REST_ELEM_MUST_LAST
				}
				return nil, p.errorAt(tok.value, &tok.begin, msg)
			}
		}

		if check {
			// `spread` or `pattern_rest` expression is legal argument:
			// `f(c, b, ...a)`
			if err := p.checkArg(arg, true); err != nil {
				return nil, err
			}
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
	return p.assignExpr(false, false)
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

		bin := &BinExpr{N_EXPR_BIN, nil, T_ILLEGAL, nil, nil, nil}
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
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), p.unParen(obj), prop, true, false, nil}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node) (Node, error) {
	p.lexer.Next()

	loc := p.loc()
	tok := p.lexer.Next()
	_, ok := tok.CanBePropKey()

	var prop Node
	if (ok && tok.value != T_NUM) || tok.value == T_NAME_PVT {
		prop = &Ident{N_NAME, p.finLoc(loc), tok.Text(), tok.value == T_NAME_PVT, tok.ContainsEscape(), nil}
	} else {
		return nil, p.errorTok(tok)
	}

	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), p.unParen(obj), prop, false, false, nil}
	return node, nil
}

func (p *Parser) primaryExpr() (Node, error) {
	tok := p.lexer.Peek()
	loc := p.locFromTok(tok)

	switch tok.value {
	case T_NUM:
		p.lexer.Next()
		return &NumLit{N_LIT_NUM, p.finLoc(loc), nil}, nil
	case T_STRING:
		p.lexer.Next()
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, p.errorAtLoc(p.finLoc(loc), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), legacyOctalEscapeSeq, nil}, nil
	case T_NULL:
		p.lexer.Next()
		return &NullLit{N_LIT_NULL, p.finLoc(loc), nil}, nil
	case T_TRUE, T_FALSE:
		p.lexer.Next()
		return &BoolLit{N_LIT_BOOL, p.finLoc(loc), tok.Text() == "true", nil}, nil
	case T_NAME:
		if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			return p.fnDec(true, true)
		}
		p.lexer.Next()

		name := tok.Text()
		if p.scope().IsKind(SPK_STRICT) && isProhibitedName(name) {
			return nil, p.errorAtLoc(p.finLoc(loc), ERR_RESERVED_WORD_IN_STRICT_MODE)
		}
		return &Ident{N_NAME, p.finLoc(loc), name, false, tok.ContainsEscape(), nil}, nil
	case T_THIS:
		p.lexer.Next()
		return &ThisExpr{N_EXPR_THIS, p.finLoc(loc), nil}, nil
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
		return &RegexpLit{N_LIT_REGEXP, p.finLoc(loc), tok.Text(), ext.Pattern(), ext.Flags(), nil}, nil
	case T_CLASS:
		return p.classDec(true)
	case T_SUPER:
		p.lexer.Next()
		return &Super{N_SUPER, p.finLoc(loc), nil}, nil
	case T_IMPORT:
		return p.importCall()
	case T_TPL_HEAD:
		return p.tplExpr(nil)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) parenExpr() (Node, error) {
	loc := p.loc()
	args, err := p.argList(false)
	if err != nil {
		return nil, err
	}

	// next is arrow-expression
	if p.lexer.Peek().value == T_ARROW {
		params, err := p.argsToFormalParams(args)
		if err != nil {
			return nil, err
		}

		p.lexer.Next()
		scope := p.symtab.EnterScope(true)

		paramNames, firstComplicated := p.collectNames(params)
		for _, paramName := range paramNames {
			ref := NewRef()
			ref.Node = paramName.(*Ident)
			ref.BindKind = BK_PARAM
			// duplicate-checking is enable in strict mode so here skip doing checking,
			// checking is delegated to below `checkParams`
			p.addLocalBinding(nil, ref, false)
		}

		var body Node
		if p.lexer.Peek().value == T_BRACE_L {
			body, err = p.fnBody()
		} else {
			body, err = p.expr(false)
		}
		if err != nil {
			return nil, err
		}

		isStrict := scope.IsKind(SPK_STRICT)
		if err := p.checkParams(paramNames, firstComplicated, isStrict); err != nil {
			return nil, err
		}
		p.symtab.LeaveScope()

		return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), false, params, body, nil}, nil
	}

	argsLen := len(args)
	if argsLen == 0 {
		return nil, p.errorAt(p.lexer.prtVal, &p.lexer.prtBegin, "")
	}

	if err := p.checkArgs(args, false); err != nil {
		return nil, err
	}

	if argsLen == 1 {
		return &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), args[0], nil}, nil
	}

	return &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), args, nil}, nil
}

func (p *Parser) unParen(expr Node) Node {
	loc := expr.Loc().Clone()
	if expr.Type() == N_EXPR_PAREN {
		sub := expr.(*ParenExpr).Expr()
		if extra, ok := sub.Extra().(*ExprExtra); ok {
			if extra == nil {
				extra = &ExprExtra{}
				sub.setExtra(extra)
			}
			extra.OuterParen = loc
		}
		return sub
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
		if node.Type() == N_PAT_REST && tok.value != T_BRACKET_R {
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
	return &ArrLit{N_LIT_ARR, p.finLoc(loc), elems, nil}, nil
}

func (p *Parser) arrElem() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(false, true)
}

func (p *Parser) spread() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	node, err := p.assignExpr(false, true)
	if err != nil {
		return nil, err
	}

	// trailing comma is legal in after spread-expr but it's syntax
	// error after rest-expr
	//
	// when paring expr like `({...obj1,} = foo)` the part `{...obj1,}`
	// is parsed as spread-expr firstly then applied a arg-to-param
	// transform to become obj-pattern, that behavior is caused by the
	// left-most paren, it leads the state of parser to fulfill the rule
	// of paren-expr and then the inner `{...obj1,}` is parsed as obj-expr
	//
	// keep the loc of tailing comma for reporting the `tailing comma after rest-expr`
	// err in the arg-to-param transform
	var trailingCommaLoc *Loc
	tok := p.lexer.Peek()
	if tok.value == T_COMMA {
		trailingCommaLoc = p.locFromTok(tok)
	}
	return &Spread{N_SPREAD, p.finLoc(loc), p.unParen(node), trailingCommaLoc, nil}, nil
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
						return nil, p.errorAtLoc(node.Loc(), ERR_REDEF_PROP)
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
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props, nil}, nil
}

func (p *Parser) objProp() (Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_DOT_TRI {
		return p.spread()
	}

	if tok.value == T_MUL {
		return p.method(nil, nil, false, PK_INIT, true, false, false, false, false)
	} else if p.aheadIsAsync(tok, true, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(nil, nil, false, PK_INIT, false, true, false, false, false)
	}
	return p.propData()
}

func (p *Parser) propData() (Node, error) {
	key, err := p.propName(false)

	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
	}

	loc := key.Loc().Clone()

	var value Node
	tok := p.lexer.Peek()
	opLoc := p.locFromTok(tok)
	assign := tok.value == T_ASSIGN
	if tok.value == T_COLON || assign {
		p.lexer.Next()
		value, err = p.assignExpr(false, true)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.method(nil, key, false, PK_INIT, false, false, false, false, false)
	}

	shorthand := assign
	if value == nil && key.Type() == N_NAME {
		shorthand = true
		value = key
	}
	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, !IsLitPropName(key), shorthand, assign, PK_INIT}, nil
}

func (p *Parser) method(loc *Loc, key Node, shorthand bool, kind PropKind, gen bool, async bool, allowNamePVT bool,
	inClass bool, static bool) (Node, error) {
	if loc == nil {
		loc = p.loc()
	}

	// depart `gen` and `async` here since below stmt is legal:
	// `class a{ async *a() {} }`
	if async {
		p.lexer.Next()
		gen = p.lexer.Peek().value == T_MUL
	}
	if gen {
		p.lexer.Next()
	}

	var err error
	if key == nil {
		if inClass {
			key, err = p.classElemName()
		} else {
			key, err = p.propName(allowNamePVT)
		}
		if err != nil {
			return nil, err
		}
	}

	fnLoc := p.loc()
	args, err := p.argList(false)
	if err != nil {
		return nil, err
	}

	params, err := p.argsToFormalParams(args)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if kind == PK_GETTER && len(params) > 0 {
		return nil, p.errorAtLoc(params[0].Loc(), ERR_GETTER_SHOULD_NO_PARAM)
	} else if kind == PK_SETTER && len(params) != 1 {
		return nil, p.errorAtLoc(fnLoc, ERR_SETTER_SHOULD_ONE_PARAM)
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

	value := &FnDec{N_EXPR_FN, p.finLoc(fnLoc), nil, gen, async, params, body, nil}
	if inClass {
		return &Method{N_METHOD, p.finLoc(loc), key, static, key.Type() != N_NAME, kind, value}, nil
	}
	return &Prop{N_PROP, p.finLoc(loc), key, nil, value, !IsLitPropName(key), shorthand, false, kind}, nil
}

func (p *Parser) advanceIfSemi(raise bool) error {
	tok := p.lexer.Peek()
	if tok.value == T_SEMI {
		p.lexer.Next()
	}
	if raise && tok.value != T_SEMI && tok.value != T_BRACE_R && !tok.afterLineTerminator && tok.value != T_EOF {
		errMsg := ERR_UNEXPECTED_TOKEN
		if tok.value == T_ILLEGAL {
			if msg, ok := tok.ext.(string); ok {
				errMsg = msg
			} else if msg, ok := tok.ext.(*LexerError); ok {
				errMsg = msg.Error()
			}
		}
		return p.errorAt(tok.value, &tok.begin, errMsg)
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

func (p *Parser) nextMustName(name string, canContainsEscape bool) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != T_NAME || tok.Text() != name {
		return nil, p.errorTok(tok)
	}
	if !canContainsEscape && tok.ContainsEscape() {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
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
		return NewParserError(fmt.Sprintf(ERR_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok.value].Name),
			p.lexer.src.path, tok.begin.line, tok.begin.col)
	}
	return NewParserError(tok.ErrMsg(), p.lexer.src.path, tok.begin.line, tok.begin.col)
}

func (p *Parser) errorAt(tok TokenValue, pos *Pos, errMsg string) *ParserError {
	if tok != T_ILLEGAL && errMsg == "" {
		return NewParserError(fmt.Sprintf(ERR_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok].Name),
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
