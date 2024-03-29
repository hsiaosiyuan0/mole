package parser

import (
	"errors"
	"fmt"

	span "github.com/hsiaosiyuan0/mole/span"
)

// parser is one-pass mode and returns the first error in which either syntax-error or semantic-error
//
// it supports below syntaxes out-of-box by setting the `ParserOpts.Feature`:
//
// - ecmascript up to 2020
// - jsx
// - typescript
//
// an AST couple with a Symtab will be constructed after the source is processed successfully
type Parser struct {
	lexer           *Lexer
	symtab          *SymTab
	feat            Feature
	imp             map[string]*Ident
	checkName       bool
	danglingPvtRefs []*Ref

	ts  bool
	dts bool

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
	ltTokens map[uint32]bool

	// the parsed decorators but have not been attached to the target nodes
	hangingDecorators []Node

	// a temporary stack which holds the loop nodes in their lexical order
	loopStk []Node

	// for resolving the `FnDec.rets`
	retsStk [][]Node

	// keep tryStmts in their lexical order
	tryStk []Node

	// the root node after process is finished
	prog Node

	// node => comments
	prevCmts map[Node][]span.Range
	postCmts map[Node][]span.Range

	errTypArgMissingGT ErrTypArgMissingGT
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
	FEAT_DYNAMIC_IMPORT | FEAT_JSON_SUPER_SET | FEAT_EXPORT_ALL_AS_NS | FEAT_JSX | FEAT_DECORATOR

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
	if on, ok := obj["typescript"]; ok {
		o.Feature = o.Feature.Turn(FEAT_TS, on == true)
	}
	if on, ok := obj["jsx"]; ok {
		o.Feature = o.Feature.Turn(FEAT_JSX, on == true)
	}
	if on, ok := obj["dts"]; ok {
		o.Feature = o.Feature.Turn(FEAT_DTS, on == true)
	}
	if on, ok := obj["strict"]; ok {
		o.Feature = o.Feature.Turn(FEAT_STRICT, on == true)
	}
}

func NewParser(src *span.Source, opts *ParserOpts) *Parser {
	parser := &Parser{}
	parser.Setup(src, opts)
	return parser
}

func (p *Parser) Setup(src *span.Source, opts *ParserOpts) {
	if opts.Feature&FEAT_ASYNC_AWAIT == 0 {
		opts.Feature = opts.Feature.Off(FEAT_GLOBAL_ASYNC)
	}
	if opts.Feature&FEAT_MODULE != 0 {
		opts.Feature = opts.Feature.On(FEAT_IMPORT_DEC).On(FEAT_EXPORT_DEC)
	}

	p.feat = opts.Feature
	p.imp = map[string]*Ident{}
	p.checkName = true
	p.danglingPvtRefs = make([]*Ref, 0)
	p.ltTokens = map[uint32]bool{}
	p.symtab = NewSymTab(opts.Externals)
	p.loopStk = []Node{}
	p.retsStk = [][]Node{}
	p.tryStk = []Node{}
	p.prevCmts = map[Node][]span.Range{}
	p.postCmts = map[Node][]span.Range{}

	p.lexer = NewLexer(src)
	p.lexer.ver = opts.Version
	p.lexer.feat = opts.Feature
	if p.feat&FEAT_TS != 0 || p.feat&FEAT_DTS != 0 {
		p.lexer.AddMode(LM_TS)
	}

	p.ts = p.feat&FEAT_TS != 0
	p.dts = p.feat&FEAT_DTS != 0
}

func (p *Parser) pushLoopStk(loopNode Node) {
	p.loopStk = append(p.loopStk, loopNode)
}

func (p *Parser) popLoopStk() {
	p.loopStk = p.loopStk[0 : len(p.loopStk)-1]
}

func (p *Parser) pushTryStk(try Node) {
	p.tryStk = append(p.tryStk, try)
}

func (p *Parser) popTryStk() {
	p.tryStk = p.tryStk[0 : len(p.tryStk)-1]
}

func (p *Parser) lastTry() Node {
	if len(p.tryStk) == 0 {
		return nil
	}
	return p.tryStk[len(p.tryStk)-1]
}

func (p *Parser) incRetsStk() {
	p.retsStk = append(p.retsStk, make([]Node, 0, 5))
}

func (p *Parser) decRetsStk() []Node {
	last, rest := p.retsStk[len(p.retsStk)-1], p.retsStk[:len(p.retsStk)-1]
	p.retsStk = rest
	return last
}

func (p *Parser) pushRetsStk(ret Node) Node {
	last := len(p.retsStk) - 1
	p.retsStk[last] = append(p.retsStk[last], ret)
	return ret
}

func (p *Parser) Symtab() *SymTab {
	return p.symtab
}

func (p *Parser) PrevCmts(stmt Node) []span.Range {
	return p.prevCmts[stmt]
}

func (p *Parser) PostCmts(stmt Node) []span.Range {
	return p.postCmts[stmt]
}

func (p *Parser) Source() *span.Source {
	return p.lexer.src
}

func (p *Parser) Lexer() *Lexer {
	return p.lexer
}

func (p *Parser) Ast() Node {
	return p.prog
}

func (p *Parser) Prog() (Node, error) {
	rng := p.rng()
	pg := &Prog{N_PROG, span.Range{}, make([]Node, 0, 20)}
	p.prog = pg

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
	rng.Hi = p.lexer.src.Ofst()
	pg.rng = rng

	if err := p.checkExp(scope.Exports); err != nil {
		return nil, err
	}

	if err := p.resolvingDanglingPvtRefs(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Parser) resolvingDanglingPvtRefs() error {
	for _, ref := range p.danglingPvtRefs {
		if ref.Typ == RDT_PVT_FIELD {
			name := "#" + ref.Id.val
			target := ref.Scope.BindingOf(name)
			if target != nil {
				target.RetainBy(ref)
			} else {
				return p.errorAtLoc(ref.Id.Range(), fmt.Sprintf(ERR_TPL_ALONE_PVT_FIELD, name))
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
func (p *Parser) checkExp(exps []*ExportDec) error {
	names := map[string]bool{}
	// check duplication
	for _, exp := range exps {
		var subnames []Node
		if !exp.def.Empty() {
			subnames = []Node{&Ident{N_NAME, exp.def, "default", false, false, span.Range{}, true, p.newTypInfo(N_NAME)}}
		} else if exp.dec != nil {
			subnames = p.namesInNode(exp.dec)
		} else {
			subnames = make([]Node, 0, len(exp.specs))
			for _, spec := range exp.specs {
				s := spec.(*ExportSpec)
				if s.id != nil && !s.tsTyp {
					subnames = append(subnames, s.id)
				}
			}
		}
		for _, sn := range subnames {
			id := sn.(*Ident)
			name := id.val
			if _, ok := names[name]; ok {
				return p.errorAtLoc(id.Range(), fmt.Sprintf(ERR_DUP_EXPORT, name))
			} else {
				names[name] = true
			}
		}
	}

	// also check definition
	// here separate the definition checking into two checks since
	// their errors needed to be reported independently - firstly report
	// the duplication then the definition
	for _, exp := range exps {
		if exp.src != nil {
			continue
		}
		for _, spec := range exp.specs {
			id := spec.(*ExportSpec).local.(*Ident)
			name := id.val
			if !p.scope().HasName(name) {
				return p.errorAtLoc(id.rng, fmt.Sprintf(ERR_EXPORT_NOT_DEFINED, name))
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
	checkHangingDec := true

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
			node, err = p.tsEnum(span.Range{Lo: 1, Hi: 0}, false)
		}
	} else if tok.value == T_BRACE_L {
		node, err = p.blockStmt(true, SPK_NONE)
	} else if ok, kind := p.aheadIsVarDec(tok); ok {
		if allowDec {
			node, err = p.varDecStmt(kind, false)
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
		}
		node, err = p.fnDec(false, tok, false)
	} else if p.aheadIsLabel(tok) {
		node, err = p.labelStmt()
	} else if p.aheadIsTsTypDec(tok, false) {
		rng := p.lexer.Next().rng
		node, err = p.tsTypDec(rng, false)
	} else if p.aheadIsTsItf(tok) {
		node, err = p.tsItf()
	} else if p.aheadIsTsNS(tok) {
		node, err = p.tsNS()
	} else if p.aheadIsModDec(tok) {
		node, err = p.tsModDec()
	} else if p.aheadIsTsDec(tok) {
		node, err = p.tsDec()
	} else if ok, itf, _ := p.tsAheadIsAbstract(tok, false, false, false); ok {
		if itf {
			return nil, p.errorAtLoc(tok.rng, ERR_ABSTRACT_AT_INVALID_POSITION)
		}
		node, err = p.classDec(false, false, false, true)
	} else if p.aheadIsDecorator(tok) {
		ds, err := p.decorators()
		if err != nil {
			return nil, err
		}
		p.hangingDecorators = ds
		checkHangingDec = false
	} else if tok.value == T_SEMI {
		node, err = p.emptyStmt()
	} else if tok.value == T_EOF {
		node, err = nil, errEof
	} else {
		node, err = p.exprStmt()
	}

	if err != nil {
		return nil, err
	} else if node == nil && checkHangingDec {
		return nil, p.errorTok(tok)
	}

	if node != nil {
		typ := node.Type()
		if scope.IsKind(SPK_INTERIM) {
			// `if (morning) function a(){}` is legal
			// `for (morning;;) function a(){}` is illegal
			if typ == N_STMT_FN && (scope.IsKind(SPK_STRICT) || scope.IsKind(SPK_LOOP_DIRECT)) {
				return nil, p.errorAtLoc(node.Range(), ERR_FN_IN_SINGLE_STMT_CTX)
			} else if typ == N_STMT_IMPORT || typ == N_STMT_EXPORT {
				return nil, p.errorAtLoc(node.Range(), ERR_IMPORT_EXPORT_SHOULD_AT_TOP_LEVEL)
			}
		}
	}

	if checkHangingDec && len(p.hangingDecorators) != 0 {
		n := p.hangingDecorators[0]
		return nil, p.errorAtLoc(n.Range(), ERR_DECORATOR_INVALID_POSITION)
	}

	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ExportDeclaration
func (p *Parser) exportDec() (Node, error) {
	rng := p.rng()
	tok := p.lexer.Next()
	if p.feat&FEAT_EXPORT_DEC == 0 {
		return nil, p.errorTok(tok)
	}

	var err error
	node := &ExportDec{N_STMT_EXPORT, span.Range{}, false, span.Range{}, nil, nil, nil, false}
	specs := make([]Node, 0, 3)
	tok = p.lexer.Peek()
	tv := tok.value
	if tv == T_MUL || tv == T_BRACE_L {
		ss, all, src, err := p.exportFrom(false)
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
		node.dec, err = p.fnDec(false, nil, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tv, tok.rng, ERR_ESCAPE_IN_KEYWORD)
		}
		node.dec, err = p.fnDec(false, tok, false)
		if err != nil {
			return nil, err
		}
	} else if tv == T_CLASS {
		node.dec, err = p.classDec(false, false, false, false)
		if err != nil {
			return nil, err
		}
	} else if tv == T_DEFAULT {
		node.def = p.lexer.Next().rng
		tok := p.lexer.Peek()
		tv = tok.value
		if tv == T_FUNC {
			node.dec, err = p.fnDec(false, nil, true)
		} else if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tv, tok.rng, ERR_ESCAPE_IN_KEYWORD)
			}
			node.dec, err = p.fnDec(false, tok, true)
		} else if tv == T_CLASS {
			node.dec, err = p.classDec(false, true, false, false)
		} else if ok, itf, _ := p.tsAheadIsAbstract(tok, false, false, false); ok {
			if itf {
				return nil, p.errorAtLoc(tok.rng, ERR_ABSTRACT_AT_INVALID_POSITION)
			}
			node.dec, err = p.classDec(false, true, false, true)
		} else if p.aheadIsTsItf(tok) {
			node.dec, err = p.tsItf()
			if err != nil {
				return nil, err
			}
		} else {
			node.dec, err = p.assignExpr(true, false, false, false)
			if err := p.advanceIfSemi(false); err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	} else if p.ts && tv == T_IMPORT {
		p.lexer.Next()
		id, err := p.ident(nil, true)
		if err != nil {
			return nil, err
		}
		n, err := p.tsImportAlias(rng, id, true)
		if err != nil {
			return nil, err
		}
		if n.Type() == N_TS_IMPORT_ALIAS {
			return n, nil
		}
		node.dec = n
	} else if ok, itf, _ := p.tsAheadIsAbstract(tok, false, false, false); ok {
		if itf {
			return nil, p.errorAtLoc(tok.rng, ERR_ABSTRACT_AT_INVALID_POSITION)
		}
		node.dec, err = p.classDec(false, false, false, true)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsItf(tok) {
		node.dec, err = p.tsItf()
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsTsTypDec(tok, false) {
		node.tsTyp = true
		rng := p.lexer.Next().rng // consume `type`
		// `export type { A };`
		if p.lexer.Peek().value == T_BRACE_L {
			ss, all, src, err := p.exportFrom(true)
			node.src = src
			node.all = all
			specs = append(specs, ss...)
			if err != nil {
				return nil, err
			}
			p.advanceIfSemi(false)
		} else {
			node.dec, err = p.tsTypDec(rng, false)
			if err != nil {
				return nil, err
			}
		}
	} else if p.aheadIsTsEnum(tok) {
		node.dec, err = p.tsEnum(span.Range{Lo: 1, Hi: 0}, false)
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
	} else if p.aheadIsModDec(tok) {
		node.dec, err = p.tsModDec()
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
		p.advanceIfSemi(false)
		return &TsExportAssign{N_TS_EXPORT_ASSIGN, p.finRng(rng), id}, nil
	} else {
		return nil, p.errorTok(tok)
	}

	node.rng = p.finRng(rng)
	node.specs = specs

	if err = p.addExp(node); err != nil {
		return nil, err
	}
	return node, nil
}

// add exp to its nearest upper scope whose type is one of:
// - top level scope
// - ts namespace
// - ts module
func (p *Parser) addExp(exp *ExportDec) error {
	// skip to record the method overloads as export
	if exp.dec != nil && exp.dec.Type() == N_STMT_FN && exp.dec.(*FnDec).body == nil {
		return nil
	}
	scope := p.scope()
	if scope.IsKind(SPK_GLOBAL) ||
		(scope.IsKind(SPK_TS_MODULE) && !scope.IsKind(SPK_TS_MODULE_INDIRECT)) {
		if scope.Exports == nil {
			scope.Exports = make([]*ExportDec, 0, 1)
		}
		scope.Exports = append(scope.Exports, exp)
		return nil
	}

	return p.errorAtLoc(exp.Range(), ERR_IMPORT_EXPORT_SHOULD_AT_TOP_LEVEL)
}

func (p *Parser) exportFrom(typ bool) ([]Node, bool, Node, error) {
	tok := p.lexer.Next()
	var specs []Node
	var err error

	ns := false
	if tok.value == T_MUL {
		ns = true
		ahead := p.lexer.Peek()
		if ahead.value == T_NAME && ahead.text == "as" && p.feat&FEAT_EXPORT_ALL_AS_NS != 0 {
			p.lexer.Next()

			id, err := p.ident(nil, false)
			if err != nil {
				return nil, false, nil, err
			}
			specs = make([]Node, 1)
			specs[0] = &ExportSpec{N_EXPORT_SPEC, p.finRng(tok.rng), true, id, nil, false}
		}
	} else {
		specs, err = p.exportNamed(typ)
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
		src = &StrLit{N_LIT_STR, p.finRng(str.rng), p.TokText(str), str.HasLegacyOctalEscapeSeq(), span.Range{}, nil}
	} else {
		// `export { default } from "a"` is legal
		// `export { default }` is illegal
		for _, spec := range specs {
			id := spec.(*ExportSpec).local.(*Ident)
			if id.kw {
				return nil, false, nil, p.errorAtLoc(id.rng, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
			}
		}
	}
	return specs, ns, src, nil
}

func (p *Parser) exportNamed(typ bool) ([]Node, error) {
	specs := make([]Node, 0, 3)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		spec, err := p.exportSpec(typ)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_COMMA {
			p.lexer.Next()
		} else if av == T_BRACE_R {
			break
		} else {
			return nil, p.errorAtLoc(ahead.rng, ERR_UNEXPECTED_TOKEN)
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
		return &Ident{N_NAME, p.finRng(ahead.rng), str, false, false, span.Range{}, true, p.newTypInfo(N_NAME)}, nil
	}
	return p.ident(scope, binding)
}

func (p *Parser) exportSpec(typ bool) (Node, error) {
	var local Node
	var err error
	rng := p.rng()
	ahead := p.lexer.Peek()
	if p.aheadIsTsTypDec(ahead, false) {
		rng := p.lexer.Next().rng // consume `type`
		typLoc := p.finRng(rng)
		if typ {
			return nil, p.errorAtLoc(rng, ERR_EXPORT_DUP_TYPE_MODIFIER)
		}
		typ = true
		local, err = p.tsTypDec(rng, true)
		if err != nil {
			return nil, err
		}

		// `export { type as as } from "./mod.js";`
		if local.Type() == N_NAME && local.(*Ident).val == "as" {
			_, canProp := p.lexer.Peek().CanBePropKey()
			ahead2nd := p.lexer.Peek2nd()

			// if `canProp2` is true, then the stmt may match:
			// `export { type as as if };`
			_, canProp2 := ahead2nd.CanBePropKey()
			if canProp && !canProp2 {
				local = &Ident{N_NAME, typLoc, "type", false, false, span.Range{}, false, p.newTypInfo(N_NAME)}
				id, err := p.identWithKw(nil, false)
				if err != nil {
					return nil, err
				}
				return &ExportSpec{N_EXPORT_SPEC, p.finRng(rng), false, local, id, false}, nil
			}
		}
	} else {
		local, err = p.identWithKw(nil, false)
		if err != nil {
			return nil, err
		}
	}

	id := local
	if p.aheadIsName("as") {
		tok := p.lexer.Next()
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
		}
		id, err = p.identWithKw(nil, false)
		if err != nil {
			return nil, err
		}
	}

	return &ExportSpec{N_EXPORT_SPEC, p.finRng(rng), false, local, id, typ}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-scripts-and-modules.html#prod-ImportDeclaration
func (p *Parser) importDec() (Node, error) {
	rng := p.rng()
	importTok := p.lexer.Peek()

	if p.feat&FEAT_IMPORT_DEC == 0 && p.feat&FEAT_DYNAMIC_IMPORT == 0 {
		return nil, p.errorTok(importTok)
	}

	ahead := p.lexer.Peek2nd()
	if ahead.value == T_PAREN_L || ahead.value == T_DOT {
		return p.exprStmt()
	}

	p.lexer.Next() // consume `import`

	specs := make([]Node, 0, 5)
	tok := p.lexer.Peek()
	node := &ImportDec{N_STMT_IMPORT, span.Range{}, specs, nil, false}

	// the second arg set to `true` for stmt like: `import type * as Types`
	typDec := p.aheadIsTsTypDec(tok, true)
	if typDec {
		ahead2nd := p.lexer.Peek2nd()
		typDec = !(ahead2nd.value == T_NAME && ahead2nd.text == "from")
	}

	if typDec {
		node.tsTyp = true
		typRng := p.lexer.Next().rng // consume `type`
		// `import type { A }`
		// `import type * as Types`
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACE_L || av == T_MUL {
			ss, err := p.importNamedOrNS(true)
			specs = append(specs, ss...)
			if err != nil {
				return nil, err
			}
		} else {
			// `import type A`
			tn, err := p.tsTypName(nil)
			if err != nil {
				return nil, err
			}

			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_ASSIGN {
				alias, err := p.tsImportAliasRsh(rng, typRng, tn, false)
				if err != nil {
					return nil, err
				}
				// `import type A = B.C;` is illegal
				// `import type a = require("a")` is legal
				if alias.Type() == N_TS_IMPORT_ALIAS {
					return nil, p.errorAtLoc(rng, ERR_IMPORT_TYPE_IN_IMPORT_ALIAS)
				}
				return alias, nil
			}

			spec := &ImportSpec{N_IMPORT_SPEC, p.finRng(tn.Range()), true, false, tn, tn, true}
			specs = append(specs, spec)

			if p.lexer.Peek().value == T_COMMA {
				return nil, p.errorAtLoc(p.lexer.Next().rng, ERR_IMPORT_TYP_MIX_NAMED)
			}
		}
	} else if tok.value != T_STRING {
		var id Node
		var err error
		if tok.value == T_NAME {
			id, err = p.ident(nil, true)
			if err != nil {
				return nil, err
			}
			spec := &ImportSpec{N_IMPORT_SPEC, p.finRng(tok.rng), true, false, id, id, false}
			specs = append(specs, spec)
		} else {
			ss, err := p.importNamedOrNS(false)
			if err != nil {
				return nil, err
			}
			specs = append(specs, ss...)
		}

		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
			ss, err := p.importNamedOrNS(false)
			if err != nil {
				return nil, err
			}
			specs = append(specs, ss...)
		}

		ahead := p.lexer.Peek()
		av := ahead.value
		if p.ts && (av == T_ASSIGN || (av == T_NAME && ahead.text != "from")) && id != nil {
			return p.tsImportAlias(rng, id, false)
		}
	}

	ahead = p.lexer.Peek()
	av := ahead.value
	if av == T_NAME && ahead.text == "from" && !ahead.ContainsEscape() {
		p.lexer.Next()
	}

	str, err := p.nextMustTok(T_STRING)
	if err != nil {
		return nil, err
	}
	legacyOctalEscapeSeq := str.HasLegacyOctalEscapeSeq()
	if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
		return nil, p.errorAtLoc(str.rng, ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
	}

	node.src = &StrLit{N_LIT_STR, p.finRng(str.rng), p.TokText(str), legacyOctalEscapeSeq, span.Range{}, nil}
	node.specs = specs
	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	if !typDec {
		for _, spec := range specs {
			s := spec.(*ImportSpec)
			ref := NewRef()
			ref.Id = s.local.(*Ident)
			ref.Dec = node
			ref.BindKind = BK_CONST
			ref.Typ = RDT_IMPORT
			if p.addLocalBinding(p.Symtab().Scopes[0], ref, true, ref.Id.val); err != nil {
				return nil, err
			}
		}
	}

	node.rng = p.finRng(rng)
	return node, nil
}

func (p *Parser) importNamedOrNS(typ bool) ([]Node, error) {
	tok := p.lexer.Peek()
	if tok.value == T_BRACE_L {
		return p.importNamed(typ)
	} else if tok.value == T_MUL {
		return p.importNS(typ)
	} else {
		return nil, p.errorTok(tok)
	}
}

func (p *Parser) importSpec(typ bool) (Node, error) {
	var binding Node
	var err error
	rng := p.rng()
	ahead := p.lexer.Peek()
	if p.aheadIsTsTypDec(ahead, false) {
		rng := p.lexer.Next().rng // consume `type`
		typLoc := p.finRng(rng)
		if typ {
			return nil, p.errorAtLoc(rng, ERR_EXPORT_DUP_TYPE_MODIFIER)
		}
		typ = true
		binding, err = p.tsTypDec(rng, true)
		if err != nil {
			return nil, err
		}

		// `export { type as as } from "./mod.js";`
		if binding.Type() == N_NAME && binding.(*Ident).val == "as" {
			_, canProp := p.lexer.Peek().CanBePropKey()
			ahead2nd := p.lexer.Peek2nd()

			// if `canProp2` is true, then the stmt may match:
			// `export { type as as if };`
			_, canProp2 := ahead2nd.CanBePropKey()
			if canProp && !canProp2 {
				binding = &Ident{N_NAME, typLoc, "type", false, false, span.Range{}, false, p.newTypInfo(N_NAME)}
				id, err := p.identWithKw(nil, false)
				if err != nil {
					return nil, err
				}
				return &ImportSpec{N_IMPORT_SPEC, p.finRng(rng), false, false, id, binding, false}, nil
			}
		}
	} else {
		binding, err = p.identWithKw(nil, false)
		if err != nil {
			return nil, err
		}
	}

	id := binding
	if p.aheadIsName("as") {
		p.lexer.Next()
		binding, err = p.ident(nil, true)
		if err != nil {
			return nil, err
		}
	} else if binding.Type() == N_NAME {
		// for statement like `import { true } from "bar"`, report `true` is a keyword
		id := binding.(*Ident)
		if id.kw {
			return nil, p.errorAtLoc(binding.Range(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, binding.(*Ident).val))
		}
	}

	return &ImportSpec{N_IMPORT_SPEC, p.finRng(rng), false, false, binding, id, typ}, nil
}

func (p *Parser) importNamed(typ bool) ([]Node, error) {
	p.lexer.Next()

	specs := make([]Node, 0, 5)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		}
		spec, err := p.importSpec(typ)
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

func (p *Parser) importNS(typ bool) ([]Node, error) {
	rng := p.rng()
	p.lexer.Next() // consume `*`
	_, err := p.nextMustName("as", false)
	if err != nil {
		return nil, err
	}

	id, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}

	specs := make([]Node, 1)
	specs[0] = &ImportSpec{N_IMPORT_SPEC, p.finRng(rng), false, true, id, nil, typ}
	return specs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-ClassDeclaration
func (p *Parser) classDec(expr bool, canNameOmitted bool, declare bool, abstract bool) (Node, error) {
	declare = declare || p.feat&FEAT_DTS != 0
	rng := p.lexer.Next().rng

	if abstract {
		p.lexer.Next()
	}

	ps := p.scope()
	// all parts of the class dec are in strict mode(include the id part)
	// here push an intermediate mode as strict to handle the id part
	p.lexer.PushMode(LM_STRICT, true)

	var id Node
	var err error
	ahead := p.lexer.Peek()
	av := ahead.value
	ti := p.newTypInfo(N_STMT_CLASS)

	var ds []Node
	if len(p.hangingDecorators) > 0 {
		ds = p.hangingDecorators
		rng = ds[0].Range()
		p.hangingDecorators = nil
		ti.decorators = ds
	}

	dec := &ClassDec{}
	if av != T_BRACE_L && av != T_EXTENDS {
		if p.ts && (av == T_LT || av == T_IMPLEMENTS) {
			// hit expr like:
			// - `(class<T> {});`
			// - `(class implements X.Y<T> {});`
		} else {
			id, err = p.identStrict(ps, true, true)
			if err != nil {
				return nil, err
			}
		}

		if ti != nil {
			typParams, err := p.tsTryTypParams()
			if err != nil {
				return nil, err
			}
			ti.SetTypParams(typParams)
		}

		if id != nil {
			ref := NewRef()
			ref.Id = id.(*Ident)
			ref.Dec = dec
			ref.BindKind = BK_CONST
			ref.Typ = RDT_CLASS | RDT_TYPE
			if err := p.addLocalBinding(ps, ref, true, ref.Id.val); err != nil {
				return nil, err
			}
		}
	}
	if !expr && !canNameOmitted && id == nil {
		return nil, p.errorAtLoc(p.rng(), ERR_CLASS_NAME_REQUIRED)
	}

	var super Node
	if p.lexer.Peek().value == T_EXTENDS {
		p.lexer.Next()
		scope := p.scope()
		scope.AddKind(SPK_CLASS_EXTEND_SUPER)
		super, err = p.lhs(false)
		scope.EraseKind(SPK_CLASS_EXTEND_SUPER)
		if err != nil {
			return nil, err
		}
		if err := p.checkDefaultVal(super, false, false, false); err != nil {
			return nil, err
		}
	}

	if p.ts {
		ti.SetAbstract(abstract)
		impl, err := p.tsImplements()
		if err != nil {
			return nil, err
		}
		ti.SetImplements(impl)
	}

	scope := p.symtab.EnterScope(true, false, true)
	p.enterStrict(true).AddKind(SPK_CLASS)
	if abstract {
		scope.AddKind(SPK_ABSTRACT_CLASS)
	} else {
		// erase the abstract flag inherits from parent scope
		scope.EraseKind(SPK_ABSTRACT_CLASS)
	}
	if super != nil {
		scope.AddKind(SPK_CLASS_HAS_SUPER)
	}

	body, err := p.classBody(declare, super != nil)
	if err != nil {
		return nil, err
	}

	s := p.symtab.LeaveScope()
	// balance the intermediate mode described above to handle
	// the id part of the class dec
	p.lexer.PopMode()

	typ := N_STMT_CLASS
	if expr {
		typ = N_EXPR_CLASS
	}

	dec.typ = typ
	dec.rng = p.finRng(rng)
	dec.id = id
	dec.super = super
	dec.body = body
	dec.declare = declare
	dec.ti = ti
	s.Node = dec
	return dec, nil
}

func (p *Parser) classBody(declare bool, hasSuper bool) (Node, error) {
	rng := p.rng()

	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}

	elems := make([]Node, 0, 3)
	hasCtor := false
	pvtNames := make(map[string]Node)
	scope := p.scope()
	for {
		tok := p.lexer.Peek()
		if tok.value == T_BRACE_R {
			break
		} else if tok.value == T_EOF {
			return nil, p.errorTok(tok)
		} else if p.aheadIsDecorator(tok) {
			ds, err := p.decorators()
			if err != nil {
				return nil, err
			}
			p.hangingDecorators = ds
			continue
		}
		if tok.value == T_SEMI {
			p.lexer.Next()
			continue
		}
		elem, err := p.classElem(declare)
		if err != nil {
			return nil, err
		}

		// attach decorators
		if len(p.hangingDecorators) > 0 {
			if wt, ok := elem.(NodeWithTypInfo); ok {
				ti := wt.TypInfo()
				if ti != nil {
					ti.decorators = p.hangingDecorators
					p.hangingDecorators = nil
				}
			}
		}

		if len(p.hangingDecorators) != 0 {
			n := p.hangingDecorators[0]
			return nil, p.errorAtLoc(n.Range(), ERR_DECORATOR_INVALID_POSITION)
		}

		if elem.Type() == N_METHOD {
			m := elem.(*Method)
			// `!m.Declare` is used to skip the constructor overloads
			if !m.Declare() && p.isName(m.key, "constructor", false, true) {
				if hasCtor {
					return nil, p.errorAtLoc(m.key.Range(), ERR_CTOR_DUP)
				}
				hasCtor = true
			}
		} else if elem.Type() == N_FIELD {
			f := elem.(*Field)
			if f.ti != nil && f.ti.Abstract() && !scope.IsKind(SPK_ABSTRACT_CLASS) {
				return nil, p.errorAtLoc(f.Range(), ERR_BARE_ABSTRACT_PROPERTY)
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
					if a.HasBody() && (a.kind == PK_SETTER || a.kind == PK_GETTER) && (b.kind == PK_SETTER || b.kind == PK_GETTER) {
						dup = a.static != b.static ||
							(a.kind == PK_GETTER && b.kind == PK_GETTER || a.kind == PK_SETTER && b.kind == PK_SETTER)
					}
				}
				if dup {
					return nil, p.errorAtLoc(key.Range(), fmt.Sprintf(ERR_TPL_ID_DUP_DEF, name))
				}
			}
			if elem.Type() == N_FIELD || elem.(*Method).HasBody() {
				pvtNames[name] = elem

				ref := NewRef()
				ref.Id = key.(*Ident)
				ref.Typ = RDT_PVT_FIELD
				ref.BindKind = BK_PVT_FIELD
				// skip check dup since getter/setter is dup but legal
				if err := p.addLocalBinding(nil, ref, false, name); err != nil {
					return nil, err
				}
			}
		}
		if p.ts {
			// do some TS related checks
			if wt, ok := elem.(NodeWithTypInfo); ok {
				override := wt.TypInfo().Override()
				if override && !hasSuper {
					return nil, p.errorAtLoc(elem.Range(), ERR_OVERRIDE_IN_NO_EXTEND)
				}
			}

			if elem.Type() == N_METHOD {
				m := elem.(*Method)
				if m.ti != nil && m.ti.Abstract() && !scope.IsKind(SPK_ABSTRACT_CLASS) {
					return nil, p.errorAtLoc(m.Range(), ERR_BARE_ABSTRACT_PROPERTY)
				}
				if m.key != nil && m.key.Type() == N_NAME && m.key.(*Ident).pvt {
					if m.ti.Abstract() {
						return nil, p.errorAtLoc(m.rng, ERR_PVT_ELEM_WITH_ABSTRACT)
					}
					if m.ti.AccMod() != ACC_MOD_NONE {
						return nil, p.errorAtLoc(m.rng, fmt.Sprintf(ERR_TPL_PVT_ELEM_WITH_ACCESS_MODIFIER, m.ti.AccMod().String()))
					}
				}
			} else if elem.Type() == N_FIELD {
				f := elem.(*Field)
				if f.ti != nil && f.ti.Abstract() && !scope.IsKind(SPK_ABSTRACT_CLASS) {
					return nil, p.errorAtLoc(f.Range(), ERR_BARE_ABSTRACT_PROPERTY)
				}
				if f.key != nil && f.key.Type() == N_NAME && f.key.(*Ident).pvt {
					if f.ti.Abstract() {
						return nil, p.errorAtLoc(f.rng, ERR_PVT_ELEM_WITH_ABSTRACT)
					}
					if f.ti.AccMod() != ACC_MOD_NONE {
						return nil, p.errorAtLoc(f.rng, fmt.Sprintf(ERR_TPL_PVT_ELEM_WITH_ACCESS_MODIFIER, f.ti.AccMod().String()))
					}
				}
			} else if elem.Type() == N_STATIC_BLOCK {
				n := elem.(*StaticBlock)
				ti := n.ti
				if ti.clsTyp.accMod != ACC_MOD_NONE || ti.clsTyp.abstract || ti.clsTyp.override {
					return nil, p.errorAtLoc(n.rng, ERR_STATIC_BLOCK_WITH_MODIFIER)
				}
			}
		}
		elems = append(elems, elem)
	}

	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}

	return &ClassBody{N_CLASS_BODY, p.finRng(rng), elems}, nil
}

func (p *Parser) modifiers() (begin, static, access, abstract, readonly, override, declare span.Range,
	isField, escape bool, name string, fieldLoc span.Range, accMod ACC_MOD, ahead *Token, mayStaticBlock bool) {
	for {
		ahead = p.lexer.Peek()
		av := ahead.value
		if av == T_STATIC && static.Empty() {
			tok := p.lexer.Next()
			static = tok.rng
			if begin.Empty() {
				begin = static
			}
			escape = tok.ContainsEscape()
			name = p.TokText(tok)
			isField, ahead = p.isField(true, false)
			fieldLoc = static
			if isField {
				static = span.Range{}
				return
			}
			mayStaticBlock = p.lexer.Peek().value == T_BRACE_L
			fieldLoc = static
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
				access = tok.rng
				if begin.Empty() {
					begin = access
				}
				escape = tok.ContainsEscape()
				name = p.TokText(tok)
				isField, ahead = p.isField(false, false)
				fieldLoc = access
				if isField {
					access = span.Range{}
					return
				}
				fieldLoc = access
			}
		} else if p.ts && abstract.Empty() && IsName(ahead, "abstract", false) {
			tok := p.lexer.Next()
			abstract = tok.rng
			if begin.Empty() {
				begin = abstract
			}
			escape = tok.ContainsEscape()
			name = p.TokText(tok)
			isField, ahead = p.isField(false, false)
			fieldLoc = abstract
			if isField {
				abstract = span.Range{}
				return
			}
			fieldLoc = abstract
		} else if p.ts && readonly.Empty() && IsName(ahead, "readonly", false) {
			tok := p.lexer.Next()
			readonly = tok.rng
			if begin.Empty() {
				begin = readonly
			}
			escape = tok.ContainsEscape()
			name = p.TokText(tok)
			isField, ahead = p.isField(false, false)
			fieldLoc = readonly
			if isField {
				readonly = span.Range{}
				return
			}
			fieldLoc = readonly
		} else if p.ts && override.Empty() && IsName(ahead, "override", false) {
			tok := p.lexer.Next()
			override = tok.rng
			if begin.Empty() {
				begin = override
			}
			escape = tok.ContainsEscape()
			name = p.TokText(tok)
			isField, ahead = p.isField(false, false)
			fieldLoc = override
			if isField {
				override = span.Range{}
				return
			}
			fieldLoc = override
		} else if p.ts && declare.Empty() && IsName(ahead, "declare", false) {
			tok := p.lexer.Next()
			declare = tok.rng
			if begin.Empty() {
				begin = declare
			}
			escape = tok.ContainsEscape()
			name = p.TokText(tok)
			isField, ahead = p.isField(false, false)
			fieldLoc = declare
			if isField {
				declare = span.Range{}
				return
			}
			fieldLoc = declare
		} else {
			break
		}
	}
	return
}

func (p *Parser) classElem(inDeclare bool) (Node, error) {
	beginLoc, staticLoc, accLoc, abstractLoc, readonlyLoc, overrideLoc, declareLoc, isField, escape, fieldName, fieldLoc, accMod, ahead, mayStaticBlock := p.modifiers()
	if err := p.tsModifierOrder(staticLoc, overrideLoc, readonlyLoc, accLoc, abstractLoc, declareLoc, accMod, mayStaticBlock); err != nil {
		return nil, err
	}

	static := !staticLoc.Empty()
	abstract := !abstractLoc.Empty()
	override := !overrideLoc.Empty()
	declare := !declareLoc.Empty()

	if static {
		if ahead.value == T_BRACE_L {
			blk, err := p.staticBlock(beginLoc)
			if err != nil {
				return nil, err
			}
			ti := blk.(*StaticBlock).ti
			if ti != nil {
				ti.SetAccMod(accMod)
				ti.SetAbstract(abstract)
				ti.SetDeclare(declare)
				ti.SetReadonly(!readonlyLoc.Empty())
				ti.SetOverride(override)
			}
			return blk, nil
		}
		if p.aheadIsArgList(ahead) {
			key := &Ident{N_NAME, fieldLoc, fieldName, false, escape, span.Range{}, true, p.newTypInfo(N_STMT_CLASS)}
			return p.method(beginLoc, key, accMod, span.Range{}, false, PK_METHOD, false, false, false, true, false, beginLoc, false, false, false, nil)
		}
	} else if isField {
		ti := p.newTypInfo(N_STMT_CLASS)
		key := &Ident{N_NAME, fieldLoc, fieldName, false, escape, span.Range{}, true, ti}
		if ti != nil {
			ti.ques, ti.not = p.tsAdvanceHook(true)
		}
		return p.field(key, beginLoc, span.Range{}, accMod, beginLoc, abstractLoc, readonlyLoc, overrideLoc, declareLoc, inDeclare)
	}

	if p.aheadIsAsync(ahead, true, true) {
		if !isField && !readonlyLoc.Empty() {
			return nil, p.errorAtLoc(readonlyLoc, ERR_METHOD_CANNOT_READONLY)
		}
		if ahead.ContainsEscape() {
			return nil, p.errorAt(ahead.value, ahead.rng, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(beginLoc, nil, accMod, span.Range{}, false, PK_METHOD, false, true, true, true, static, beginLoc, declare, abstract, override, nil)
	} else if p.aheadIsArgList(ahead) {
		key := &Ident{N_NAME, fieldLoc, fieldName, false, escape, span.Range{}, true, p.newTypInfo(N_STMT_CLASS)}
		return p.method(beginLoc, key, accMod, span.Range{}, false, PK_METHOD, false, false, false, true, false, beginLoc, false, false, false, nil)
	} else if ahead.value == T_MUL {
		if !isField && !readonlyLoc.Empty() {
			return nil, p.errorAtLoc(readonlyLoc, ERR_METHOD_CANNOT_READONLY)
		}
		if p.feat&FEAT_ASYNC_GENERATOR == 0 {
			return nil, p.errorTok(ahead)
		}
		return p.method(beginLoc, nil, accMod, span.Range{}, false, PK_METHOD, true, false, true, true, static, beginLoc, declare, abstract, override, nil)
	}

	propRng := ahead.rng
	kw := ahead.IsKw()
	if ahead.value == T_NAME || ahead.value == T_STRING || kw {
		ahead = p.lexer.Next()
		name := p.TokText(ahead)

		var key Node
		if ahead.value == T_STRING {
			key = &StrLit{N_LIT_STR, p.finRng(propRng), name, ahead.HasLegacyOctalEscapeSeq(), span.Range{}, nil}
		} else {
			key = &Ident{N_NAME, p.finRng(propRng), name, false, ahead.ContainsEscape(), span.Range{}, kw, nil}
		}

		ti := p.newTypInfo(N_STMT_CLASS)
		if ti != nil {
			ques, not := p.tsAdvanceHook(true)
			if !ques.Empty() && key.Type() == N_TS_THIS {
				return nil, p.errorAtLoc(ques, ERR_THIS_CANNOT_BE_OPTIONAL)
			}
			ti.SetQues(ques)
			ti.SetNot(not)
			key.(NodeWithTypInfo).SetTypInfo(ti)
		}

		isField, ahead = p.isField(false, name == "get" || name == "set")
		if isField {
			return p.field(key, beginLoc, staticLoc, accMod, beginLoc, abstractLoc, readonlyLoc, overrideLoc, declareLoc, inDeclare)
		}

		if !beginLoc.Empty() {
			propRng = beginLoc
		}

		kd := PK_INIT
		if p.aheadIsArgList(ahead) {
			kd = PK_METHOD
			if name == "constructor" {
				kd = PK_CTOR
			}

			if !isField && !readonlyLoc.Empty() {
				return nil, p.errorAtLoc(readonlyLoc, ERR_METHOD_CANNOT_READONLY)
			}
			if !declareLoc.Empty() {
				return nil, p.errorAtLoc(beginLoc, ERR_ILLEGAL_DECLARE_IN_CLASS)
			}
			return p.method(propRng, key, accMod, span.Range{}, false, kd, false, false, true, true, static, beginLoc, !declareLoc.Empty(), !abstractLoc.Empty(), !overrideLoc.Empty(), nil)
		}

		if name == "get" {
			kd = PK_GETTER
		} else if name == "set" {
			kd = PK_SETTER
		} else {
			return nil, p.errorTok(ahead)
		}

		if !declareLoc.Empty() {
			return nil, p.errorAtLoc(beginLoc, ERR_ILLEGAL_DECLARE_IN_CLASS)
		}

		if !isField && !readonlyLoc.Empty() {
			return nil, p.errorAtLoc(readonlyLoc, ERR_METHOD_CANNOT_READONLY)
		}
		return p.method(propRng, nil, accMod, span.Range{}, false, kd, false, false, true, true, static, beginLoc, !declareLoc.Empty(), !abstractLoc.Empty(), !overrideLoc.Empty(), nil)
	}

	return p.field(nil, beginLoc, staticLoc, accMod, beginLoc, abstractLoc, readonlyLoc, overrideLoc, declareLoc, inDeclare)
}

func (p *Parser) isName(node Node, name string, canContainsEscape bool, str bool) bool {
	nv := node.Type()
	hasEscape := false
	ns := ""
	if nv == N_LIT_STR && str {
		s := node.(*StrLit)
		ns = s.val
	} else if nv == N_NAME {
		id := node.(*Ident)
		ns = id.val
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

func (p *Parser) field(key Node, begin, static span.Range, accMod ACC_MOD, beginLoc, abstract, readonly, override, declare span.Range, inDeclare bool) (Node, error) {
	var rng span.Range
	var err error
	var compute span.Range
	var ti *TypInfo
	if key == nil {
		key, compute, err = p.classElemName()
		if err != nil {
			return nil, err
		}
		ti = p.newTypInfo(N_FIELD)
		if ti != nil {
			// computed prop can also be optional in ts, eg:
			// `class C { [Symbol.iterator]?(): void; }`
			ti.ques, ti.not = p.tsAdvanceHook(true)
		}
	} else if p.ts {
		if wt, ok := key.(NodeWithTypInfo); ok {
			// use clone here to decouple the typInfo of key and method
			ti = wt.TypInfo().Clone()
		}
	}

	if !begin.Empty() {
		rng = begin
	} else if !compute.Empty() {
		rng = compute
	} else {
		rng = key.Range()
	}

	if ti != nil {
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}

		ti.SetAccMod(accMod)
		ti.SetTypAnnot(typAnnot)
		ti.SetAbstract(!abstract.Empty())
		ti.SetDeclare(!declare.Empty())
		ti.SetReadonly(!readonly.Empty())
		ti.SetOverride(!override.Empty())
	}

	var value Node
	tok := p.lexer.Peek()
	var assignLoc span.Range
	if tok.value == T_ASSIGN {
		assignLoc = p.rng()
		p.lexer.Next()
		if !abstract.Empty() {
			return nil, p.errorAtLoc(assignLoc, ERR_ABSTRACT_PROP_WITH_INIT)
		}
		value, err = p.assignExpr(true, false, false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsArgList(tok) {
		if !static.Empty() {
			rng = static
		}
		return p.method(rng, key, accMod, compute, false, PK_METHOD, false, false, true, true, !static.Empty(), beginLoc, !declare.Empty(), !abstract.Empty(), !override.Empty(), ti)
	}
	p.advanceIfSemi(false)

	if p.ts && !compute.Empty() && ti.typAnnot != nil {
		if !abstract.Empty() {
			return nil, p.errorAtLoc(abstract, fmt.Sprintf(ERR_TPL_IDX_SIG_CANNOT_HAVE_MODIFIER, "abstract"))
		} else if !declare.Empty() {
			return nil, p.errorAtLoc(declare, fmt.Sprintf(ERR_TPL_IDX_SIG_CANNOT_HAVE_MODIFIER, "declare"))
		} else if !override.Empty() {
			return nil, p.errorAtLoc(override, fmt.Sprintf(ERR_TPL_IDX_SIG_CANNOT_HAVE_MODIFIER, "override"))
		} else if accMod != ACC_MOD_NONE {
			return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_IDX_SIG_CANNOT_HAVE_ACCESS, accMod.String()))
		}
	}

	if value != nil {
		p.checkName = true
		if err := p.checkDefaultVal(value, false, false, true); err != nil {
			return nil, err
		}
		p.checkName = false

		if p.ts {
			if inDeclare || !declare.Empty() || p.scope().IsKind(SPK_TS_DECLARE) {
				return nil, p.errorAtLoc(assignLoc, ERR_INIT_IN_ALLOWED_CTX)
			}
		}
	}

	isStatic := !static.Empty()
	if isStatic && p.isName(key, "prototype", false, true) {
		return nil, p.errorAtLoc(key.Range(), ERR_STATIC_PROP_PROTOTYPE)
	} else if p.isName(key, "constructor", false, true) {
		return nil, p.errorAtLoc(key.Range(), ERR_CTOR_CANNOT_BE_Field)
	}

	var ds []Node
	if len(p.hangingDecorators) > 0 {
		ds = p.hangingDecorators
		rng = ds[0].Range()
		p.hangingDecorators = nil
		ti.decorators = ds
	}
	return &Field{N_FIELD, p.finRng(rng), key, isStatic, !compute.Empty(), value, ti}, nil
}

func (p *Parser) classElemName() (Node, span.Range, error) {
	return p.propName(true, false, false)
}

func (p *Parser) staticBlock(static span.Range) (Node, error) {
	block, err := p.blockStmt(true, SPK_NONE)
	if err != nil {
		return nil, err
	}
	return &StaticBlock{N_STATIC_BLOCK, p.finRng(static), block.body, p.newTypInfo(N_STATIC_BLOCK)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-EmptyStatement
func (p *Parser) emptyStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()
	return &EmptyStmt{N_STMT_EMPTY, p.finRng(rng)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WithStatement
func (p *Parser) withStmt() (Node, error) {
	rng := p.rng()
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
		return nil, p.errorAtLoc(p.finRng(rng), ERR_WITH_STMT_IN_STRICT)
	}

	return &WithStmt{N_STMT_WITH, p.finRng(rng), expr, body}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DebuggerStatement
func (p *Parser) debugStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	return &DebugStmt{N_STMT_DEBUG, p.finRng(rng)}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-TryStatement
func (p *Parser) tryStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tryStmt := &TryStmt{N_STMT_TRY, span.Range{}, nil, nil, nil}
	p.pushTryStk(tryStmt)

	try, err := p.blockStmt(true, SPK_TRY)
	if err != nil {
		return nil, err
	}

	tok := p.lexer.Peek()
	if tok.value != T_CATCH && tok.value != T_FINALLY {
		return nil, p.errorTok(tok)
	}

	p.popTryStk()
	var catch Node
	if tok.value == T_CATCH {
		rng := p.rng()
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
			p.tsNodeTypAnnot(param, typAnnot, ACC_MOD_NONE, span.Range{}, false, false, false, false, span.Range{})

			if _, err := p.nextMustTok(T_PAREN_R); err != nil {
				return nil, err
			}
		} else if p.feat&FEAT_OPT_CATCH_PARAM == 0 {
			return nil, p.errorTok(ahead)
		}

		scope := p.symtab.EnterScope(false, false, true)
		scope.AddKind(SPK_CATCH)

		if param != nil {
			names, _, _ := p.collectNames([]Node{param})
			for _, nameNode := range names {
				id := nameNode.(*Ident)
				if ok := p.isProhibitedName(nil, id.val, true, true, false, false); ok {
					return nil, p.errorAtLoc(id.Range(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
				}
				ref := NewRef()
				ref.Id = id
				ref.BindKind = BK_LET
				if err := p.addLocalBinding(nil, ref, true, id.val); err != nil {
					return nil, err
				}
			}
		}

		body, err := p.blockStmt(false, SPK_NONE)
		s := p.symtab.LeaveScope()

		if err != nil {
			return nil, err
		}

		catch = &Catch{N_CATCH, p.finRng(rng), param, body}
		s.Node = catch
	}

	var fin Node
	if p.lexer.Peek().value == T_FINALLY {
		p.lexer.Next()
		fin, err = p.blockStmt(true, SPK_NONE)
		if err != nil {
			return nil, err
		}
	}

	tryStmt.rng = p.finRng(rng)
	tryStmt.try = try
	tryStmt.catch = catch
	tryStmt.fin = fin
	return tryStmt, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ThrowStatement
func (p *Parser) throwStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tok := p.lexer.Peek()
	var arg Node
	var err error
	if tok.value != T_ILLEGAL && tok.value != T_EOF && !tok.afterLineTerm {
		arg, err = p.expr()
		if err != nil {
			return nil, err
		}
	} else {
		if tok.afterLineTerm {
			return nil, p.errorAtLoc(rng, ERR_ILLEGAL_NEWLINE_AFTER_THROW)
		}
		return nil, p.errorTok(tok)
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	try := p.lastTry()
	if try == nil {
		try = p.prog
	}
	return &ThrowStmt{N_STMT_THROW, p.finRng(rng), arg, try}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ReturnStatement
func (p *Parser) retStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tok := p.lexer.Peek()
	var arg Node
	var err error

	closed := false
	if tok.value == T_SEMI {
		p.lexer.Next()
		closed = true
	} else if tok.value != T_ILLEGAL &&
		tok.value != T_BRACE_R &&
		tok.value != T_PAREN_R &&
		tok.value != T_BRACKET_R &&
		tok.value != T_EOF && !tok.afterLineTerm {
		arg, err = p.expr()
		if err != nil {
			return nil, err
		}
	}

	if err := p.advanceIfSemi(!closed); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_FUNC) && !scope.IsKind(SPK_FUNC_INDIRECT) {
		return nil, p.errorAtLoc(rng, ERR_ILLEGAL_RETURN)
	}

	return p.pushRetsStk(&RetStmt{N_STMT_RET, p.finRng(rng), arg}), nil
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
	rng := p.rng()
	label, err := p.ident(nil, false)
	if err != nil {
		return nil, err
	}

	// advance `:`
	p.lexer.Next()

	scope := p.scope()
	labelName := label.val
	if scope.HasLabel(labelName) {
		return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_DUP_LABEL, labelName))
	}

	node := &LabelStmt{N_STMT_LABEL, span.Range{}, label, nil, false}
	scope.uniqueLabels[labelName] = node
	scope.Labels = append(scope.Labels, node)

	scope.AddKind(SPK_INTERIM)
	body, err := p.stmt()
	scope.EraseKind(SPK_INTERIM)
	if err != nil {
		return nil, err
	}

	node.rng = p.finRng(rng)
	node.body = body
	// reset to check next label chain
	scope.uniqueLabels = make(map[string]Node)
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-BreakStatement
func (p *Parser) brkStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tok := p.lexer.Peek()
	if tok.value == T_NAME && !tok.afterLineTerm {
		label, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}

		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}

		target := p.scope().GetLabel(label.val)
		if target == nil {
			return nil, p.errorAtLoc(label.rng, fmt.Sprintf(ERR_UNDEF_LABEL, label.val))
		} else {
			target.(*LabelStmt).used = true
		}

		return &BrkStmt{N_STMT_BRK, p.finRng(rng), label, target}, nil
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_LOOP_DIRECT) &&
		!scope.IsKind(SPK_LOOP_INDIRECT) &&
		!scope.IsKind(SPK_SWITCH) &&
		!scope.IsKind(SPK_SWITCH_INDIRECT) {
		return nil, p.errorAtLoc(rng, ERR_ILLEGAL_BREAK)
	}
	return &BrkStmt{N_STMT_BRK, p.finRng(rng), nil, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ContinueStatement
func (p *Parser) contStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tok := p.lexer.Peek()
	if tok.value == T_NAME && !tok.afterLineTerm {
		label, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}

		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}

		ln := label.val
		target := p.scope().GetLabel(ln)
		if target == nil {
			return nil, p.errorAtLoc(label.rng, fmt.Sprintf(ERR_UNDEF_LABEL, label.val))
		} else {
			target.(*LabelStmt).used = true
		}

		return &ContStmt{N_STMT_CONT, p.finRng(rng), label, target}, nil
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	scope := p.scope()
	if !scope.IsKind(SPK_LOOP_DIRECT) && !scope.IsKind(SPK_LOOP_INDIRECT) {
		return nil, p.errorAtLoc(rng, ERR_ILLEGAL_CONTINUE)
	}

	return &ContStmt{N_STMT_CONT, p.finRng(rng), nil, p.loopStk[len(p.loopStk)-1]}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-SwitchStatement
func (p *Parser) switchStmt() (Node, error) {
	rng := p.rng()
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

	cases := make([]Node, 0, 3)
	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}
	metDefault := false

	scope := p.symtab.EnterScope(false, false, true)
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
			return nil, p.errorAt(tok.value, tok.rng, ERR_MULTI_DEFAULT)
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

	sw := &SwitchStmt{N_STMT_SWITCH, p.finRng(rng), test, cases}
	s := p.symtab.LeaveScope()
	s.Node = sw
	return sw, nil
}

func (p *Parser) switchCase(tok *Token) (*SwitchCase, error) {
	rng := p.lexer.Next().rng

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

	cons := make([]Node, 0, 3)
	for {
		tok := p.lexer.PeekStmtBegin()
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
	return &SwitchCase{N_SWITCH_CASE, p.finRng(rng), test, cons}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-IfStatement
func (p *Parser) ifStmt() (Node, error) {
	rng := p.rng()
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
	return &IfStmt{N_STMT_IF, p.finRng(rng), test, cons, alt}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-DoWhileStatement
func (p *Parser) doWhileStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	tok := p.lexer.PeekStmtBegin()

	loop := &DoWhileStmt{N_STMT_DO_WHILE, p.finRng(rng), nil, nil}
	p.pushLoopStk(loop)

	scope := p.symtab.EnterScope(false, false, true)
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

	s := p.symtab.LeaveScope()
	s.Node = loop
	p.popLoopStk()

	loop.rng = p.finRng(rng)
	loop.test = test
	loop.body = body
	return loop, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-WhileStatement
func (p *Parser) whileStmt() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	loop := &WhileStmt{N_STMT_WHILE, span.Range{}, nil, nil}
	p.pushLoopStk(loop)

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

	scope := p.symtab.EnterScope(false, false, true)
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

	s := p.symtab.LeaveScope()
	s.Node = loop
	p.popLoopStk()

	loop.rng = p.finRng(rng)
	loop.test = test
	loop.body = body
	return loop, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-ForStatement
func (p *Parser) forStmt() (Node, error) {
	rng := p.rng()
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

	scope := p.symtab.EnterScope(false, false, true)
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
			return nil, p.errorAt(T_IN, tok.rng, "")
		}
		return nil, p.errorTok(tok)
	} else if isOf && !await && init.Type() == N_NAME && init.(*Ident).val == "async" {
		return nil, p.errorAtLoc(init.Range(), ERR_LHS_OF_FOR_OF_CANNOT_ASYNC)
	}

	if isIn || isOf {
		loop := &ForInOfStmt{N_STMT_FOR_IN_OF, span.Range{}, false, false, nil, nil, nil}
		p.pushLoopStk(loop)

		if init == nil {
			return nil, p.errorTok(tok)
		}

		it := init.Type()
		if it != N_STMT_VAR_DEC {
			if !p.isSimpleLVal(init, true, false, true, false) {
				return nil, p.errorAtLoc(init.Range(), ERR_ASSIGN_TO_RVALUE)
			}

			// do the `argToParam` check only if the type of init is LitObj or LitArr otherwise
			// just check their simplicity
			if it == N_LIT_OBJ || it == N_LIT_ARR {
				if init, err = p.argToParam(init, 0, false, true, false); err != nil {
					return nil, err
				}
			} else if !p.isSimpleLVal(init, true, false, true, false) {
				return nil, p.errorAtLoc(init.Range(), ERR_ASSIGN_TO_RVALUE)
			}
		} else if it == N_STMT_VAR_DEC {
			varDec := init.(*VarDecStmt)
			if len(varDec.decList) > 1 {
				return nil, p.errorAtLoc(varDec.decList[1].Range(), ERR_DUP_BINDING)
			}
			if p.scope().IsKind(SPK_STRICT) {
				for _, dec := range varDec.decList {
					d := dec.(*VarDec)
					if d.init != nil {
						et := ERR_FOR_OF_LOOP_HAS_INIT
						if isIn {
							et = ERR_FOR_IN_LOOP_HAS_INIT
						}
						return nil, p.errorAtLoc(varDec.rng, et)
					}
				}
			}
		}

		revise := T_IN
		if !isIn {
			revise = T_OF
		}
		p.lexer.NextRevise(revise)

		var right Node
		if isOf {
			right, err = p.assignExpr(true, false, false, false)
		} else {
			right, err = p.expr()
		}
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

		s := p.symtab.LeaveScope()
		s.Node = loop

		p.popLoopStk()
		loop.rng = p.finRng(rng)
		loop.in = isIn
		loop.await = await
		loop.left = init
		loop.right = right
		loop.body = body
		return loop, nil
	}

	loopNode := &ForStmt{N_STMT_FOR, span.Range{}, nil, nil, nil, nil}
	p.pushLoopStk(loopNode)

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

	s := p.symtab.LeaveScope()
	s.Node = loopNode

	p.popLoopStk()
	loopNode.rng = p.finRng(rng)
	loopNode.init = init
	loopNode.test = test
	loopNode.update = update
	loopNode.body = body
	return loopNode, nil
}

func (p *Parser) aheadIsAsync(tok *Token, prop bool, pvt bool) bool {
	if p.feat&FEAT_ASYNC_AWAIT != 0 && IsName(tok, "async", true) {
		ahead := p.lexer.Peek2nd()
		if ahead.afterLineTerm {
			return false
		}
		if ahead.value == T_FUNC ||
			(p.aheadIsArgList(ahead) && !prop) ||
			ahead.value == T_MUL {
			return true
		}
		_, canProp := ahead.CanBePropKey()
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
	rng := p.rng()

	// below value cache is needed since token is saved in the ring-buffer
	// so the advanced token maybe override by next token peek
	asyncHasEscape := false
	var asyncLoc span.Range
	if async != nil {
		asyncHasEscape = async.ContainsEscape()
		rng = async.rng
		asyncLoc = async.rng

		p.lexer.Next()
		p.finRng(asyncLoc)
	}
	tok := p.lexer.Peek()
	fn := false
	if tok.value == T_FUNC {
		fn = true
		p.lexer.Next()
	}

	tok = p.lexer.Peek()
	generator := tok.value == T_MUL
	genLoc := tok.rng
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
			fnRef.Id = id.(*Ident)
			fnRef.BindKind = BK_VAR
			fnRef.Typ = RDT_FN
			if ps.IsKind(SPK_STRICT) {
				fnRef.BindKind = BK_LET
			}
			if ps.IsKind(SPK_TS_DECLARE) {
				fnRef.Typ = fnRef.Typ.On(RDT_TYPE)
			}
			if err := p.addLocalBinding(ps, fnRef, ps.IsKind(SPK_STRICT), fnRef.Id.val); err != nil {
				return nil, err
			}
		}
	}
	if fn && !expr && !canNameOmitted && id == nil {
		return nil, p.errorTok(tok)
	}

	scope := p.symtab.EnterScope(true, false, true)
	if async != nil {
		scope.AddKind(SPK_ASYNC)
		p.lexer.AddMode(LM_ASYNC)
	}
	// 'yield' as function names
	if generator {
		p.scope().AddKind(SPK_GENERATOR)
	}
	p.incRetsStk()

	var args []Node
	var typArgs Node
	ahead := p.lexer.Peek()
	if id != nil && ahead.value == T_ARROW && !ahead.afterLineTerm {
		// async a => {}
		args = make([]Node, 1)
		args[0] = id
	} else if fn {
		args, typArgs, _, err = p.paramList(false, PK_INIT, true)
		if err != nil {
			return nil, err
		}
	} else {
		// the arg check is skipped here, the correctness of args is guaranteed by
		// below `argsToFormalParams`
		p.checkName = false
		scope.AddKind(SPK_FORMAL_PARAMS)
		args, _, typArgs, _, err = p.argList(false, false, asyncLoc, false)
		scope.EraseKind(SPK_FORMAL_PARAMS)
		p.checkName = true
		if err != nil {
			return nil, err
		}
		if typArgs != nil && typArgs.Type() != N_TS_PARAM_INST {
			if err := p.advanceIfSemi(true); err != nil {
				return nil, err
			}
			return &ExprStmt{N_STMT_EXPR, p.finRng(rng), typArgs, false}, nil
		}
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	ti := p.newTypInfo(N_STMT_FN)
	if ti != nil {
		ti.SetTypAnnot(typAnnot)
	}

	if generator {
		p.lexer.AddMode(LM_GENERATOR)
	}

	tok = p.lexer.Peek()
	arrow := false
	var arrowLoc span.Range
	if tok.value == T_ARROW && !tok.afterLineTerm {
		if !fn {
			arrowLoc = p.lexer.Next().rng
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
	var firstComplicated span.Range
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
			ref.Id = paramName.(*Ident)
			ref.BindKind = BK_PARAM
			// duplicate-checking for params is enable in strict and delegated to below `checkParams`
			p.addLocalBinding(nil, ref, false, ref.Id.val)
		}

		if tok.value == T_BRACE_L {
			if body, err = p.fnBody(); err != nil {
				return nil, err
			}
		} else if expr || arrow {
			if body, err = p.expr(); err != nil {
				return nil, err
			}
		} else if fn && !expr && tok.value == T_FUNC && tok.afterLineTerm && p.ts {
			// ts func overloads:
			// `function f(a:number)`
			// `function f(): any {}`
			if err = p.tsIsFnSigValid(fnRef.Id.val); err != nil {
				return nil, err
			}
		} else if (tok.value == T_SEMI || tok.afterLineTerm || tok.value == T_EOF) && p.ts {
			// AmbientFunctionDeclaration
			// `declare function a();`
		} else {
			return nil, p.errorTok(tok)
		}
	} else if async != nil {
		// this branch means the input is callExpr like:
		// `async ({a: b = c})` callExpr
		// `async* ({a: b = c})` binExpr
		lhs := &Ident{N_NAME, asyncLoc, "async", false, asyncHasEscape, span.Range{}, true, p.newTypInfo(N_NAME)}

		var exp Node
		if generator {
			var rhs Node
			argsLen := len(args)
			if argsLen == 0 {
				return nil, p.errorTok(tok)
			} else if argsLen == 1 {
				rhs = args[0]
			} else {
				rhs = &SeqExpr{N_EXPR_SEQ, p.finRng(rng), args, span.Range{}}
			}
			exp = &BinExpr{N_EXPR_BIN, p.finRng(rng), T_MUL, genLoc, lhs, rhs, span.Range{}}
		} else {
			if err := p.checkArgs(args, false, true); err != nil {
				return nil, err
			}
			ti := p.newTypInfo(N_EXPR_CALL)
			if ti != nil {
				// typArgs is produced by `argList` in this branch, so it's required
				// to do a typeParam to typeArg transformation
				if err = p.tsCheckTypArgs(typArgs); err != nil {
					return nil, err
				}
				ti.SetTypArgs(typArgs)
			}
			exp = &CallExpr{N_EXPR_CALL, p.finRng(rng), lhs, args, false, span.Range{}, ti}
		}

		if !expr {
			binExpr, err := p.binExpr(exp, 0, false, false, false, false)
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
			return &ExprStmt{N_STMT_EXPR, p.finRng(rng), seq, false}, nil
		}
		return exp, nil
	} else {
		return nil, p.errorTok(tok)
	}

	if generator {
		p.lexer.PopMode()
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)
	if id != nil {
		name := id.(*Ident).val
		if p.isProhibitedName(idScope, name, isStrict, true, false, false) {
			return nil, p.errorAtLoc(id.Range(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
		}
	}

	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	if p.ts {
		if body == nil {
			if fnRef != nil {
				if !scope.IsKind(SPK_TS_DECLARE) {
					// suppress the dup-checking of the binding name on the func overloads
					ps.DelLocal(fnRef)
				}
			}
			p.advanceIfSemi(false)
		}

		if (scope.IsKind(SPK_TS_DECLARE) || p.feat&FEAT_DTS != 0) && body != nil {
			return nil, p.errorAtLoc(body.Range(), ERR_IMPL_IN_AMBIENT_CTX)
		}

		opts := NewTsCheckParamOpts()
		opts.impl = body != nil
		if err = p.tsCheckParams(params, opts); err != nil {
			return nil, err
		}
	}

	s := p.symtab.LeaveScope()

	if ti != nil {
		if !fn {
			typArgs, err = p.tsTypArgsToTypParams(typArgs)
			if err != nil {
				return nil, err
			}
		}
		ti.SetTypParams(typArgs)
	}

	if arrow {
		fn := &ArrowFn{N_EXPR_ARROW, p.finRng(rng), arrowLoc, async != nil, params, body, body.Type() != N_STMT_BLOCK, p.decRetsStk(), span.Range{}, ti}
		s.Node = fn
		if expr {
			return fn, nil
		}
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
		return &ExprStmt{N_STMT_EXPR, p.finRng(rng), fn, false}, nil
	}

	typ := N_STMT_FN
	if expr {
		typ = N_EXPR_FN
	}

	fnDec := &FnDec{typ, p.finRng(rng), id, generator, async != nil, params, body, p.decRetsStk(), span.Range{}, ti}
	if fnRef != nil {
		fnRef.Dec = fnDec
	}
	s.Node = fnDec

	if expr && p.lexer.Peek().value == T_PAREN_L {
		node, _, err := p.callExpr(fnDec, true, false, span.Range{}, false)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	if !expr && p.ts {
		if body == nil {
			p.lastTsFnSig = fnDec
		} else if err = p.tsIsFnImplValid(id); err != nil {
			return nil, err
		}
	}
	return fnDec, nil
}

func (p *Parser) collectNames(nodes []Node) (names []Node, firstComplicated span.Range, err error) {
	names = make([]Node, 0, 5)
	var ns []Node
	for _, param := range nodes {
		if firstComplicated.Empty() && param.Type() != N_NAME {
			firstComplicated = param.Range()
		}
		ns, err = p.namesInPattern(param, false)
		if err != nil {
			return nil, span.Range{}, err
		}
		names = append(names, ns...)
	}
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#sec-parameter-lists-static-semantics-early-errors
// `isSimpleParamList` should be true if function body directly contains `use strict` directive
func (p *Parser) checkParams(names []Node, firstComplicated span.Range, isStrict bool, directStrict bool) error {
	var dupLoc span.Range
	unique := make(map[string]bool)
	for _, id := range names {
		name := id.(*Ident).val
		if p.isProhibitedName(nil, name, isStrict, true, false, false) {
			return p.errorAtLoc(id.Range(), fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, name))
		}

		if dupLoc.Empty() {
			if _, ok := unique[name]; ok {
				dupLoc = id.Range()

			} else {
				unique[name] = true
			}
		}
	}

	if directStrict && !firstComplicated.Empty() {
		return p.errorAtLoc(firstComplicated, ERR_STRICT_DIRECTIVE_AFTER_NOT_SIMPLE)
	}

	if !dupLoc.Empty() {
		if isStrict {
			return p.errorAtLoc(dupLoc, ERR_DUP_PARAM_NAME)
		}
		if !firstComplicated.Empty() {
			return p.errorAtLoc(dupLoc, ERR_DUP_PARAM_NAME)
		}
	}
	return nil
}

// `kw` means whether the name of node can be keyword or not
func (p *Parser) namesInPattern(node Node, kw bool) ([]Node, error) {
	out := make([]Node, 0, 10)
	if node == nil {
		return out, nil
	}
	switch node.Type() {
	case N_NAME:
		id := node.(*Ident)
		if !kw && id.kw {
			em := ERR_TPL_UNEXPECTED_TOKEN_TYPE
			if id.val == "eval" {
				em = ERR_TPL_BINDING_RESERVED_WORD
			}
			return nil, p.errorAtLoc(node.Range(), fmt.Sprintf(em, id.val))
		}
		out = append(out, node)
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
	case N_PAT_REST:
		names, errLoc := p.namesInPattern(node.(*RestPat).arg, kw)
		if errLoc != nil {
			return nil, errLoc
		}
		out = append(out, names...)
	}
	return out, nil
}

func (p *Parser) fnBody() (Node, error) {
	return p.blockStmt(false, SPK_NONE)
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
		str := p.RngText(expr.Range())
		if str == "\"use strict\"" || str == "'use strict'" {
			return true, true
		}
		return false, true
	}
	return false, false
}

func (p *Parser) enterStrict(lex bool) *Scope {
	if lex {
		p.lexer.AddMode(LM_STRICT)
	}
	return p.scope().AddKind(SPK_STRICT).AddKind(SPK_STRICT_DIR)
}

func (p *Parser) stmts(terminal TokenValue) ([]Node, error) {
	stmts := make([]Node, 0, 20)
	prologue := 0 // the index in above `stmts` contains the last stmt in Directive Prologue

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
		cmts := p.lexer.takeStmtCmts()
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}

		if len(cmts) > 0 {
			p.prevCmts[stmt] = cmts
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

					// lookbehind to check that exprs before the 'use strict' directive
					if prologue > 0 {
						for i := 0; i < prologue; i++ {
							stmt := stmts[i]
							if stmt.Type() == N_STMT_EXPR {
								expr := stmts[i].(*ExprStmt).expr
								if expr.Type() == N_LIT_STR && expr.(*StrLit).loSeq {
									return nil, p.errorAtLoc(expr.Range(), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
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

			cmts := p.lexer.takeExprCmts()
			if len(cmts) > 0 {
				p.postCmts[stmt] = cmts
			}
			stmts = append(stmts, stmt)
		}
	}
	return stmts, nil
}

func (p *Parser) blockStmt(newScope bool, scopeKind ScopeKind) (*BlockStmt, error) {
	tok, err := p.nextMustTok(T_BRACE_L)
	if err != nil {
		return nil, err
	}

	var scope *Scope
	if newScope {
		fn := scopeKind == SPK_TS_MODULE
		scope = p.symtab.EnterScope(fn, false, true)
		scope.AddKind(scopeKind)
	}
	rng := tok.rng

	stmts, err := p.stmts(T_BRACE_R)
	if err != nil {
		return nil, err
	}

	if newScope {
		if scope.IsKind(SPK_GLOBAL) ||
			(scope.IsKind(SPK_TS_MODULE) && !scope.IsKind(SPK_TS_MODULE_INDIRECT)) {
			if scope.Exports != nil {
				if err := p.checkExp(scope.Exports); err != nil {
					return nil, err
				}
			}
		}
		p.symtab.LeaveScope()
	}
	return &BlockStmt{N_STMT_BLOCK, p.finRng(rng), stmts, newScope}, nil
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
		if !ahead.afterLineTerm && (av == T_NAME ||
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
		return p.errorAtLoc(ref.Id.rng, fmt.Sprintf(ERR_TPL_ID_DUP_DEF, name))
	}
	return nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#prod-VariableStatement
func (p *Parser) varDecStmt(kind TokenValue, asExpr bool) (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	node := &VarDecStmt{N_STMT_VAR_DEC, span.Range{}, T_ILLEGAL, make([]Node, 0, 5), nil}

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
		return p.tsEnum(rng, isConst)
	}

	lvs := make([]Node, 0, 5)
	for {
		dec, err := p.varDec(bindKind != BK_VAR)
		if err != nil {
			return nil, err
		}
		lvs = append(lvs, dec.id)

		if !p.dts && isConst && dec.init == nil && !p.scope().IsKind(SPK_NOT_IN) && !(p.ts && p.scope().IsKind(SPK_TS_DECLARE)) {
			return nil, p.errorAtLoc(dec.rng, ERR_CONST_DEC_INIT_REQUIRED)
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
			return nil, p.errorAtLoc(id.rng, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
		}

		ref := NewRef()
		ref.Id = id
		ref.Dec = node
		ref.BindKind = bindKind
		if err := p.addLocalBinding(nil, ref, true, ref.Id.val); err != nil {
			return nil, err
		}
	}

	if !asExpr {
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
	}
	node.rng = p.finRng(rng)
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
	rng := binding.Range()
	scope.EraseKind(SPK_LEXICAL_DEC)

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.tsNodeTypAnnot(binding, typAnnot, ACC_MOD_NONE, span.Range{}, false, false, false, false, span.Range{})

	var init Node
	if p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		init, err = p.assignExpr(true, false, false, false)
		if err != nil {
			return nil, err
		}
	}

	if binding.Type() != N_NAME && init == nil && !p.scope().IsKind(SPK_NOT_IN) && !(p.ts && p.scope().IsKind(SPK_TS_DECLARE)) {
		return nil, p.errorAtLoc(p.rng(), ERR_COMPLEX_BINDING_MISSING_INIT)
	}

	return &VarDec{N_VAR_DEC, p.finRng(rng), binding, init}, nil
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

func (p *Parser) identStrict(scope *Scope, forceStrict bool, binding bool) (Node, error) {
	if scope == nil {
		scope = p.scope()
	}

	tok := p.lexer.Next()
	tv := tok.value
	if tv != T_NAME && !(tv > T_CTX_KEYWORD_BEGIN && tv < T_CTX_KEYWORD_END) {
		return nil, p.errorTok(tok)
	}

	name := p.TokText(tok)
	rng := p.finRng(tok.rng)

	if p.isProhibitedName(scope, name, true, false, false, forceStrict) {
		if binding {
			return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, name))
		}
		return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
	}

	// for reporting `'let' is disallowed as a lexically bound name` for stmt like `let let`
	if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) && !tok.ContainsEscape() {
		if name == "let" || name == "const" {
			return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_FORBIDDEN_LEXICAL_NAME, name))
		}
	}

	return &Ident{N_NAME, rng, name, false, tok.ContainsEscape(), span.Range{}, tok.IsKw(), p.newTypInfo(N_NAME)}, nil
}

func (p *Parser) ident(scope *Scope, binding bool) (*Ident, error) {
	id, err := p.identStrict(scope, false, binding)
	if err != nil {
		return nil, err
	}
	return id.(*Ident), nil
}

func (p *Parser) accMod() (accMod ACC_MOD, accLoc span.Range, abstractLoc, readonlyLoc, overrideLoc, declareLoc span.Range, beginLoc span.Range,
	isField, escape bool, name string, fieldLoc span.Range, err error) {
	if !p.ts {
		return
	}

	var staticLoc span.Range
	var mayStaticBlock bool
	beginLoc, staticLoc, accLoc, abstractLoc, readonlyLoc, overrideLoc, declareLoc, isField, escape, name, fieldLoc, accMod, _, mayStaticBlock = p.modifiers()
	if err = p.tsModifierOrder(staticLoc, overrideLoc, readonlyLoc, accLoc, abstractLoc, declareLoc, accMod, mayStaticBlock); err != nil {
		return
	}

	if !staticLoc.Empty() {
		err = p.errorAtLoc(staticLoc, ERR_UNEXPECTED_TOKEN)
		return
	}
	if !abstractLoc.Empty() {
		err = p.errorAtLoc(abstractLoc, ERR_UNEXPECTED_TOKEN)
		return
	}
	return
}

func (p *Parser) roughParam(ctor bool) (Node, error) {
	accMod, accLoc, abstract, readonlyLoc, overrideLoc, declareLoc, _, isField, escape, fieldName, fieldLoc, err := p.accMod()
	if err != nil {
		return nil, err
	}

	var name Node
	if isField {
		name = &Ident{N_NAME, fieldLoc, fieldName, false, escape, span.Range{}, false, p.newTypInfo(N_NAME)}
	} else {
		if p.ts && !ctor {
			if accMod != ACC_MOD_NONE {
				return nil, p.errorAtLoc(accLoc, ERR_ILLEGAL_PARAMETER_MODIFIER)
			}
			if !readonlyLoc.Empty() {
				return nil, p.errorAtLoc(readonlyLoc, ERR_ILLEGAL_PARAMETER_MODIFIER)
			}
		}
		name, err = p.tsTyp(true, false, false)
		if err != nil {
			return nil, err
		}
	}

	ques, _ := p.tsAdvanceHook(false)
	if !ques.Empty() && name.Type() == N_TS_THIS {
		return nil, p.errorAtLoc(ques, ERR_THIS_CANNOT_BE_OPTIONAL)
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}

	ti := p.newTypInfo(N_NAME)
	if ti != nil {
		ti.SetTypAnnot(typAnnot)
		ti.SetAccMod(accMod)
		ti.SetReadonly(!readonlyLoc.Empty())
		ti.SetOverride(!overrideLoc.Empty())
		ti.SetDeclare(!declareLoc.Empty())
		ti.SetAbstract(!abstract.Empty())
		ti.SetQues(ques)
	}

	var colonLoc span.Range
	if typAnnot != nil {
		colonLoc = typAnnot.Range()
	}
	return &TsRoughParam{N_TS_ROUGH_PARAM, p.finRng(name.Range()), name, colonLoc, ti}, nil
}

func (p *Parser) param(methodKind PropKind) (Node, error) {
	accMod, accLoc, abstract, readonly, override, declare, beginLoc, isField, escape, fieldName, fieldLoc, err := p.accMod()
	if err != nil {
		return nil, err
	}

	var binding Node
	var this bool
	if isField {
		binding = &Ident{N_NAME, fieldLoc, fieldName, false, escape, span.Range{}, false, p.newTypInfo(N_NAME)}
	} else {
		if p.ts {
			if accMod != ACC_MOD_NONE {
				ahead := p.lexer.Peek()
				av := ahead.value
				if av == T_BRACE_L || av == T_BRACKET_L {
					return nil, p.errorAtLoc(accLoc, ERR_PARAM_PROP_WITH_BINDING_PATTERN)
				}
			}
			if methodKind != PK_CTOR {
				if accMod != ACC_MOD_NONE {
					return nil, p.errorAtLoc(accLoc, ERR_ILLEGAL_PARAMETER_MODIFIER)
				}
				if !readonly.Empty() {
					return nil, p.errorAtLoc(readonly, ERR_ILLEGAL_PARAMETER_MODIFIER)
				}
			}
		}

		ahead := p.lexer.Peek()
		this = p.scope().IsKind(SPK_METHOD) && ahead.value == T_THIS
		if this {
			if methodKind == PK_GETTER || methodKind == PK_SETTER {
				return nil, p.errorAt(ahead.value, ahead.rng, ERR_GETTER_SETTER_WITH_THIS_PARAM)
			}
			rng := p.lexer.Next().rng
			binding = &Ident{N_NAME, p.finRng(rng), "this", false, false, span.Range{}, true, p.newTypInfo(N_NAME)}
		} else {
			binding, err = p.bindingPattern()
			if err != nil {
				return nil, err
			}
		}
	}

	if err != nil {
		return nil, err
	}
	rng := binding.Range()

	if ques, _ := p.tsAdvanceHook(false); !ques.Empty() {
		if wt, ok := binding.(NodeWithTypInfo); ok {
			wt.TypInfo().SetQues(ques)
		}
		if this {
			return nil, p.errorAtLoc(ques, ERR_THIS_CANNOT_BE_OPTIONAL)
		}
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.tsNodeTypAnnot(binding, typAnnot, accMod, beginLoc, !abstract.Empty(), !readonly.Empty(), !override.Empty(), !declare.Empty(), span.Range{})

	// default value
	if !this && p.lexer.Peek().value == T_ASSIGN {
		p.lexer.Next()
		value, err := p.assignExpr(true, false, false, false)
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
			return nil, p.errorAtLoc(r.arg.Range(), ERR_REST_CANNOT_SET_DEFAULT)
		}
		ap := &AssignPat{
			typ: N_PAT_ASSIGN,
			rng: p.finRng(rng),
			lhs: binding,
			rhs: value,
			ti:  p.newTypInfo(N_PAT_ASSIGN),
		}
		if p.ts {
			ap.hoistTypInfo()
		}
		binding = ap
	}

	return binding, nil
}

// `ctor` indicates this method is called when processing the constructor method of class,
// in that case the access modifier is needed to be considered if TS is enabled
func (p *Parser) paramList(firstRough bool, methodKind PropKind, typParams bool) ([]Node, Node, span.Range, error) {
	scope := p.scope()
	p.checkName = false
	ctor := methodKind == PK_CTOR
	scope.AddKind(SPK_FORMAL_PARAMS)

	var err error
	var tp Node
	if typParams {
		tp, err = p.tsTryTypParams()
		if err != nil {
			return nil, nil, span.Range{}, err
		}
		if ctor && tp != nil {
			return nil, nil, span.Range{}, p.errorAtLoc(tp.Range(), ERR_CTOR_CANNOT_WITH_TYPE_PARAMS)
		}
	}

	parenL, err := p.nextMustTok(T_PAREN_L)
	if err != nil {
		return nil, nil, span.Range{}, err
	}
	parenLoc := parenL.rng

	params := make([]Node, 0, 5)
	i := 0
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, nil, span.Range{}, p.errorTok(tok)
		}

		var param Node
		var err error
		if firstRough && i == 0 {
			param, err = p.roughParam(ctor)
		} else {
			ahead := p.lexer.Peek()
			var ds []Node
			var err error
			if scope.IsKind(SPK_METHOD) && scope.IsKind(SPK_FORMAL_PARAMS) && p.aheadIsDecorator(ahead) {
				ds, err = p.decorators()
				if err != nil {
					return nil, nil, span.Range{}, err
				}
			}

			param, err = p.param(methodKind)
			if err != nil {
				return nil, nil, span.Range{}, err
			}

			if len(ds) > 0 {
				if wt, ok := param.(NodeWithTypInfo); ok {
					ti := wt.TypInfo()
					if ti != nil {
						ti.decorators = ds
					}
				} else {
					return nil, nil, span.Range{}, p.errorAtLoc(ds[0].Range(), ERR_DECORATOR_INVALID_POSITION)
				}
			}
		}

		if err != nil {
			return nil, nil, span.Range{}, err
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
				return nil, nil, span.Range{}, p.errorAt(tok.value, tok.rng, msg)
			}
		} else if av != T_PAREN_R {
			return nil, nil, span.Range{}, p.errorTok(ahead)
		}
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, nil, span.Range{}, err
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
	rng := p.rng()
	p.lexer.Next()

	props := make([]Node, 0, 5)
	for {
		if p.lexer.Peek().value == T_BRACE_R {
			p.lexer.Next()
			break
		}

		node, err := p.patternProp()
		if err != nil {
			return nil, err
		}
		if node.Type() == N_PROP {
			prop := node.(*Prop)
			if prop.method {
				return nil, p.errorAtLoc(prop.rng, ERR_INVALID_DESTRUCTING_TARGET)
			}
		}

		tok := p.lexer.Peek()
		if node.Type() == N_PAT_REST && tok.value != T_BRACE_R {
			if tok.value == T_COMMA {
				return nil, p.errorAt(tok.value, tok.rng, ERR_REST_ELEM_MUST_LAST)
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
	return &ObjPat{N_PAT_OBJ, p.finRng(rng), props, span.Range{}, p.newTypInfo(N_PAT_OBJ)}, nil
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

	rng := compute
	if rng.Empty() {
		rng = key.Range()
	}

	tok := p.lexer.Peek()
	opLoc := tok.rng
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

	return &Prop{N_PROP, p.finRng(rng), key, opLoc, value, !compute.Empty(), false, shorthand, assign, PK_INIT, ACC_MOD_NONE}, nil
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
		av == T_BRACE_R ||
		av == T_PAREN_R ||
		(p.ts && (av == T_HOOK || av == T_NOT))

	if isField {
		return true, ahead
	}

	if getter {
		return !TokenKinds[av].StartExpr && av != T_NAME_PVT && av != T_BRACKET_L && !ahead.IsKw(), ahead
	}

	if static {
		return isField, ahead
	}

	return ahead.afterLineTerm, ahead
}

func (p *Parser) propName(allowNamePVT bool, maybeMethod bool, tsRough bool) (Node, span.Range, error) {
	var key Node
	tok := p.lexer.Next()
	rng := tok.rng
	kw, ok := tok.CanBePropKey()
	keyName := p.TokText(tok)

	scope := p.scope()
	var computeLoc span.Range
	tv := tok.value
	if allowNamePVT && tv == T_NAME_PVT {
		key = &Ident{N_NAME, p.finRng(rng), p.TokText(tok), true, tok.ContainsEscape(), span.Range{}, false, p.newTypInfo(N_NAME)}
	} else if tv == T_STRING {
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, span.Range{}, p.errorAtLoc(tok.rng, ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		key = &StrLit{N_LIT_STR, p.finRng(rng), p.TokText(tok), tok.HasLegacyOctalEscapeSeq(), span.Range{}, p.newTypInfo(N_LIT_STR)}
	} else if tv == T_NUM {
		key = &NumLit{N_LIT_NUM, p.finRng(rng), span.Range{}}
	} else if tv == T_BRACKET_L {
		computeLoc = tok.rng
		scope.AddKind(SPK_PROP_NAME)
		name, err := p.assignExpr(true, false, false, false)
		scope.EraseKind(SPK_PROP_NAME)
		if err != nil {
			return nil, span.Range{}, err
		}
		_, err = p.nextMustTok(T_BRACKET_R)
		if err != nil {
			return nil, span.Range{}, err
		}
		key = name
	} else if ok {
		if !kw && p.isProhibitedName(nil, keyName, true, false, false, false) {
			kw = true
		}
		// stmt `let { let } = {}` will raise error `let is disallowed as a lexically bound name` in sloppy mode
		if !scope.IsKind(SPK_STRICT) && scope.IsKind(SPK_LEXICAL_DEC) {
			if !tok.ContainsEscape() && (keyName == "let" || keyName == "const") {
				return nil, span.Range{}, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_FORBIDDEN_LEXICAL_NAME, keyName))
			}
		}
		key = &Ident{N_NAME, p.finRng(rng), keyName, false, tok.ContainsEscape(), span.Range{}, kw, p.newTypInfo(N_NAME)}
	} else {
		return nil, span.Range{}, p.errorTok(tok)
	}

	getter := keyName == "get" || keyName == "set"
	isField, ahead := p.isField(false, getter)
	if isField || !maybeMethod {
		return key, computeLoc, nil
	}

	kd := PK_INIT
	if getter {
		if tok.ContainsEscape() {
			return nil, span.Range{}, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
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

	m, err := p.method(rng, key, ACC_MOD_NONE, computeLoc, false, kd, false, false, false, false, false, span.Range{}, false, false, false, nil)
	if err != nil {
		return nil, span.Range{}, err
	}
	return m, span.Range{}, nil
}

func (p *Parser) patternArr() (Node, error) {
	rng := p.rng()
	p.lexer.Next()

	elems := make([]Node, 0, 5)
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
	return &ArrPat{N_PAT_ARRAY, p.finRng(rng), elems, span.Range{}, p.newTypInfo(N_PAT_ARRAY)}, nil
}

func (p *Parser) elision() []Node {
	ret := make([]Node, 0, 5)
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
	var opLoc span.Range
	if p.lexer.Peek().value == T_ASSIGN {
		tok := p.lexer.Next()
		opLoc = tok.rng
		init, err = p.assignExpr(true, false, false, false)
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

	rng := ident.Range()
	val := &AssignPat{N_PAT_ASSIGN, p.finRng(rng), ident, init, span.Range{}, p.newTypInfo(N_PAT_ASSIGN)}
	if !asProp {
		return val, nil
	}
	return &Prop{N_PROP, p.finRng(rng), val.lhs, opLoc, val, false, false, true, true, PK_INIT, ACC_MOD_NONE}, nil
}

// `arrPat` indicates whether `restExpr` is in array-pattern or not
func (p *Parser) patternRest(arrPat bool, allowNotLast bool) (Node, error) {
	rng := p.rng()
	tok := p.lexer.Next()

	if p.feat&FEAT_BINDING_REST_ELEM == 0 {
		return nil, p.errorTok(tok)
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	if av != T_NAME && (!arrPat || av != T_BRACKET_L && av != T_BRACE_L) {
		if av == T_BRACKET_L || av == T_BRACE_L {
			return nil, p.errorAt(ahead.value, ahead.rng, ERR_REST_ARG_NOT_BINDING_PATTERN)
		}
		return nil, p.errorTok(ahead)
	}

	arg, err := p.bindingPattern()
	if err != nil {
		return nil, err
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.tsNodeTypAnnot(arg, typAnnot, ACC_MOD_NONE, span.Range{}, false, false, false, false, span.Range{})

	// always allow the trailing comma in typescript
	if !allowNotLast && !p.ts {
		tok = p.lexer.Peek()
		if tok.value == T_COMMA {
			return nil, p.errorAt(tok.value, tok.rng, ERR_REST_ELEM_MUST_LAST)
		}
	}

	rest := &RestPat{N_PAT_REST, p.finRng(rng), arg, span.Range{}, p.newTypInfo(N_PAT_REST)}
	if p.ts {
		rest.hoistTypInfo()
	}
	return rest, nil
}

func (p *Parser) exprStmt() (Node, error) {
	rng := p.rng()
	stmt := &ExprStmt{N_STMT_EXPR, span.Range{}, nil, false}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	stmt.expr = expr

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	stmt.rng = p.finRng(rng)
	return stmt, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-Expression
func (p *Parser) expr() (Node, error) {
	return p.seqExpr(nil, false)
}

// `notGT` is `true` tells the later subroutine does NOT treat the `>` symbol as the greatThen operator
func (p *Parser) seqExpr(expr Node, notGT bool) (Node, error) {
	var err error
	if expr == nil {
		expr, err = p.assignExpr(true, notGT, false, false)
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

	rng := expr.Range()
	exprs := make([]Node, 0, 5)
	exprs = append(exprs, expr)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_COMMA {
			p.lexer.Next()
			expr, err = p.assignExpr(true, notGT, false, false)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.errorAt(p.lexer.PrevTok(), p.lexer.PrevTokRng(), "")
			}
			exprs = append(exprs, expr)
		} else {
			break
		}
	}
	return &SeqExpr{N_EXPR_SEQ, p.finRng(rng), exprs, span.Range{}}, nil
}

func (p *Parser) aheadIsYield() bool {
	if !p.scope().IsKind(SPK_GENERATOR) {
		return false
	}
	ahead := p.lexer.Peek()
	av := ahead.value
	return av == T_YIELD || av == T_NAME && p.RngText(ahead.rng) == "yield"
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-YieldExpression
func (p *Parser) yieldExpr() (Node, error) {
	rng := p.rng()
	tok := p.lexer.Next()

	if p.scope().IsKind(SPK_FORMAL_PARAMS) {
		return nil, p.errorAt(tok.value, tok.rng, ERR_YIELD_IN_FORMAL_PARAMS)
	}

	tok = p.lexer.Peek()
	kind := TokenKinds[tok.value]
	tv := tok.value
	startExpr := kind.StartExpr || p.feat&FEAT_JSX != 0 && tv == T_LT
	if tok.afterLineTerm || !startExpr && tv != T_MUL {
		return &YieldExpr{N_EXPR_YIELD, p.finRng(rng), false, nil, span.Range{}}, nil
	}

	delegate := false
	if p.lexer.Peek().value == T_MUL {
		p.lexer.Next()
		delegate = true
	}

	arg, err := p.assignExpr(true, false, false, true)
	if err != nil {
		return nil, err
	}
	return &YieldExpr{N_EXPR_YIELD, p.finRng(rng), delegate, arg, span.Range{}}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-AssignmentExpression
func (p *Parser) assignExpr(checkLhs bool, notGT bool, notHook bool, notColon bool) (Node, error) {
	if p.aheadIsYield() {
		return p.yieldExpr()
	}

	lhs, err := p.condExpr(nil, notGT, notHook, notColon)
	if err != nil {
		return nil, err
	}
	rng := lhs.Range()

	if !notColon && p.tsExprHasTypAnnot(lhs) {
		ques, _ := p.tsAdvanceHook(false)
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		p.tsNodeTypAnnot(lhs, typAnnot, ACC_MOD_NONE, span.Range{}, false, false, false, false, ques)
	}

	tok := p.lexer.Peek()
	if lhs.Type() == N_NAME && tok.value == T_ARROW && !tok.afterLineTerm {
		fn, err := p.arrowFn(rng, []Node{lhs}, nil, nil)
		if err != nil {
			return nil, err
		}
		lhs = fn
	}

	assign := p.advanceIfTokIn(T_ASSIGN_BEGIN, T_ASSIGN_END)
	if assign == nil {
		cmts := p.lexer.takeExprCmts()
		if len(cmts) > 0 {
			p.prevCmts[lhs] = cmts
		}
		return lhs, nil
	}
	op := assign.value
	opLoc := assign.rng

	rhs, err := p.assignExpr(checkLhs, notGT, false, false)
	if err != nil {
		return nil, err
	}

	cmts := p.lexer.takeExprCmts()
	if len(cmts) > 0 {
		p.prevCmts[rhs] = cmts
	}

	// set `depth` to 1 to permit expr like `i + 2 = 42`
	// and so just do the arg to param transform silently
	lhs, err = p.argToParam(lhs, 1, false, true, false)
	if err != nil {
		return nil, err
	}

	if checkLhs && !p.isSimpleLVal(lhs, true, false, true, op != T_ASSIGN) {
		return nil, p.errorAtLoc(lhs.Range(), ERR_ASSIGN_TO_RVALUE)
	}

	if err := p.checkArg(rhs, false, false); err != nil {
		return nil, err
	}

	node := &AssignExpr{N_EXPR_ASSIGN, p.finRng(rng), op, opLoc, lhs, rhs, span.Range{}, p.newTypInfo(N_EXPR_ASSIGN)}
	return node, nil
}

// https://tc39.es/ecma262/multipage/syntax-directed-operations.html#sec-static-semantics-assignmenttargettype
// `pat` indicates whether to treat the pattern syntax as legal or not
// `member` indicates whether the member expr can be treated as legal or not
// `optAssign` indicates whether the expr is the lhs of the op-assign expr
func (p *Parser) isSimpleLVal(expr Node, pat bool, inParen bool, member bool, optAssign bool) bool {
	switch expr.Type() {
	case N_NAME:
		node := expr.(*Ident)
		scope := p.scope()
		if scope.IsKind(SPK_ASYNC) && node.val == "await" {
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

func (p *Parser) advanceIfHook() *Token {
	tok := p.lexer.Peek()
	if !p.ts {
		if tok.value != T_HOOK {
			return nil
		}
		return p.lexer.Next()
	}
	hook := tok.value == T_HOOK
	if hook {
		ahead := p.lexer.Peek2nd()
		av := ahead.value
		if ahead.Kind().StartExpr || (p.feat&FEAT_JSX != 0 && av == T_LT) {
			return p.lexer.Next()
		}
	}
	return nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ConditionalExpression
func (p *Parser) condExpr(test Node, notGT bool, notHook bool, notColon bool) (Node, error) {
	var err error
	test, err = p.binExpr(test, 0, false, false, notGT, notColon)
	if err != nil {
		return nil, err
	}
	rng := test.Range()

	if notHook {
		return test, nil
	}

	hook := p.advanceIfHook()
	if hook == nil {
		return test, nil
	}

	// colon after ques maybe the leading of the typAnnot: `async (x?: number): any => x;`
	typAnnot, _ := p.tsTypAnnot()
	if typAnnot != nil {
		if wt, ok := test.(NodeWithTypInfo); ok {
			ti := wt.TypInfo()
			ti.SetQues(hook.rng)
			ti.SetTypAnnot(typAnnot)
		}
		return test, nil
	}

	cons, err := p.assignExpr(true, notGT, false, true)
	if err != nil {
		return nil, err
	}

	_, err = p.nextMustTok(T_COLON)
	if err != nil {
		return nil, err
	}

	alt, err := p.assignExpr(true, notGT, false, false)
	if err != nil {
		return nil, err
	}

	node := &CondExpr{N_EXPR_COND, p.finRng(rng), test, cons, alt, span.Range{}}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-functions-and-classes.html#prod-AwaitExpression
func (p *Parser) awaitExpr(tok *Token) (Node, error) {
	rng := tok.rng
	scope := p.scope()

	ahead := p.lexer.Peek()
	if !TokenKinds[ahead.value].StartExpr {
		if p.feat&FEAT_MODULE != 0 {
			// report friendly message for expr like: `async function foo(await) {}`
			if ahead.value == T_PAREN_R || ahead.value == T_COMMA {
				return nil, p.errorAtLoc(rng, fmt.Sprintf(ERR_TPL_BINDING_RESERVED_WORD, "await"))
			} else if !scope.IsKind(SPK_ASYNC) {
				return nil, p.errorAt(tok.value, tok.rng, ERR_AWAIT_OUTSIDE_ASYNC)
			}
			return nil, p.errorTok(ahead)
		}
		return &Ident{N_NAME, p.finRng(rng), "await", false, tok.ContainsEscape(), span.Range{}, true, p.newTypInfo(N_NAME)}, nil
	}

	if scope.IsKind(SPK_FORMAL_PARAMS) {
		return nil, p.errorAt(tok.value, tok.rng, ERR_AWAIT_IN_FORMAL_PARAMS)
	}
	if !scope.IsKind(SPK_ASYNC) {
		return nil, p.errorAt(tok.value, tok.rng, ERR_AWAIT_OUTSIDE_ASYNC)
	}

	arg, err := p.unaryExpr(nil, span.Range{}, false)
	if err != nil {
		return nil, err
	}
	return &UnaryExpr{N_EXPR_UNARY, p.finRng(rng), T_AWAIT, arg, span.Range{}}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-UnaryExpression
func (p *Parser) unaryExpr(typArgs Node, typArgsLoc span.Range, notColon bool) (Node, error) {
	var err error
	if typArgs == nil {
		typArgs, err = p.tsTryTypArgs(span.Range{}, false)
		if err != nil {
			if err == p.errTypArgMissingGT {
				return nil, p.errorAtLoc(p.rng(), ERR_UNEXPECTED_TOKEN)
			}
			return nil, err
		}
		if typArgs != nil && typArgs.Type() == N_JSX_ELEM {
			return typArgs, nil
		}
	}

	tok := p.lexer.Peek()
	rng := tok.rng
	op := tok.value
	if tok.IsUnary() || op == T_ADD || op == T_SUB || (op == T_LT && p.ts && p.feat&FEAT_JSX == 0) {
		if op != T_LT {
			p.lexer.Next()
		}
		arg, err := p.unaryExpr(nil, span.Range{}, notColon)
		if err != nil {
			return nil, err
		}

		// current es grammar does not allow the arrowFn to be the arg of
		// the unaryExpr such as `typeof () => {}` will raise an exception
		// `Malformed arrow function parameter list`
		if arg.Type() == N_EXPR_ARROW {
			return nil, p.errorAtLoc(arg.Range(), ERR_MALFORMED_ARROW_PARAM)
		}

		scope := p.scope()
		if scope.IsKind(SPK_STRICT) && tok.value == T_DELETE && arg.Type() == N_NAME {
			return nil, p.errorAtLoc(arg.Range(), ERR_DELETE_LOCAL_IN_STRICT)
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
					return nil, p.errorAtLoc(prop.rng, ERR_DELETE_PVT_FIELD)
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

		return &UnaryExpr{N_EXPR_UNARY, p.finRng(rng), op, arg, span.Range{}}, nil
	}

	if tok.value == T_AWAIT {
		if p.feat&FEAT_ASYNC_AWAIT == 0 {
			return nil, p.errorTok(tok)
		}
		p.lexer.Next()
		return p.awaitExpr(tok)
	}

	return p.updateExpr(typArgs, typArgsLoc, notColon)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-UpdateExpression
func (p *Parser) updateExpr(typArgs Node, typArgsLoc span.Range, notColon bool) (Node, error) {
	rng := p.rng()
	tok := p.lexer.Peek()
	if tok.value == T_INC || tok.value == T_DEC {
		p.lexer.Next()
		arg, err := p.unaryExpr(nil, span.Range{}, notColon)
		if err != nil {
			return nil, err
		}
		if !p.isSimpleLVal(arg, true, false, true, false) {
			return nil, p.errorAtLoc(arg.Range(), ERR_ASSIGN_TO_RVALUE)
		}
		ud := &UpdateExpr{N_EXPR_UPDATE, p.finRng(rng), tok.value, true, arg, span.Range{}}
		arg, err = p.tsTypAssert(ud, typArgs)
		if err != nil {
			return nil, err
		}
		return arg, nil
	}

	arg, err := p.lhs(notColon)
	if err != nil {
		return nil, err
	}

	tok = p.lexer.Peek()
	postfix := !tok.afterLineTerm && (tok.value == T_INC || tok.value == T_DEC)
	if !postfix {
		// for the type info before the arrow fn in this stmt `let a = <T, R>(a: T): void => { a++ }`,
		// the type info `<T, R>` is the typeParams of the arrowFn rather than typeAssert
		if arg.Type() == N_EXPR_ARROW {
			ti := arg.(NodeWithTypInfo).TypInfo()
			if ti != nil && typArgs != nil {
				typArgs, err = p.tsTypArgsToTypParams(typArgs)
				if err != nil {
					return nil, err
				}
				ti.SetTypParams(typArgs)
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
		return nil, p.errorAtLoc(arg.Range(), ERR_ASSIGN_TO_RVALUE)
	}

	p.lexer.Next()

	ud := &UpdateExpr{N_EXPR_UPDATE, p.finRng(rng), tok.value, false, arg, span.Range{}}
	ta, err := p.tsTypAssert(ud, typArgs)
	if err != nil {
		return nil, err
	}
	return ta, nil
}

func (p *Parser) lhs(notColon bool) (Node, error) {
	tok := p.lexer.Peek()
	var node Node
	var err error
	if tok.value == T_NEW {
		node, err = p.newExpr()
	} else {
		node, _, err = p.callExpr(nil, true, false, span.Range{}, notColon)
	}
	if err != nil {
		return nil, err
	}
	node = p.tsNoNull(node)
	if p.lexer.Peek().value == T_DOT {
		node, _, err = p.callExpr(node, true, false, span.Range{}, notColon)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-NewExpression
func (p *Parser) newExpr() (Node, error) {
	rng := p.rng()
	new := p.lexer.Next()

	var expr Node
	var err error

	scope := p.scope()
	tok := p.lexer.Peek()
	if tok.value == T_DOT && p.feat&FEAT_META_PROPERTY != 0 {
		meta := &Ident{N_NAME, p.finRng(new.rng), "new", false, new.ContainsEscape(), span.Range{}, true, p.newTypInfo(N_NAME)}
		p.lexer.Next() // consume dot

		id, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		if id.val != "target" {
			return nil, p.errorAtLoc(id.Range(), ERR_INVALID_META_PROP)
		}
		if !(scope.IsKind(SPK_CLASS) || scope.IsKind(SPK_CLASS_INDIRECT) ||
			scope.IsKind(SPK_ARROW) && scope.IsKind(SPK_FUNC_INDIRECT) ||
			(!scope.IsKind(SPK_ARROW) && (scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_FUNC_INDIRECT)))) {
			return nil, p.errorAtLoc(rng, ERR_META_PROP_OUTSIDE_FN)
		}

		expr = &MetaProp{N_META_PROP, p.finRng(rng), meta, id}
		return expr, nil
	}

	var optLoc span.Range
	expr, optLoc, err = p.memberExpr(nil, false, true, span.Range{})
	if err != nil {
		return nil, err
	}
	if !optLoc.Empty() {
		return nil, p.errorAtLoc(optLoc, ERR_OPT_EXPR_IN_NEW)
	}
	if expr.Type() == N_IMPORT_CALL {
		return nil, p.errorAtLoc(expr.Range(), ERR_DYNAMIC_IMPORT_CANNOT_NEW)
	}

	var args []Node
	var typArgs Node
	ti := p.newTypInfo(N_EXPR_NEW)
	ahead := p.lexer.Peek()
	if p.aheadIsArgList(ahead) {
		p.pushState()
		args, _, typArgs, _, err = p.argList(true, true, span.Range{}, false)
		if err != nil {
			// `new A < T`
			if err == p.errTypArgMissingGT {
				p.popState()

				e := err.(ErrTypArgMissingGT)
				p.addLtTok(e.rng.Lo)
				return &NewExpr{N_EXPR_NEW, p.finRng(rng), expr, nil, span.Range{}, ti}, nil
			}
			return nil, err
		}
		p.discardState()

		if ti != nil {
			if err = p.tsCheckTypArgs(typArgs); err != nil {
				return nil, err
			}
			ti.SetTypArgs(typArgs)
		}

		// below is newExpr with callee tplExpr
		// ```
		// new C``
		// ```
		if len(args) == 0 && p.lexer.Peek().value == T_TPL_HEAD {
			if wt, ok := expr.(NodeWithTypInfo); ok {
				wt.SetTypInfo(ti)
			}
			expr, err = p.tplExpr(expr, false)
			if err != nil {
				return nil, err
			}
		}
	}

	var ret Node
	ret = &NewExpr{N_EXPR_NEW, p.finRng(rng), expr, args, span.Range{}, ti}
	root := true
	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if p.aheadIsArgList(tok) {
			if ret, _, err = p.callExpr(ret, root, false, span.Range{}, false); err != nil {
				return nil, err
			}
		} else if tv == T_BRACKET_L || tv == T_DOT || tv == T_OPT_CHAIN {
			if tv == T_OPT_CHAIN {
				optLoc = tok.rng
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

func (p *Parser) checkCallee(callee Node, nextLoc span.Range) error {
	scope := p.scope()
	switch callee.Type() {
	case N_EXPR_ARROW:
		if !scope.IsKind(SPK_PAREN) {
			return p.errorAtLoc(nextLoc, ERR_UNEXPECTED_TOKEN)
		}
	case N_SUPER:
		if !scope.IsKind(SPK_CTOR) || !scope.IsKind(SPK_CLASS_HAS_SUPER) {
			return p.errorAtLoc(callee.Range(), ERR_SUPER_CALL_OUTSIDE_CTOR)
		}
	case N_EXPR_MEMBER:
		n := callee.(*MemberExpr)
		if n.obj.Type() != N_SUPER {
			return p.checkCallee(n.obj, nextLoc)
		}
	}
	return nil
}

func (p *Parser) isLtTok(ofst uint32) bool {
	_, ok := p.ltTokens[ofst]
	return ok
}

// records the `T_LT` at given position should be considered
// as `less then` operator
func (p *Parser) addLtTok(ofst uint32) {
	p.ltTokens[ofst] = true
}

func (p *Parser) pushState() {
	p.lexer.PushState()
	p.lexer.src.PushState()
}

func (p *Parser) discardState() {
	p.lexer.DiscardState()
	p.lexer.src.DiscardState()
}

func (p *Parser) popState() {
	p.lexer.src.PopState()
	p.lexer.PopState()
}

func (p *Parser) aheadIsArgList(tok *Token) bool {
	tv := tok.value
	return tv == T_PAREN_L || (tv == T_LT && p.ts && !p.isLtTok(tok.rng.Lo))
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-CallExpression
func (p *Parser) callExpr(callee Node, root bool, directOpt bool, opt span.Range, notColon bool) (Node, span.Range, error) {
	var rng span.Range
	var err error
	if callee == nil {
		rng = p.rng()
		callee, err = p.primaryExpr(notColon)
		if err != nil {
			return nil, span.Range{}, err
		}
		callee = p.tsNoNull(callee)
	} else {
		rng = callee.Range()
	}

	firstOpt := opt
	var fo span.Range
	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if p.aheadIsArgList(tok) {
			aheadLoc := tok.rng

			// below pair `pushState` and `popState` is used to dealing with
			// the ambiguity between `<` in typArgs and `<` operator in binExpr
			lt := tok.value == T_LT
			if lt {
				p.pushState()
			}
			// `superTypArgs` is used to represent expr like `(class extends f()<T> {})`
			args, _, typArgs, superTypArgs, err := p.argList(true, true, span.Range{}, true)
			if err != nil {
				if err == p.errTypArgMissingGT && firstOpt.Empty() {
					p.popState()

					e := err.(ErrTypArgMissingGT)
					if p.lexer.maybeLsh(e.rng.Lo) {
						p.lexer.tryLsh(e.rng.Lo)
					} else {
						p.addLtTok(e.rng.Lo)
					}
					return callee, span.Range{}, nil
				}
				return nil, span.Range{}, err
			}
			if lt {
				p.discardState()
			}

			if superTypArgs != nil {
				if wt, ok := callee.(NodeWithTypInfo); ok {
					wt.TypInfo().SetSuperTypArgs(superTypArgs)
				}
				return callee, span.Range{}, nil
			}

			ti := p.newTypInfo(N_EXPR_CALL)
			if ti != nil {
				if err = p.tsCheckTypArgs(typArgs); err != nil {
					return nil, span.Range{}, err
				}
				ti.SetTypArgs(typArgs)
			}

			if err = p.checkCallee(callee, aheadLoc); err != nil {
				return nil, span.Range{}, err
			}

			// ```
			// f<T>``;
			// ```
			if args == nil && typArgs != nil {
				if wt, ok := callee.(NodeWithTypInfo); ok {
					wt.SetTypInfo(ti)
				} else {
					return nil, span.Range{}, p.errorAtLoc(typArgs.Range(), ERR_UNEXPECTED_TOKEN)
				}
			} else {
				callee = &CallExpr{N_EXPR_CALL, p.finRng(rng), callee, args, directOpt, span.Range{}, ti}
			}
		} else if tv == T_BRACKET_L || tv == T_DOT || tv == T_OPT_CHAIN {
			callee, fo, err = p.memberExpr(callee, true, root, firstOpt)
			if err != nil {
				return nil, span.Range{}, err
			}
			if firstOpt.Empty() {
				firstOpt = fo
			}
		} else if tv == T_TPL_HEAD {
			callee, err = p.tplExpr(callee, false)
			if err != nil {
				return nil, span.Range{}, err
			}
		} else {
			break
		}
	}

	ct := callee.Type()
	if root && !firstOpt.Empty() && (ct != N_NAME && ct != N_EXPR_CHAIN) {
		return &ChainExpr{N_EXPR_CHAIN, callee.Range(), callee}, firstOpt, nil
	}

	return callee, firstOpt, nil
}

// https://262.ecma-international.org/12.0/#prod-ImportCall
func (p *Parser) importCall(tok *Token) (Node, error) {
	if tok == nil {
		tok = p.lexer.Next()
	}
	rng := tok.rng

	meta := &Ident{N_NAME, p.finRng(tok.rng), p.TokText(tok), false, tok.ContainsEscape(), span.Range{}, false, p.newTypInfo(N_NAME)}

	ahead := p.lexer.Peek()
	if ahead.value == T_DOT && p.feat&FEAT_META_PROPERTY != 0 {
		p.lexer.Next()
		prop, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		if prop.val != "meta" {
			return nil, p.errorAtLoc(prop.rng, ERR_ILLEGAL_IMPORT_PROP)
		} else if prop.ContainsEscape() {
			return nil, p.errorAtLoc(prop.rng, ERR_META_PROP_CONTAINS_ESCAPE)
		}

		mp := &MetaProp{N_META_PROP, p.finRng(rng), meta, prop}
		ahead = p.lexer.Peek()
		av := ahead.value
		if av == T_PAREN_L {
			node, _, err := p.callExpr(mp, true, false, span.Range{}, false)
			return node, err
		} else if av == T_BRACKET_L || av == T_DOT || av == T_OPT_CHAIN {
			node, _, err := p.memberExpr(mp, true, true, span.Range{})
			return node, err
		}
		return &MetaProp{N_META_PROP, p.finRng(rng), meta, prop}, nil
	}

	if ahead.value == T_PAREN_L && p.feat&FEAT_DYNAMIC_IMPORT != 0 {
		p.lexer.Next()
		src, err := p.assignExpr(true, false, false, false)
		if err != nil {
			return nil, err
		}
		_, err = p.nextMustTok(T_PAREN_R)
		if err != nil {
			return nil, err
		}

		call := &ImportCall{N_IMPORT_CALL, p.finRng(rng), src, span.Range{}}
		ahead := p.lexer.Peek()
		semi := ahead.value == T_SEMI
		if semi || ahead.afterLineTerm {
			if semi {
				p.lexer.Next()
			}
			return call, nil
		}

		return p.condExpr(call, false, false, false)
	}
	return nil, p.errorTok(ahead)
}

func (p *Parser) tplExpr(tag Node, ts bool) (Node, error) {

	if tag != nil {
		if tag.Type() == N_EXPR_CHAIN {
			return nil, p.errorAtLoc(p.lexer.Peek().rng, ERR_OPT_EXPR_IN_TAG)
		}
	}

	var rng span.Range
	elems := make([]Node, 0, 5)
	for {
		tok := p.lexer.Peek()
		if rng.Empty() {
			rng = tok.rng
		}
		if tok.value >= T_TPL_HEAD && tok.value <= T_TPL_TAIL {
			cooked := ""
			ext := tok.ext.(*TokExtTplSpan)
			if ext.IllegalEscape != nil {
				// raise error for bad escape sequence if the template is not tagged
				if tag == nil || p.feat&FEAT_BAD_ESCAPE_IN_TAGGED_TPL == 0 {
					return nil, p.errorAt(tok.value, ext.IllegalEscape.Rng, ext.IllegalEscape.Err)
				}
			} else {
				cooked = ext.str
			}

			p.lexer.Next()
			str := &StrLit{N_LIT_STR, span.Range{Lo: ext.strRng.Lo, Hi: ext.strRng.Hi}, cooked, false, span.Range{}, nil}
			elems = append(elems, str)

			if tok.value == T_TPL_TAIL || tok.IsPlainTpl() {
				break
			}

			var expr Node
			var err error
			if !ts {
				expr, err = p.expr()
			} else {
				expr, err = p.tsTyp(false, false, true)
			}
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
		} else {
			return nil, p.errorTok(tok)
		}
	}

	return &TplExpr{N_EXPR_TPL, p.finRng(rng), tag, elems, span.Range{}}, nil
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

// `yield` indicates whether the yield-expr is permitted
func (p *Parser) checkDefaultVal(val Node, yield bool, destruct bool, field bool) error {
	switch val.Type() {
	case N_EXPR_YIELD:
		scope := p.scope()
		if !yield || !scope.IsKind(SPK_GENERATOR) {
			return p.errorAtLoc(val.Range(), ERR_YIELD_CANNOT_BE_DEFAULT_VALUE)
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
			return p.errorAtLoc(n.rng, ERR_AWAIT_AS_DEFAULT_VALUE)
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
		name := val.(*Ident).val
		if p.checkName && p.isProhibitedName(nil, name, true, false, field, false) {
			return p.errorAtLoc(id.rng, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
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
			rng:   n.rng,
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
			rng:   n.rng,
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
				return nil, p.errorAtLoc(arg.Range(), ERR_OBJ_PATTERN_CANNOT_FN)
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
			return nil, p.errorAtLoc(n.key.Range(), ERR_UNEXPECTED_TOKEN)
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
				rng: n.rng,
				lhs: n.key,
				rhs: n.value,
				ti:  p.newTypInfo(N_PAT_ASSIGN),
			}
		}
		return n, nil
	case N_EXPR_ASSIGN:
		n := arg.(*AssignExpr)
		if n.op != T_ASSIGN {
			return nil, p.errorAtLoc(n.opLoc, ERR_UNEXPECTED_TOKEN)
		}

		if pn, ok := n.lhs.(InParenNode); ok {
			inParen = !pn.OuterParen().Empty()
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
			rng: n.rng,
			lhs: lhs,
			rhs: n.rhs,
		}
		return p, nil
	case N_NAME:
		id := arg.(*Ident)
		name := arg.(*Ident).val
		if p.checkName && p.isProhibitedName(nil, name, true, true, false, false) {
			et := ERR_TPL_BINDING_RESERVED_WORD
			if destruct {
				et = ERR_TPL_ASSIGN_TO_RESERVED_WORD_IN_STRICT_MODE
			}
			return nil, p.errorAtLoc(id.rng, fmt.Sprintf(et, name))
		}
		return arg, nil
	case N_PAT_REST:
		return arg, nil
	case N_SPREAD:
		n := arg.(*Spread)
		if !n.tcLoc.Empty() {
			return nil, p.errorAtLoc(n.tcLoc, ERR_REST_ELEM_MUST_LAST)
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
					if !n.OuterParen().Empty() {
						return nil, p.errorAtLoc(n.OuterParen(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
					}
				}
			}
			if _, err := p.argToParam(n.arg, depth, prop, destruct, inParen); err != nil {
				return nil, err
			}
		} else if at == N_EXPR_ASSIGN {
			return nil, p.errorAtLoc(n.arg.Range(), ERR_REST_CANNOT_SET_DEFAULT)
		} else if at == N_EXPR_PAREN {
			if destruct {
				arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen)
				if err != nil {
					return nil, err
				}
				n.arg = arg
			} else {
				return nil, p.errorAtLoc(n.arg.Range(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
			}
		} else if p.feat&FEAT_BINDING_REST_ELEM_NESTED != 0 && (at == N_LIT_ARR || at == N_LIT_OBJ) {
			arg, err := p.argToParam(n.arg, depth, prop, destruct, inParen)
			if err != nil {
				return nil, err
			}
			if prop && !p.isSimpleLVal(arg, false, false, true, true) {
				return nil, p.errorAtLoc(arg.Range(), ERR_REST_ARG_NOT_SIMPLE)
			}
			n.arg = arg
		} else {
			if !prop && p.feat&FEAT_BINDING_REST_ELEM_NESTED == 0 {
				nested := UnParen(n.arg)
				if nested.Type() != N_NAME {
					return nil, p.errorAtLoc(nested.Range(), ERR_REST_ARG_NOT_BINDING_PATTERN)
				}
			}

			return nil, p.errorAtLoc(n.arg.Range(), ERR_REST_ARG_NOT_SIMPLE)
		}

		rest := &RestPat{
			typ: N_PAT_REST,
			rng: n.rng,
			arg: n.arg,
			ti:  n.TypInfo(),
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
					return nil, p.errorAtLoc(sub.Range(), ERR_ASSIGN_TO_RVALUE)
				}
				return nil, p.errorAtLoc(arg.Range(), ERR_INVALID_PAREN_ASSIGN_PATTERN)
			}
		}
		arg, err := p.argToParam(sub, depth, prop, destruct, true)
		if err != nil {
			return nil, err
		}
		if pn, ok := arg.(InParenNode); ok {
			pn.SetOuterParen(sub.Range())
		}
		return arg, nil
	case N_TS_TYP_ASSERT:
		n := arg.(*TsTypAssert)
		if destruct && !inParen {
			if depth < 2 {
				// `[a as number] = [42];` is legal
				// `<string>foo = '100';` is illegal
				return nil, p.errorAtLoc(n.rng, ERR_ASSIGN_TO_RVALUE)
			}
			// transform the arg at first: `<number>(a)`
			arg, err := p.argToParam(n.arg, depth, prop, destruct, true)
			if err != nil {
				return nil, err
			}

			// the transformed arg should be `NodeWithTypInfo` since we need to attach the
			// `des` of TsTypAssert as typAnnot of it
			if wt, ok := arg.(NodeWithTypInfo); ok {
				wt.TypInfo().SetTypAnnot(n.des)
			} else {
				return nil, p.errorAtLoc(n.Range(), ERR_UNEXPECTED_TOKEN)
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
				return nil, p.errorAtLoc(n.rng, ERR_ASSIGN_TO_RVALUE)
			}
			arg, err := p.argToParam(n.lhs, depth, prop, destruct, true)
			if err != nil {
				return nil, err
			}
			if wt, ok := arg.(NodeWithTypInfo); ok {
				wt.TypInfo().SetTypAnnot(n.rhs)
			} else {
				return nil, p.errorAtLoc(n.Range(), ERR_UNEXPECTED_TOKEN)
			}
			return arg, nil
		}
	}
	if depth == 0 {
		return nil, p.errorAtLoc(arg.Range(), ERR_UNEXPECTED_TOKEN)
	}
	// `([a.a]) => 42` is illegal since the `a.a` is not permitted to occur
	// `[a.r] = b` is legal since `a.r` is permitted to occur in destruct
	if !p.isSimpleLVal(arg, true, inParen, destruct, false) {
		return nil, p.errorAtLoc(arg.Range(), ERR_ASSIGN_TO_RVALUE)
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

	priv := key.Type() == N_NAME && key.(*Ident).pvt
	propName := p.nameOfNode(key)
	return propName, key, priv
}

// check the `arg` is legal as argument
// `spread` means whether the spread is permitted
// `simplicity` means whether to check simplicity of lhs of the assignExpr
func (p *Parser) checkArg(arg Node, spread bool, simplicity bool) error {
	if wt, ok := arg.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		// report the error of hook in expr like: `async(x?)`
		if ti != nil && !ti.Ques().Empty() {
			return p.errorAtLoc(ti.Ques(), ERR_UNEXPECTED_TOKEN)
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
						return p.errorAtLoc(pp.Range(), ERR_REDEF_PROP)
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
			return p.errorAtLoc(arg.Range(), ERR_UNEXPECTED_TOKEN)
		}
	case N_EXPR_ASSIGN:
		if simplicity {
			n := arg.(*AssignExpr)
			if n.op != T_ASSIGN && !p.isSimpleLVal(n.lhs, false, false, true, false) {
				return p.errorAtLoc(n.lhs.Range(), ERR_ASSIGN_TO_RVALUE)
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
			return p.errorAtLoc(arg.Range(), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
		}
		// for reporting `(a:b)` is illegal in ts
		if id.ti != nil && id.ti.TypAnnot() != nil {
			return p.errorAtLoc(id.ti.TypAnnot().Range(), ERR_UNEXPECTED_TYPE_ANNOTATION)
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

func (p *Parser) argList(check bool, incall bool, asyncLoc span.Range, noJsx bool) ([]Node, span.Range, Node, Node, error) {
	typArgs, err := p.tsTryTypArgs(asyncLoc, noJsx)
	if err != nil {
		return nil, span.Range{}, nil, nil, err
	}
	if typArgs != nil {
		tt := typArgs.Type()
		if tt != N_TS_PARAM_INST && tt != N_TS_PARAM_DEC {
			return nil, span.Range{}, typArgs, nil, nil
		}
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	if av != T_PAREN_L {
		// ```
		// (class extends f()<T> {}     // isExtending
		// f<T>``;                      // av == T_TPL_HEAD
		// ```
		isExtending := p.scope().IsKind(SPK_CLASS_EXTEND_SUPER)
		if p.ts && isExtending {
			// returns `typArgs` as `superTypArgs`
			return nil, span.Range{}, nil, typArgs, nil
		} else if av == T_TPL_HEAD {
			return nil, span.Range{}, typArgs, nil, nil
		}
		return nil, span.Range{}, nil, nil, p.errorTok(ahead)
	}
	p.lexer.Next()

	var tailingComma span.Range
	args := make([]Node, 0, 5)
	for {
		tok := p.lexer.Peek()
		if tok.value == T_PAREN_R {
			break
		} else if tok.value == T_EOF {
			return nil, span.Range{}, nil, nil, p.errorTok(tok)
		}
		arg, err := p.arg()
		if err != nil {
			return nil, span.Range{}, nil, nil, err
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
				return nil, tailingComma, nil, nil, p.errorAt(tok.value, tok.rng, msg)
			}
			if tailingComma.Empty() && ahead.value == T_PAREN_R {
				tailingComma = tok.rng
			}
		} else if av != T_PAREN_R {
			return nil, span.Range{}, nil, nil, p.errorTok(ahead)
		}

		if check {
			// `spread` or `pattern_rest` expression is legal argument:
			// `f(c, b, ...a)`
			if err := p.checkArg(arg, true, false); err != nil {
				return nil, span.Range{}, nil, nil, err
			}
		}

		args = append(args, arg)
	}

	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, span.Range{}, nil, nil, err
	}
	return args, tailingComma, typArgs, nil, nil
}

// consider below exprs:
// `(a,b)`
// `(a,b) =>`
// we cannot judge `(a,b)` is a parenExpr or the formalParamsList of
// an arrayExpr before we see the `=>` token, for avoiding to rollback
// the parsing state, we firstly parse `(a,b)` as parenExpr which children
// is parsed by this method and then convert the parsed subtree to formalParamList
// by using `argToParam` when required
func (p *Parser) arg() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(false, false, false, false)
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
func (p *Parser) binExpr(lhs Node, minPcd int, logic bool, nullish bool, notGT bool, notColon bool) (Node, error) {
	var err error
	if lhs == nil {
		if lhs, err = p.unaryExpr(nil, span.Range{}, notColon); err != nil {
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
			return nil, p.errorAtLoc(ahead.rng, ERR_NULLISH_MIXED_WITH_LOGIC)
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
		opLoc := p.lexer.Next().rng

		var rhs Node
		if op != T_TS_AS {
			rhs, err = p.unaryExpr(nil, span.Range{}, false)
		} else {
			rhs, err = p.tsTyp(false, true, false)
		}
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		aheadOp := ahead.IsBin(notIn, ts)
		kind = TokenKinds[aheadOp]
		for aheadOp != T_ILLEGAL && (kind.Pcd > pcd || kind.Pcd == pcd && kind.RightAssoc) {
			pcd = kind.Pcd
			rhs, err = p.binExpr(rhs, pcd, logic, nullish, notGT, false)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			aheadOp = ahead.IsBin(notIn, ts)
			kind = TokenKinds[aheadOp]
		}

		// deal with expr like: `console.log( -2 ** 4 )`
		if lhs.Type() == N_EXPR_UNARY && op == T_POW {
			n := lhs.(*UnaryExpr)
			return nil, p.errorAtLoc(UnParen(lhs.(*UnaryExpr).arg).Range(), fmt.Sprintf(ERR_TPL_UNARY_IMMEDIATELY_BEFORE_POW, n.OpText()))
		}

		// deal with expr like: `4 + async() => 2`
		if rhs.Type() == N_EXPR_ARROW {
			return nil, p.errorAtLoc(rhs.(*ArrowFn).arrowLoc, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, "=>"))
		}

		bin := &BinExpr{N_EXPR_BIN, span.Range{}, T_ILLEGAL, span.Range{}, nil, nil, span.Range{}}
		bin.rng = p.finRng(lhs.Range())
		bin.op = op
		bin.opLoc = opLoc
		bin.lhs = lhs
		bin.rhs = rhs
		lhs = bin
	}
	return lhs, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-MemberExpression
func (p *Parser) memberExpr(obj Node, call bool, root bool, optLoc span.Range) (Node, span.Range, error) {
	var err error
	if obj == nil {
		if p.lexer.Peek().value == T_NEW {
			if obj, err = p.newExpr(); err != nil {
				return nil, span.Range{}, err
			}
		} else if obj, err = p.memberExprObj(); err != nil {
			return nil, span.Range{}, err
		}
	}

	for {
		tok := p.lexer.Peek()
		tv := tok.value
		if tv == T_OPT_CHAIN {
			if optLoc.Empty() {
				optLoc = tok.rng
			}
			p.lexer.Next()

			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_BRACKET_L { // a?.[b]
				if obj, err = p.memberExprPropSubscript(obj, true); err != nil {
					return nil, span.Range{}, err
				}
			} else if p.aheadIsArgList(ahead) { // a?.()
				if obj, _, err = p.callExpr(obj, false, true, optLoc, false); err != nil {
					return nil, span.Range{}, err
				}
			} else {
				// a?.b
				if obj, err = p.memberExprPropDot(obj, true); err != nil {
					return nil, span.Range{}, err
				}
			}
		} else if tv == T_BRACKET_L {
			if obj, err = p.memberExprPropSubscript(obj, false); err != nil {
				return nil, span.Range{}, err
			}
		} else if tv == T_DOT {
			p.lexer.Next()
			if obj, err = p.memberExprPropDot(obj, false); err != nil {
				return nil, span.Range{}, err
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
				return nil, span.Range{}, p.errorAtLoc(prop.rng, ERR_UNEXPECTED_PVT_FIELD)
			}
		}
	}

	if call && p.aheadIsArgList(p.lexer.Peek()) {
		return p.callExpr(obj, root, false, optLoc, false)
	}

	if root && !optLoc.Empty() && obj.Type() != N_NAME {
		return &ChainExpr{N_EXPR_CHAIN, obj.Range(), obj}, optLoc, nil
	}
	return obj, optLoc, nil
}

func (p *Parser) memberExprObj() (Node, error) {
	obj, err := p.primaryExpr(false)
	if err != nil {
		return nil, err
	}
	obj = p.tsNoNull(obj)
	if p.lexer.Peek().value == T_TPL_HEAD {
		return p.tplExpr(obj, false)
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
	node := &MemberExpr{N_EXPR_MEMBER, p.finRng(obj.Range()), obj, prop, true, opt, span.Range{}}
	return node, nil
}

func (p *Parser) memberExprPropDot(obj Node, opt bool) (Node, error) {
	loc := p.rng()
	tok := p.lexer.Next()
	tv := tok.value
	kw, ok := tok.CanBePropKey()

	var prop Node
	if (ok && tv != T_NUM) || tv == T_NAME_PVT {
		pvt := tv == T_NAME_PVT
		id := &Ident{N_NAME, p.finRng(loc), p.TokText(tok), pvt, tok.ContainsEscape(), span.Range{}, kw, p.newTypInfo(N_NAME)}
		if pvt {
			scope := p.scope().UpperCls()
			if scope == nil {
				return nil, p.errorAtLoc(loc, fmt.Sprintf(ERR_TPL_ALONE_PVT_FIELD, "#"+p.TokText(tok)))
			}
			ref := NewRef()
			ref.Id = id
			ref.Typ = RDT_PVT_FIELD
			ref.Scope = scope
			p.danglingPvtRefs = append(p.danglingPvtRefs, ref)
		}
		prop = id
	} else {
		return nil, p.errorTok(tok)
	}

	node := &MemberExpr{N_EXPR_MEMBER, p.finRng(obj.Range()), obj, prop, false, opt, span.Range{}}
	return node, nil
}

func (p *Parser) primaryExpr(notColon bool) (Node, error) {
	tok := p.lexer.Peek()

	switch tok.value {
	case T_NUM:
		loc := tok.rng
		p.lexer.Next()
		return &NumLit{N_LIT_NUM, p.finRng(loc), span.Range{}}, nil
	case T_STRING:
		loc := tok.rng
		p.lexer.Next()
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, p.errorAtLoc(p.finRng(loc), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		return &StrLit{N_LIT_STR, p.finRng(loc), p.TokText(tok), legacyOctalEscapeSeq, span.Range{}, nil}, nil
	case T_NULL:
		loc := tok.rng
		p.lexer.Next()
		return &NullLit{N_LIT_NULL, p.finRng(loc), span.Range{}, nil}, nil
	case T_TRUE, T_FALSE:
		loc := tok.rng
		p.lexer.Next()
		return &BoolLit{N_LIT_BOOL, p.finRng(loc), p.TokText(tok) == "true", span.Range{}, nil}, nil
	case T_NAME:
		loc := tok.rng
		if p.aheadIsAsync(tok, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
			}
			return p.fnDec(true, tok, false)
		}
		if ok, itf, _ := p.tsAheadIsAbstract(tok, false, false, false); ok {
			if itf {
				return nil, p.errorAtLoc(tok.rng, ERR_ABSTRACT_AT_INVALID_POSITION)
			}
			return p.classDec(true, false, false, true)
		}

		p.lexer.Next()
		name := p.TokText(tok)
		ahead := p.lexer.Peek()
		// `ahead.value != T_ARROW` is used to skip checking name when it appears in the param list of arrow expr
		// for `eval => 42` we should report binding-reserved-word error instead of unexpected-reserved-word error
		if p.checkName && ahead.value != T_ARROW && !ahead.afterLineTerm && p.isProhibitedName(nil, name, true, false, false, false) {
			if tok.ContainsEscape() {
				return nil, p.errorAtLoc(p.finRng(loc), ERR_ESCAPE_IN_KEYWORD)
			}
			return nil, p.errorAtLoc(p.finRng(loc), fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, name))
		}
		kw := p.isProhibitedName(nil, name, true, false, false, false)
		return &Ident{N_NAME, p.finRng(loc), name, false, tok.ContainsEscape(), span.Range{}, kw, p.newTypInfo(N_NAME)}, nil
	case T_THIS:
		loc := tok.rng
		p.lexer.Next()
		return &ThisExpr{N_EXPR_THIS, p.finRng(loc), span.Range{}, nil}, nil
	case T_PAREN_L:
		return p.parenExpr(nil, notColon)
	case T_BRACKET_L:
		return p.arrLit()
	case T_BRACE_L:
		cmts := p.lexer.takeExprCmts()
		node, err := p.objLit()
		if err != nil {
			return nil, err
		}
		if len(cmts) > 0 {
			p.prevCmts[node] = cmts
		}
		return node, nil
	case T_FUNC:
		return p.fnDec(true, nil, false)
	case T_REGEXP:
		loc := tok.rng
		p.lexer.Next()
		ext := tok.ext.(*TokExtRegexp)
		return &RegLit{N_LIT_REGEXP, p.finRng(loc), p.TokText(tok), p.RngText(ext.pattern), p.RngText(ext.flags), span.Range{}, nil}, nil
	case T_CLASS:
		return p.classDec(true, false, false, false)
	case T_SUPER:
		loc := tok.rng
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
		return &Super{N_SUPER, p.finRng(loc), span.Range{}, nil}, nil
	case T_IMPORT:
		return p.importCall(nil)
	case T_TPL_HEAD:
		return p.tplExpr(nil, false)
	case T_LT:
		if p.feat&FEAT_JSX != 0 && p.feat&FEAT_TS == 0 {
			return p.jsx(true, false)
		}
		if p.feat&FEAT_TS != 0 {
			typArgs, err := p.tsTryTypArgs(span.Range{}, false)
			if err != nil {
				return nil, err
			}
			if typArgs != nil && typArgs.Type() == N_JSX_ELEM {
				return typArgs, nil
			}
			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_PAREN_L {
				return p.parenExpr(typArgs, false)
			}
			return p.unaryExpr(typArgs, typArgs.Range(), false)
		}
		return nil, p.errorTok(tok)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) arrowFn(rng span.Range, args []Node, params []Node, ti *TypInfo) (Node, error) {
	var err error
	if params == nil {
		params, err = p.argsToParams(args)
		if err != nil {
			return nil, err
		}
	}

	arrowLoc := p.lexer.Next().rng
	ps := p.scope()
	scope := p.symtab.EnterScope(true, true, true)
	p.incRetsStk()

	paramNames, firstComplicated, err := p.collectNames(params)
	if err != nil {
		return nil, err
	}

	for _, paramName := range paramNames {
		ref := NewRef()
		ref.Id = paramName.(*Ident)
		ref.BindKind = BK_PARAM
		// duplicate-checking is enable in strict mode by below `checkParams`
		p.addLocalBinding(nil, ref, false, ref.Id.val)
	}

	if ti == nil {
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		ti = p.newTypInfo(N_EXPR_ARROW)
		if ti != nil {
			ti.SetTypAnnot(typAnnot)
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
		body, err = p.assignExpr(true, false, false, false)
		if err != nil {
			return nil, err
		}
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)
	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	s := p.symtab.LeaveScope()
	n := &ArrowFn{N_EXPR_ARROW, p.finRng(rng), arrowLoc, false, params, body, body.Type() != N_STMT_BLOCK, p.decRetsStk(), span.Range{}, ti}
	s.Node = n
	return n, nil
}

func (p *Parser) parenExpr(typArgs Node, notColon bool) (Node, error) {
	rng := p.rng()

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
	scope := p.symtab.EnterScope(false, false, false)
	scope.AddKind(SPK_PAREN)
	p.checkName = false
	args, tailingComma, ta, _, err := p.argList(false, false, span.Range{}, false)
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
	if allowTypAnnot && !notColon {
		typAnnot, err = p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
	}

	ti := p.newTypInfo(N_EXPR_PAREN)
	if ti != nil {
		ti.SetTypAnnot(typAnnot)
		ti.SetTypArgs(typArgs)
	}

	// next is arrow-expression
	ahead := p.lexer.Peek()
	if ahead.value == T_ARROW && !ahead.afterLineTerm {
		if paramsErr != nil {
			return nil, paramsErr
		}
		return p.arrowFn(rng, nil, params, ti)
	}

	// `():number` is illegal
	if typAnnot != nil {
		return nil, p.errorAtLoc(typAnnot.Range(), ERR_UNEXPECTED_TOKEN)
	}

	// for report expr like: `(a,)`
	if !tailingComma.Empty() {
		return nil, p.errorAtLoc(tailingComma, ERR_TRAILING_COMMA)
	}

	argsLen := len(args)
	if argsLen == 0 {
		return nil, p.errorAt(p.lexer.state.prtVal, p.lexer.state.prtRng, "")
	}

	if err := p.checkArgs(args, false, true); err != nil {
		return nil, err
	}

	if argsLen == 1 {
		pe := &ParenExpr{N_EXPR_PAREN, p.finRng(rng), args[0], span.Range{}}
		if ti == nil {
			return pe, nil
		}
		node, err := p.tsTypAssert(pe, typArgs)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	seqLoc := args[0].Range()
	end := args[argsLen-1].Range()
	seqLoc.Hi = end.Hi
	seq := &SeqExpr{N_EXPR_SEQ, seqLoc, args, span.Range{}}
	pe := &ParenExpr{N_EXPR_PAREN, p.finRng(rng), seq, span.Range{}}
	if ti == nil {
		return pe, nil
	}
	node, err := p.tsTypAssert(pe, typArgs)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func UnParen(expr Node) Node {
	if expr.Type() == N_EXPR_PAREN {
		loc := expr.Range()
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
	loc := p.rng()
	p.lexer.Next()

	elems := make([]Node, 0, 5)
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
	return &ArrLit{N_LIT_ARR, p.finRng(loc), elems, span.Range{}, p.newTypInfo(N_LIT_ARR)}, nil
}

func (p *Parser) arrElem() (Node, error) {
	if p.lexer.Peek().value == T_DOT_TRI {
		return p.spread()
	}
	return p.assignExpr(true, false, false, false)
}

func (p *Parser) spread() (Node, error) {
	loc := p.rng()
	tok := p.lexer.Next()

	if p.feat&FEAT_SPREAD == 0 {
		return nil, p.errorTok(tok)
	}

	node, err := p.assignExpr(true, false, false, false)
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
	// of paren-expr and therefor the inner `{...obj1,}` was parsed as obj-expr
	//
	// keep the loc of tailing comma for reporting the `tailing comma after rest-expr`
	// err in the arg-to-param transform
	var trailingCommaLoc span.Range
	tok = p.lexer.Peek()
	if tok.value == T_COMMA {
		trailingCommaLoc = tok.rng
	}
	return &Spread{N_SPREAD, p.finRng(loc), node, trailingCommaLoc, span.Range{}, nil}, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html#prod-ObjectLiteral
func (p *Parser) objLit() (Node, error) {
	loc := p.rng()
	p.lexer.Next()

	props := make([]Node, 0, 5)
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
	return &ObjLit{N_LIT_OBJ, p.finRng(loc), props, span.Range{}, p.newTypInfo(N_LIT_OBJ)}, nil
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
		return p.method(span.Range{}, nil, ACC_MOD_NONE, span.Range{}, false, PK_INIT, true, false, false, false, false, span.Range{}, false, false, false, nil)
	} else if p.aheadIsAsync(tok, true, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
		}
		return p.method(span.Range{}, nil, ACC_MOD_NONE, span.Range{}, false, PK_INIT, false, true, false, false, false, span.Range{}, false, false, false, nil)
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
	if loc.Empty() {
		loc = key.Range()
	}

	var value Node
	tok := p.lexer.Peek()
	opLoc := tok.rng
	assign := tok.value == T_ASSIGN
	if tok.value == T_COLON || assign {
		p.lexer.Next()
		value, err = p.assignExpr(true, false, false, false)
		if err != nil {
			return nil, err
		}
	} else if p.aheadIsArgList(tok) {
		return p.method(loc, key, ACC_MOD_NONE, compute, false, PK_INIT, false, false, false, false, false, span.Range{}, false, false, false, nil)
	} else if !compute.Empty() {
		return nil, p.errorAt(tok.value, tok.rng, ERR_COMPUTE_PROP_MISSING_INIT)
	}

	shorthand := assign
	if value == nil && key.Type() == N_NAME {
		id := key.(*Ident)
		name := id.val
		if id.kw && name != "eval" && name != "arguments" {
			return nil, p.errorAtLoc(id.rng, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, id.val))
		}
		shorthand = true
		value = key
	}
	return &Prop{N_PROP, p.finRng(loc), key, opLoc, value, !compute.Empty(), false, shorthand, assign, PK_INIT, ACC_MOD_NONE}, nil
}

func (p *Parser) method(rng span.Range, key Node, accMode ACC_MOD, compute span.Range, shorthand bool, kind PropKind,
	gen, async, allowNamePVT, inclass, static bool, beginLoc span.Range, declare, abstract, override bool, ti *TypInfo) (Node, error) {

	if !inclass && !compute.Empty() {
		rng = compute
	}
	if rng.Empty() {
		rng = p.rng()
	}

	scope := p.symtab.EnterScope(true, false, true)
	scope.AddKind(SPK_METHOD)
	if kind == PK_CTOR && !static {
		scope.AddKind(SPK_CTOR)
	}
	p.incRetsStk()

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
	if ti == nil {
		ti = p.newTypInfo(N_METHOD)
	}
	if key == nil {
		if inclass {
			key, compute, err = p.classElemName()
		} else {
			key, compute, err = p.propName(allowNamePVT, false, false)
		}
		if err != nil {
			return nil, err
		}
		if ti != nil {
			ques, not := p.tsAdvanceHook(false)
			ti.SetQues(ques)
			ti.SetNot(not)
		}
	} else if ti != nil {
		if wt, ok := key.(NodeWithTypInfo); ok {
			tti := wt.TypInfo()
			if tti != nil {
				ti.SetQues(tti.Ques())
			}
		}
	}

	if ti != nil {
		ti.SetAccMod(accMode)
		ti.SetAbstract(abstract)
		ti.SetDeclare(declare)
		ti.SetOverride(override)
	}

	ctor := false
	if p.isName(key, "constructor", false, true) && compute.Empty() {
		if kind == PK_GETTER || kind == PK_SETTER {
			return nil, p.errorAtLoc(key.Range(), ERR_CTOR_CANNOT_HAVE_MODIFIER)
		} else if async {
			return nil, p.errorAtLoc(key.Range(), ERR_CTOR_CANNOT_BE_ASYNC)
		} else if gen {
			return nil, p.errorAtLoc(key.Range(), ERR_CTOR_CANNOT_BE_GENERATOR)
		}
		ctor = true
	}

	fnLoc := p.rng()
	params, typParams, _, err := p.paramList(false, kind, true)
	if err != nil {
		return nil, err
	}

	paramNames, firstComplicated, err := p.collectNames(params)
	if err != nil {
		return nil, err
	}

	for _, paramName := range paramNames {
		ref := NewRef()
		ref.Id = paramName.(*Ident)
		ref.BindKind = BK_PARAM
		// duplicate-checking is enable in strict mode so here skip doing checking,
		// checking is delegated to below `checkParams`
		p.addLocalBinding(nil, ref, false, ref.Id.val)
	}

	// the return type of method
	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	if ti != nil {
		if ctor && ti.Override() {
			return nil, p.errorAtLoc(key.Range(), ERR_OVERRIDE_ON_CTOR)
		}
		ti.SetTypAnnot(typAnnot)
		ti.SetTypParams(typParams)
	}

	if kind == PK_GETTER && len(params) > 0 {
		return nil, p.errorAtLoc(params[0].Range(), ERR_GETTER_SHOULD_NO_PARAM)
	}
	if kind == PK_SETTER {
		if len(params) != 1 {
			return nil, p.errorAtLoc(fnLoc, ERR_SETTER_SHOULD_ONE_PARAM)
		}
		if params[0].Type() == N_PAT_REST {
			return nil, p.errorAtLoc(params[0].Range(), ERR_REST_IN_SETTER)
		}
	}

	if gen {
		p.lexer.AddMode(LM_GENERATOR)
	}

	var body Node
	ahead := p.lexer.Peek()
	opt := false
	if ahead.value == T_BRACE_L {
		if abstract {
			return nil, p.errorAtLoc(rng, ERR_ABSTRACT_METHOD_WITH_IMPL)
		}
		body, err = p.fnBody()
		if gen {
			p.lexer.EraseMode(LM_GENERATOR)
		}
		if err != nil {
			return nil, err
		}
	}

	isStrict := scope.IsKind(SPK_STRICT)
	directStrict := scope.IsKind(SPK_STRICT_DIR)

	// `isProhibitedName` is not needed here since `keyword` as method name is permitted
	if err := p.checkParams(paramNames, firstComplicated, isStrict, directStrict); err != nil {
		return nil, err
	}

	s := p.symtab.LeaveScope()

	if p.ts {
		if body == nil {
			p.advanceIfSemi(false)
		}

		opts := NewTsCheckParamOpts()
		opts.impl = body != nil
		if err = p.tsCheckParams(params, opts); err != nil {
			return nil, err
		}
	}

	rets := p.decRetsStk()
	value := &FnDec{N_EXPR_FN, p.finRng(fnLoc), nil, gen, async, params, body, rets, span.Range{}, ti}
	s.Node = value
	if body == nil {
		if p.ts {
			// the method body can be emitted in typescript to represent the
			// method override, eg:
			// ```ts
			// class A { async?(): void }
			// ```
			// above code is legal in ts, it means class A contains a method
			// named `async` as well as optional, in other words the implementation
			// of method `async` can be emitted
			//
			// however, if the method is not flagged as optional then it must be followed
			// by the method definition which is the last subsequent statement in principle
			opt = ti != nil && !ti.Ques().Empty()
			if !opt {
				// memberExpr can also be the content of the key of the computed prop, eg:
				// `class C { [Symbol.iterator]?(): void; }`
				if name := p.nameOfNode(key); name == "" && key.Type() != N_EXPR_MEMBER {
					return nil, p.errorAtLoc(fnLoc, ERR_OVERRIDE_METHOD_DYNAMIC_NAME)
				}

				if !abstract {
					value.id = key
					sig := value
					p.lastTsFnSig = sig
				}
			}
		} else {
			return nil, p.errorAt(ahead.value, ahead.rng, ERR_UNEXPECTED_TOKEN)
		}
	} else if err = p.tsIsFnImplValid(key); err != nil {
		return nil, err
	}

	if inclass {
		if static && p.isName(key, "prototype", false, true) {
			return nil, p.errorAtLoc(key.Range(), ERR_STATIC_PROP_PROTOTYPE)
		}

		var ds []Node
		if len(p.hangingDecorators) > 0 {
			ds = p.hangingDecorators
			rng = ds[0].Range()
			p.hangingDecorators = nil
			ti.decorators = ds
		}
		return &Method{N_METHOD, p.finRng(rng), key, static, !compute.Empty(), kind, value, ti}, nil
	}
	return &Prop{N_PROP, p.finRng(rng), key, span.Range{}, value, !compute.Empty(), true, shorthand, false, kind, accMode}, nil
}

func (p *Parser) aheadIsDecorator(tok *Token) bool {
	return tok.value == T_AT
}

func (p *Parser) decorator() (Node, error) {
	tok := p.lexer.Peek()
	if !p.aheadIsDecorator(tok) {
		return nil, nil
	}
	loc := p.lexer.Next().rng
	expr, _, err := p.callExpr(nil, true, false, span.Range{}, false)
	if err != nil {
		return nil, err
	}
	return &Decorator{N_DECORATOR, p.finRng(loc), expr}, nil
}

func (p *Parser) decorators() ([]Node, error) {
	ds := make([]Node, 0, 5)
	for {
		d, err := p.decorator()
		if err != nil {
			return nil, err
		}
		if d == nil {
			break
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (p *Parser) nameOfNode(node Node) string {
	if node == nil {
		return ""
	}
	switch node.Type() {
	case N_NAME:
		id := node.(*Ident)
		v := id.val
		if id.pvt {
			return "#" + v
		}
		return v
	case N_LIT_STR:
		return node.(*StrLit).val
	case N_LIT_NULL:
		return "null"
	case N_LIT_BOOL:
		n := node.(*BoolLit)
		if n.val {
			return "true"
		}
		return "false"
	case N_LIT_NUM:
		n := node.(*NumLit)
		return p.NodeText(n)
	}
	return ""
}

func (p *Parser) isExprOpening(raise bool) (*Token, error) {
	tok := p.lexer.PeekStmtBegin()
	tv := tok.value
	if raise && tv != T_SEMI && tv != T_BRACE_R && tv != T_COMMA && tv != T_PAREN_R && tv != T_COLON && !tok.afterLineTerm && tv != T_EOF {
		errMsg := ERR_UNEXPECTED_TOKEN
		if tok.value == T_ILLEGAL {
			if msg, ok := tok.ext.(string); ok {
				errMsg = msg
			}
		}
		return nil, p.errorAt(tok.value, tok.rng, errMsg)
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
		return tok, p.errorTok(tok)
	}
	return tok, nil
}

func (p *Parser) nextMustName(name string, canContainsEscape bool) (*Token, error) {
	tok := p.lexer.Next()
	if tok.value != T_NAME || p.TokText(tok) != name {
		return nil, p.errorTok(tok)
	}
	if !canContainsEscape && tok.ContainsEscape() {
		return nil, p.errorAt(tok.value, tok.rng, ERR_ESCAPE_IN_KEYWORD)
	}
	return tok, nil
}

func (p *Parser) aheadIsName(name string) bool {
	tok := p.lexer.Peek()
	return tok.value == T_NAME && tok.text == name
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

func (p *Parser) rng() span.Range {
	return p.lexer.Rng()
}

func (p *Parser) RngText(rng span.Range) string {
	return p.lexer.src.RngText(rng)
}

func (p *Parser) TokText(t *Token) string {
	return TokText(t, p.lexer.src)
}

func (p *Parser) NodeText(n Node) string {
	return NodeText(n, p.lexer.src)
}

func (p *Parser) finRng(rng span.Range) span.Range {
	return p.lexer.FinRng(rng)
}

func (p *Parser) errorTok(tok *Token) *ParserError {
	if tok.value != T_ILLEGAL {
		return newParserError(p, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok.value].Name),
			p.lexer.src.Path, tok.rng.Lo)
	}
	return newParserError(p, tok.ErrMsg(), p.lexer.src.Path, tok.rng.Lo)
}

func (p *Parser) errorAt(tok TokenValue, pos span.Range, errMsg string) *ParserError {
	if tok != T_ILLEGAL && errMsg == "" {
		return newParserError(p, fmt.Sprintf(ERR_TPL_UNEXPECTED_TOKEN_TYPE, TokenKinds[tok].Name),
			p.lexer.src.Path, pos.Lo)
	}
	return newParserError(p, errMsg, p.lexer.src.Path, pos.Lo)
}

func (p *Parser) errorAtLoc(rng span.Range, errMsg string) *ParserError {
	return newParserError(p, errMsg, p.lexer.src.Path, rng.Lo)
}
