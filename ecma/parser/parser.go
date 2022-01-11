package parser

import (
	"errors"
	"fmt"
)

type Parser struct {
	lexer           *Lexer
	symtab          *SymTab
	feat            Feature
	imp             map[string]*Ident
	exp             []*ExportDec
	checkName       bool
	danglingPvtRefs []*Ref

	ts bool

	// the ts func sig cannot stand alone:
	// `function f(a:number)` is illegal unless it's followed by a
	// func dec with the same id:
	// ```
	// function f(a:number)
	// function f(): any {}
	// ```
	lastTsFnSig *FnDec

	// stores the `<` tokens which are resoled as `LT` operator and
	// identified by their `[line, column]`
	ltTokens [][2]int
}

type ParserOpts struct {
	Externals []string
	Version   ESVersion
	Feature   Feature
}

const defaultFeatures Feature = FEAT_MODULE | FEAT_GLOBAL_ASYNC | FEAT_STRICT | FEAT_LET_CONST |
	FEAT_BINDING_PATTERN | FEAT_BINDING_REST_ELEM | FEAT_BINDING_REST_ELEM_NESTED |
	FEAT_SPREAD | FEAT_META_PROPERTY | FEAT_ASYNC_AWAIT | FEAT_ASYNC_ITERATION | FEAT_ASYNC_GENERATOR |
	FEAT_POW | FEAT_CLASS_PRV | FEAT_CLASS_PUB_FIELD | FEAT_CLASS_PRIV_FIELD | FEAT_OPT_EXPR | FEAT_OPT_CATCH_PARAM |
	FEAT_NULLISH | FEAT_BAD_ESCAPE_IN_TAGGED_TPL | FEAT_BIGINT | FEAT_NUM_SEP | FEAT_LOGIC_ASSIGN |
	FEAT_DYNAMIC_IMPORT | FEAT_JSON_SUPER_SET | FEAT_EXPORT_ALL_AS_NS | FEAT_JSX

func NewParserOpts() *ParserOpts {
	return &ParserOpts{
		Externals: make([]string, 0),
		Feature:   defaultFeatures,
	}
}

func (o *ParserOpts) Clone() *ParserOpts {
	return &ParserOpts{
		Externals: o.Externals,
		Version:   o.Version,
		Feature:   o.Feature,
	}
}

func (o *ParserOpts) MergeJson(obj map[string]interface{}) {
	if moduleType, ok := obj["sourceType"]; ok {
		if moduleType == "module" {
			o.Feature = o.Feature.On(FEAT_MODULE)
		}
	}
	if ts, ok := obj["typescript"]; ok && ts == true {
		o.Feature = o.Feature.On(FEAT_TS)
	}
	if jsx, ok := obj["jsx"]; ok && jsx == true {
		o.Feature = o.Feature.On(FEAT_JSX)
	}
}

func NewParser(src *Source, opts *ParserOpts) *Parser {
	parser := &Parser{}
	parser.Setup(src, opts)
	return parser
}

func (p *Parser) Setup(src *Source, opts *ParserOpts) {
	if opts.Feature&FEAT_ASYNC_AWAIT == 0 {
		opts.Feature = opts.Feature.Off(FEAT_GLOBAL_ASYNC)
	}
	if opts.Feature&FEAT_MODULE != 0 {
		opts.Feature = opts.Feature.On(FEAT_IMPORT_DEC).On(FEAT_EXPORT_DEC)
	}

	p.feat = opts.Feature
	p.imp = map[string]*Ident{}
	p.exp = []*ExportDec{}
	p.checkName = true
	p.danglingPvtRefs = make([]*Ref, 0)

	p.symtab = NewSymTab(opts.Externals)

	p.lexer = NewLexer(src)
	p.lexer.ver = opts.Version
	p.lexer.feat = opts.Feature

	p.ts = p.feat&FEAT_TS != 0
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

	if err := p.resolvingDanglingPvtRefs(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Parser) resolvingDanglingPvtRefs() error {
	for _, ref := range p.danglingPvtRefs {
		if ref.TargetType == TT_PVT_FIELD {
			name := "#" + ref.Node.Text()
			target := ref.Scope.BindingOf(name)
			if target != nil {
				target.RetainBy(ref)
			} else {
				return p.errorAtLoc(ref.Node.Loc(), fmt.Sprintf(ERR_TPL_ALONE_PVT_FIELD, name))
			}
		}
	}
	p.danglingPvtRefs = nil
	return nil
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
			subnames = []Node{&Ident{N_NAME, exp.def, "default", false, false, nil, true, p.newTypInfo()}}
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
			node, err = p.fnDec(false, nil, false, false)
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
				node, err = p.classDec(false, false, false, false)
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
		if p.aheadIsTsEnum(tok) {
			node, err = p.tsEnum(nil)
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
		node, err = p.fnDec(false, tok, false, false)
	} else if p.aheadIsLabel(tok) {
		node, err = p.labelStmt()
	} else if p.aheadIsTsTypDec(tok) {
		node, err = p.tsTypDec()
	} else if p.aheadIsTsItf(tok) {
		node, err = p.tsItf()
	} else if p.aheadIsTsNS(tok) {
		node, err = p.tsNS()
	} else if p.aheadIsTsDec(tok) {
		node, err = p.tsDec()
	} else if p.tsAheadIsAbstract(tok, false, false) {
		node, err = p.classDec(false, false, false, true)
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
	if p.feat&FEAT_EXPORT_DEC == 0 {
		return nil, p.errorTok(tok)
	}

	var err error
	node := &ExportDec{N_STMT_EXPORT, nil, false, nil, nil, nil, nil}
	specs := make([]Node, 0)
	tok = p.lexer.Peek()
	tv := tok.value
	if tv == T_MUL || tv == T_BRACE_L {
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
	} else if tv == T_FUNC {
		node.dec, err = p.fnDec(false, nil, false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tv, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		node.dec, err = p.fnDec(false, tok, false, false)
		if err != nil {
			return nil, err
		}
	} else if tv == T_CLASS {
		node.dec, err = p.classDec(false, false, false, false)
		if err != nil {
			return nil, err
		}
	} else if tv == T_DEFAULT {
		def := p.lexer.Next()
		tok := p.lexer.Peek()
		tv = tok.value
		node.def = p.locFromTok(def)
		if tv == T_FUNC {
			node.dec, err = p.fnDec(false, nil, true, false)
		} else if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tv, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			node.dec, err = p.fnDec(false, tok, true, false)
		} else if tv == T_CLASS {
			node.dec, err = p.classDec(false, true, false, false)
		} else if p.tsAheadIsAbstract(tok, false, false) {
			node.dec, err = p.classDec(false, true, false, true)
		} else {
			node.dec, err = p.assignExpr(true, false, false)
			if err := p.advanceIfSemi(false); err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	} else if p.ts && tv == T_IMPORT {
		loc := p.locFromTok(p.lexer.Next())
		id, err := p.ident(nil, true)
		if err != nil {
			return nil, err
		}
		n, err := p.tsImportAlias(loc, id, true)
		if err != nil {
			return nil, err
		}
		if n.Type() == N_TS_IMPORT_ALIAS {
			return n, nil
		}
		node.dec = n
	} else if p.tsAheadIsAbstract(tok, false, false) {
		node.dec, err = p.classDec(false, true, false, true)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsItf(tok) {
		node.dec, err = p.tsItf()
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsTypDec(tok) {
		node.dec, err = p.tsTypDec()
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsEnum(tok) {
		node.dec, err = p.tsEnum(nil)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsNS(tok) {
		node.dec, err = p.tsNS()
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsDec(tok) {
		node.dec, err = p.tsDec()
		if err != nil {
			return nil, err
		}
	} else if p.ts && tv == T_ASSIGN {
		p.lexer.Next()
		// ExportAssignment: `export = a`
		id, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		loc := p.locFromTok(tok)
		return &TsExportAssign{N_TS_EXPORT_ASSIGN, p.finLoc(loc), id}, nil
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
		if ahead.value == T_NAME && ahead.Text() == "as" && p.feat&FEAT_EXPORT_ALL_AS_NS != 0 {
			p.lexer.Next()

			id, err := p.ident(nil, false)
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
		src = &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), str.HasLegacyOctalEscapeSeq(), nil, nil}
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

func (p *Parser) identWithKw(scope *Scope, binding bool) (Node, error) {
	ahead := p.lexer.Peek()
	if ahead.IsKw() {
		p.lexer.Next()
		str := TokenKinds[ahead.value].Name
		return &Ident{N_NAME, p.finLoc(p.locFromTok(ahead)), str, false, false, nil, true, p.newTypInfo()}, nil
	}
	return p.ident(scope, binding)
}

func (p *Parser) exportSpec() (Node, error) {
	loc := p.loc()
	local, err := p.identWithKw(nil, false)
	if err != nil {
		return nil, err
	}

	id := local
	if p.aheadIsName("as") {
		tok := p.lexer.Next()
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		id, err = p.identWithKw(nil, false)
		if err != nil {
			return nil, err
		}
	}

	return &ExportSpec{N_EXPORT_SPEC, p.finLoc(loc), false, local, id}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ImportDeclaration
func (p *Parser) importDec() (Node, error) {
	loc := p.loc()
	ipt := p.lexer.Next()
	if p.feat&FEAT_IMPORT_DEC == 0 && p.feat&FEAT_DYNAMIC_IMPORT == 0 {
		return nil, p.errorTok(ipt)
	}

	specs := make([]Node, 0)
	tok := p.lexer.Peek()
	if tok.value != T_STRING {
		var id Node
		var err error
		if tok.value == T_NAME {
			id, err = p.ident(nil, true)
			if err != nil {
				return nil, err
			}
			spec := &ImportSpec{N_IMPORT_SPEC, p.finLoc(p.locFromTok(tok)), true, false, id, id}
			specs = append(specs, spec)
		} else if tok.value == T_PAREN_L || tok.value == T_DOT {
			expr, err := p.importCall(ipt)
			if err != nil {
				return nil, err
			}
			return &ExprStmt{N_STMT_EXPR, expr.Loc().Clone(), expr, false}, nil
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

		if p.ts && p.lexer.Peek().value == T_ASSIGN && id != nil {
			return p.tsImportAlias(p.locFromTok(tok), id, false)
		}

		_, err = p.nextMustName("from", false)
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
	src := &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(str)), str.Text(), legacyOctalEscapeSeq, nil, nil}

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
	binding, err := p.identWithKw(nil, true)
	if err != nil {
		return nil, err
	}

	id := binding
	if p.aheadIsName("as") {
		p.lexer.Next()
		binding, err = p.ident(nil, true)
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

	id, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}

	specs := make([]Node, 1)
	specs[0] = &ImportSpec{N_IMPORT_SPEC, p.finLoc(loc), false, true, id, nil}
	return specs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec(expr bool, canNameOmitted bool, declare bool, abstract bool) (Node, error) {
	loc := p.locFromTok(p.lexer.Next())

	if abstract {
		p.lexer.Next()
	}

	ps := p.scope()
	// all parts of the class dec are in strict mode(include the id part)
	// here push an intermidate mode as strict to handle the id part
	p.lexer.pushMode(LM_STRICT, true)

	var id Node
	var err error
	ahead := p.lexer.Peek()
	if ahead.value != T_BRACE_L && ahead.value != T_EXTENDS {
		id, err = p.identStrict(ps, true, true, false)
		if err != nil {
			return nil, err
		}
		ti := p.newTypInfo()
		if ti != nil {
			typParams, err := p.tsTryTypParams()
			if err != nil {
				return nil, err
			}
			ti.typParams = typParams
			id.(NodeWithTypInfo).SetTypInfo(ti)
		}
		ref := NewRef()
		ref.Node = id.(*Ident)
		ref.BindKind = BK_CONST
		if err := p.addLocalBinding(ps, ref, true, ref.Node.Text()); err != nil {
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
		if err := p.checkDefaultVal(super, false, false, false); err != nil {
			return nil, err
		}
	}

	scope := p.symtab.EnterScope(true, false)
	p.enterStrict(true).AddKind(SPK_CLASS)
	if super != nil {
		scope.AddKind(SPK_CLASS_HAS_SUPER)
	}

	body, err := p.classBody(declare)
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
	return &ClassDec{typ, p.finLoc(loc), id, super, body, abstract}, nil
}

func (p *Parser) classBody(declare bool) (Node, error) {
	loc := p.loc()

	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}

	elems := make([]Node, 0, 3)
	hasCtor := false
	pvtNames := make(map[string]Node)
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
		elem, err := p.classElem(declare)
		if err != nil {
			return nil, err
		}
		if elem.Type() == N_METHOD {
			m := elem.(*Method)
			if hasCtor {
				return nil, p.errorAtLoc(m.key.Loc(), ERR_CTOR_DUP)
			}
			if p.isName(m.key, "constructor", false, true) {
				hasCtor = true
			}
		}
		if name, key, pvt := p.nameOfProp(elem); pvt {
			dup := true
			if prev, ok := pvtNames[name]; ok {
				// a pair of getter/setter is valid
				// prev is nil means a pair of getter/setter already occurred therefore the
				// node is eliminated
				if prev != nil && prev.Type() == N_METHOD && elem.Type() == N_METHOD {
					a := elem.(*Method)
					b := prev.(*Method)
					if (a.kind == PK_SETTER || a.kind == PK_GETTER) && (b.kind == PK_SETTER || b.kind == PK_GETTER) {
						dup = a.static != b.static ||
							(a.kind == PK_GETTER && b.kind == PK_GETTER || a.kind == PK_SETTER && b.kind == PK_SETTER)
					}
				}
				if dup {
					return nil, p.errorAtLoc(key.Loc(), fmt.Sprintf(ERR_TPL_ID_ALREADY_DEF, name))
				}
			}
			pvtNames[name] = elem

			ref := NewRef()
			ref.Node = key.(*Ident)
			ref.TargetType = TT_PVT_FIELD
			ref.BindKind = BK_PVT_FIELD
			// skip check dup since getter/setter is dup but legal
			if err := p.addLocalBinding(nil, ref, false, name); err != nil {
				return nil, err
			}
		}
		elems = append(elems, elem)

	}

	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}

	return &ClassBody{N_ClASS_BODY, p.finLoc(loc), elems}, nil
}

func (p *Parser) classModifier() (begin, static, access, abstract *Loc, isField, escape bool, name string, fieldLoc *Loc, accMod ACC_MOD, ahead *Token) {
	for {
		ahead = p.lexer.Peek()
		av := ahead.value
		if av == T_STATIC && static == nil {
			tok := p.lexer.Next()
			static = p.locFromTok(tok)
			if begin == nil {
				begin = static.Clone()
			}
			escape = tok.ContainsEscape()
			name = tok.Text()
			isField, ahead = p.isField(true, false)
			fieldLoc = static
			if isField {
				static = nil
				return
			}
			fieldLoc = static.Clone()
		} else if p.ts && accMod == ACC_MOD_NONE && (av == T_PUBLIC || av == T_PRIVATE || av == T_PROTECTED) {
			switch av {
			case T_PUBLIC:
				accMod = ACC_MOD_PUB
			case T_PRIVATE:
				accMod = ACC_MOD_PRI
			case T_PROTECTED:
				accMod = ACC_MOD_PRO
			}
			if accMod != ACC_MOD_NONE {
				tok := p.lexer.Next()
				access = p.locFromTok(tok)
				if begin == nil {
					begin = access.Clone()
				}
				escape = tok.ContainsEscape()
				name = tok.Text()
				isField, ahead = p.isField(true, false)
				fieldLoc = access
				if isField {
					access = nil
					return
				}
				fieldLoc = access.Clone()
			}
		} else if p.ts && abstract == nil && IsName(ahead, "abstract", false) {
			tok := p.lexer.Next()
			abstract = p.locFromTok(tok)
			if begin == nil {
				begin = abstract.Clone()
			}
			escape = tok.ContainsEscape()
			name = tok.Text()
			isField, ahead = p.isField(true, false)
			fieldLoc = abstract
			if isField {
				abstract = nil
				return
			}
			fieldLoc = abstract.Clone()
		} else {
			break
		}
	}
	return
}

func (p *Parser) classElem(declare bool) (Node, error) {
	beginLoc, staticLoc, accessLoc, abstractLoc, isField, escape, fieldName, fieldLoc, accMod, ahead := p.classModifier()

	static := staticLoc != nil
	abstract := abstractLoc != nil
	if static || accessLoc != nil || abstractLoc != nil {
		if static && abstract {
			return nil, p.errorAtLoc(abstractLoc, ERR_ABSTRACT_MIXED_WITH_STATIC)
		}

		if ahead.value == T_BRACE_L {
			return p.staticBlock(beginLoc)
		}
		if p.aheadIsArgList(ahead) {
			key := &Ident{N_NAME, fieldLoc, fieldName, false, escape, nil, true, p.newTypInfo()}
			return p.method(beginLoc, key, accMod, nil, false, PK_METHOD, false, false, false, true, false, declare, abstract)
		}
	} else if isField {
		key := &Ident{N_NAME, fieldLoc, fieldName, false, escape, nil, true, p.newTypInfo()}
		return p.field(key, nil, accMod, false)
	}

	if p.aheadIsAsync(ahead, true, true) {
		if ahead.ContainsEscape() {
			return nil, p.errorAt(ahead.value, &ahead.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(beginLoc, nil, accMod, nil, false, PK_METHOD, false, true, true, true, static, declare, abstract)
	}
	if ahead.value == T_MUL {
		if p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(ahead)
		}
		return p.method(beginLoc, nil, accMod, nil, false, PK_METHOD, true, false, true, true, static, declare, abstract)
	}

	propLoc := p.locFromTok(ahead)
	kw := ahead.IsKw()
	if ahead.value == T_NAME || ahead.value == T_STRING || kw {
		p.lexer.Next()

		name := ahead.Text()

		var key Node
		if ahead.value == T_STRING {
			key = &StrLit{N_LIT_STR, p.finLoc(propLoc.Clone()), name, ahead.HasLegacyOctalEscapeSeq(), nil, nil}
		} else {
			key = &Ident{N_NAME, p.finLoc(propLoc.Clone()), name, false, ahead.ContainsEscape(), nil, kw, p.newTypInfo()}
		}

		isField, ahead = p.isField(false, name == "get" || name == "set")
		if isField {
			return p.field(key, beginLoc, accMod, abstract)
		}

		if static || abstract || accMod != ACC_MOD_NONE {
			propLoc = beginLoc
		}

		kd := PK_INIT
		if p.aheadIsArgList(ahead) {
			kd = PK_METHOD
			if name == "constructor" {
				kd = PK_CTOR
			}
			return p.method(propLoc, key, accMod, nil, false, kd, false, false, true, true, static, declare, abstract)
		}

		if name == "get" {
			kd = PK_GETTER
		} else if name == "set" {
			kd = PK_SETTER
		} else {
			return nil, p.errorTok(ahead)
		}

		return p.method(propLoc, nil, accMod, nil, false, kd, false, false, true, true, static, declare, abstract)
	}

	if declare && ahead.value == T_BRACKET_L {
		return p.tsIdxSig(nil)
	}

	return p.field(nil, staticLoc, accMod, abstract)
}

func (p *Parser) isName(node Node, name string, canContainsEscape bool, str bool) bool {
	nv := node.Type()
	hasEscape := false
	ns := ""
	if nv == N_LIT_STR && str {
		s := node.(*StrLit)
		ns = s.Text()
	} else if nv == N_NAME {
		id := node.(*Ident)
		ns = id.Text()
		hasEscape = id.containsEscape
	}

	if ns != name {
		return false
	}
	if !canContainsEscape {
		return !hasEscape
	}
	return true
}

func (p *Parser) field(key Node, static *Loc, accMode ACC_MOD, abstract bool) (Node, error) {
	var loc *Loc
	var err error
	var compute *Loc
	if key == nil {
		key, compute, err = p.classElemName()
		if err != nil {
			return nil, err
		}
	}
	if static != nil {
		loc = static.Clone()
	} else if compute != nil {
		loc = compute.Clone()
	} else {
		loc = key.Loc().Clone()
	}

	ti := p.newTypInfo()
	if ti != nil {
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		ti.typAnnot = typAnnot
		key.(NodeWithTypInfo).SetTypInfo(ti)
	}

	var value Node
	tok := p.lexer.Peek()
	if tok.value == T_ASSIGN {
		p.lexer.Next()
		value, err = p.assignExpr(true, false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsArgList(tok) {
		if static != nil {
			loc = static
		}
		return p.method(loc, key, accMode, compute, false, PK_METHOD, false, false, true, true, static != nil, false, abstract)
	}
	p.advanceIfSemi(false)

	if value != nil {
		p.checkName = true
		if err := p.checkDefaultVal(value, false, false, true); err != nil {
			return nil, err
		}
		p.checkName = false
	}

	staticField := static != nil
	if staticField && p.isName(key, "prototype", false, true) {
		return nil, p.errorAtLoc(key.Loc(), ERR_STATIC_PROP_PROTOTYPE)
	} else if p.isName(key, "constructor", false, true) {
		return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_BE_Field)
	}

	return &Field{N_FIELD, p.finLoc(loc), key, staticField, compute != nil, value, accMode, abstract}, nil
}

func (p *Parser) classElemName() (Node, *Loc, error) {
	return p.propName(true, false, false)
}

func (p *Parser) staticBlock(static *Loc) (Node, error) {
	block, err := p.blockStmt(true)
	if err != nil {
		return nil, err
	}
	return &StaticBlock{N_STATIC_BLOCK, p.finLoc(static), block.body}, nil
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

		ahead := p.lexer.Peek()
		var param Node
		if ahead.value == T_PAREN_L {
			p.lexer.Next()

			param, err = p.bindingPattern()
			if err != nil {
				return nil, err
			}

			typAnnot, err := p.tsTypAnnot()
			if err != nil {
				return nil, err
			}
			p.tsNodeTypAnnot(param, typAnnot, ACC_MOD_NONE, nil)

			if _, err := p.nextMustTok(T_PAREN_R); err != nil {
				return nil, err
			}
		} else if p.feat&FEAT_OPT_CATCH_PARAM == 0 {
			return nil, p.errorTok(ahead)
		}

		scope := p.symtab.EnterScope(false, false)
		scope.AddKind(SPK_CATCH)

		if param != nil {
			names, _, _ := p.collectNames([]Node{param})
			for _, nameNode := range names {
				id := nameNode.(*Ident)
				if ok := p.isProhibitedName(nil, id.val, true, true, false, false); ok {
					return nil, p.errorAtLoc(id.loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
				}
				ref := NewRef()
				ref.Node = id
				ref.BindKind = BK_LET
				if err := p.addLocalBinding(nil, ref, true, id.Text()); err != nil {
					return nil, err
				}
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
	label, err := p.ident(nil, false)
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
		label, err := p.ident(nil, false)
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
		label, err := p.ident(nil, false)
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

	cases := make([]Node, 0)
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
				if init, err = p.argToParam(init, 0, false, true, false); err != nil {
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
					d := dec.(*VarDec)
					if d.init != nil {
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
		p.lexer.NextAndRevise(revise)

		right, err := p.assignExpr(true, false, false)
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
			(p.aheadIsArgList(ahead) && !prop) ||
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
func (p *Parser) fnDec(expr bool, async *Token, canNameOmitted bool, declare bool) (Node, error) {
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
	}
	tok := p.lexer.Peek()
	fn := false
	if tok.value == T_FUNC {
		fn = true
		p.lexer.Next()
	}

	tok = p.lexer.Peek()
	generator := tok.value == T_MUL
	genLoc := p.locFromTok(tok)
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
	var fnRef *Ref
	if tok.value != T_PAREN_L && tok.value != T_LT {
		id, err = p.ident(idScope, true)
		if err != nil {
			return nil, err
		}

		// name of the function expression will not add a ref record
		if !expr {
			fnRef = NewRef()
			fnRef.Node = id.(*Ident)
			fnRef.BindKind = BK_VAR
			fnRef.TargetType = TT_FN
			if ps.IsKind(SPK_STRICT) {
				fnRef.BindKind = BK_LET
			}
			if err := p.addLocalBinding(ps, fnRef, ps.IsKind(SPK_STRICT), fnRef.Node.Text()); err != nil {
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
		p.lexer.addMode(LM_ASYNC)
	}
	// 'yield' as function names
	if generator {
		p.scope().AddKind(SPK_GENERATOR)
	}

	var args []Node
	var typArgs Node
	ahead := p.lexer.Peek()
	if id != nil && ahead.value == T_ARROW && !ahead.afterLineTerminator {
		// async a => {}
		args = make([]Node, 1)
		args[0] = id
	} else if fn {
		args, typArgs, _, err = p.paramList(false, false, true)
		if err != nil {
			return nil, err
		}
	} else {
		// the arg check is skipped here, the correctness of args is guaranteed by
		// below `argsToFormalParams`
		p.checkName = false
		scope.AddKind(SPK_FORMAL_PARAMS)
		args, _, typArgs, err = p.argList(false, false, asyncLoc)
		scope.EraseKind(SPK_FORMAL_PARAMS)
		p.checkName = true
		if err != nil {
			return nil, err
		}
		if typArgs != nil && typArgs.Type() != N_TS_PARAM_INST {
			if err := p.advanceIfSemi(true); err != nil {
				return nil, err
			}
			return &ExprStmt{N_STMT_EXPR, p.finLoc(loc.Clone()), typArgs, false}, nil
		}
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	ti := p.newTypInfo()
	if ti != nil {
		ti.typAnnot = typAnnot
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
		if fn {
			params = args
		} else {
			params, err = p.argsToParams(args)
			if err != nil {
				return nil, err
			}
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
			p.addLocalBinding(nil, ref, false, ref.Node.Text())
		}

		if tok.value == T_BRACE_L {
			if body, err = p.fnBody(); err != nil {
				return nil, err
			}
		} else if expr || arrow {
			if body, err = p.expr(); err != nil {
				return nil, err
			}
		} else if fn && !expr && tok.value == T_FUNC && tok.afterLineTerminator && p.ts {
			// ts func overloads:
			// `function f(a:number)`
			// `function f(): any {}`
			ps.DelLocal(fnRef) // suppress the dup-checking of binding name
			if err = p.tsIsFnSigValid(fnRef.Node.Text()); err != nil {
				return nil, err
			}
		} else if (tok.value == T_SEMI || tok.afterLineTerminator) && declare {
			// AmbientFunctionDeclaration
			// `declare function a();`
		} else {
			return nil, p.errorTok(tok)
		}
	} else if async != nil {
		// this branch means the input is callExpr like:
		// `async ({a: b = c})` callExpr
		// `async* ({a: b = c})` binExpr
		lhs := &Ident{N_NAME, asyncLoc, "async", false, asyncHasEscape, nil, true, p.newTypInfo()}

		var exp Node
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
			exp = &BinExpr{N_EXPR_BIN, p.finLoc(loc), T_MUL, genLoc, lhs, rhs, nil}
		} else {
			if err := p.checkArgs(args, false, true); err != nil {
				return nil, err
			}
			ti := p.newTypInfo()
			if ti != nil {
				// typArgs is produced by `argList` in this branch, so it's required
				// to do a typeParam to typeArg transformation
				if err = p.tsCheckTypArgs(typArgs); err != nil {
					return nil, err
				}
				ti.typArgs = typArgs
			}
			exp = &CallExpr{N_EXPR_CALL, p.finLoc(loc), lhs, args, false, nil, ti}
		}

		if !expr {
			binExpr, err := p.binExpr(exp, 0, false, false, false)
			if err != nil {
				return nil, err
			}

			seq, err := p.seqExpr(binExpr, false)
			if err != nil {
				return nil, err
			}
			if err = p.advanceIfSemi(true); err != nil {
				return nil, err
			}
			return &ExprStmt{N_STMT_EXPR, p.finLoc(loc.Clone()), seq, false}, nil
		}
		return exp, nil
	} else {
		return nil, p.errorTok(tok)
	}

	if generator {
		p.lexer.popMode()
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)
	if id != nil {
		name := id.(*Ident).Text()
		if p.isProhibitedName(idScope, name, isStrict, true, false, false) {
			return nil, p.errorAtLoc(id.Loc(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
		}
	}

	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()

	if ti != nil {
		if !fn {
			typArgs, err = p.tsTypArgsToTypParams(typArgs)
			if err != nil {
				return nil, err
			}
		}
		ti.typParams = typArgs
	}

	if arrow {
		fn := &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), arrowLoc, async != nil, params, body, body.Type() != N_STMT_BLOCK, nil, ti}
		if expr {
			return fn, nil
		}
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
		return &ExprStmt{N_STMT_EXPR, p.finLoc(loc.Clone()), fn, false}, nil
	}

	typ := N_STMT_FN
	if expr {
		typ = N_EXPR_FN
	}

	fnDec := &FnDec{typ, p.finLoc(loc), id, generator, async != nil, params, body, nil, ti}
	if !expr && p.ts {
		if body == nil {
			p.lastTsFnSig = fnDec
		} else if err = p.tsIsFnImplValid(id); err != nil {
			return nil, err
		}
	}
	return fnDec, nil
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
		if p.isProhibitedName(nil, name, isStrict, true, false, false) {
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

func (p *Parser) isDirective(stmt Node) (bool, bool) {
	if stmt.Type() != N_STMT_EXPR {
		return false, false
	}
	expr := stmt.(*ExprStmt).expr
	if expr.Type() == N_LIT_STR {
		str := expr.(*StrLit).Raw()
		if str == "\"use strict\"" || str == "'use strict'" {
			return true, true
		}
		return false, true
	}
	return false, false
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
			if prologue != -1 && (scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_GLOBAL)) {
				strict, dir := p.isDirective(stmt)
				if !dir {
					if prologue == 0 {
						prologue = -1
					}
				} else {
					stmt.(*ExprStmt).dir = true
				}

				if strict {
					// lexer will automatically pop it's mode when the `T_BRACE_R` is met
					// here we use `ahead.value != T_BRACE_R` to prevent accidentally change
					// the upper lexer mode
					p.enterStrict(ahead.value != T_BRACE_R)

					// lookbehind to check that exprs before the 'use strcit' directive
					if prologue > 0 {
						for i := 0; i < prologue; i++ {
							stmt := stmts[i]
							if stmt.Type() == N_STMT_EXPR {
								expr := stmts[i].(*ExprStmt).expr
								if expr.Type() == N_LIT_STR && expr.(*StrLit).legacyOctalEscapeSeq {
									return nil, p.errorAtLoc(expr.Loc(), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
								}
							}
						}
						prologue = -1
					}
				}
				if dir {
					prologue += 1
				}
			}
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

func (p *Parser) addLocalBinding(s *Scope, ref *Ref, checkDup bool, name string) error {
	if s == nil {
		s = p.scope()
	}
	ok := s.AddLocal(ref, name, checkDup)
	if ok {
		return nil
	}
	if !ok {
		return p.errorAtLoc(ref.Node.loc, fmt.Sprintf(ERR_ID_DUP_DEF, name))
	}
	return nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-VariableStatement
func (p *Parser) varDecStmt(kind TokenValue, asExpr bool) (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	node := &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]Node, 0, 1), nil}

	isConst := false
	node.kind = kind
	bindKind := BK_VAR
	if kind == T_LET {
		bindKind = BK_LET
	} else if kind == T_CONST {
		isConst = true
		bindKind = BK_CONST
	}

	if p.aheadIsTsEnum(nil) {
		return p.tsEnum(loc)
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
		if ok := p.isProhibitedName(nil, id.val, true, true, false, false); ok {
			return nil, p.errorAtLoc(id.loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
		}

		ref := NewRef()
		ref.Node = id
		ref.BindKind = bindKind
		if err := p.addLocalBinding(nil, ref, true, ref.Node.Text()); err != nil {
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

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.tsNodeTypAnnot(binding, typAnnot, ACC_MOD_NONE, nil)

	var init Node
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr(true, false, false)
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
	"yield":      true,
	"implements": true,
	"interface":  true,
	"let":        true,
	"const":      true,
	"package":    true,
	"protected":  true,
	"public":     true,
	"static":     true,
	"import":     true,
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#sec-identifiers-static-semantics-early-errors
func (p *Parser) isProhibitedName(scope *Scope, name string, withStrict bool, lVal bool, field bool, forceStrict bool) bool {
	if scope == nil {
		scope = p.scope()
	}
	strict := withStrict && (scope.IsKind(SPK_STRICT) || forceStrict)
	_, ok := prohibitedNames[name]
	if strict && ok {
		return true
	}
	if (strict && lVal || field) && (name == "eval" || name == "arguments") {
		return true
	}
	return scope.IsKind(SPK_ASYNC) && name == "await"
}

func (p *Parser) identStrict(scope *Scope, forceStrict bool, binding bool, jsx bool) (Node, error) {
	if scope == nil {
		scope = p.scope()
	}

	tok := p.lexer.Next()
	tv := tok.value
	if tv != T_NAME && !(tv > T_CTX_KEYWORD_BEGIN && tv < T_CTX_KEYWORD_END) {
		return nil, p.errorTok(tok)
	}

	name := tok.Text()
	loc := p.finLoc(p.locFromTok(tok))

	if p.isProhibitedName(scope, name, true, false, false, forceStrict) {
		if binding {
			return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, name))
		}
		return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
	}

	// for resporting `'let' is disallowed as a lexically bound name` for stmt like `let let`
	if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) && !tok.ContainsEscape() {
		if name == "let" || name == "const" {
			return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_FORBIDED_LEXICAL_NAME, name))
		}
	}

	if !jsx {
		return &Ident{N_NAME, loc, name, false, tok.ContainsEscape(), nil, tok.IsKw(), p.newTypInfo()}, nil
	}

	return &JsxIdent{N_JSX_ID, loc, name, nil}, nil
}

func (p *Parser) ident(scope *Scope, binding bool) (*Ident, error) {
	id, err := p.identStrict(scope, false, binding, false)
	if err != nil {
		return nil, err
	}
	return id.(*Ident), nil
}

func (p *Parser) accMod() (ACC_MOD, *Loc) {
	if !p.ts {
		return ACC_MOD_NONE, nil
	}

	var loc *Loc
	ahead := p.lexer.Peek()
	mod := ACC_MOD_NONE
	switch ahead.value {
	case T_PUBLIC:
		mod = ACC_MOD_PUB
		loc = p.locFromTok(p.lexer.Next())
	case T_PRIVATE:
		mod = ACC_MOD_PRI
		loc = p.locFromTok(p.lexer.Next())
	case T_PROTECTED:
		mod = ACC_MOD_PRO
		loc = p.locFromTok(p.lexer.Next())
	}
	return mod, loc
}

func (p *Parser) roughParam(ctor bool) (Node, error) {
	accMod, accLoc := p.accMod()
	if accLoc != nil && !ctor {
		return nil, p.errorAtLoc(accLoc, ERR_ILLEGAL_PARAMETER_MODIFIER)
	}

	name, err := p.tsTyp(true)
	if err != nil {
		return nil, err
	}

	ques := p.tsAdvanceHook()
	if ques != nil && name.Type() == N_TS_THIS {
		return nil, p.errorAtLoc(ques, ERR_THIS_CANNOT_BE_OPTIONAL)
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	ti := p.newTypInfo()
	ti.typAnnot = typAnnot
	ti.accMod = accMod
	ti.ques = ques
	var colonLoc *Loc
	if typAnnot != nil {
		colonLoc = typAnnot.Loc().Clone()
	}
	return &TsRoughParam{N_TS_ROUGH_PARAM, p.finLoc(name.Loc().Clone()), name, colonLoc, ti}, nil
}

func (p *Parser) param(ctor bool) (Node, error) {
	accMod, accLoc := p.accMod()
	if accLoc != nil && !ctor {
		return nil, p.errorAtLoc(accLoc, ERR_ILLEGAL_PARAMETER_MODIFIER)
	}

	binding, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}
	loc := binding.Loc().Clone()

	if ques := p.tsAdvanceHook(); ques != nil {
		if binding.Type() == N_NAME {
			binding.(*Ident).ti.ques = ques
		} else {
			return nil, p.errorAtLoc(ques, ERR_UNEXPECTED_TOKEN)
		}
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.tsNodeTypAnnot(binding, typAnnot, accMod, nil)

	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		value, err := p.assignExpr(true, false, false)
		if err != nil {
			return nil, err
		}

		if err = p.checkDefaultVal(value, false, false, true); err != nil {
			return nil, err
		}
		if err = p.checkArg(value, false, false); err != nil {
			return nil, err
		}

		if binding.Type() == N_PAT_REST {
			r := binding.(*RestPat)
			return nil, p.errorAtLoc(r.arg.Loc(), ERR_REST_CANNOT_SET_DEFAULT)
		}
		binding = &AssignPat{
			typ: N_PAT_ASSIGN,
			loc: p.finLoc(loc),
			lhs: binding,
			rhs: value,
		}
	}

	return binding, nil
}

// `ctor` indicates this method is called when processing the constructor method of class,
// in that case the access modifier is needed to be considered as long as TS is enabled
func (p *Parser) paramList(firstRough bool, ctor bool, typParams bool) ([]Node, Node, *Loc, error) {
	scope := p.scope()
	p.checkName = false
	scope.AddKind(SPK_FORMAL_PARAMS)

	var tp Node
	var err error
	if typParams {
		tp, err = p.tsTryTypParams()
		if err != nil {
			return nil, nil, nil, err
		}
		if ctor && tp != nil {
			return nil, nil, nil, p.errorAtLoc(tp.Loc(), ERR_CTOR_CANNOT_WITH_TYPE_PARAMS)
		}
	}

	parenL, err := p.nextMustTok(T_PAREN_L)
	if err != nil {
		return nil, nil, nil, err
	}
	parenLoc := p.locFromTok(parenL)

	params := make([]Node, 0)
	i := 0
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, nil, nil, p.errorTok(tok)
		}

		var param Node
		var err error
		if firstRough && i == 0 {
			param, err = p.roughParam(ctor)
		} else {
			param, err = p.param(ctor)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		params = append(params, param)
		i += 1

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_COMMA {
			tok := p.lexer.Next()
			ahead := p.lexer.Peek()
			if param.Type() == N_SPREAD {
				msg := ERR_REST_ELEM_MUST_LAST
				if ahead.value != T_PAREN_R {
					msg = ERR_REST_ELEM_MUST_LAST
				}
				return nil, nil, nil, p.errorAt(tok.value, &tok.begin, msg)
			}
		} else if av != T_PAREN_R {
			return nil, nil, nil, p.errorTok(ahead)
		}
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, nil, nil, err
	}

	scope.EraseKind(SPK_FORMAL_PARAMS)
	p.checkName = true
	return params, tp, parenLoc, nil
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
	} else if tv == T_DOT_TRI {
		binding, err = p.patternRest(p.feat&FEAT_BINDING_REST_ELEM_NESTED != 0, false)
	} else {
		binding, err = p.ident(nil, true)
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
		if node.Type() == N_PROP {
			prop := node.(*Prop)
			if prop.method {
				return nil, p.errorAtLoc(prop.loc, ERR_INVALID_DESTRUCTING_TARGET)
			}
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PAT_REST && tok.value != T_BRACE_R {
			if tok.value == T_COMMA {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_REST_ELEM_MUST_LAST)
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
	return &ObjPat{N_PAT_OBJ, p.finLoc(loc), props, nil, p.newTypInfo()}, nil
}

func (p *Parser) patternProp() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		binding, err := p.patternRest(false, false)
		if err != nil {
			return nil, err
		}
		return binding, nil
	}

	key, compute, err := p.propName(false, true, false)
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

	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, compute != nil, false, shorthand, assign, PK_INIT, ACC_MOD_NONE}, nil
}

// test whether current place is filed or method, this method just lookahead
// the caller should ensure current place can be property key, eg: `T_NAME` or `T_STATIC`
//
// `static` indicates whether current place is `static` or not
// `getter` indicates whether current place is `get/set` or not
func (p *Parser) isField(static bool, getter bool) (bool, *Token) {
	ahead := p.lexer.Peek()
	av := ahead.value
	if p.feat&FEAT_CLASS_PUB_FIELD == 0 {
		return false, ahead
	}

	isField := av == T_COLON ||
		av == T_ASSIGN ||
		av == T_SEMI ||
		av == T_COMMA ||
		av == T_BRACE_R

	if isField {
		return true, ahead
	}

	if getter {
		return !TokenKinds[av].StartExpr && av != T_NAME_PVT && av != T_BRACKET_L && !ahead.IsKw(), ahead
	}

	if static {
		return isField, ahead
	}

	return ahead.afterLineTerminator, ahead
}

func (p *Parser) propName(allowNamePVT bool, maybeMethod bool, tsRough bool) (Node, *Loc, error) {
	var key Node
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	keyName, kw, ok := tok.CanBePropKey()

	scope := p.scope()
	var computeLoc *Loc
	tv := tok.value
	if allowNamePVT && tv == T_NAME_PVT {
		key = &Ident{N_NAME, p.finLoc(loc), tok.Text(), true, tok.ContainsEscape(), nil, false, p.newTypInfo()}
	} else if tv == T_STRING {
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, nil, p.errorAtLoc(p.locFromTok(tok), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		key = &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil, nil}
	} else if tv == T_NUM {
		key = &NumLit{N_LIT_NUM, p.finLoc(loc), nil}
	} else if tv == T_BRACKET_L {
		computeLoc = p.locFromTok(tok)
		scope.AddKind(SPK_PROP_NAME)
		name, err := p.assignExpr(true, false, false)
		scope.EraseKind(SPK_PROP_NAME)
		if err != nil {
			return nil, nil, err
		}
		_, err = p.nextMustTok(T_BRACKET_R)
		if err != nil {
			return nil, nil, err
		}
		key = name
	} else if ok {
		if !kw && p.isProhibitedName(nil, keyName, true, false, false, false) {
			kw = true
		}
		// stmt `let { let } = {}` will raise error `let is disallowed as a lexically bound name` in sloppy mode
		if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) {
			if !tok.ContainsEscape() && (keyName == "let" || keyName == "const") {
				return nil, nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_FORBIDED_LEXICAL_NAME, keyName))
			}
		}
		key = &Ident{N_NAME, p.finLoc(loc), keyName, false, tok.ContainsEscape(), nil, kw, p.newTypInfo()}
	} else {
		return nil, nil, p.errorTok(tok)
	}

	getter := keyName == "get" || keyName == "set"
	isField, ahead := p.isField(false, getter)
	if isField || !maybeMethod {
		return key, computeLoc, nil
	}

	kd := PK_INIT
	loc = loc.Clone()
	if getter {
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

	m, err := p.method(loc, key, ACC_MOD_NONE, computeLoc, false, kd, false, false, false, false, false, false, false)
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
	return &ArrPat{N_PAT_ARRAY, p.finLoc(loc), elems, nil, p.newTypInfo()}, nil
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
		binding, err = p.patternRest(!asProp, false)
	} else {
		binding, err = p.ident(nil, true)
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
		init, err = p.assignExpr(true, false, false)
		if err != nil {
			return nil, err
		}

		if err = p.checkArg(init, true, false); err != nil {
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
	return &Prop{N_PROP, p.finLoc(loc.Clone()), val.lhs, opLoc, val, false, false, true, true, PK_INIT, ACC_MOD_NONE}, nil
}

// `arrPat` indicats whether `restExpr` is in array-pattern or not
func (p *Parser) patternRest(arrPat bool, allowNotLast bool) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()

	if p.feat&FEAT_BINDING_REST_ELEM == 0 {
		return nil, p.errorTok(tok)
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	if av != T_NAME && (!arrPat || av != T_BRACKET_L && av != T_BRACE_L) {
		if av == T_BRACKET_L || av == T_BRACE_L {
			return nil, p.errorAt(ahead.value, &ahead.begin, ERR_REST_ARG_NOT_BINDING_PATTERN)
		}
		return nil, p.errorTok(ahead)
	}

	arg, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}

	if !allowNotLast {
		tok = p.lexer.Peek()
		if tok.value == T_COMMA {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_REST_ELEM_MUST_LAST)
		}
	}

	rest := &RestPat{N_PAT_REST, p.finLoc(loc), arg, nil, nil}
	if p.ts {
		rest.hoistTypInfo()
	}
	return rest, nil
}

func (p *Parser) exprStmt() (Node, error) {
	loc := p.loc()
	stmt := &ExprStmt{N_STMT_EXPR, &Loc{}, nil, false}
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
	return p.seqExpr(nil, false)
}

func (p *Parser) seqExpr(expr Node, notGT bool) (Node, error) {
	loc := p.loc()

	var err error
	if expr == nil {
		expr, err = p.assignExpr(true, notGT, false)
		if err != nil {
			return nil, err
		}
		// reports the illegal typAnnot usage in expr like `[a:b];` and `[x?]`
		if err = p.checkArg(expr, false, true); err != nil {
			return nil, err
		}
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
			expr, err = p.assignExpr(true, notGT, false)
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
	tok := p.lexer.Next()

	if p.scope().IsKind(SPK_FORMAL_PARAMS) {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_YIELD_IN_FORMAL_PARAMS)
	}

	tok = p.lexer.Peek()
	kind := TokenKinds[tok.value]
	tv := tok.value
	startExpr := kind.StartExpr || p.feat&FEAT_JSX != 0 && tv == T_LT
	if tok.afterLineTerminator || !startExpr && tv != T_MUL {
		return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), false, nil, nil}, nil
	}

	delegate := false
	if p.lexer.Peek().value == T_MUL {
		p.lexer.Next()
		delegate = true
	}

	arg, err := p.assignExpr(true, false, false)
	if err != nil {
		return nil, err
	}
	return &YieldExpr{N_EXPR_YIELD, p.finLoc(loc), delegate, arg, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-AssignmentExpression
func (p *Parser) assignExpr(checkLhs bool, notGT bool, notHook bool) (Node, error) {
	if p.aheadIsYield() {
		return p.yieldExpr()
	}

	lhs, err := p.condExpr(notGT, notHook)
	if err != nil {
		return nil, err
	}
	loc := lhs.Loc().Clone()

	if p.tsExprHasTypAnnot(lhs) {
		ques := p.tsQues()
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		p.tsNodeTypAnnot(lhs, typAnnot, ACC_MOD_NONE, ques)
	}

	tok := p.lexer.Peek()
	if lhs.Type() == N_NAME && tok.value == T_ARROW && !tok.afterLineTerminator {
		fn, err := p.arrowFn(loc, []Node{lhs}, nil, nil)
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

	rhs, err := p.assignExpr(checkLhs, notGT, false)
	if err != nil {
		return nil, err
	}

	// set `depth` to 1 to permit expr like `i + 2 = 42`
	// and so just do the arg to param transform silently
	lhs, err = p.argToParam(lhs, 1, false, true, false)
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
		scope := p.scope()
		if scope.IsKind(SPK_ASYNC) && node.Text() == "await" {
			return false
		} else if scope.IsKind(SPK_STRICT) && (node.val == "eval" || node.val == "arguments") {
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
	case N_EXPR_BIN:
		node := expr.(*BinExpr)
		return node.op == T_TS_AS
	case N_TS_NO_NULL, N_TS_TYP_ASSERT:
		return true
	}
	return false
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ConditionalExpression
func (p *Parser) condExpr(notGT bool, notHook bool) (Node, error) {
	loc := p.loc()
	test, err := p.binExpr(nil, 0, false, false, notGT)
	if err != nil {
		return nil, err
	}

	if notHook {
		return test, nil
	}

	hook := p.advanceIfTok(T_HOOK)
	if hook == nil {
		return test, nil
	}

	typAnnot, _ := p.tsTypAnnot()
	if typAnnot != nil {
		// `async (x?: number): any => x;`
		if wt, ok := test.(NodeWithTypInfo); ok {
			ti := wt.TypInfo()
			ti.ques = p.locFromTok(hook)
			ti.typAnnot = typAnnot
		}
		return test, nil
	}

	cons, err := p.assignExpr(true, notGT, false)
	if err != nil {
		return nil, err
	}

	_, err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr(true, notGT, false)
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_COND, p.finLoc(loc), test, cons, alt, nil}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-AwaitExpression
func (p *Parser) awaitExpr(tok *Token) (Node, error) {
	loc := p.locFromTok(tok)

	scope := p.scope()

	ahead := p.lexer.Peek()
	if !TokenKinds[ahead.value].StartExpr {
		if p.feat&FEAT_MODULE != 0 {
			// report friendly message for expr like: `async function foo(await) {}`
			if ahead.value == T_PAREN_R || ahead.value == T_COMMA {
				return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, "await"))
			} else if !scope.IsKind(SPK_ASYNC) {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_AWAIT_OUTSIDE_ASYNC)
			}
			return nil, p.errorTok(ahead)
		}
		return &Ident{N_NAME, p.finLoc(loc), "await", false, tok.ContainsEscape(), nil, true, p.newTypInfo()}, nil
	}

	if scope.IsKind(SPK_FORMAL_PARAMS) {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_AWAIT_IN_FORMAL_PARAMS)
	}
	if !scope.IsKind(SPK_ASYNC) {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_AWAIT_OUTSIDE_ASYNC)
	}

	arg, err := p.unaryExpr(nil, nil)
	if err != nil {
		return nil, err
	}
	return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), T_AWAIT, arg, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-UnaryExpression
func (p *Parser) unaryExpr(typArgs Node, typArgsLoc *Loc) (Node, error) {
	var err error
	if typArgs == nil {
		typArgs, err = p.tsTryTypArgs(nil)
		if err != nil {
			return nil, err
		}
	}

	tok := p.lexer.Peek()
	loc := p.locFromTok(tok)
	op := tok.value
	if tok.IsUnary() || op == T_ADD || op == T_SUB || (op == T_LT && p.ts && p.feat&FEAT_JSX == 0) {
		if op != T_LT {
			p.lexer.Next()
		}
		arg, err := p.unaryExpr(nil, nil)
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

		if op == T_DELETE {
			var m *MemberExpr
			at := arg.Type()
			if at == N_EXPR_MEMBER {
				m = arg.(*MemberExpr)
			} else if at == N_EXPR_CHAIN {
				m = arg.(*ChainExpr).expr.(*MemberExpr)
			}

			if m != nil && m.prop.Type() == N_NAME {
				prop := m.prop.(*Ident)
				if prop.pvt {
					return nil, p.errorAtLoc(prop.loc, ERR_DELETE_PVT_FIELD)
				}
			}
		}

		if op == T_LT {
			arg, err = p.tsTypAssert(arg, typArgs)
			if err != nil {
				return nil, err
			}

			return arg, nil
		}

		return &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), op, arg, nil}, nil
	}

	if tok.value == T_AWAIT {
		if p.feat&FEAT_ASYNC_AWAIT == 0 {
			return nil, p.errorTok(tok)
		}
		p.lexer.Next()
		return p.awaitExpr(tok)
	}
	return p.updateExpr(typArgs, typArgsLoc)
}

func (p *Parser) updateExpr(typArgs Node, typArgsLoc *Loc) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Peek()
	if tok.value == T_INC || tok.value == T_DEC {
		p.lexer.Next()
		arg, err := p.unaryExpr(nil, nil)
		if err != nil {
			return nil, err
		}
		if !p.isSimpleLVal(arg, true, false, true, false) {
			return nil, p.errorAtLoc(arg.Loc(), ERR_ASSIGN_TO_RVALUE)
		}
		ud := &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, true, arg, nil}
		arg, err = p.tsTypAssert(ud, typArgs)
		if err != nil {
			return nil, err
		}
		return arg, nil
	}

	arg, err := p.lhs()
	if err != nil {
		return nil, err
	}

	tok = p.lexer.Peek()
	postfix := !tok.afterLineTerminator && (tok.value == T_INC || tok.value == T_DEC)
	if !postfix {
		// for the type info before the arrow fn in this stmt `let a = <T, R>(a: T): void => { a++ }`,
		// it's typeParams of the arrowFn rather than typeAssert
		if arg.Type() == N_EXPR_ARROW {
			ti := arg.(NodeWithTypInfo).TypInfo()
			if ti != nil {
				typArgs, err = p.tsTypArgsToTypParams(typArgs)
				if err != nil {
					return nil, err
				}
				ti.typParams = typArgs
			}
			return arg, nil
		}

		arg, err = p.tsTypAssert(arg, typArgs)
		if err != nil {
			return nil, err
		}
		return arg, nil
	}

	if !p.isSimpleLVal(arg, true, false, true, false) {
		return nil, p.errorAtLoc(arg.Loc(), ERR_ASSIGN_TO_RVALUE)
	}

	p.lexer.Next()

	ud := &UpdateExpr{N_EXPR_UPDATE, p.finLoc(loc), tok.value, false, arg, nil}
	ta, err := p.tsTypAssert(ud, typArgs)
	if err != nil {
		return nil, err
	}
	return ta, nil
}

func (p *Parser) lhs() (Node, error) {
	tok := p.lexer.Peek()
	var node Node
	var err error
	if tok.value == T_NEW {
		node, err = p.newExpr()
	} else {
		node, _, err = p.callExpr(nil, true, false, nil)
	}
	node = p.tsNoNull(node)
	if err != nil {
		return nil, err
	}
	return node, nil
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
		meta := &Ident{N_NAME, p.finLoc(p.locFromTok(new)), "new", false, new.ContainsEscape(), nil, true, p.newTypInfo()}
		p.lexer.Next() // consume dot

		id, err := p.ident(nil, false)
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

	var optLoc *Loc
	expr, optLoc, err = p.memberExpr(nil, false, true, nil)
	if err != nil {
		return nil, err
	}
	if optLoc != nil {
		return nil, p.errorAtLoc(optLoc, ERR_OPT_EXPR_IN_NEW)
	}
	if expr.Type() == N_IMPORT_CALL {
		return nil, p.errorAtLoc(expr.Loc(), ERR_DYNAMIC_IMPORT_CANNOT_NEW)
	}

	var args []Node
	var typArgs Node
	if p.aheadIsArgList(p.lexer.Peek()) {
		args, _, typArgs, err = p.argList(true, true, nil)
		if err != nil {
			return nil, err
		}
	}

	var ret Node
	ti := p.newTypInfo()
	if ti != nil {
		if err = p.tsCheckTypArgs(typArgs); err != nil {
			return nil, err
		}
		ti.typArgs = typArgs
	}
	ret = &NewExpr{N_EXPR_NEW, p.finLoc(loc), expr, args, nil, ti}
	root := true
	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if p.aheadIsArgList(tok) {
			if ret, _, err = p.callExpr(ret, root, false, nil); err != nil {
				return nil, err
			}
		} else if tv == T_BRACKET_L || tv == T_DOT || tv == T_OPT_CHAIN {
			if tv == T_OPT_CHAIN {
				optLoc = p.locFromTok(tok)
			}
			if ret, _, err = p.memberExpr(ret, true, root, optLoc); err != nil {
				return nil, err
			}
		} else {
			break
		}
		if root {
			root = false
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
	case N_EXPR_MEMBER:
		n := callee.(*MemberExpr)
		if n.obj.Type() != N_SUPER {
			return p.checkCallee(n.obj, nextLoc)
		}
	}
	return nil
}

func (p *Parser) isLtTok(tok *Token) bool {
	line := tok.begin.line
	col := tok.begin.col
	for _, lc := range p.ltTokens {
		if lc[0] == line && lc[1] == col {
			return true
		}
	}
	return false
}

func (p *Parser) addLtTok(line, col int) {
	p.ltTokens = append(p.ltTokens, [2]int{line, col})
}

func (p *Parser) pushState() {
	p.lexer.pushState()
	p.lexer.src.pushState()
}

func (p *Parser) discardState() {
	p.lexer.discardState()
	p.lexer.src.discardState()
}

func (p *Parser) popState() {
	p.lexer.src.popState()
	p.lexer.popState()
}

func (p *Parser) aheadIsArgList(tok *Token) bool {
	tv := tok.value
	return tv == T_PAREN_L || (tv == T_LT && p.ts && !p.isLtTok(tok))
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-CallExpression
func (p *Parser) callExpr(callee Node, root bool, directOpt bool, opt *Loc) (Node, *Loc, error) {
	var loc *Loc
	var err error
	if callee == nil {
		loc = p.loc()
		callee, err = p.primaryExpr()
		if err != nil {
			return nil, nil, err
		}
		callee = p.tsNoNull(callee)
	} else {
		loc = callee.Loc().Clone()
	}

	firstOpt := opt
	var fo *Loc
	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if p.aheadIsArgList(tok) {
			aheadLoc := p.locFromTok(tok)

			// below pair `pushState` and `popState` is used to dealing with
			// the ambiguity between `<` in typArgs and `<` operator in binExpr
			lt := tok.value == T_LT
			var line, col int
			if lt {
				line = tok.begin.line
				col = tok.begin.col
				p.pushState()
			}
			args, _, typArgs, err := p.argList(true, true, nil)
			if err != nil {
				if err == errTypArgMissingGT && firstOpt == nil {
					p.popState()
					p.addLtTok(line, col)
					return callee, nil, nil
				}
				return nil, nil, err
			}
			if lt {
				p.discardState()
			}

			ti := p.newTypInfo()
			if ti != nil {
				if err = p.tsCheckTypArgs(typArgs); err != nil {
					return nil, nil, err
				}
				ti.typArgs = typArgs
			}

			if err = p.checkCallee(callee, aheadLoc); err != nil {
				return nil, nil, err
			}

			callee = &CallExpr{N_EXPR_CALL, p.finLoc(loc), callee, args, directOpt, nil, ti}
		} else if tv == T_BRACKET_L || tv == T_DOT || tv == T_OPT_CHAIN {
			callee, fo, err = p.memberExpr(callee, true, root, firstOpt)
			if err != nil {
				return nil, nil, err
			}
			if firstOpt == nil {
				firstOpt = fo
			}
		} else if tv == T_TPL_HEAD {
			callee, err = p.tplExpr(callee)
			if err != nil {
				return nil, nil, err
			}
		} else {
			break
		}
	}

	ct := callee.Type()
	if root && firstOpt != nil && (ct != N_NAME && ct != N_EXPR_CHAIN) {
		return &ChainExpr{N_EXPR_CHAIN, callee.Loc().Clone(), callee}, firstOpt, nil
	}

	return callee, firstOpt, nil
}

// https://262.ecma-international.org/12.0/#prod-ImportCall
func (p *Parser) importCall(tok *Token) (Node, error) {
	if tok == nil {
		tok = p.lexer.Next()
	}
	loc := p.locFromTok(tok)

	meta := &Ident{N_NAME, p.finLoc(p.locFromTok(tok)), tok.Text(), false, tok.ContainsEscape(), nil, false, p.newTypInfo()}

	ahead := p.lexer.Peek()
	if ahead.value == T_DOT && p.feat&FEAT_META_PROPERTY != 0 {
		p.lexer.Next()
		prop, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		if prop.Text() != "meta" {
			return nil, p.errorAtLoc(prop.loc, ERR_ILLEGAL_IMPORT_PROP)
		} else if prop.ContainsEscape() {
			return nil, p.errorAtLoc(prop.loc, ERR_META_PROP_CONTAINS_ESCAPE)
		}

		mp := &MetaProp{N_META_PROP, p.finLoc(loc), meta, prop}
		ahead = p.lexer.Peek()
		av := ahead.value
		if av == T_PAREN_L {
			node, _, err := p.callExpr(mp, true, false, nil)
			return node, err
		} else if av == T_BRACKET_L || av == T_DOT || av == T_OPT_CHAIN {
			node, _, err := p.memberExpr(mp, true, true, nil)
			return node, err
		}
		return &MetaProp{N_META_PROP, p.finLoc(loc), meta, prop}, nil
	}

	if ahead.value == T_PAREN_L && p.feat&FEAT_DYNAMIC_IMPORT != 0 {
		p.lexer.Next()
		src, err := p.assignExpr(true, false, false)
		if err != nil {
			return nil, err
		}
		_, err = p.nextMustTok(T_PAREN_R)
		if err != nil {
			return nil, err
		}
		return &ImportCall{N_IMPORT_CALL, p.finLoc(loc), src, nil}, nil
	}
	return nil, p.errorTok(ahead)
}

func (p *Parser) tplExpr(tag Node) (Node, error) {
	loc := p.loc()
	if tag != nil {
		tl := tag.Loc()
		loc.begin = tl.end.Clone()
		loc.rng.start = tl.rng.end

		if tag.Type() == N_EXPR_CHAIN {
			return nil, p.errorAtLoc(p.locFromTok(p.lexer.Peek()), ERR_OPT_EXPR_IN_TAG)
		}
	}

	elems := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value >= T_TPL_HEAD && tok.value <= T_TPL_TAIL {
			cooked := ""
			ext := tok.ext.(*TokExtTplSpan)
			if ext.IllegalEscape != nil {
				// raise error for bad escape sequence if the template is not tagged
				if tag == nil || p.feat&FEAT_BAD_ESCAPE_IN_TAGGED_TPL == 0 {
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
			str := &StrLit{N_LIT_STR, loc, cooked, false, nil, nil}
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

func (p *Parser) argsToParams(args []Node) ([]Node, error) {
	params := make([]Node, len(args))
	var err error
	for i, arg := range args {
		if arg != nil {
			params[i], err = p.argToParam(arg, 0, false, false, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return params, nil
}

// `yield` indicates whether is yield-expr is permitted
func (p *Parser) checkDefaultVal(val Node, yield bool, destruct bool, field bool) error {
	switch val.Type() {
	case N_EXPR_YIELD:
		scope := p.scope()
		if !yield || !scope.IsKind(SPK_GENERATOR) {
			return p.errorAtLoc(val.Loc(), ERR_YIELD_CANNOT_BE_DEFAULT_VALUE)
		}
		return nil
	case N_EXPR_BIN:
		n := val.(*BinExpr)
		if err := p.checkDefaultVal(n.lhs, yield, destruct, field); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.rhs, yield, destruct, field); err != nil {
			return err
		}
	case N_EXPR_PAREN:
		n := val.(*ParenExpr)
		return p.checkDefaultVal(n.expr, yield, destruct, field)
	case N_EXPR_UNARY:
		n := val.(*UnaryExpr)
		// `{a = await b} = obj` is legal
		// `({a = await b}) => obj` is illegal
		if n.op == T_AWAIT && !destruct {
			return p.errorAtLoc(n.loc, ERR_AWAIT_AS_DEFAULT_VALUE)
		}
		return p.checkDefaultVal(n.arg, yield, destruct, field)
	case N_EXPR_UPDATE:
		n := val.(*UpdateExpr)
		return p.checkDefaultVal(n.arg, yield, destruct, field)
	case N_EXPR_COND:
		n := val.(*CondExpr)
		if err := p.checkDefaultVal(n.test, yield, destruct, field); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.cons, yield, destruct, field); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.alt, yield, destruct, field); err != nil {
			return err
		}
	case N_PAT_ASSIGN:
		n := val.(*AssignPat)
		if err := p.checkDefaultVal(n.lhs, yield, destruct, field); err != nil {
			return err
		}
		if err := p.checkDefaultVal(n.rhs, yield, destruct, field); err != nil {
			return err
		}
	case N_LIT_ARR:
		n := val.(*ArrLit)
		for _, elem := range n.elems {
			if err := p.checkDefaultVal(elem, yield, destruct, field); err != nil {
				return err
			}
		}
	case N_LIT_OBJ:
		n := val.(*ObjLit)
		for _, prop := range n.props {
			if err := p.checkDefaultVal(prop, yield, destruct, field); err != nil {
				return err
			}
		}
	case N_PROP:
		n := val.(*Prop)
		if err := p.checkDefaultVal(n.value, yield, destruct, field); err != nil {
			return err
		}
	case N_SPREAD:
		n := val.(*Spread)
		return p.checkDefaultVal(n.arg, yield, destruct, field)
	case N_NAME:
		id := val.(*Ident)
		name := val.(*Ident).Text()
		if p.checkName && p.isProhibitedName(nil, name, true, false, field, false) {
			return p.errorAtLoc(id.loc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
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
func (p *Parser) argToParam(arg Node, depth int, prop bool, destruct bool, inParen bool) (Node, error) {
	switch arg.Type() {
	case N_LIT_ARR:
		n := arg.(*ArrLit)
		pat := &ArrPat{
			typ:   N_PAT_ARRAY,
			loc:   n.loc,
			elems: make([]Node, len(n.elems)),
			ti:    n.ti,
		}
		var err error
		for i, node := range n.elems {
			// elem maybe nil in expr like `([a, , b]) => 42`
			if node != nil {
				pat.elems[i], err = p.argToParam(node, depth+1, false, destruct, inParen)
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
			ti:    n.ti,
		}
		for i, prop := range n.props {
			pp, err := p.argToParam(prop, depth+1, true, destruct, inParen)
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
			_, err := p.argToParam(prop, depth+1, true, destruct, inParen)
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
				if n.key, err = p.argToParam(n.key, depth+1, prop, destruct, inParen); err != nil {
					return nil, err
				}
				if err = p.checkDefaultVal(n.value, destruct, destruct, false); err != nil {
					return nil, err
				}
			} else {
				// the correctness of the value should be checked account for
				// using it as an alias
				val, err := p.argToParam(n.value, depth+1, prop, destruct, inParen)
				if err != nil {
					return nil, err
				}
				n.value = val
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
			err = p.checkDefaultVal(n.value, destruct, destruct, false)
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

		if pn, ok := n.lhs.(InParenNode); ok {
			inParen = pn.OuterParen() != nil
		}
		lhs, err := p.argToParam(n.lhs, depth+1, false, destruct, inParen)
		if err != nil {
			return nil, err
		}

		// also check the default value
		err = p.checkDefaultVal(n.rhs, destruct, destruct, false)
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
		if p.checkName && p.isProhibitedName(nil, name, true, true, false, false) {
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
			return nil, p.errorAtLoc(n.trailingCommaLoc, ERR_REST_ELEM_MUST_LAST)
		}

		at := n.arg.Type()
		if at == N_NAME {
			// `({...(obj)} = foo)` raises error`Parenthesized pattern` in acorn
			// however it's legal in babel-parser, chrome and firefox
			//
			// use `destruct` to require the caller to indicate the parsing state
			// is in destructing or not
			if !destruct {
				if n, ok := n.arg.(InParenNode); ok {
					if n.OuterParen() != nil {
						return nil, p.errorAtLoc(n.OuterParen(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
					}
				}
			}
			if _, err := p.argToParam(n.arg, depth, prop, destruct, inParen); err != nil {
				return nil, err
			}
		} else if at == N_EXPR_ASSIGN {
			return nil, p.errorAtLoc(n.arg.Loc(), ERR_REST_CANNOT_SET_DEFAULT)
		} else if at == N_EXPR_PAREN {
			if destruct {
				arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen)
				if err != nil {
					return nil, err
				}
				n.arg = arg
			} else {
				return nil, p.errorAtLoc(n.arg.Loc(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
			}
		} else if p.feat&FEAT_BINDING_REST_ELEM_NESTED != 0 && (at == N_LIT_ARR || at == N_LIT_OBJ) {
			arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen)
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

		rest := &RestPat{
			typ: N_PAT_REST,
			loc: n.loc,
			arg: n.arg,
			ti:  p.newTypInfo(),
		}
		if p.ts {
			rest.hoistTypInfo()
		}
		return rest, nil
	case N_EXPR_PAREN:
		sub := arg.(*ParenExpr).expr
		if !destruct || !p.isPrimitive(sub) && !p.isTsLhs(sub) {
			st := sub.Type()
			if !(destruct && st == N_EXPR_BIN && sub.(*BinExpr).op == T_TS_AS) {
				if st != N_LIT_ARR && st != N_LIT_OBJ && st != N_NAME {
					return nil, p.errorAtLoc(sub.Loc(), ERR_ASSIGN_TO_RVALUE)
				}
				return nil, p.errorAtLoc(arg.Loc(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
			}
		}
		arg, err := p.argToParam(sub, depth, prop, destruct, true)
		if err != nil {
			return nil, err
		}
		if pn, ok := arg.(InParenNode); ok {
			pn.SetOuterParen(sub.Loc().Clone())
		}
		return arg, nil
	case N_TS_TYP_ASSERT:
		n := arg.(*TsTypAssert)
		if destruct && !inParen {
			if depth < 2 {
				// `[a as number] = [42];` is legal
				// `<string>foo = '100';` is illegal
				return nil, p.errorAtLoc(n.loc, ERR_ASSIGN_TO_RVALUE)
			}
			// transform the arg at first: `<number>(a)`
			arg, err := p.argToParam(n.arg, depth, prop, destruct, true)
			if err != nil {
				return nil, err
			}

			// the transformed arg should be `NodeWithTypInfo` since we need to attach the
			// `des` of TsTypAssert as typAnnot of it
			if wt, ok := arg.(NodeWithTypInfo); ok {
				wt.TypInfo().typAnnot = n.des
			} else {
				return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
			}
			return arg, nil
		}
	case N_EXPR_BIN:
		n := arg.(*BinExpr)
		if destruct && depth > 0 && n.op == T_TS_AS {
			if inParen {
				// however `foo as any = 10;` is illegal
				return n, nil
			} else if depth < 2 {
				// opposite to above true case, `[a as number] = [42];` is legal
				return nil, p.errorAtLoc(n.loc, ERR_ASSIGN_TO_RVALUE)
			}
			arg, err := p.argToParam(n.lhs, depth, prop, destruct, true)
			if err != nil {
				return nil, err
			}
			if wt, ok := arg.(NodeWithTypInfo); ok {
				wt.TypInfo().typAnnot = n.rhs
			} else {
				return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
			}
			return arg, nil
		}
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

func (p *Parser) nameOfProp(propOrField Node) (string, Node, bool) {
	var key Node
	pt := propOrField.Type()
	if pt == N_PROP {
		key = propOrField.(*Prop).key
	} else if pt == N_FIELD {
		key = propOrField.(*Field).key
	} else if pt == N_METHOD {
		key = propOrField.(*Method).key
	} else {
		return "", nil, false
	}

	var propName string
	priv := false
	switch key.Type() {
	case N_NAME:
		id := key.(*Ident)
		propName = id.Text()
		priv = id.pvt
		if priv {
			propName = "#" + propName
		}
	case N_LIT_STR:
		propName = key.(*StrLit).Text()
	case N_LIT_NUM:
		propName = key.(*NumLit).Text()
	case N_LIT_BOOL:
		propName = key.(*BoolLit).Text()
	case N_LIT_NULL:
		propName = key.(*NullLit).Text()
	}
	return propName, key, priv
}

// check the `arg` is legal as argument
// `spread` means whether the spread is permitted
// `simplicity` means whether check simplicity of lhs of the assignExpr
func (p *Parser) checkArg(arg Node, spread bool, simplicity bool) error {
	if wt, ok := arg.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		// report the error of hook in expr like: `async(x?)`
		if ti != nil && ti.ques != nil {
			return p.errorAtLoc(ti.ques, ERR_UNEXPECTED_TOKEN)
		}
	}

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
				pn, _, _ := p.nameOfProp(pp)
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
			notIn := p.scope().IsKind(SPK_NOT_IN)
			// assign is legal in expr like `for ({x = 0} in arr);`
			if !notIn {
				return p.errorAtLoc(n.opLoc, ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT)
			}
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
	case N_LIT_ARR:
		n := arg.(*ArrLit)
		for _, el := range n.elems {
			if el == nil {
				continue
			}
			if err := p.checkArg(el, true, simplicity); err != nil {
				return err
			}
		}
	case N_NAME:
		id := arg.(*Ident)
		if id.kw && p.scope().IsKind(SPK_STRICT) {
			return p.errorAtLoc(arg.Loc(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.Text()))
		}
		// for reporting `(a:b)` is illegal in ts
		if id.ti != nil && id.ti.typAnnot != nil {
			return p.errorAtLoc(id.ti.typAnnot.Loc(), ERR_UNEXPECTED_TYPE_ANNOTATION)
		}
	case N_TS_NO_NULL:
		return nil
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

func (p *Parser) argList(check bool, incall bool, asyncLoc *Loc) ([]Node, *Loc, Node, error) {
	typArgs, err := p.tsTryTypArgs(asyncLoc)
	if err != nil {
		return nil, nil, nil, err
	}
	if typArgs != nil {
		tt := typArgs.Type()
		if tt != N_TS_PARAM_INST && tt != N_TS_PARAM_DEC {
			return nil, nil, typArgs, nil
		}
	}

	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, nil, nil, err
	}

	var tailingComma *Loc
	args := make([]Node, 0)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, nil, nil, p.errorTok(tok)
		}
		arg, err := p.arg()
		if err != nil {
			return nil, nil, nil, err
		}

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_COMMA {
			tok := p.lexer.Next()
			// trailing comma is need to be checked when it's in
			// parenExpr, this snippet `...a,` is illegal as parenExpr: `(...a,)`
			// however it is legal as arguments `foo(...a,)`
			ahead := p.lexer.Peek()
			if !incall && arg.Type() == N_SPREAD {
				msg := ERR_REST_ELEM_MUST_LAST
				if ahead.value != T_PAREN_R {
					msg = ERR_REST_ELEM_MUST_LAST
				}
				return nil, tailingComma, nil, p.errorAt(tok.value, &tok.begin, msg)
			}
			if tailingComma == nil && ahead.value == T_PAREN_R {
				tailingComma = p.locFromTok(tok)
			}
		} else if av != T_PAREN_R {
			return nil, nil, nil, p.errorTok(ahead)
		}

		if check {
			// `spread` or `pattern_rest` expression is legal argument:
			// `f(c, b, ...a)`
			if err := p.checkArg(arg, true, false); err != nil {
				return nil, nil, nil, err
			}
		}

		args = append(args, arg)
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, nil, nil, err
	}
	return args, tailingComma, typArgs, nil
}

// consider below exprs:
// `(a,b)`
// `(a,b) =>`
// we cannot judge `(a,b)` is a parenExpr or the formalParamsList of
// an arrayExpr before we see the `=>` token, for avoding to rollback
// the parsing state, we firstly parse `(a,b)` as parenExpr which children
// is parsed by this method and then convert the parsed subtree to formalParamList
// by using `argToParam` when required
func (p *Parser) arg() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(false, false, true)
}

func (p *Parser) checkOp(tok *Token) error {
	if tok.value == T_POW && p.feat&FEAT_POW == 0 {
		return p.errorTok(tok)
	}
	return nil
}

// `logic` indicates that there is at least one logic expr in previous lhs group
// `nullish` indicates that there is at least one nullish expr in previous lhs group
// bypassing above two params to do the `nullish-can-not-along-with-logic` syntax check as
// well as in one parse - avoiding to traverse the sub-tree later
func (p *Parser) binExpr(lhs Node, minPcd int, logic bool, nullish bool, notGT bool) (Node, error) {
	var err error
	if lhs == nil {
		if lhs, err = p.unaryExpr(nil, nil); err != nil {
			return nil, err
		}
	}

	ts := p.ts
	notIn := p.scope().IsKind(SPK_NOT_IN)
	for {
		ahead := p.lexer.Peek()
		op := ahead.IsBin(notIn, ts)
		if op == T_ILLEGAL {
			break
		}

		kind := TokenKinds[op]
		pcd := kind.Pcd
		if pcd < minPcd || (op == T_GT && notGT) {
			break
		}

		if logic && op == T_NULLISH || nullish && (op == T_AND || op == T_OR) {
			return nil, p.errorAtLoc(p.locFromTok(ahead), ERR_NULLISH_MIXED_WITH_LOGIC)
		}

		if op == T_AND || op == T_OR {
			if !logic {
				logic = true
			}
		} else if op == T_NULLISH {
			if !nullish {
				nullish = true
			}
		}

		if err = p.checkOp(ahead); err != nil {
			return nil, err
		}
		opLoc := p.locFromTok(p.lexer.Next())

		var rhs Node
		if op != T_TS_AS {
			rhs, err = p.unaryExpr(nil, nil)
		} else {
			rhs, err = p.tsTyp(false)
		}
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		aheadOp := ahead.IsBin(notIn, ts)
		kind = TokenKinds[aheadOp]
		for aheadOp != T_ILLEGAL && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			pcd = kind.Pcd
			rhs, err = p.binExpr(rhs, pcd, logic, nullish, notGT)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			aheadOp = ahead.IsBin(notIn, ts)
			kind = TokenKinds[aheadOp]
		}

		// deal with expr like: `console.log( -2 ** 4 )`
		if lhs.Type() == N_EXPR_UNARY && op == T_POW {
			return nil, p.errorAtLoc(p.UnParen(lhs.(*UnaryExpr).arg).Loc(), ERR_UNARY_OPERATOR_IMMEDIATELY_BEFORE_POW)
		}

		// deal with expr like: `4 + async() => 2`
		if rhs.Type() == N_EXPR_ARROW {
			return nil, p.errorAtLoc(rhs.(*ArrowFn).arrowLoc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, "=>"))
		}

		bin := &BinExpr{N_EXPR_BIN, nil, T_ILLEGAL, nil, nil, nil, nil}
		bin.loc = p.finLoc(lhs.Loc().Clone())
		bin.op = op
		bin.opLoc = opLoc
		bin.lhs = lhs
		bin.rhs = rhs
		lhs = bin
	}
	return lhs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-MemberExpression
func (p *Parser) memberExpr(obj Node, call bool, root bool, optLoc *Loc) (Node, *Loc, error) {
	var err error
	if obj == nil {
		if p.lexer.Peek().value == T_NEW {
			if obj, err = p.newExpr(); err != nil {
				return nil, nil, err
			}
		} else if obj, err = p.memberExprObj(); err != nil {
			return nil, nil, err
		}
	}

	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if tv == T_OPT_CHAIN {
			if optLoc == nil {
				optLoc = p.locFromTok(tok)
			}
			p.lexer.Next()

			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_BRACKET_L { // a?.[b]
				if obj, err = p.memberExprPropSubscript(obj, true); err != nil {
					return nil, nil, err
				}
			} else if p.aheadIsArgList(ahead) { // a?.()
				if obj, _, err = p.callExpr(obj, false, true, optLoc); err != nil {
					return nil, nil, err
				}
			} else {
				// a?.b
				if obj, err = p.memberExprPropDot(obj, true); err != nil {
					return nil, nil, err
				}
			}
		} else if tv == T_BRACKET_L {
			if obj, err = p.memberExprPropSubscript(obj, false); err != nil {
				return nil, nil, err
			}
		} else if tv == T_DOT {
			p.lexer.Next()
			if obj, err = p.memberExprPropDot(obj, false); err != nil {
				return nil, nil, err
			}
		} else {
			break
		}
	}

	// `super.#aaa` is illegal since the direct pvt access
	// `super.c.#aaa` is legal and fall into runtime-check
	if obj.Type() == N_EXPR_MEMBER {
		m := obj.(*MemberExpr)
		if m.obj.Type() == N_SUPER && m.prop.Type() == N_NAME {
			prop := m.prop.(*Ident)
			if prop.pvt {
				return nil, nil, p.errorAtLoc(prop.loc, ERR_UNEXPECTED_PVT_FIELD)
			}
		}
	}

	if call && p.aheadIsArgList(p.lexer.Peek()) {
		return p.callExpr(obj, root, false, optLoc)
	}

	if root && optLoc != nil && obj.Type() != N_NAME {
		return &ChainExpr{N_EXPR_CHAIN, obj.Loc().Clone(), obj}, optLoc, nil
	}
	return obj, optLoc, nil
}

func (p *Parser) memberExprObj() (Node, error) {
	obj, err := p.primaryExpr()
	if err != nil {
		return nil, err
	}
	obj = p.tsNoNull(obj)
	if p.lexer.Peek().value == T_TPL_HEAD {
		return p.tplExpr(obj)
	}
	return obj, nil
}

func (p *Parser) memberExprPropSubscript(obj Node, opt bool) (Node, error) {
	p.lexer.Next()
	prop, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}
	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, true, opt, nil}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node, opt bool) (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()
	tv := tok.value
	_, kw, ok := tok.CanBePropKey()

	var prop Node
	if (ok && tv != T_NUM) || tv == T_NAME_PVT {
		pvt := tv == T_NAME_PVT
		id := &Ident{N_NAME, p.finLoc(loc), tok.Text(), pvt, tok.ContainsEscape(), nil, kw, p.newTypInfo()}
		if pvt {
			scope := p.scope().UpperCls()
			if scope == nil {
				return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_ALONE_PVT_FIELD, "#"+tok.Text()))
			}
			ref := NewRef()
			ref.Node = id
			ref.TargetType = TT_PVT_FIELD
			ref.Scope = scope
			p.danglingPvtRefs = append(p.danglingPvtRefs, ref)
		}
		prop = id
	} else {
		return nil, p.errorTok(tok)
	}

	node := &MemberExpr{N_EXPR_MEMBER, p.finLoc(obj.Loc().Clone()), obj, prop, false, opt, nil}
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
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), legacyOctalEscapeSeq, nil, nil}, nil
	case T_NULL:
		p.lexer.Next()
		return &NullLit{N_LIT_NULL, p.finLoc(loc), nil, nil}, nil
	case T_TRUE, T_FALSE:
		p.lexer.Next()
		return &BoolLit{N_LIT_BOOL, p.finLoc(loc), tok.Text() == "true", nil, nil}, nil
	case T_NAME:
		if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
			}
			return p.fnDec(true, tok, false, false)
		} else if p.tsAheadIsAbstract(tok, false, false) {
			return p.classDec(true, false, false, true)
		}

		p.lexer.Next()
		name := tok.Text()
		ahead := p.lexer.Peek()
		// `ahead.value != T_ARROW` is used to skip checking name when it appears in the param list of arrow expr
		// for `eval => 42` we should report binding-reserved-word error instead of unexpected-reserved-word error
		if p.checkName && ahead.value != T_ARROW && !ahead.afterLineTerminator && p.isProhibitedName(nil, name, true, false, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAtLoc(p.finLoc(loc), ERR_ESCAPE_IN_KEYWORD)
			}
			return nil, p.errorAtLoc(p.finLoc(loc), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
		}
		kw := p.isProhibitedName(nil, name, true, false, false, false)
		return &Ident{N_NAME, p.finLoc(loc), name, false, tok.ContainsEscape(), nil, kw, p.newTypInfo()}, nil
	case T_THIS:
		p.lexer.Next()
		return &ThisExpr{N_EXPR_THIS, p.finLoc(loc), nil, nil}, nil
	case T_PAREN_L:
		return p.parenExpr(nil)
	case T_BRACKET_L:
		return p.arrLit()
	case T_BRACE_L:
		return p.objLit()
	case T_FUNC:
		return p.fnDec(true, nil, false, false)
	case T_REGEXP:
		p.lexer.Next()
		ext := tok.ext.(*TokExtRegexp)
		return &RegLit{N_LIT_REGEXP, p.finLoc(loc), tok.Text(), ext.Pattern(), ext.Flags(), nil, nil}, nil
	case T_CLASS:
		return p.classDec(true, false, false, false)
	case T_SUPER:
		scope := p.scope()
		sup := p.lexer.Next()

		ahead := p.lexer.Peek()
		if !scope.IsKind(SPK_CLASS) && !scope.IsKind(SPK_CLASS_INDIRECT) && !scope.IsKind(SPK_METHOD) ||
			scope.IsKind(SPK_PROP_NAME) {
			em := ERR_SUPER_OUTSIDE_CLASS
			if ahead.value == T_PAREN_L {
				em = ERR_SUPER_CALL_OUTSIDE_CTOR
			}
			return nil, p.errorAtLoc(loc, em)
		}

		if ahead.value != T_DOT && ahead.value != T_PAREN_L {
			return nil, p.errorTok(sup)
		}
		return &Super{N_SUPER, p.finLoc(loc), nil, nil}, nil
	case T_IMPORT:
		return p.importCall(nil)
	case T_TPL_HEAD:
		return p.tplExpr(nil)
	case T_LT:
		if p.feat&FEAT_JSX != 0 {
			return p.jsx(true, false)
		} else if p.feat&FEAT_TS != 0 {
			typArgs, err := p.tsTryTypArgs(nil)
			if err != nil {
				return nil, err
			}
			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_PAREN_L {
				return p.parenExpr(typArgs)
			}
			return p.unaryExpr(typArgs, typArgs.Loc().Clone())
		}
		return nil, p.errorTok(tok)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) arrowFn(loc *Loc, args []Node, params []Node, ti *TypInfo) (Node, error) {
	var err error
	if params == nil {
		params, err = p.argsToParams(args)
		if err != nil {
			return nil, err
		}
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
		p.addLocalBinding(nil, ref, false, ref.Node.Text())
	}

	if ti == nil {
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		ti = p.newTypInfo()
		if ti != nil {
			ti.typAnnot = typAnnot
		}
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

	return &ArrowFn{N_EXPR_ARROW, p.finLoc(loc), arrowLoc, false, params, body, body.Type() != N_STMT_BLOCK, nil, ti}, nil
}

func (p *Parser) parenExpr(typArgs Node) (Node, error) {
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
	args, tailingComma, ta, err := p.argList(false, false, nil)
	p.checkName = true
	p.symtab.LeaveScope()

	if ta != nil {
		typArgs = ta
	}

	if err != nil {
		return nil, err
	}

	params, paramsErr := p.argsToParams(args)
	allowTypAnnot := paramsErr == nil

	var typAnnot Node
	if allowTypAnnot {
		typAnnot, err = p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
	}

	ti := p.newTypInfo()
	if ti != nil {
		ti.typAnnot = typAnnot
		ti.typArgs = typArgs
	}

	// next is arrow-expression
	ahead := p.lexer.Peek()
	if ahead.value == T_ARROW && !ahead.afterLineTerminator {
		if paramsErr != nil {
			return nil, paramsErr
		}
		return p.arrowFn(loc, nil, params, ti)
	}

	// `():number` is illegal
	if typAnnot != nil {
		return nil, p.errorAtLoc(typAnnot.Loc(), ERR_UNEXPECTED_TOKEN)
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
		pe := &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), args[0], nil}
		if ti == nil {
			return pe, nil
		}
		node, err := p.tsTypAssert(pe, typArgs)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	seqLoc := args[0].Loc().Clone()
	end := args[argsLen-1].Loc()
	seqLoc.rng.end = end.rng.end
	seqLoc.end = end.end.Clone()
	seq := &SeqExpr{N_EXPR_SEQ, seqLoc, args, nil}
	pe := &ParenExpr{N_EXPR_PAREN, p.finLoc(loc), seq, nil}
	if ti == nil {
		return pe, nil
	}
	node, err := p.tsTypAssert(pe, typArgs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) UnParen(expr Node) Node {
	if expr.Type() == N_EXPR_PAREN {
		loc := expr.Loc().Clone()
		sub := expr.(*ParenExpr).Expr()
		if n, ok := sub.(InParenNode); ok {
			n.SetOuterParen(loc)
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
	return &ArrLit{N_LIT_ARR, p.finLoc(loc), elems, nil, p.newTypInfo()}, nil
}

func (p *Parser) arrElem() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(true, false, false)
}

func (p *Parser) spread() (Node, error) {
	loc := p.loc()
	tok := p.lexer.Next()

	if p.feat&FEAT_SPREAD == 0 {
		return nil, p.errorTok(tok)
	}

	node, err := p.assignExpr(true, false, false)
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
	return &Spread{N_SPREAD, p.finLoc(loc), node, trailingCommaLoc, nil, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ObjectLiteral
func (p *Parser) objLit() (Node, error) {
	loc := p.loc()
	p.lexer.Next()

	props := make([]Node, 0, 1)
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
	return &ObjLit{N_LIT_OBJ, p.finLoc(loc), props, nil, p.newTypInfo()}, nil
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
		return p.method(nil, nil, ACC_MOD_NONE, nil, false, PK_INIT, true, false, false, false, false, false, false)
	} else if p.aheadIsAsync(tok, true, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(nil, nil, ACC_MOD_NONE, nil, false, PK_INIT, false, true, false, false, false, false, false)
	}
	return p.propData()
}

func (p *Parser) propData() (Node, error) {
	key, compute, err := p.propName(false, true, false)

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
		value, err = p.assignExpr(true, false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsArgList(tok) {
		return p.method(loc, key, ACC_MOD_NONE, compute, false, PK_INIT, false, false, false, false, false, false, false)
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
	return &Prop{N_PROP, p.finLoc(loc), key, opLoc, value, compute != nil, false, shorthand, assign, PK_INIT, ACC_MOD_NONE}, nil
}

func (p *Parser) method(loc *Loc, key Node, accMode ACC_MOD, compute *Loc, shorthand bool, kind PropKind,
	gen bool, async bool, allowNamePVT bool, inclass bool, static bool, declare bool, abstract bool) (Node, error) {

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
			key, compute, err = p.propName(allowNamePVT, false, false)
		}
		if err != nil {
			return nil, err
		}
	}

	ctor := false
	if p.isName(key, "constructor", false, true) && compute == nil {
		if kind == PK_GETTER || kind == PK_SETTER {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_HAVE_MODIFIER)
		} else if async {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_BE_ASYNC)
		} else if gen {
			return nil, p.errorAtLoc(key.Loc(), ERR_CTOR_CANNOT_BE_GENERATOR)
		}
		ctor = true
	}

	fnLoc := p.loc()
	params, typParams, _, err := p.paramList(false, ctor && p.ts, true)
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
		p.addLocalBinding(nil, ref, false, ref.Node.Text())
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	ti := p.newTypInfo()
	if ti != nil {
		ti.typAnnot = typAnnot
		ti.typParams = typParams
	}

	if kind == PK_GETTER && len(params) > 0 {
		return nil, p.errorAtLoc(params[0].Loc(), ERR_GETTER_SHOULD_NO_PARAM)
	}
	if kind == PK_SETTER {
		if len(params) != 1 {
			return nil, p.errorAtLoc(fnLoc, ERR_SETTER_SHOULD_ONE_PARAM)
		}
		if params[0].Type() == N_PAT_REST {
			return nil, p.errorAtLoc(params[0].Loc(), ERR_REST_IN_SETTER)
		}
	}

	if gen {
		p.lexer.addMode(LM_GENERATOR)
	}

	var body Node
	ahead := p.lexer.Peek()
	if ahead.value == T_BRACE_L {
		body, err = p.fnBody()
		if gen {
			p.lexer.popMode()
		}
		if err != nil {
			return nil, err
		}
	} else if !declare && !abstract {
		return nil, p.errorAt(ahead.value, &ahead.begin, ERR_UNEXPECTED_TOKEN)
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)

	// `isProhibitedName` is not needed here since `keyword` as method name is permitted
	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	p.symtab.LeaveScope()

	if declare || abstract {
		p.advanceIfSemi(false)
	}

	value := &FnDec{N_EXPR_FN, p.finLoc(fnLoc), nil, gen, async, params, body, nil, ti}
	if inclass {
		if static && p.isName(key, "prototype", false, true) {
			return nil, p.errorAtLoc(key.Loc(), ERR_STATIC_PROP_PROTOTYPE)
		}

		return &Method{N_METHOD, p.finLoc(loc), key, static, compute != nil, kind, value, accMode, abstract}, nil
	}
	return &Prop{N_PROP, p.finLoc(loc), key, nil, value, compute != nil, true, shorthand, false, kind, accMode}, nil
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
		end:   tok.end.Clone(),
		rng:   &Range{tok.raw.lo, tok.raw.hi},
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
