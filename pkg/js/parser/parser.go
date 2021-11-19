package parser

import (
	"errors"
	"fmt"
)

type Parser struct {
	lexer     *Lexer
	symtab    *SymTab
	ver       ESVersion
	srcTyp    SourceType
	feat      Feature
	imp       map[string]*Ident
	exp       []*ExportDec
	checkName bool
}

type ParserOpts struct {
	Externals  []string
	Version    ESVersion
	SourceType SourceType
	Feature    Feature
}

const defaultFeatures Feature = FEAT_GLOBAL_ASYNC | FEAT_STRICT | FEAT_LET_CONST |
	FEAT_BINDING_PATTERN | FEAT_BINDING_REST_ELEM | FEAT_BINDING_REST_ELEM_NESTED |
	FEAT_SPREAD | FEAT_MODULE | FEAT_META_PROPERTY | FEAT_ASYNC_AWAIT | FEAT_ASYNC_ITERATION | FEAT_ASYNC_GENERATOR |
	FEAT_POW

func NewParserOpts() *ParserOpts {
	return &ParserOpts{
		Externals:  make([]string, 0),
		Version:    ES12,
		SourceType: ST_MODULE,
		Feature:    defaultFeatures,
	}
}

func NewParser(src *Source, opts *ParserOpts) *Parser {
	parser := &Parser{}
	parser.ver = opts.Version
	parser.srcTyp = opts.SourceType
	parser.feat = opts.Feature
	parser.imp = map[string]*Ident{}
	parser.exp = []*ExportDec{}
	parser.checkName = true

	parser.symtab = NewSymTab(opts.Externals)

	parser.lexer = NewLexer(src)
	parser.lexer.ver = opts.Version
	parser.lexer.feature = opts.Feature
	return parser
}

func (p *Parser) Prog() (Node, error) {
	loc := p.loc()
	pg := NewProg()

	scope := p.scope()
	scope.AddKind(SPK_GLOBAL)
	if p.feat&FEAT_GLOBAL_ASYNC != 0 {
		scope.AddKind(SPK_ASYNC)
	}
	if p.feat&FEAT_STRICT != 0 {
		p.enterStrict(true)
	}

	stmts, err := p.stmts(T_ILLEGAL)
	if err != nil {
		return nil, err
	}

	pg.stmts = stmts
	pg.loc = p.finLoc(loc)
	pg.loc.end.line = p.lexer.src.line
	pg.loc.end.col = p.lexer.src.col
	pg.loc.rng.end = p.lexer.src.Ofst()

	if err := p.checkExp(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Parser) namesInNode(node Node) []Node {
	switch node.Type() {
	case N_STMT_VAR_DEC:
		n := node.(*VarDecStmt)
		return n.names
	case N_STMT_FN, N_EXPR_FN:
		n := node.(*FnDec)
		if n.id != nil {
			return []Node{n.id}
		}
	case N_STMT_EXPR:
		return p.namesInNode(node.(*ExprStmt).expr)
	case N_STMT_CLASS, N_EXPR_CLASS:
		n := node.(*ClassDec)
		if n.id != nil {
			return []Node{n.id}
		}
	}
	return []Node{}
}

// check the exports
func (p *Parser) checkExp() error {
	names := map[string]bool{}
	// check duplication
	for _, exp := range p.exp {
		var subnames []Node
		if exp.def != nil {
			subnames = []Node{&Ident{N_NAME, exp.def, "default", false, false, nil, true}}
		} else if exp.dec != nil {
			subnames = p.namesInNode(exp.dec)
		} else {
			subnames = make([]Node, 0, len(exp.specs))
			for _, spec := range exp.specs {
				s := spec.(*ExportSpec)
				if s.id != nil {
					subnames = append(subnames, s.id)
				}
			}
		}
		for _, sn := range subnames {
			id := sn.(*Ident)
			name := id.Text()
			if _, ok := names[name]; ok {
				return p.errorAtLoc(id.Loc(), fmt.Sprintf(ERR_DUP_EXPORT, name))
			} else {
				names[name] = true
			}
		}
	}

	// also check definition
	// here separate the definition cheking into two checks since
	// their errors needed to be reported independently - firstly report
	// the duplication then the definition
	for _, exp := range p.exp {
		if exp.src != nil {
			continue
		}
		for _, spec := range exp.specs {
			id := spec.(*ExportSpec).local.(*Ident)
			name := id.Text()
			if !p.scope().HasName(name) {
				return p.errorAtLoc(id.loc, fmt.Sprintf(ERR_EXPORT_NOT_DEFINED, name))
			}
		}
	}
	return nil
}

var errEof = errors.New("eof")

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-Statement
func (p *Parser) stmt() (node Node, err error) {
	tok := p.lexer.PeekStmtBegin()

	scope := p.scope()
	allowDec := !scope.IsKind(SPK_INTERIM)

	if tok.value > T_KEYWORD_BEGIN && tok.value < T_KEYWORD_END {
		switch tok.value {
		case T_VAR:
			node, err = p.varDecStmt(T_VAR, false)
		case T_FUNC:
			node, err = p.fnDec(false, nil, false)
		case T_IF:
			node, err = p.ifStmt()
		case T_FOR:
			node, err = p.forStmt()
		case T_RETURN:
			node, err = p.retStmt()
		case T_WHILE:
			node, err = p.whileStmt()
		case T_CLASS:
			if allowDec {
				node, err = p.classDec(false, false)
			}
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
	} else if ok, kind := p.aheadIsVarDec(tok); ok {
		if allowDec {
			node, err = p.varDecStmt(kind, false)
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		node, err = p.fnDec(false, tok, false)
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
	} else if node == nil {
		return nil, p.errorTok(tok)
	}

	typ := node.Type()
	if scope.IsKind(SPK_INTERIM) {
		// `if (morning) function a(){}` is legal
		// `for (morning;;) function a(){}` is illegal
		if typ == N_STMT_FN && (scope.IsKind(SPK_STRICT) || scope.IsKind(SPK_LOOP_DIRECT)) {
			return nil, p.errorAtLoc(node.Loc(), ERR_FN_IN_SINGLE_STMT_CTX)
		} else if typ == N_STMT_IMPORT || typ == N_STMT_EXPORT {
			return nil, p.errorAtLoc(node.Loc(), ERR_IMPORT_EXPORT_SHOULD_AT_TOP_LEVEL)
		}
	}

	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ExportDeclaration
func (p *Parser) exportDec() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	if p.feat&FEAT_MODULE == 0 {
		return nil, p.errorTok(tok)
	}

	var err error
	node := &ExportDec{N_STMT_EXPORT, nil, false, nil, nil, nil, nil}
	specs := make([]Node, 0)
	tok = p.lexer.Peek()
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
	} else if ok, kind := p.aheadIsVarDec(tok); ok {
		node.dec, err = p.varDecStmt(kind, false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_FUNC {
		node.dec, err = p.fnDec(false, nil, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		node.dec, err = p.fnDec(false, tok, false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_CLASS {
		node.dec, err = p.classDec(false, false)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_DEFAULT {
		def := p.lexer.Next()
		tok := p.lexer.Peek()
		node.def = p.locFromTok(def)
		if tok.value == T_FUNC {
			node.dec, err = p.fnDec(false, nil, true)
		} else if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			node.dec, err = p.fnDec(false, tok, true)
		} else if tok.value == T_CLASS {
			node.dec, err = p.classDec(false, true)
		} else {
			node.dec, err = p.assignExpr(true)
			if err := p.advanceIfSemi(false); err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, p.errorTok(tok)
	}

	node.loc = p.finLoc(loc)
	node.specs = specs
	p.exp = append(p.exp, node)
	return node, nil
}

func (p *Parser) exportFrom() ([]Node, bool, Node, error) {
	tok := p.lexer.Next()
	var specs []Node
	var err error

	ns := false
	if tok.value == T_MUL {
		ns = true
		ahead := p.lexer.Peek()
		if ahead.value == T_NAME && ahead.Text() == "as" {
			p.lexer.Next()

			id, err := p.ident(nil)
			if err != nil {
				return nil, false, nil, err
			}
			specs = make([]Node, 1)
			specs[0] = &ExportSpec{N_EXPORT_SPEC, p.finLoc(p.locFromTok(tok)), true, id, nil}
		}
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
	} else {
		// `export { default } from "a"` is legal
		// `export { default }` is illegal
		for _, spec := range specs {
			id := spec.(*ExportSpec).local.(*Ident)
			if id.kw {
				return nil, false, nil, p.errorAtLoc(id.loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.Text()))
			}
		}
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

func (p *Parser) identWithKw(scope *Scope) (Node, error) {
	ahead := p.lexer.Peek()
	if ahead.IsKw() {
		p.lexer.Next()
		str := TokenKinds[ahead.value].Name
		return &Ident{N_NAME, p.finLoc(p.locFromTok(ahead)), str, false, false, nil, true}, nil
	}
	return p.ident(scope)
}

func (p *Parser) exportSpec() (Node, error) {
	loc := p.loc()
	local, err := p.identWithKw(nil)
	if err != nil {
		return nil, err
	}

	id := local
	if p.aheadIsName("as") {
		tok := p.lexer.Next()
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		id, err = p.identWithKw(nil)
		if err != nil {
			return nil, err
		}
	}

	return &ExportSpec{N_EXPORT_SPEC, p.finLoc(loc), false, local, id}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ImportDeclaration
func (p *Parser) importDec() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	if p.feat&FEAT_MODULE == 0 {
		return nil, p.errorTok(tok)
	}

	specs := make([]Node, 0)
	tok = p.lexer.Peek()
	if tok.value != T_STRING {
		if tok.value == T_NAME {
			id, err := p.ident(nil)
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
	binding, err := p.identWithKw(nil)
	if err != nil {
		return nil, err
	}

	id := binding
	if p.aheadIsName("as") {
		p.lexer.Next()
		binding, err = p.ident(nil)
		if err != nil {
			return nil, err
		}
	} else if binding.Type() == N_NAME {
		// for statemtnt `import { true } from "bar"`, report `true` is a keyword
		id := binding.(*Ident)
		if id.kw {
			return nil, p.errorAtLoc(binding.Loc(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, binding.(*Ident).Text()))
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

	id, err := p.ident(nil)
	if err != nil {
		return nil, err
	}

	specs := make([]Node, 1)
	specs[0] = &ImportSpec{N_IMPORT_SPEC, p.finLoc(loc), false, true, id, nil}
	return specs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec(expr bool, canNameOmitted bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	ps := p.scope()
	// all parts of the class dec are in strict mode(include the id part)
	// here push an intermidate mode as strict to handle the id part
	p.lexer.pushMode(LM_STRICT, true)
	scope := p.symtab.EnterScope(false, false)
	p.enterStrict(true).AddKind(SPK_CLASS)

	var id Node
	var err error
	ahead := p.lexer.Peek()
	if ahead.value != T_BRACE_L && ahead.value != T_EXTENDS {
		id, err = p.ident(nil)
		if err != nil {
			return nil, err
		}
		ref := NewRef()
		ref.Node = id.(*Ident)
		ref.BindKind = BK_CONST
		if err := p.addLocalBinding(ps, ref, true); err != nil {
			return nil, err
		}
	}
	if !expr && !canNameOmitted && id == nil {
		return nil, p.errorAtLoc(p.loc(), ERR_CLASS_NAME_REQUIRED)
	}

	var super Node
	if p.lexer.Peek().value == T_EXTENDS {
		p.lexer.Next()
		super, err = p.lhs()
		if err != nil {
			return nil, err
		}
		scope.AddKind(SPK_CLASS_HAS_SUPER)

		if err := p.checkDefaultVal(super, false, false); err != nil {
			return nil, err
		}
	}

	body, err := p.classBody()
	if err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()
	// balance the intermediate mode described above to handle
	// the id part of the class dec
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
	hasCtor := false
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		if tok.value == T_SEMI {
			p.lexer.Next()
			continue
		}
		elem, err := p.classElem()
		if err != nil {
			return nil, err
		}
		if elem.Type() == N_METHOD {
			m := elem.(*Method)
			if hasCtor {
				return nil, p.errorAtLoc(m.key.Loc(), ERR_CTOR_DUP)
			}
			if p.isName(m.key, "constructor", false) {
				hasCtor = true
			}
		}
		elems = append(elems, elem)
	}

	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}

	return &ClassBody{N_ClASS_BODY, p.finLoc(loc), elems}, nil
}

func (p *Parser) classElem() (Node, error) {
	var staticLoc *Loc
	static := false

	ahead := p.lexer.Peek()
	var isField bool
	if ahead.value == T_STATIC {
		staticLoc = p.loc()
		tok := p.lexer.Next()
		static = true

		isField, ahead = p.isField()
		if isField {
			key := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), "static", false, tok.ContainsEscape(), nil, true}
			return p.field(key, nil)
		} else if ahead.value == T_BRACE_L {
			return p.staticBlock()
		} else if ahead.value == T_PAREN_L {
			key := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), "static", false, tok.ContainsEscape(), nil, true}
			return p.method(staticLoc, key, nil, false, PK_METHOD, false, false, false, true, false)
		}
	}

	if p.aheadIsAsync(ahead, true, true) {
		if ahead.ContainsEscape() {
			return nil, p.errorAt(ahead.value, &ahead.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(staticLoc, nil, nil, false, PK_METHOD, false, true, true, true, static)
	} else if ahead.value == T_MUL {
		if p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(ahead)
		}
		return p.method(staticLoc, nil, nil, false, PK_METHOD, true, false, true, true, static)
	}

	tok := p.lexer.Peek()
	propLoc := p.locFromTok(tok)
	kw := tok.IsKw()
	if tok.value == T_NAME || tok.value == T_STRING || kw {
		p.lexer.Next()

		name := tok.text

		var key Node
		if tok.value == T_STRING {
			key = &StrLit{N_LIT_STR, p.finLoc(propLoc.Clone()), name, tok.HasLegacyOctalEscapeSeq(), nil}
		} else {
			key = &Ident{N_NAME, p.finLoc(propLoc.Clone()), name, false, tok.ContainsEscape(), nil, kw}
		}

		isField, ahead = p.isField()
		if isField {
			return p.field(key, staticLoc)
		}

		if static {
			propLoc = staticLoc
		}

		kd := PK_INIT
		if ahead.value == T_PAREN_L {
			kd = PK_METHOD
			if name == "constructor" {
				kd = PK_CTOR
			}
			return p.method(propLoc, key, nil, false, kd, false, false, true, true, static)
		}

		if name == "get" {
			kd = PK_GETTER
		} else if name == "set" {
			kd = PK_SETTER
		} else {
			return nil, p.errorTok(tok)
		}

		return p.method(propLoc, nil, nil, false, kd, false, false, true, true, static)
	}

	return p.field(nil, staticLoc)
}

func (p *Parser) isName(node Node, name string, canContainsEscape bool) bool {
	if node.Type() != N_NAME {
		return false
	}
	id := node.(*Ident)
	if id.Text() != name {
		return false
	}
	if !canContainsEscape {
		return !id.ContainsEscape()
	}
	return true
}

func (p *Parser) field(key Node, static *Loc) (Node, error) {
	var loc *Loc
	var err error
	var compute *Loc
	if key == nil {
		loc = p.loc()
		key, compute, err = p.classElemName()
		if err != nil {
			return nil, err
		}
	} else {
		loc = key.Loc().Clone()
	}

	var value Node
	tok := p.lexer.Peek()
	if tok.value == T_ASSIGN {
		p.lexer.Next()
		value, err = p.assignExpr(true)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		if static != nil {
			loc = static
		}
		return p.method(loc, key, compute, false, PK_METHOD, false, false, true, true, static != nil)
	}
	p.advanceIfSemi(false)

	staticField := static != nil
	if staticField && p.isName(key, "prototype", false) {
		return nil, p.errorAtLoc(key.Loc(), ERR_STATIC_PROP_PROTOTYPE)
	}

	return &Field{N_FIELD, p.finLoc(loc), key, staticField, compute != nil, value}, nil
}

func (p *Parser) classElemName() (Node, *Loc, error) {
	return p.propName(true, false)
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
	expr, err := p.expr()
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

		scope := p.symtab.EnterScope(false, false)
		scope.AddKind(SPK_CATCH)

		names, _, _ := p.collectNames([]Node{param})
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
		arg, err = p.expr()
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
		arg, err = p.expr()
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
	label, err := p.ident(nil)
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

	scope.AddKind(SPK_INTERIM)
	body, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)
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
		label, err := p.ident(nil)
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
		label, err := p.ident(nil)
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
	test, err := p.expr()
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

	scope := p.symtab.EnterScope(false, false)
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

	p.symtab.LeaveScope()

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

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-IfStatement
func (p *Parser) ifStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	scope := p.scope()
	tok := p.lexer.PeekStmtBegin()

	scope.AddKind(SPK_INTERIM)
	cons, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)

	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	var alt Node
	if p.lexer.Peek().value == T_ELSE {
		p.lexer.Next()
		tok := p.lexer.PeekStmtBegin()

		scope.AddKind(SPK_INTERIM)
		alt, err = p.stmt()
		scope.EraseKind(SPK_INTERIM)

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

	scope := p.symtab.EnterScope(false, false)
	scope.AddKind(SPK_LOOP_DIRECT).AddKind(SPK_INTERIM)
	body, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)

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
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	if err := p.advanceIfSemi(false); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()
	return &DoWhileStmt{N_STMT_DO_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WhileStatement
func (p *Parser) whileStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	test, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}

	scope := p.symtab.EnterScope(false, false)
	scope.AddKind(SPK_LOOP_DIRECT)

	tok := p.lexer.PeekStmtBegin()

	scope.AddKind(SPK_INTERIM)
	body, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)

	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	p.symtab.LeaveScope()
	return &WhileStmt{N_STMT_WHILE, p.finLoc(loc), test, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ForStatement
func (p *Parser) forStmt() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	await := false
	ps := p.scope()
	tok := p.lexer.Peek()
	if ps.IsKind(SPK_ASYNC) && tok.value == T_AWAIT {
		if p.feat&FEAT_ASYNC_ITERATION == 0 {
			return nil, p.errorTok(tok)
		}
		await = true
		p.lexer.Next()
	}

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}

	scope := p.symtab.EnterScope(false, false)
	scope.AddKind(SPK_LOOP_DIRECT)
	tok = p.lexer.Peek()

	var init Node
	var err error
	scope.AddKind(SPK_NOT_IN)
	if ok, kind := p.aheadIsVarDec(tok); ok {
		init, err = p.varDecStmt(kind, true)
		if err != nil {
			return nil, err
		}
	} else if tok.value != T_SEMI {
		init, err = p.expr()
		if err != nil {
			return nil, err
		}
	}
	scope.EraseKind(SPK_NOT_IN)

	tok = p.lexer.Peek()
	isIn := IsName(tok, "in", false)
	isOf := IsName(tok, "of", false)

	if await && !isOf {
		if isIn {
			return nil, p.errorAt(T_IN, &tok.begin, "")
		}
		return nil, p.errorTok(tok)
	} else if isOf && !await && init.Type() == N_NAME && init.(*Ident).val == "async" {
		return nil, p.errorAtLoc(init.Loc(), ERR_LHS_OF_FOR_OF_CANNOT_ASYNC)
	}

	if isIn || isOf {
		if init == nil {
			return nil, p.errorTok(tok)
		}

		it := init.Type()
		if it != N_STMT_VAR_DEC {
			if !p.isSimpleLVal(init, true, false, true, false) {
				return nil, p.errorAtLoc(init.Loc(), ERR_ASSIGN_TO_RVALUE)
			}

			// do the `argToParam` check only if the type of init is LitObj or LitArr otherwise
			// just check their simplicity
			if it == N_LIT_OBJ || it == N_LIT_ARR {
				if init, err = p.argToParam(init, 0, false, true, false, false); err != nil {
					return nil, err
				}
			} else if !p.isSimpleLVal(init, true, false, true, false) {
				return nil, p.errorAtLoc(init.Loc(), ERR_ASSIGN_TO_RVALUE)
			}
		} else if it == N_STMT_VAR_DEC {
			varDec := init.(*VarDecStmt)
			if len(varDec.decList) > 1 {
				return nil, p.errorAtLoc(varDec.decList[1].Loc(), ERR_DUP_BINDING)
			}
			if p.scope().IsKind(SPK_STRICT) {
				for _, dec := range varDec.decList {
					if dec.init != nil {
						et := ERR_FOR_OF_LOOP_HAS_INIT
						if isIn {
							et = ERR_FOR_IN_LOOP_HAS_INIT
						}
						return nil, p.errorAtLoc(varDec.loc, et)
					}
				}
			}
		}

		revise := T_IN
		if !isIn {
			revise = T_OF
		}
		p.lexer.NextRevise(revise)

		right, err := p.assignExpr(true)
		if err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_PAREN_R); err != nil {
			return nil, err
		}
		tok := p.lexer.PeekStmtBegin()

		scope.AddKind(SPK_INTERIM)
		body, err := p.stmt()
		scope.EraseKind(SPK_INTERIM)

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
		test, err = p.expr()
		if err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_SEMI); err != nil {
			return nil, err
		}
	}

	var update Node
	if p.lexer.Peek().value != T_PAREN_R {
		update, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}
	tok = p.lexer.PeekStmtBegin()

	scope.AddKind(SPK_INTERIM)
	body, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)

	if err != nil {
		if err == errEof {
			return nil, p.errorTok(tok)
		}
		return nil, err
	}

	p.symtab.LeaveScope()

	return &ForStmt{N_STMT_FOR, p.finLoc(loc), init, test, update, body}, nil
}

func (p *Parser) aheadIsAsync(tok *Token, prop bool, pvt bool) bool {
	if p.feat&FEAT_ASYNC_AWAIT != 0 && IsName(tok, "async", true) {
		ahead := p.lexer.PeekN(2)
		if ahead.afterLineTerminator {
			return false
		}
		if ahead.value == T_FUNC ||
			(ahead.value == T_PAREN_L && !prop) ||
			ahead.value == T_MUL {
			return true
		}
		_, _, canProp := ahead.CanBePropKey()
		if prop && (ahead.value == T_BRACKET_L || ahead.value == T_NAME || ahead.value == T_STRING || canProp) {
			return true
		}
		if pvt && ahead.value == T_NAME_PVT {
			return true
		}
		if ahead.value == T_NAME {
			if p.scope().IsKind(SPK_NOT_IN) && (IsName(ahead, "in", false) || IsName(ahead, "of", false)) {
				return false
			}
			return true
		}
	}
	return false
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-HoistableDeclaration
func (p *Parser) fnDec(expr bool, async *Token, canNameOmitted bool) (Node, error) {
	loc := p.loc()

	// below value cache is needed since token is saved in the ring-buffer
	// so maybe override by next token peek
	asyncHasEscape := false
	var asyncLoc *Loc
	if async != nil {
		asyncHasEscape = async.ContainsEscape()
		loc = p.locFromTok(async)
		asyncLoc = p.locFromTok(async)

		p.lexer.Next()
		p.finLoc(asyncLoc)
		p.lexer.addMode(LM_ASYNC)
	}
	tok := p.lexer.Peek()
	fn := false
	if tok.value == T_FUNC {
		fn = true
		p.lexer.Next()
	}

	tok = p.lexer.Peek()
	generator := tok.value == T_MUL
	if generator {
		if async != nil && p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(tok)
		}
		p.lexer.Next()
	}

	ps := p.scope()
	idScope := ps
	if !ps.IsKind(SPK_STRICT) {
		idScope = ps.OuterFn()
	}

	var id Node
	var err error
	tok = p.lexer.Peek()
	if tok.value != T_PAREN_L {
		id, err = p.ident(idScope)
		if err != nil {
			return nil, err
		}

		// name of the function expression will not add a ref record
		if !expr {
			ref := NewRef()
			ref.Node = id.(*Ident)
			ref.BindKind = BK_VAR
			ref.TargetType = TT_FN
			if ps.IsKind(SPK_STRICT) {
				ref.BindKind = BK_LET
			}
			if err := p.addLocalBinding(ps, ref, ps.IsKind(SPK_STRICT)); err != nil {
				return nil, err
			}
		}
	}
	if fn && !expr && !canNameOmitted && id == nil {
		return nil, p.errorTok(tok)
	}

	scope := p.symtab.EnterScope(true, false)
	if async != nil {
		scope.AddKind(SPK_ASYNC)
	}
	// 'yield' as function names
	if generator {
		p.scope().AddKind(SPK_GENERATOR)
	}

	var args []Node
	// async a => {}
	ahead := p.lexer.Peek()
	if id != nil && ahead.value == T_ARROW && !ahead.afterLineTerminator {
		args = make([]Node, 1)
		args[0] = id
	} else {
		// the arg check is skipped here, the correctness of args is guaranteed by
		// below `argsToFormalParams`
		p.checkName = false
		args, _, err = p.argList(false, false)
		p.checkName = true
		if err != nil {
			return nil, err
		}
	}

	if generator {
		p.lexer.addMode(LM_GENERATOR)
	}

	tok = p.lexer.Peek()
	arrow := false
	var arrowLoc *Loc
	if tok.value == T_ARROW && !tok.afterLineTerminator {
		if !fn {
			arrowLoc = p.locFromTok(p.lexer.Next())
			arrow = true
			scope.AddKind(SPK_ARROW)
		} else {
			return nil, p.errorTok(tok)
		}
	}

	var body Node
	tok = p.lexer.Peek()

	var params []Node
	var paramNames []Node
	var firstComplicated *Loc
	if fn || tok.value == T_BRACE_L || arrow {
		params, err = p.argsToParams(args, false)
		if err != nil {
			return nil, err
		}

		paramNames, firstComplicated, err = p.collectNames(params)
		if err != nil {
			return nil, err
		}
		for _, paramName := range paramNames {
			ref := NewRef()
			ref.Node = paramName.(*Ident)
			ref.BindKind = BK_PARAM
			// duplicate-checking for params is enable in strict and delegated to below `checkParams`
			p.addLocalBinding(nil, ref, false)
		}

		if tok.value == T_BRACE_L {
			if body, err = p.fnBody(); err != nil {
				return nil, err
			}
		} else if expr || arrow {
			if body, err = p.expr(); err != nil {
				return nil, err
			}
		} else {
			return nil, p.errorTok(tok)
		}
	} else if async != nil {
		// this branch means the input is callExpr like:
		// `async ({a: b = c})` callExpr
		// `async* ({a: b = c})` binExpr
		lhs := &Ident{N_NAME, asyncLoc, "async", false, asyncHasEscape, nil, true}

		var expr Node
		if generator {
			var rhs Node
			argsLen := len(args)
			if argsLen == 0 {
				return nil, p.errorTok(tok)
			} else if argsLen == 1 {
				rhs = args[0]
			} else {
				rhs = &SeqExpr{N_EXPR_SEQ, p.finLoc(loc), args, nil}
			}
			expr = &BinExpr{N_EXPR_BIN, p.finLoc(loc), T_MUL, lhs, rhs, nil}
		} else {
			if err := p.checkArgs(args, false, true); err != nil {
				return nil, err
			}
			expr = &CallExpr{N_EXPR_CALL, p.finLoc(loc), lhs, args, nil}
		}

		return &ExprStmt{N_STMT_EXPR, p.finLoc(loc.Clone()), expr}, nil
	} else {
		return nil, p.errorTok(tok)
	}

	if generator {
		p.lexer.popMode()
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)
	if id != nil && p.isProhibitedName(idScope, id.(*Ident).Text(), isStrict) {
		return nil, p.errorAtLoc(id.Loc(), ERR_RESERVED_WORD_IN_STRICT_MODE)
	}

	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()

	if arrow {
		fn := &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), arrowLoc, async != nil, params, body, body.Type() != N_STMT_BLOCK, nil}
		if expr {
			return fn, nil
		}
		return &ExprStmt{N_STMT_EXPR, p.finLoc(loc.Clone()), fn}, nil
	}

	typ := N_STMT_FN
	if expr {
		typ = N_EXPR_FN
	}
	return &FnDec{typ, p.finLoc(loc), id, generator, async != nil, params, body, nil}, nil
}

func (p *Parser) collectNames(nodes []Node) (names []Node, firstComplicated *Loc, err error) {
	names = make([]Node, 0)
	var ns []Node
	for _, param := range nodes {
		if firstComplicated == nil && param.Type() != N_NAME {
			firstComplicated = param.Loc()
		}
		ns, err = p.namesInPattern(param, false)
		if err != nil {
			return nil, nil, err
		}
		names = append(names, ns...)
	}
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#sec-parameter-lists-static-semantics-early-errors
// `isSimpleParamList` should be true if function body directly contains `use strict` directive
func (p *Parser) checkParams(names []Node, firstComplicated *Loc, isStrict bool, directStrict bool) error {
	var dupLoc *Loc
	unique := make(map[string]bool)
	for _, id := range names {
		name := id.(*Ident).Text()
		if p.isProhibitedName(nil, name, isStrict) {
			return p.errorAtLoc(id.Loc(), fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, name))
		}

		if dupLoc == nil {
			if _, ok := unique[name]; ok {
				dupLoc = id.Loc()

			} else {
				unique[name] = true
			}
		}
	}

	if directStrict && firstComplicated != nil {
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

func (p *Parser) namesInPattern(node Node, kw bool) ([]Node, error) {
	out := make([]Node, 0)
	if node == nil {
		return out, nil
	}
	switch node.Type() {
	case N_PAT_ARRAY:
		elems := node.(*ArrPat).elems
		for _, node := range elems {
			names, err := p.namesInPattern(node, false)
			if err != nil {
				return nil, err
			}
			out = append(out, names...)
		}
	case N_PAT_ASSIGN:
		names, err := p.namesInPattern(node.(*AssignPat).lhs, kw)
		if err != nil {
			return nil, err
		}
		out = append(out, names...)
	case N_PAT_OBJ:
		props := node.(*ObjPat).props
		for _, node := range props {
			var names []Node
			var err error
			if node.Type() == N_NAME {
				names, err = p.namesInPattern(node, false)
			} else if node.Type() == N_PROP {
				val := node.(*Prop).value
				names, err = p.namesInPattern(val, false)
			} else {
				names, err = p.namesInPattern(node, false)
			}
			if err != nil {
				return nil, err
			}
			out = append(out, names...)
		}
	case N_PAT_REST:
		names, errLoc := p.namesInPattern(node.(*RestPat).arg, kw)
		if errLoc != nil {
			return nil, errLoc
		}
		out = append(out, names...)
	case N_NAME:
		id := node.(*Ident)
		if !kw && id.kw {
			em := ERR_TPL_UNEXPECTED_TOKEN_TYPE
			if id.Text() == "eval" {
				em = ERR_TPL_BINDING_RESERVED_WORD
			}
			return nil, p.errorAtLoc(node.Loc(), fmt.Sprintf(em, id.Text()))
		}
		out = append(out, node)
	}
	return out, nil
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

func (p *Parser) enterStrict(lex bool) *Scope {
	if lex {
		p.lexer.addMode(LM_STRICT)
	}
	return p.scope().AddKind(SPK_STRICT).AddKind(SPK_STRICT_DIR)
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

		ahead := p.lexer.PeekStmtBegin()
		if stmt != nil {
			// StrictDirective processing
			if (scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_GLOBAL)) && stmt.Type() == N_STMT_EXPR {
				if p.isStrictDirective(stmt) {
					// lexer will automatically pop it's mode when the `T_BRACE_R` is met
					// here we use `ahead.value != T_BRACE_R` to prevent accidentally change
					// the upper lexer mode
					p.enterStrict(ahead.value != T_BRACE_R)

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
		p.symtab.EnterScope(false, false)
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

func (p *Parser) aheadIsVarDec(tok *Token) (bool, TokenValue) {
	if tok.value == T_VAR {
		return true, T_VAR
	}
	if p.feat&FEAT_LET_CONST != 0 {
		var ok bool
		var v TokenValue

		if tok.value == T_LET || tok.value == T_CONST {
			ok = true
			v = tok.value
		} else if IsName(tok, "let", false) {
			ok = true
			v = T_LET
		} else if IsName(tok, "const", false) {
			ok = true
			v = T_CONST
		}

		if !ok {
			return false, T_ILLEGAL
		}

		if p.scope().IsKind(SPK_STRICT) {
			return true, v
		}

		// an additional lookahead is needed to judge the various:
		// - `let + 1`
		// - `let a`
		ahead := p.lexer.PeekGrow()
		av := ahead.value
		if !ahead.afterLineTerminator && (av == T_NAME ||
			(av > T_CTX_KEYWORD_BEGIN && av < T_CTX_KEYWORD_END) ||
			av == T_BRACE_L || av == T_BRACKET_L) {
			return true, v
		}
	}
	return false, T_ILLEGAL
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
func (p *Parser) varDecStmt(kind TokenValue, asExpr bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	node := &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]*VarDec, 0, 1), nil}

	isConst := false
	node.kind = kind
	bindKind := BK_VAR
	if kind == T_LET {
		bindKind = BK_LET
	} else if kind == T_CONST {
		isConst = true
		bindKind = BK_CONST
	}

	lvs := make([]Node, 0)
	for {
		dec, err := p.varDec(bindKind != BK_VAR)
		if err != nil {
			return nil, err
		}
		lvs = append(lvs, dec.id)

		if isConst && dec.init == nil && !p.scope().IsKind(SPK_NOT_IN) {
			return nil, p.errorAtLoc(dec.loc, ERR_CONST_DEC_INIT_REQUIRED)
		}

		node.decList = append(node.decList, dec)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		} else {
			break
		}
	}

	names, _, err := p.collectNames(lvs)
	if err != nil {
		return nil, err
	}
	node.names = names

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

func (p *Parser) varDec(lexical bool) (*VarDec, error) {
	scope := p.scope()
	if lexical {
		scope.AddKind(SPK_LEXICAL_DEC)
	}

	binding, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}
	loc := binding.Loc().Clone()
	scope.EraseKind(SPK_LEXICAL_DEC)

	var init Node
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr(true)
		if err != nil {
			return nil, err
		}
	}

	if binding.Type() != N_NAME && init == nil && !p.scope().IsKind(SPK_NOT_IN) {
		return nil, p.errorAtLoc(p.loc(), ERR_COMPLEX_BINDING_MISSING_INIT)
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
func (p *Parser) isProhibitedName(scope *Scope, name string, withStrict bool) bool {
	if scope == nil {
		scope = p.scope()
	}
	strict := withStrict && scope.IsKind(SPK_STRICT)
	_, ok := prohibitedNames[name]
	if strict && ok {
		return true
	}
	return scope.IsKind(SPK_ASYNC) && name == "await"
}

func (p *Parser) ident(scope *Scope) (*Ident, error) {
	if scope == nil {
		scope = p.scope()
	}

	tok, err := p.nextMustTok(T_NAME)
	if err != nil {
		return nil, err
	}

	name := tok.Text()
	ident := &Ident{N_NAME, nil, "", false, tok.ContainsEscape(), nil, false}
	ident.loc = p.finLoc(p.locFromTok(tok))
	ident.val = name

	if p.isProhibitedName(scope, ident.val, true) {
		return nil, p.errorAtLoc(ident.loc, ERR_RESERVED_WORD_IN_STRICT_MODE)
	}

	// for resporting `'let' is disallowed as a lexically bound name` for stmt like `let let`
	if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) && !tok.ContainsEscape() {
		if name == "let" || name == "const" {
			return nil, p.errorAtLoc(ident.loc, fmt.Sprintf(ERR_TPL_FORBIDED_LEXICAL_NAME, name))
		}
	}

	return ident, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#sec-destructuring-binding-patterns
func (p *Parser) bindingPattern() (Node, error) {
	tok := p.lexer.Peek()
	tv := tok.value

	var binding Node
	var err error

	if p.feat&FEAT_BINDING_PATTERN == 0 && (tv == T_BRACE_L || tv == T_BRACKET_L) {
		return nil, p.errorTok(tok)
	}

	if tv == T_BRACE_L {
		binding, err = p.patternObj()
	} else if tv == T_BRACKET_L {
		binding, err = p.patternArr()
	} else {
		binding, err = p.ident(nil)
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
		binding, err := p.patternRest(false)
		if err != nil {
			return nil, err
		}
		return binding, nil
	}

	key, compute, err := p.propName(false, true)
	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
	}

	loc := compute
	if loc == nil {
		loc = key.Loc().Clone()
	}

	tok := p.lexer.Peek()
	opLoc := p.locFromTok(tok)
	assign := tok.value == T_ASSIGN
	var value Node
	if tok.value == T_COLON {
		p.lexer.Next()
		value, err = p.bindingElem(false)
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

	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, compute != nil, false, shorthand, assign, PK_INIT}, nil
}

// test whether current place is filed or method
func (p *Parser) isField() (bool, *Token) {
	ahead := p.lexer.Peek()
	isField := ahead.value == T_COLON ||
		ahead.value == T_ASSIGN ||
		ahead.value == T_SEMI ||
		ahead.value == T_COMMA ||
		ahead.value == T_BRACE_R ||
		ahead.afterLineTerminator
	return isField, ahead
}

func (p *Parser) propName(allowNamePVT bool, maybeMethod bool) (Node, *Loc, error) {
	var key Node
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	keyName, kw, ok := tok.CanBePropKey()

	scope := p.scope()
	var computeLoc *Loc
	if allowNamePVT && tok.value == T_NAME_PVT {
		key = &Ident{N_NAME, p.finLoc(loc), tok.Text(), true, tok.ContainsEscape(), nil, false}
	} else if tok.value == T_STRING {
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, nil, p.errorAtLoc(p.locFromTok(tok), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		key = &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil}
	} else if tok.value == T_NUM {
		key = &NumLit{N_LIT_NUM, p.finLoc(loc), nil}
	} else if tok.value == T_BRACKET_L {
		computeLoc = p.locFromTok(tok)
		name, err := p.assignExpr(true)
		if err != nil {
			return nil, nil, err
		}
		_, err = p.nextMustTok(T_BRACKET_R)
		if err != nil {
			return nil, nil, err
		}
		key = name
	} else if ok {
		if !kw && p.isProhibitedName(nil, keyName, true) {
			kw = true
		}
		// stmt `let { let } = {}` will raise error `let is disallowed as a lexically bound name` in sloppy mode
		if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) {
			if !tok.ContainsEscape() && (keyName == "let" || keyName == "const") {
				return nil, nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_FORBIDED_LEXICAL_NAME, keyName))
			}
		}
		key = &Ident{N_NAME, p.finLoc(loc), keyName, false, tok.ContainsEscape(), nil, kw}
	} else {
		return nil, nil, p.errorTok(tok)
	}

	isField, ahead := p.isField()
	if isField || !maybeMethod {
		return key, computeLoc, nil
	}

	kd := PK_INIT
	loc = loc.Clone()
	if keyName == "get" || keyName == "set" {
		if tok.ContainsEscape() {
			return nil, nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		if ahead.value != T_PAREN_L {
			key = nil
			if keyName == "get" {
				kd = PK_GETTER
			} else {
				kd = PK_SETTER
			}
		}
	}

	m, err := p.method(loc, key, computeLoc, false, kd, false, false, false, false, false)
	if err != nil {
		return nil, nil, err
	}
	return m, nil, nil
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
		binding, err = p.patternRest(true)
	} else {
		binding, err = p.ident(nil)
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
		init, err = p.assignExpr(true)
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
	return &Prop{N_PROP, p.finLoc(loc.Clone()), val.lhs, opLoc, val, false, false, true, true, PK_INIT}, nil
}

func (p *Parser) patternRest(arrPat bool) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()

	if p.feat&FEAT_BINDING_REST_ELEM == 0 {
		return nil, p.errorTok(tok)
	}

	ahead := p.lexer.Peek()
	if ahead.value != T_NAME && (!arrPat || ahead.value != T_BRACKET_L) {
		return nil, p.errorTok(ahead)
	}

	arg, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}

	tok = p.lexer.Peek()
	if tok.value == T_COMMA {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_REST_TRAILING_COMMA)
	}

	return &RestPat{N_PAT_REST, p.finLoc(loc), arg, nil}, nil
}

func (p *Parser) exprStmt() (Node, error) {
	loc := p.loc()
	stmt := &ExprStmt{N_STMT_EXPR, &Loc{}, nil}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	stmt.expr = expr

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	stmt.loc = p.finLoc(loc)
	return stmt, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-Expression
func (p *Parser) expr() (Node, error) {
	return p.seqExpr()
}

func (p *Parser) seqExpr() (Node, error) {
	loc := p.loc()
	expr, err := p.assignExpr(true)
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
			expr, err = p.assignExpr(true)
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
	return p.lexer.Peek().value == T_YIELD
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-YieldExpression
func (p *Parser) yieldExpr() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	tok := p.lexer.Peek()
	kind := TokenKinds[tok.value]
	if tok.afterLineTerminator || !kind.StartExpr && tok.value != T_MUL {
		return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), false, nil, nil}, nil
	}

	delegate := false
	if p.lexer.Peek().value == T_MUL {
		p.lexer.Next()
		delegate = true
	}

	arg, err := p.assignExpr(true)
	if err != nil {
		return nil, err
	}
	return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), delegate, arg, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-AssignmentExpression
func (p *Parser) assignExpr(checkLhs bool) (Node, error) {
	if p.aheadIsYield() {
		return p.yieldExpr()
	}

	lhs, err := p.condExpr()
	if err != nil {
		return nil, err
	}
	loc := lhs.Loc().Clone()

	tok := p.lexer.Peek()
	if lhs.Type() == N_NAME && tok.value == T_ARROW && !tok.afterLineTerminator {
		fn, err := p.arrowFn(loc, []Node{lhs})
		if err != nil {
			return nil, err
		}
		lhs = fn
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == nil {
		return lhs, nil
	}
	op := assign.value
	opLoc := p.locFromTok(assign)

	rhs, err := p.assignExpr(checkLhs)
	if err != nil {
		return nil, err
	}

	// set `depth` to 1 to permit expr like `i + 2 = 42`
	// and so just do the arg to param transform silently
	lhs, err = p.argToParam(lhs, 1, false, true, false, false)
	if err != nil {
		return nil, err
	}

	if checkLhs && !p.isSimpleLVal(lhs, true, false, true, op != T_ASSIGN) {
		return nil, p.errorAtLoc(lhs.Loc(), ERR_ASSIGN_TO_RVALUE)
	}

	if err := p.checkArg(rhs, false, false); err != nil {
		return nil, err
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finLoc(loc), op, opLoc, lhs, rhs, nil}
	return node, nil
}

// https://tc39.es/ecma262/multipage/syntax-directed-operations.html#sec-static-semantics-assignmenttargettype
// `pat` indicates whether to treat the pattern syntax as legal or not
// `member` indicates whether the member expr can be treated as legal or not
// `optAssign` indicats whether the expr is the lhs of the op-assign expr
func (p *Parser) isSimpleLVal(expr Node, pat bool, inParen bool, member bool, optAssign bool) bool {
	switch expr.Type() {
	case N_NAME:
		node := expr.(*Ident)
		if p.scope().IsKind(SPK_ASYNC) && node.Text() == "await" {
			return false
		}
		return true
	case N_PAT_ASSIGN, N_PAT_REST:
		if inParen || optAssign {
			return false
		}
	case N_PAT_OBJ, N_PAT_ARRAY, N_LIT_ARR, N_LIT_OBJ:
		if optAssign {
			return false
		}
		return pat
	case N_EXPR_MEMBER:
		return member
	case N_EXPR_PAREN:
		node := expr.(*ParenExpr)
		return p.isSimpleLVal(node.expr, pat, true, false, optAssign)
	}
	return false
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

	cons, err := p.assignExpr(true)
	if err != nil {
		return nil, err
	}

	_, err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr(true)
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_COND, p.finLoc(loc), test, cons, alt, nil}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-AwaitExpression
func (p *Parser) awaitExpr(tok *Token) (Node, error) {
	loc := p.locFromTok(tok)
	ahead := p.lexer.Peek()
	if !TokenKinds[ahead.value].StartExpr {
		// report friendly message for expr like: `async function foo(await) {}`
		if ahead.value == T_PAREN_R || ahead.value == T_COMMA {
			return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, "await"))
		}
		return nil, p.errorTok(ahead)
	}
	arg, err := p.unaryExpr()
	if err != nil {
		return nil, err
	}
	return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), T_AWAIT, arg, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-UnaryExpression
func (p *Parser) unaryExpr() (Node, error) {
	tok := p.lexer.Peek()
	loc := p.locFromTok(tok)
	op := tok.value
	if tok.IsUnary() || tok.value == T_ADD || tok.value == T_SUB {
		p.lexer.Next()
		arg, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		// current es grammar does not allow the arrowFn to be the arg of
		// the unaryExpr such as `typeof () => {}` will raise an exception
		// `Malformed arrow function parameter list`
		if arg.Type() == N_EXPR_ARROW {
			return nil, p.errorAtLoc(arg.Loc(), ERR_MALFORMED_ARROW_PARAM)
		}

		scope := p.scope()
		if scope.IsKind(SPK_STRICT) && tok.value == T_DELETE && arg.Type() == N_NAME {
			return nil, p.errorAtLoc(arg.Loc(), ERR_DELETE_LOCAL_IN_STRICT)
		}

		return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), op, arg, nil}, nil
	} else if tok.value == T_AWAIT {
		if !p.scope().IsKind(SPK_ASYNC) {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_AWAIT_OUTSIDE_ASYNC)
		}
		p.lexer.Next()
		return p.awaitExpr(tok)
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
		if !p.isSimpleLVal(arg, true, false, true, false) {
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

	if !p.isSimpleLVal(arg, true, false, true, false) {
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
	new := p.lexer.Next()

	var expr Node
	var err error

	scope := p.scope()
	tok := p.lexer.Peek()
	if tok.value == T_DOT && p.feat&FEAT_META_PROPERTY != 0 {
		meta := &Ident{N_NAME, p.finLoc(p.locFromTok(new)), "new", false, new.ContainsEscape(), nil, true}
		p.lexer.Next() // consume dot

		id, err := p.ident(nil)
		if err != nil {
			return nil, err
		}
		if id.Text() != "target" {
			return nil, p.errorAtLoc(id.loc, ERR_INVALID_META_PROP)
		}
		if !(scope.IsKind(SPK_CLASS) || scope.IsKind(SPK_CLASS_INDIRECT) ||
			(!scope.IsKind(SPK_ARROW) && scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_FUNC_INDIRECT))) {
			return nil, p.errorAtLoc(loc, ERR_META_PROP_OUTSIDE_FN)
		}

		expr = &MetaProp{N_META_PROP, p.finLoc(loc), meta, id}
		return expr, nil
	}

	expr, err = p.memberExpr(nil, false)
	if err != nil {
		return nil, err
	}

	var args []Node
	if p.lexer.Peek().value == T_PAREN_L {
		args, _, err = p.argList(true, true)
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

func (p *Parser) checkCallee(callee Node, nextLoc *Loc) error {
	scope := p.scope()
	switch callee.Type() {
	case N_EXPR_FN, N_EXPR_ARROW:
		if !scope.IsKind(SPK_PAREN) {
			return p.errorAtLoc(nextLoc, ERR_UNEXPECTED_TOKEN)
		}
	case N_SUPER:
		if !scope.IsKind(SPK_CTOR) || !scope.IsKind(SPK_CLASS_HAS_SUPER) {
			return p.errorAtLoc(callee.Loc(), ERR_SUPER_CALL_OUTSIDE_CTOR)
		}
	}
	return nil
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

	ahead := p.lexer.Peek()
	if ahead.value == T_PAREN_L {
		if err = p.checkCallee(callee, p.locFromTok(ahead)); err != nil {
			return nil, err
		}
	}

	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_L {
			args, _, err := p.argList(true, true)
			if err != nil {
				return nil, err
			}
			callee = &CallExpr{N_EXPR_CALL, p.finLoc(loc), callee, args, nil}
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
	meta := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), tok.Text(), false, tok.ContainsEscape(), nil, false}

	if p.lexer.Peek().value == T_DOT {
		p.lexer.Next()
		prop, err := p.ident(nil)
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
	src, err := p.assignExpr(true)
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
	}

	elems := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value >= T_TPL_HEAD || tok.value <= T_TPL_TAIL {
			cooked := ""
			ext := tok.ext.(*TokExtTplSpan)
			if ext.IllegalEscape != nil {
				// raise error for bad escape sequence if the template is not tagged
				if tag == nil {
					return nil, p.errorAt(tok.value, ext.IllegalEscape.Loc.begin, ext.IllegalEscape.Err)
				}
			} else {
				cooked = ext.str
			}

			loc := &Loc{
				src:   ext.strRng.src,
				begin: ext.strBegin.Clone(),
				end:   ext.strEnd.Clone(),
				rng:   &Range{ext.strRng.lo, ext.strRng.hi},
			}
			p.lexer.Next()
			str := &StrLit{N_LIT_STR, loc, cooked, false, nil}
			elems = append(elems, str)

			if tok.value == T_TPL_TAIL || tok.IsPlainTpl() {
				break
			}

			expr, err := p.expr()
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		} else {
			return nil, p.errorTok(tok)
		}
	}

	loc = p.finLoc(loc)
	return &TplExpr{N_EXPR_TPL, loc, tag, elems, nil}, nil
}

func (p *Parser) argsToParams(args []Node, setter bool) ([]Node, error) {
	params := make([]Node, len(args))
	var err error
	for i, arg := range args {
		if arg != nil {
			params[i], err = p.argToParam(arg, 0, false, false, false, setter)
			if err != nil {
				return nil, err
			}
		}
	}
	return params, nil
}

// `yield` indicates whether is yield-expr is permitted
func (p *Parser) checkDefaultVal(val Node, yield bool, destruct bool) error {
	switch val.Type() {
	case N_EXPR_YIELD:
		scope := p.scope()
		if !yield || !scope.IsKind(SPK_GENERATOR) {
			return p.errorAtLoc(val.Loc(), ERR_YIELD_CANNOT_BE_DEFAULT_VALUE)
		}
		return nil
	case N_EXPR_BIN:
		n := val.(*BinExpr)
		if err := p.checkDefaultVal(n.lhs, yield, destruct); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.rhs, yield, destruct); err != nil {
			return err
		}
	case N_EXPR_PAREN:
		n := val.(*ParenExpr)
		return p.checkDefaultVal(n.expr, yield, destruct)
	case N_EXPR_UNARY:
		n := val.(*UnaryExpr)
		// `{a = await b} = obj` is legal
		// `({a = await b}) => obj` is illegal
		if n.op == T_AWAIT && !destruct {
			return p.errorAtLoc(n.loc, ERR_AWAIT_AS_DEFAULT_VALUE)
		}
		return p.checkDefaultVal(n.arg, yield, destruct)
	case N_EXPR_UPDATE:
		n := val.(*UpdateExpr)
		return p.checkDefaultVal(n.arg, yield, destruct)
	case N_EXPR_COND:
		n := val.(*CondExpr)
		if err := p.checkDefaultVal(n.test, yield, destruct); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.cons, yield, destruct); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.alt, yield, destruct); err != nil {
			return err
		}
	case N_PAT_ASSIGN:
		n := val.(*AssignPat)
		if err := p.checkDefaultVal(n.lhs, yield, destruct); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.rhs, yield, destruct); err != nil {
			return err
		}
	case N_LIT_ARR:
		n := val.(*ArrLit)
		for _, elem := range n.elems {
			if err := p.checkDefaultVal(elem, yield, destruct); err != nil {
				return err
			}
		}
	case N_SPREAD:
		n := val.(*Spread)
		return p.checkDefaultVal(n.arg, yield, destruct)
	case N_NAME:
		id := val.(*Ident)
		name := val.(*Ident).Text()
		if p.checkName && p.isProhibitedName(nil, name, true) {
			return p.errorAtLoc(id.loc, ERR_RESERVED_WORD_IN_STRICT_MODE)
		}
	}
	return nil
}

func (p *Parser) isPrimitive(node Node) bool {
	switch node.Type() {
	case N_NAME, N_LIT_BOOL, N_LIT_NUM, N_LIT_REGEXP:
		return true
	}
	return false
}

// `destruct` indicate whether the parsing state is in destructing assignment or not
func (p *Parser) argToParam(arg Node, depth int, prop bool, destruct bool, inParen bool, setter bool) (Node, error) {
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
			// elem maybe nil in expr like `([a, , b]) => 42`
			if node != nil {
				pat.elems[i], err = p.argToParam(node, depth+1, false, destruct, inParen, setter)
				if err != nil {
					return nil, err
				}
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
		for i, prop := range n.props {
			pp, err := p.argToParam(prop, depth+1, true, destruct, inParen, setter)
			if err != nil {
				return nil, err
			}
			pat.props[i] = pp
		}
		return pat, nil
	case N_PAT_OBJ:
		// `function* wrap() {({a = yield b} = obj) }` is legal
		// `function* wrap() {({a = yield b} = obj) => a }` is illegal
		// so the firstly `argToParam` to make the `objectPat` of lhs of the assignExpr
		// does not raise the error `Yield expression cannot be a default value` for `yield b`
		// it's the duty of second `argToParam` to raise that error after `=>` is consumed
		n := arg.(*ObjPat)
		for _, prop := range n.props {
			_, err := p.argToParam(prop, depth+1, true, destruct, inParen, setter)
			if err != nil {
				return nil, err
			}
		}
		return n, nil
	case N_PROP:
		n := arg.(*Prop)
		var err error
		if n.value != nil {
			if n.value.Type() == N_EXPR_FN && depth > 0 {
				return nil, p.errorAtLoc(arg.Loc(), ERR_OBJ_PATTERN_CANNOT_FN)
			}

			if n.assign {
				// the key is needed to be checked as a legal binding name
				if n.key, err = p.argToParam(n.key, depth+1, prop, destruct, inParen, setter); err != nil {
					return nil, err
				}
				if err = p.checkDefaultVal(n.value, destruct, destruct); err != nil {
					return nil, err
				}
			} else {
				// the correctness of the value should be checked account for
				// using it as an alias
				if n.value, err = p.argToParam(n.value, depth+1, prop, destruct, inParen, setter); err != nil {
					return nil, err
				}
			}
		} else if n.key.Type() != N_NAME {
			// raise syntax error for stmt like `({ 5 }) => {}`
			return nil, p.errorAtLoc(n.key.Loc(), ERR_UNEXPECTED_TOKEN)
		}
		if n.assign {
			if !prop && depth > 0 {
				return nil, p.errorAtLoc(n.opLoc, ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT)
			}

			// also check the default value
			err = p.checkDefaultVal(n.value, destruct, destruct)
			if err != nil {
				return nil, err
			}

			if n.value.Type() == N_PAT_ASSIGN {
				return n, nil
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

		lhs, err := p.argToParam(n.lhs, depth+1, false, destruct, inParen, setter)
		if err != nil {
			return nil, err
		}

		// also check the default value
		err = p.checkDefaultVal(n.rhs, destruct, destruct)
		if err != nil {
			return nil, err
		}

		err = p.checkArg(n.rhs, false, false)
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
	case N_NAME:
		id := arg.(*Ident)
		name := arg.(*Ident).Text()
		if p.checkName && p.isProhibitedName(nil, name, true) {
			et := ERR_TPL_BINDING_RESERVED_WORD
			if destruct {
				et = ERR_TPL_ASSIGN_TO_RESERVED_WORD_IN_STRICT_MODE
			}
			return nil, p.errorAtLoc(id.loc, fmt.Sprintf(et, name))
		}
		return arg, nil
	case N_PAT_REST:
		return arg, nil
	case N_SPREAD:
		n := arg.(*Spread)
		if n.trailingCommaLoc != nil {
			return nil, p.errorAtLoc(n.trailingCommaLoc, ERR_REST_TRAILING_COMMA)
		}

		if setter {
			return nil, p.errorAtLoc(n.loc, ERR_REST_IN_SETTER)
		}

		at := n.arg.Type()
		if at == N_NAME {
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
		} else if at == N_EXPR_ASSIGN {
			return nil, p.errorAtLoc(n.arg.Loc(), ERR_REST_CANNOT_SET_DEFAULT)
		} else if at == N_EXPR_PAREN {
			if destruct {
				arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen, setter)
				if err != nil {
					return nil, err
				}
				n.arg = arg
			} else {
				return nil, p.errorAtLoc(n.arg.Loc(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
			}
		} else if p.feat&FEAT_BINDING_REST_ELEM_NESTED != 0 && (at == N_LIT_ARR || at == N_LIT_OBJ) {
			arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen, setter)
			if err != nil {
				return nil, err
			}
			if prop && !p.isSimpleLVal(arg, false, false, true, true) {
				return nil, p.errorAtLoc(arg.Loc(), ERR_REST_ARG_NOT_SIMPLE)
			}
			n.arg = arg
		} else {
			if !prop && p.feat&FEAT_BINDING_REST_ELEM_NESTED == 0 {
				nested := p.UnParen(n.arg)
				if nested.Type() != N_NAME {
					return nil, p.errorAtLoc(nested.Loc(), ERR_REST_ARG_NOT_BINDING_PATTERN)
				}
			}

			return nil, p.errorAtLoc(n.arg.Loc(), ERR_REST_ARG_NOT_SIMPLE)
		}

		return &RestPat{
			typ: N_PAT_REST,
			loc: n.loc,
			arg: n.arg,
		}, nil
	case N_EXPR_PAREN:
		sub := arg.(*ParenExpr).expr
		if !destruct || !p.isPrimitive(sub) {
			st := sub.Type()
			if st != N_LIT_ARR && st != N_LIT_OBJ && st != N_NAME {
				return nil, p.errorAtLoc(sub.Loc(), ERR_ASSIGN_TO_RVALUE)
			}
			return nil, p.errorAtLoc(arg.Loc(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
		}
		arg, err := p.argToParam(sub, depth, prop, destruct, true, setter)
		if err != nil {
			return nil, err
		}
		return arg, nil
	}
	if depth == 0 {
		return nil, p.errorAtLoc(arg.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	// `([a.a]) => 42` is illegal since the `a.a` is not permitted to occur
	// `[a.r] = b` is legal since `a.r` is permitted to occur in destruct
	if !p.isSimpleLVal(arg, true, inParen, destruct, false) {
		return nil, p.errorAtLoc(arg.Loc(), ERR_ASSIGN_TO_RVALUE)
	}
	return arg, nil
}

func (p *Parser) nameOfProp(prop *Prop) string {
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
	return propName
}

// check the `arg` is legal as argument
// `spread` means whether the spread is permitted
// `simplicity` means whether check simplicity of lhs of the assignExpr
func (p *Parser) checkArg(arg Node, spread bool, simplicity bool) error {
	switch arg.Type() {
	case N_LIT_OBJ:
		n := arg.(*ObjLit)
		hasProto := false
		for _, prop := range n.props {
			err := p.checkArg(prop, true, false)
			if err != nil {
				return err
			}

			if prop.Type() == N_PROP {
				pp := prop.(*Prop)
				pn := p.nameOfProp(pp)
				if !pp.computed && !pp.method && !pp.shorthand && pn == "__proto__" {
					if hasProto {
						return p.errorAtLoc(pp.Loc(), ERR_REDEF_PROP)
					}
					hasProto = true
				}
			}
		}
	case N_PROP:
		n := arg.(*Prop)
		if n.assign && n.shorthand {
			return p.errorAtLoc(n.opLoc, ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT)
		}
		var err error
		if n.value != nil {
			err = p.checkArg(n.value, true, false)
			if err != nil {
				return err
			}
		}
	case N_PAT_REST, N_SPREAD:
		if !spread {
			return p.errorAtLoc(arg.Loc(), ERR_UNEXPECTED_TOKEN)
		}
	case N_EXPR_ASSIGN:
		if simplicity {
			n := arg.(*AssignExpr)
			if n.op != T_ASSIGN && !p.isSimpleLVal(n.lhs, false, false, true, false) {
				return p.errorAtLoc(n.lhs.Loc(), ERR_ASSIGN_TO_RVALUE)
			}
		}
	case N_EXPR_PAREN:
		return p.checkArg(arg.(*ParenExpr).expr, spread, simplicity)
	case N_NAME:
		id := arg.(*Ident)
		if id.kw && p.scope().IsKind(SPK_STRICT) {
			return p.errorAtLoc(arg.Loc(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.Text()))
		}
	}
	return nil
}

func (p *Parser) checkArgs(args []Node, spread bool, simplicity bool) error {
	for _, arg := range args {
		err := p.checkArg(arg, spread, simplicity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) argList(check bool, incall bool) ([]Node, *Loc, error) {
	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, nil, err
	}

	var tailingComma *Loc
	args := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, nil, p.errorTok(tok)
		}
		arg, err := p.arg()
		if err != nil {
			return nil, nil, err
		}

		ahead := p.lexer.Peek()
		if ahead.value == T_COMMA {
			tok := p.lexer.Next()
			// trailing comma is need to be checked when it's in
			// parenExpr, this snippet `...a,` is illegal as parenExpr: `(...a,)`
			// however it is legal as arguments `foo(...a,)`
			ahead := p.lexer.Peek()
			if !incall && arg.Type() == N_SPREAD {
				msg := ERR_REST_TRAILING_COMMA
				if ahead.value != T_PAREN_R {
					msg = ERR_REST_ELEM_MUST_LAST
				}
				return nil, tailingComma, p.errorAt(tok.value, &tok.begin, msg)
			}
			if tailingComma == nil && ahead.value == T_PAREN_R {
				tailingComma = p.locFromTok(tok)
			}
		} else if ahead.value != T_PAREN_R {
			return nil, nil, p.errorTok(ahead)
		}

		if check {
			// `spread` or `pattern_rest` expression is legal argument:
			// `f(c, b, ...a)`
			if err := p.checkArg(arg, true, false); err != nil {
				return nil, nil, err
			}
		}

		args = append(args, arg)
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, nil, err
	}
	return args, tailingComma, nil
}

func (p *Parser) arg() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(false)
}

func (p *Parser) checkOp(tok *Token) error {
	if tok.value == T_POW && p.feat&FEAT_POW == 0 {
		return p.errorTok(tok)
	}
	return nil
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
	notIn := p.scope().IsKind(SPK_NOT_IN)
	for {
		op := ahead.IsBin(notIn)
		if op == T_ILLEGAL || pcd < minPcd {
			break
		}
		if err = p.checkOp(ahead); err != nil {
			return nil, err
		}
		p.lexer.Next()

		rhs, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		kind = ahead.Kind()
		for ahead.IsBin(notIn) != T_ILLEGAL && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			rhs, err = p.binExpr(rhs, pcd)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			kind = ahead.Kind()
		}
		pcd = kind.Pcd

		// deal with expr like: `console.log( -2 ** 4 )`
		if lhs.Type() == N_EXPR_UNARY && op == T_POW {
			return nil, p.errorAtLoc(p.UnParen(lhs.(*UnaryExpr).arg).Loc(), ERR_UNARY_OPERATOR_IMMEDIATELY_BEFORE_POW)
		}

		// deal with expr like: `4 + async() => 2`
		if rhs.Type() == N_EXPR_ARROW {
			return nil, p.errorAtLoc(rhs.(*ArrowFn).arrowLoc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, "=>"))
		}

		bin := &BinExpr{N_EXPR_BIN, nil, T_ILLEGAL, nil, nil, nil}
		bin.loc = p.finLoc(lhs.Loc().Clone())
		bin.op = op
		bin.lhs = lhs
		bin.rhs = rhs
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
	obj, err := p.primaryExpr()
	if err != nil {
		return nil, err
	}
	if p.lexer.Peek().value == T_TPL_HEAD {
		return p.tplExpr(obj)
	}
	return obj, nil
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
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, true, false, nil}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node) (Node, error) {
	p.lexer.Next()

	loc := p.loc()
	tok := p.lexer.Next()
	_, kw, ok := tok.CanBePropKey()

	var prop Node
	if (ok && tok.value != T_NUM) || tok.value == T_NAME_PVT {
		prop = &Ident{N_NAME, p.finLoc(loc), tok.Text(), tok.value == T_NAME_PVT, tok.ContainsEscape(), nil, kw}
	} else {
		return nil, p.errorTok(tok)
	}

	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, false, nil}
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
			return p.fnDec(true, tok, false)
		}
		p.lexer.Next()

		name := tok.Text()
		ahead := p.lexer.Peek()
		// `ahead.value != T_ARROW` is used to skip checking name when it appears in the param list of arrow expr
		// for `eval => 42` we should report binding-reserved-word error instead of unexpected-reserved-word error
		if p.checkName && ahead.value != T_ARROW && !ahead.afterLineTerminator && p.isProhibitedName(nil, name, true) {
			if tok.ContainsEscape() {
				return nil, p.errorAtLoc(p.finLoc(loc), ERR_ESCAPE_IN_KEYWORD)
			}
			return nil, p.errorAtLoc(p.finLoc(loc), ERR_RESERVED_WORD_IN_STRICT_MODE)
		}
		return &Ident{N_NAME, p.finLoc(loc), name, false, tok.ContainsEscape(), nil, p.isProhibitedName(nil, name, true)}, nil
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
		return p.fnDec(true, nil, false)
	case T_REGEXP:
		p.lexer.Next()
		ext := tok.ext.(*TokExtRegexp)
		return &RegexpLit{N_LIT_REGEXP, p.finLoc(loc), tok.Text(), ext.Pattern(), ext.Flags(), nil}, nil
	case T_CLASS:
		return p.classDec(true, false)
	case T_SUPER:
		scope := p.scope()
		sup := p.lexer.Next()

		ahead := p.lexer.Peek()
		if !scope.IsKind(SPK_CLASS) && !scope.IsKind(SPK_CLASS_INDIRECT) && !scope.IsKind(SPK_METHOD) {
			em := ERR_SUPER_OUTSIDE_CLASS
			if ahead.value == T_PAREN_L {
				em = ERR_SUPER_CALL_OUTSIDE_CTOR
			}
			return nil, p.errorAtLoc(loc, em)
		}

		if ahead.value != T_DOT && ahead.value != T_PAREN_L {
			return nil, p.errorTok(sup)
		}
		return &Super{N_SUPER, p.finLoc(loc), nil}, nil
	case T_IMPORT:
		return p.importCall()
	case T_TPL_HEAD:
		return p.tplExpr(nil)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) arrowFn(loc *Loc, args []Node) (Node, error) {
	params, err := p.argsToParams(args, false)
	if err != nil {
		return nil, err
	}

	arrowLoc := p.locFromTok(p.lexer.Next())
	ps := p.scope()
	scope := p.symtab.EnterScope(true, true)

	paramNames, firstComplicated, err := p.collectNames(params)
	if err != nil {
		return nil, err
	}

	for _, paramName := range paramNames {
		ref := NewRef()
		ref.Node = paramName.(*Ident)
		ref.BindKind = BK_PARAM
		// duplicate-checking is enable in strict mode so here skip doing checking,
		// checking is delegated to below `checkParams`
		p.addLocalBinding(nil, ref, false)
	}

	var body Node
	scope.AddKind(SPK_ARROW)
	if p.lexer.Peek().value == T_BRACE_L {
		body, err = p.fnBody()
		if err != nil {
			return nil, err
		}

		if _, err := p.isExprOpening(true); err != nil {
			return nil, err
		}
	} else {
		if ps.IsKind(SPK_NOT_IN) {
			scope.AddKind(SPK_NOT_IN)
		}
		body, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)
	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}
	p.symtab.LeaveScope()

	return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), arrowLoc, false, params, body, body.Type() != N_STMT_BLOCK, nil}, nil
}

func (p *Parser) parenExpr() (Node, error) {
	loc := p.loc()

	// the fnExpr and/or arrowExpr can not be followed by a pair of parens:
	// `() => {}()` is illegal, therefor it should be encapsulated in parens
	// to become the well-known IIFE - `(() => {})()`
	//
	// however, that's legal if the bad expr described above directly appear
	// in parens, eg. `(() => {}())`
	//
	// for dealing with above situations, enter a new scope and flag it as paren
	// to let the nested states can tell if they are in parenExpr to judge whether
	// to raise the syntax-error exception or not
	scope := p.symtab.EnterScope(false, false)
	scope.AddKind(SPK_PAREN)
	p.checkName = false
	args, tailingComma, err := p.argList(false, false)
	p.checkName = true
	p.symtab.LeaveScope()

	if err != nil {
		return nil, err
	}

	// next is arrow-expression
	ahead := p.lexer.Peek()
	if ahead.value == T_ARROW && !ahead.afterLineTerminator {
		return p.arrowFn(loc, args)
	}

	// for report expr like: `(a,)`
	if tailingComma != nil {
		return nil, p.errorAtLoc(tailingComma, ERR_TRAILING_COMMA)
	}

	argsLen := len(args)
	if argsLen == 0 {
		return nil, p.errorAt(p.lexer.prtVal, &p.lexer.prtBegin, "")
	}

	if err := p.checkArgs(args, false, true); err != nil {
		return nil, err
	}

	if argsLen == 1 {
		return &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), args[0], nil}, nil
	}

	seqLoc := args[0].Loc().Clone()
	end := args[argsLen-1].Loc()
	seqLoc.rng.end = end.rng.end
	seqLoc.end = end.end.Clone()
	seq := &SeqExpr{N_EXPR_SEQ, seqLoc, args, nil}
	return &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), seq, nil}, nil
}

func (p *Parser) UnParen(expr Node) Node {
	if expr.Type() == N_EXPR_PAREN {
		loc := expr.Loc().Clone()
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
	return p.assignExpr(true)
}

func (p *Parser) spread() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()

	if p.feat&FEAT_SPREAD == 0 {
		return nil, p.errorTok(tok)
	}

	node, err := p.assignExpr(true)
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
	tok = p.lexer.Peek()
	if tok.value == T_COMMA {
		trailingCommaLoc = p.locFromTok(tok)
	}
	return &Spread{N_SPREAD, p.finLoc(loc), node, trailingCommaLoc, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ObjectLiteral
func (p *Parser) objLit() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	props := make([]Node, 0, 1)
	// hasProto := false
	p.checkName = false
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
			return nil, p.errorTok(tok)
		}
	}
	p.checkName = true
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props, nil}, nil
}

func (p *Parser) objProp() (Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_DOT_TRI {
		return p.spread()
	}

	if tok.value == T_MUL {
		if p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(tok)
		}
		return p.method(nil, nil, nil, false, PK_INIT, true, false, false, false, false)
	} else if p.aheadIsAsync(tok, true, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(nil, nil, nil, false, PK_INIT, false, true, false, false, false)
	}
	return p.propData()
}

func (p *Parser) propData() (Node, error) {
	key, compute, err := p.propName(false, true)

	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
	}

	loc := compute
	if loc == nil {
		loc = key.Loc().Clone()
	}

	var value Node
	tok := p.lexer.Peek()
	opLoc := p.locFromTok(tok)
	assign := tok.value == T_ASSIGN
	if tok.value == T_COLON || assign {
		p.lexer.Next()
		value, err = p.assignExpr(true)
		if err != nil {
			return nil, err
		}
	} else if tok.value == T_PAREN_L {
		return p.method(loc, key, compute, false, PK_INIT, false, false, false, false, false)
	} else if compute != nil {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_COMPUTE_PROP_MISSING_INIT)
	}

	shorthand := assign
	if value == nil && key.Type() == N_NAME {
		id := key.(*Ident)
		name := id.Text()
		if id.kw && name != "eval" && name != "arguments" {
			return nil, p.errorAtLoc(id.loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.Text()))
		}
		shorthand = true
		value = key
	}
	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, compute != nil, false, shorthand, assign, PK_INIT}, nil
}

func (p *Parser) method(loc *Loc, key Node, compute *Loc, shorthand bool, kind PropKind,
	gen bool, async bool, allowNamePVT bool, inclass bool, static bool) (Node, error) {

	if !inclass && compute != nil {
		loc = compute
	}
	if loc == nil {
		loc = p.loc()
	}

	scope := p.symtab.EnterScope(true, false)
	scope.AddKind(SPK_METHOD)
	if kind == PK_CTOR && !static {
		scope.AddKind(SPK_CTOR)
	}

	// depart `gen` and `async` here since below stmt is legal:
	// `class a{ async *a() {} }`
	if async {
		p.lexer.Next()
		scope.AddKind(SPK_ASYNC)

		ahead := p.lexer.Peek()
		gen = ahead.value == T_MUL
		if gen && p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(ahead)
		}
	}
	if gen {
		scope.AddKind(SPK_GENERATOR)
		p.lexer.Next()
	}

	var err error
	if key == nil {
		if inclass {
			key, compute, err = p.classElemName()
		} else {
			key, compute, err = p.propName(allowNamePVT, false)
		}
		if err != nil {
			return nil, err
		}
	}

	if p.isName(key, "constructor", false) {
		if kind == PK_GETTER || kind == PK_SETTER {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_HAVE_MODIFIER)
		} else if async {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_BE_ASYNC)
		} else if gen {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_BE_GENERATOR)
		}
	}

	fnLoc := p.loc()
	p.checkName = false
	args, _, err := p.argList(false, false)
	p.checkName = true
	if err != nil {
		return nil, err
	}

	params, err := p.argsToParams(args, kind == PK_SETTER)
	if err != nil {
		return nil, err
	}

	paramNames, firstComplicated, err := p.collectNames(params)
	if err != nil {
		return nil, err
	}

	for _, paramName := range paramNames {
		ref := NewRef()
		ref.Node = paramName.(*Ident)
		ref.BindKind = BK_PARAM
		// duplicate-checking is enable in strict mode so here skip doing checking,
		// checking is delegated to below `checkParams`
		p.addLocalBinding(nil, ref, false)
	}

	if kind == PK_GETTER && len(params) > 0 {
		return nil, p.errorAtLoc(params[0].Loc(), ERR_GETTER_SHOULD_NO_PARAM)
	} else if kind == PK_SETTER && len(params) != 1 {
		return nil, p.errorAtLoc(fnLoc, ERR_SETTER_SHOULD_ONE_PARAM)
	}

	if gen {
		p.lexer.addMode(LM_GENERATOR)
	}

	body, err := p.fnBody()
	if gen {
		p.lexer.popMode()
	}
	if err != nil {
		return nil, err
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)

	// `isProhibitedName` is not needed here since `keyword` as method name is permitted
	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()

	value := &FnDec{N_EXPR_FN, p.finLoc(fnLoc), nil, gen, async, params, body, nil}
	if inclass {
		if static && p.isName(key, "prototype", false) {
			return nil, p.errorAtLoc(key.Loc(), ERR_STATIC_PROP_PROTOTYPE)
		}

		return &Method{N_METHOD, p.finLoc(loc), key, static, compute != nil, kind, value}, nil
	}
	return &Prop{N_PROP, p.finLoc(loc), key, nil, value, compute != nil, true, shorthand, false, kind}, nil
}

func (p *Parser) isExprOpening(raise bool) (*Token, error) {
	tok := p.lexer.Peek()
	tv := tok.value
	if raise && tv != T_SEMI && tv != T_BRACE_R && tv != T_COMMA && tv != T_PAREN_R && !tok.afterLineTerminator && tv != T_EOF {
		errMsg := ERR_UNEXPECTED_TOKEN
		if tok.value == T_ILLEGAL {
			if msg, ok := tok.ext.(string); ok {
				errMsg = msg
			} else if msg, ok := tok.ext.(*LexerError); ok {
				errMsg = msg.Error()
			}
		}
		return nil, p.errorAt(tok.value, &tok.begin, errMsg)
	}
	return tok, nil
}

func (p *Parser) advanceIfSemi(raise bool) error {
	tok, err := p.isExprOpening(raise)
	if err != nil {
		return err
	}

	if tok.value == T_SEMI {
		p.lexer.Next()
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
	return p.lexer.finLoc(loc)
}

func (p *Parser) errorTok(tok *Token) *ParserError {
	if tok.value != T_ILLEGAL {
		return NewParserError(fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok.value].Name),
			p.lexer.src.path, tok.begin.line, tok.begin.col)
	}
	return NewParserError(tok.ErrMsg(), p.lexer.src.path, tok.begin.line, tok.begin.col)
}

func (p *Parser) errorAt(tok TokenValue, pos *Pos, errMsg string) *ParserError {
	if tok != T_ILLEGAL && errMsg == "" {
		return NewParserError(fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok].Name),
			p.lexer.src.path, pos.line, pos.col)
	}
	return NewParserError(errMsg, p.lexer.src.path, pos.line, pos.col)
}

func (p *Parser) errorAtLoc(loc *Loc, errMsg string) *ParserError {
	return NewParserError(errMsg, p.lexer.src.path, loc.begin.line, loc.begin.col)
}
