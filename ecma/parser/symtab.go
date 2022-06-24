package parser

type ScopeKind uint64

const (
	SPK_NONE        ScopeKind = 0
	SPK_LOOP_DIRECT ScopeKind = 1 << iota
	SPK_LOOP_INDIRECT
	SPK_SWITCH
	SPK_SWITCH_INDIRECT
	SPK_STRICT
	SPK_STRICT_DIR
	SPK_TRY
	SPK_TRY_INDIRECT
	SPK_CATCH
	SPK_BLOCK
	SPK_GLOBAL
	SPK_INTERIM
	SPK_FUNC
	SPK_FUNC_INDIRECT
	SPK_ARROW
	SPK_ASYNC
	SPK_GENERATOR
	SPK_PAREN
	SPK_CLASS
	SPK_CLASS_INDIRECT
	SPK_CLASS_EXTEND_SUPER
	SPK_CLASS_HAS_SUPER
	SPK_CTOR
	SPK_LEXICAL_DEC
	SPK_SHORTHAND_PROP
	SPK_METHOD
	SPK_NOT_IN
	SPK_PROP_NAME
	SPK_FORMAL_PARAMS
	SPK_ABSTRACT_CLASS
	SPK_TS_DECLARE
	SPK_TS_MODULE
	SPK_TS_MODULE_INDIRECT
	SPK_TS_INTERFACE
	SPK_TS_MAY_INTRINSIC
)

type BindKind uint8

const (
	BK_NONE BindKind = iota
	BK_VAR
	BK_PARAM
	BK_LET
	BK_CONST
	BK_PVT_FIELD
)

type RefDefType uint32

const (
	RDT_NONE       RefDefType = 0
	RDT_FN         RefDefType = 1 << iota
	RDT_PVT_FIELD  RefDefType = 1 << iota
	RDT_CLASS      RefDefType = 1 << iota
	RDT_ENUM       RefDefType = 1 << iota
	RDT_CONST_ENUM RefDefType = 1 << iota
	RDT_ITF        RefDefType = 1 << iota
	RDT_NS         RefDefType = 1 << iota
	RDT_TYPE       RefDefType = 1 << iota
)

func (t RefDefType) On(flag RefDefType) RefDefType {
	return t | flag
}

func (t RefDefType) Off(flag RefDefType) RefDefType {
	return t & ^flag
}

func (t RefDefType) IsTyp() bool {
	return t&RDT_TYPE != 0
}

func (t RefDefType) IsPureTyp() bool {
	return t&RDT_TYPE != 0 &&
		t&RDT_CLASS == 0 &&
		t&RDT_ENUM == 0 &&
		t&RDT_CONST_ENUM == 0
}

func (t RefDefType) IsVal() bool {
	return !t.IsPureTyp()
}

func (t RefDefType) IsPureVal() bool {
	return !t.IsPureTyp() && !t.IsTyp()
}

type Ref struct {
	Scope *Scope
	Def   *Ident
	Typ   RefDefType

	// points to the ref referenced by this one, eg:
	//
	// ```
	// A -> B // A points to B
	// ```
	//
	// `B` is the value of `Forward` of `A`
	Forward *Ref

	// ref with bind kind not none means it's a variable binding
	BindKind BindKind
	Props    map[string][]*Ref
	Refs     []*Ref
}

func (r *Ref) RetainBy(ref *Ref) {
	ref.Forward = r
	r.Refs = append(r.Refs, ref)
}

func NewRef() *Ref {
	return &Ref{
		Props: make(map[string][]*Ref),
		Refs:  make([]*Ref, 0),
	}
}

type Scope struct {
	// an auto-increment number which is generated according
	// the depth-first walk over the entire AST
	Id   int
	Kind ScopeKind

	// the Node which introduces this Scope, only usable after the scope is post-process since
	// the Node has not been resolved before the scope leave, eg. `p.fnDec`` can also return `ArrowFn`
	Node

	Up   *Scope
	Down []*Scope

	// label should be unique in its label chain
	uniqueLabels map[string]Node
	// labels can be redefined in their defined scope
	// so slice type used here
	Labels []Node

	// `IsBind` of the elems of the `scope.Refs` are all `true`,
	// `IsBind` of their children are `true` means `rebind`
	Refs map[string]*Ref

	// exports declared at this scope
	Exports []*ExportDec
}

func NewScope() *Scope {
	scope := &Scope{
		Id:   0,
		Down: make([]*Scope, 0),

		uniqueLabels: make(map[string]Node),
		Labels:       make([]Node, 0),

		Refs: make(map[string]*Ref),
	}
	return scope
}

func (s *Scope) IsKind(kind ScopeKind) bool {
	return s.Kind&kind != 0
}

func (s *Scope) AddKind(kind ScopeKind) *Scope {
	s.Kind |= kind
	return s
}

func (s *Scope) EraseKind(kind ScopeKind) *Scope {
	s.Kind &= ^kind
	return s
}

func (s *Scope) Local(name string) *Ref {
	return s.Refs[name]
}

func (s *Scope) UpperFn() *Scope {
	scope := s
	for scope != nil {
		if scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_GLOBAL) {
			return scope
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) UpperCls() *Scope {
	scope := s
	for scope != nil {
		if scope.IsKind(SPK_CLASS) {
			return scope
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) OuterFn() *Scope {
	scope := s.OuterScope()
	for scope != nil {
		if scope.IsKind(SPK_FUNC) || scope.IsKind(SPK_GLOBAL) {
			return scope
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) OuterScope() *Scope {
	if s.IsKind(SPK_GLOBAL) {
		return s
	}
	return s.Up
}

func (s *Scope) UpperLoop() *Scope {
	scope := s
	for scope != nil {
		if scope.IsKind(SPK_LOOP_DIRECT) {
			return scope
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) UpperScope(kinds *[]ScopeKind) *Scope {
	scope := s
	for scope != nil {
		for _, k := range *kinds {
			if scope.IsKind(k) {
				return scope
			}
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) AddLocal(ref *Ref, name string, checkDup bool) bool {
	cur := s
	local := s.Local(name)

	// `try {} catch (foo) { let foo; }` is illegal
	if cur.IsKind(SPK_CATCH) && local != nil {
		return false
	}

	// register binding to parent fn scope if it's `BK_VAR`
	if ref.BindKind == BK_VAR {
		ps := s.UpperFn()
		localInPs := ps.Refs[name]

		if checkDup && localInPs != nil &&
			((localInPs.BindKind != BK_VAR && localInPs.Typ != RDT_FN) ||
				(localInPs.Scope.Id == 0 && localInPs.Typ == RDT_FN)) {
			return CheckRefDup(localInPs, ref)
		}
		ps.Refs[name] = ref
	}

	if !checkDup {
		ref.Scope = s
		s.Refs[name] = ref
		return true
	}

	bindKind := ref.BindKind
	if local != nil && ((local.BindKind != BK_VAR && local.Typ != RDT_FN && local.Scope.Id != 0) || bindKind != BK_VAR) {
		return CheckRefDup(local, ref)
	}

	ref.Scope = s
	s.Refs[name] = ref
	return true
}

// caller should ensure both `r1` and `r2` have the same name,
// return `true` - `r1` and `r2` are diff
// return `false` - `r1` and `r2` are dup
func CheckRefDup(r1, r2 *Ref) bool {
	if IsCallableClass(r1, r2) {
		return true
	}
	if IsBothFnTypDec(r1, r2) {
		return true
	}
	if IsClsAndIft(r1, r2) {
		return true
	}
	if IsBothEnum(r1, r2) {
		return true
	}
	if IsBothItf(r1, r2) {
		return true
	}
	return !IsBothTyp(r1, r2) && !IsBothVal(r1, r2)
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
//
// this method is used to allow stmts:
//
// ```
// declare class C { }
// function C() { }
// ```
func IsCallableClass(r1, r2 *Ref) bool {
	var fn *Ref
	var cls *Ref
	if r1.Typ&RDT_FN != 0 {
		fn = r1
		if r2.Typ&RDT_CLASS != 0 && r2.Typ&RDT_TYPE != 0 {
			cls = r2
		}
	} else if r2.Typ&RDT_FN != 0 {
		fn = r2
		if r1.Typ&RDT_CLASS != 0 && r1.Typ&RDT_TYPE != 0 {
			cls = r1
		}
	}
	return fn != nil && cls != nil
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
//
// this method is used to allow stmts:
//
// ```
// declare function f(): void;
// declare function f<T>(): T;
// ```
func IsBothFnTypDec(r1, r2 *Ref) bool {
	return r1.Typ&RDT_FN != 0 &&
		r1.Typ&RDT_TYPE != 0 &&
		r2.Typ&RDT_FN != 0 &&
		r2.Typ&RDT_TYPE != 0
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
func IsBothTyp(r1, r2 *Ref) bool {
	typ := 0
	if r1.Typ.IsTyp() {
		typ += 1
	}
	if r2.Typ.IsTyp() {
		typ += 1
	}
	return typ == 2
}

// caller should ensure both `r1` and `r2` match below rules:
// - be def rather than ref
// - have the same name
func IsBothVal(r1, r2 *Ref) bool {
	typ := 0
	if r1.Typ.IsVal() {
		typ += 1
	}
	if r2.Typ.IsVal() {
		typ += 1
	}
	return typ == 2
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
//
// this method is used to allow stmts:
//
// ```
// class A {}
// interface A {}
// ```
func IsClsAndIft(r1, r2 *Ref) bool {
	var cls *Ref
	var itf *Ref
	if r1.Typ&RDT_CLASS != 0 {
		cls = r1
		if r2.Typ&RDT_ITF != 0 {
			itf = r2
		}
	} else if r2.Typ&RDT_CLASS != 0 {
		cls = r2
		if r1.Typ&RDT_ITF != 0 {
			itf = r1
		}
	}
	return cls != nil && itf != nil
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
//
// this method is used to allow stmts:
//
// ```
// const enum Foo {}
// const enum Foo {}
// ```
func IsBothEnum(r1, r2 *Ref) bool {
	return (r1.Typ&RDT_ENUM != 0 && r2.Typ&RDT_ENUM != 0) ||
		(r1.Typ&RDT_CONST_ENUM != 0 && r2.Typ&RDT_CONST_ENUM != 0)
}

// caller should ensure both `r1` and `r2` match below rules:
// - are def rather than ref
// - have the same name
//
// this method is used to allow stmts:
//
// ```
// interface A {}
// interface A {}
// ```
func IsBothItf(r1, r2 *Ref) bool {
	return r1.Typ&RDT_ITF != 0 && r2.Typ&RDT_ITF != 0
}

func (s *Scope) DelLocal(ref *Ref) {
	s.Refs[ref.Def.Text()] = nil
}

func (s *Scope) BindingOf(name string) *Ref {
	scope := s
	for scope != nil {
		ref := scope.Local(name)
		if ref != nil {
			return ref
		}
		scope = scope.Up
	}
	return nil
}

func (s *Scope) HasName(name string) bool {
	ref := s.BindingOf(name)
	return ref != nil
}

func (s *Scope) HasLabel(name string) bool {
	_, ok := s.uniqueLabels[name]
	return ok
}

func (s *Scope) GetLabel(name string) Node {
	scope := s
	for scope != nil {
		lb := scope.uniqueLabels[name]
		if lb != nil {
			return lb
		}
		if scope.IsKind(SPK_FUNC) {
			break
		}
		scope = scope.Up
	}
	return nil
}

type SymTab struct {
	Externals []string
	Scopes    map[int]*Scope
	Root      *Scope
	Cur       *Scope

	scopeIdSeed int // the seed of scope id
}

func NewSymTab(externals []string) *SymTab {
	scope := NewScope()

	symtab := &SymTab{
		Externals: externals,
		Scopes:    make(map[int]*Scope),
		Root:      scope,
		Cur:       scope,
	}
	symtab.Scopes[scope.Id] = scope
	return symtab
}

// `settled` to increase the scope id, otherwise the new entered scope will be
// treated as a temporary one
func (s *SymTab) EnterScope(fn bool, arrow bool, settled bool) *Scope {
	scope := NewScope()

	if settled {
		s.scopeIdSeed += 1
	}
	scope.Id = s.scopeIdSeed

	if fn {
		scope.Kind = SPK_FUNC
	} else {
		// inherit labels from parent scope
		for k, v := range s.Cur.uniqueLabels {
			scope.uniqueLabels[k] = v
		}
		scope.Kind = SPK_BLOCK
	}
	// inherit scope kind
	if !fn {
		if s.Cur.IsKind(SPK_LOOP_DIRECT) || s.Cur.IsKind(SPK_LOOP_INDIRECT) {
			scope.Kind |= SPK_LOOP_INDIRECT
		}
		if s.Cur.IsKind(SPK_SWITCH) || s.Cur.IsKind(SPK_SWITCH_INDIRECT) {
			scope.Kind |= SPK_SWITCH_INDIRECT
		}
		if s.Cur.IsKind(SPK_TRY) || s.Cur.IsKind(SPK_TRY_INDIRECT) {
			scope.Kind |= SPK_SWITCH_INDIRECT
		}
		if s.Cur.IsKind(SPK_FORMAL_PARAMS) {
			scope.Kind |= SPK_FORMAL_PARAMS
		}
		if s.Cur.IsKind(SPK_GENERATOR) {
			scope.Kind |= SPK_GENERATOR
		}
		if s.Cur.IsKind(SPK_ASYNC) {
			scope.Kind |= SPK_ASYNC
		}
	}

	if s.Cur.IsKind(SPK_FUNC) || s.Cur.IsKind(SPK_FUNC_INDIRECT) {
		scope.Kind |= SPK_FUNC_INDIRECT
	}
	if s.Cur.IsKind(SPK_CLASS) || s.Cur.IsKind(SPK_CLASS_INDIRECT) {
		scope.Kind |= SPK_CLASS_INDIRECT
	}
	if s.Cur.IsKind(SPK_CLASS_HAS_SUPER) {
		scope.Kind |= SPK_CLASS_HAS_SUPER
	}
	if s.Cur.IsKind(SPK_STRICT) {
		scope.Kind |= SPK_STRICT
	}
	if s.Cur.IsKind(SPK_ABSTRACT_CLASS) {
		scope.Kind |= SPK_ABSTRACT_CLASS
	}
	if s.Cur.IsKind(SPK_TS_DECLARE) {
		scope.Kind |= SPK_TS_DECLARE
	}
	if s.Cur.IsKind(SPK_TS_MODULE) || s.Cur.IsKind(SPK_TS_MODULE_INDIRECT) {
		scope.Kind |= SPK_TS_MODULE_INDIRECT
	}

	// `(class A extends B { constructor() { (() => { super() }); } })` is legal
	// `(class A extends B { constructor() { function f() { super() } } })` is illegal
	// it requires the `SPK_CTOR` to be inherited if new scope is arrow fn
	if !fn || (fn && arrow) {
		if s.Cur.IsKind(SPK_CTOR) {
			scope.Kind |= SPK_CTOR
		}
	}

	s.Scopes[scope.Id] = scope

	scope.Up = s.Cur
	s.Cur.Down = append(s.Cur.Down, scope)

	s.Cur = scope
	return scope
}

func (s *SymTab) LeaveScope() *Scope {
	cur := s.Cur
	// prevent the scope being overlaid by its tmp child
	s.Scopes[s.Cur.Up.Id] = s.Cur.Up
	s.Cur = s.Cur.Up
	return cur
}

func (s *SymTab) HasExternal(name string) bool {
	for _, ext := range s.Externals {
		if ext == name {
			return true
		}
	}
	return false
}
