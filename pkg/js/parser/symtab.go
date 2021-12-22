package parser

type ScopeKind uint64

const (
	SPK_NONE            ScopeKind = 0
	SPK_LOOP_DIRECT     ScopeKind = 1 << iota
	SPK_LOOP_INDIRECT   ScopeKind = 1 << iota
	SPK_SWITCH          ScopeKind = 1 << iota
	SPK_STRICT          ScopeKind = 1 << iota
	SPK_STRICT_DIR      ScopeKind = 1 << iota
	SPK_CATCH           ScopeKind = 1 << iota
	SPK_BLOCK           ScopeKind = 1 << iota
	SPK_GLOBAL          ScopeKind = 1 << iota
	SPK_INTERIM         ScopeKind = 1 << iota
	SPK_FUNC            ScopeKind = 1 << iota
	SPK_FUNC_INDIRECT   ScopeKind = 1 << iota
	SPK_ARROW           ScopeKind = 1 << iota
	SPK_ASYNC           ScopeKind = 1 << iota
	SPK_GENERATOR       ScopeKind = 1 << iota
	SPK_PAREN           ScopeKind = 1 << iota
	SPK_CLASS           ScopeKind = 1 << iota
	SPK_CLASS_INDIRECT  ScopeKind = 1 << iota
	SPK_CLASS_HAS_SUPER ScopeKind = 1 << iota
	SPK_CTOR            ScopeKind = 1 << iota
	SPK_LEXICAL_DEC     ScopeKind = 1 << iota
	SPK_SHORTHAND_PROP  ScopeKind = 1 << iota
	SPK_METHOD          ScopeKind = 1 << iota
	SPK_NOT_IN          ScopeKind = 1 << iota
	SPK_PROP_NAME       ScopeKind = 1 << iota
	SPK_FORMAL_PARAMS   ScopeKind = 1 << iota
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

type TargetType int

const (
	TT_NONE      TargetType = 0
	TT_FN        TargetType = 1 << iota
	TT_PVT_FIELD TargetType = 1 << iota
)

type Ref struct {
	Node  *Ident
	Scope *Scope
	// points to the ref referenced by this one
	Target     *Ref
	TargetType TargetType
	// ref with bind kind not none means it's a variable binding
	BindKind BindKind
	Props    map[string][]*Ref
	Refs     []*Ref
}

func (r *Ref) RetainBy(ref *Ref) {
	ref.Target = r
	r.Refs = append(r.Refs, ref)
}

func NewRef() *Ref {
	return &Ref{
		Node:     nil,
		Scope:    nil,
		Target:   nil,
		BindKind: BK_NONE,
		Props:    make(map[string][]*Ref),
		Refs:     make([]*Ref, 0),
	}
}

type Scope struct {
	// an auto-increment number which is generated according
	// the depth-first walk over the entire AST
	Id   uint
	Kind ScopeKind

	Up   *Scope
	Down []*Scope

	// label should be unique in its label chain
	uniqueLabels map[string]int
	// labels can be redefined in their defined scope
	// so slice type used here
	Labels []Node

	// `IsBind` of the elems of the `scope.Refs` are all `true`,
	// `IsBind` of their children are `true` means `rebind`
	Refs map[string]*Ref
}

func NewScope() *Scope {
	scope := &Scope{
		Id:   0,
		Down: make([]*Scope, 0),

		uniqueLabels: make(map[string]int),
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
		if localInPs != nil && localInPs.BindKind != BK_VAR {
			return false
		}
		ps.Refs[name] = ref
	}

	if !checkDup {
		ref.Scope = s
		s.Refs[name] = ref
		return true
	}

	bindKind := ref.BindKind
	if local != nil && (local.BindKind != BK_VAR || bindKind != BK_VAR) {
		return false
	}

	ref.Scope = s
	s.Refs[name] = ref
	return true
}

func (s *Scope) DelLocal(ref *Ref) {
	s.Refs[ref.Node.Text()] = nil
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

type SymTab struct {
	Externals []string
	Scopes    map[uint]*Scope
	Root      *Scope
	Cur       *Scope
}

func NewSymTab(externals []string) *SymTab {
	scope := NewScope()

	symtab := &SymTab{
		Externals: externals,
		Scopes:    make(map[uint]*Scope),
		Root:      scope,
		Cur:       scope,
	}
	symtab.Scopes[scope.Id] = scope
	return symtab
}

func (s *SymTab) EnterScope(fn bool, arrow bool) *Scope {
	scope := NewScope()
	scope.Id = s.Cur.Id + 1

	if fn {
		scope.Kind = SPK_FUNC
	} else {
		// inherit labels from parent scope
		for k := range s.Cur.uniqueLabels {
			scope.uniqueLabels[k] = 1
		}
		scope.Kind = SPK_BLOCK
	}
	// inherit scope kind
	if s.Cur.IsKind(SPK_LOOP_DIRECT) || s.Cur.IsKind(SPK_LOOP_INDIRECT) && !fn {
		scope.Kind |= SPK_LOOP_INDIRECT
	}
	if s.Cur.IsKind(SPK_FUNC) || s.Cur.IsKind(SPK_FUNC_INDIRECT) {
		scope.Kind |= SPK_FUNC_INDIRECT
	}
	if s.Cur.IsKind(SPK_CLASS) || s.Cur.IsKind(SPK_CLASS_INDIRECT) {
		scope.Kind |= SPK_CLASS_INDIRECT
	}
	if s.Cur.IsKind(SPK_STRICT) {
		scope.Kind |= SPK_STRICT
	}
	if s.Cur.IsKind(SPK_CLASS_HAS_SUPER) {
		scope.Kind |= SPK_CLASS_HAS_SUPER
	}
	if s.Cur.IsKind(SPK_FORMAL_PARAMS) && !fn {
		scope.Kind |= SPK_FORMAL_PARAMS
	}

	// `(class A extends B { constructor() { (() => { super() }); } })` is legal
	// `(class A extends B { constructor() { function f() { super() } } })` is illegal
	// it requires the `SPK_CTOR` to be inherited if new scope is arrow fn
	if !fn || (fn && arrow) {
		if s.Cur.IsKind(SPK_CTOR) {
			scope.Kind |= SPK_CTOR
		}
	}

	if s.Cur.IsKind(SPK_GENERATOR) && !fn {
		scope.Kind |= SPK_GENERATOR
	}
	if s.Cur.IsKind(SPK_ASYNC) && !fn {
		scope.Kind |= SPK_ASYNC
	}

	s.Scopes[scope.Id] = scope

	scope.Up = s.Cur
	s.Cur.Down = append(s.Cur.Down, scope)

	s.Cur = scope
	return scope
}

func (s *SymTab) LeaveScope() {
	s.Cur = s.Cur.Up
}

func (s *SymTab) HasExternal(name string) bool {
	for _, ext := range s.Externals {
		if ext == name {
			return true
		}
	}
	return false
}
