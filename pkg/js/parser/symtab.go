package parser

type ScopeKind int

const (
	SPK_NONE            ScopeKind = 0
	SPK_LOOP_DIRECT               = 1 << 0
	SPK_LOOP_INDIRECT             = 1 << 1
	SPK_SWITCH                    = 1 << 2
	SPK_STRICT                    = 1 << 3
	SPK_STRICT_DIR                = 1 << 4
	SPK_CATCH                     = 1 << 5
	SPK_BLOCK                     = 1 << 6
	SPK_GLOBAL                    = 1 << 7
	SPK_FUNC                      = 1 << 8
	SPK_FUNC_INDIRECT             = 1 << 9
	SPK_ARROW                     = 1 << 10
	SPK_ASYNC                     = 1 << 11
	SPK_GENERATOR                 = 1 << 12
	SPK_PAREN                     = 1 << 13
	SPK_CLASS                     = 1 << 14
	SPK_CLASS_HAS_SUPER           = 1 << 15
	SPK_CTOR                      = 1 << 16
)

type BindKind uint8

const (
	BK_NONE  BindKind = 0
	BK_VAR   BindKind = 1
	BK_PARAM BindKind = 2
	BK_LET   BindKind = 3
	BK_CONST BindKind = 4
)

type Ref struct {
	Node     *Ident
	Scope    *Scope
	Orig     *Ref
	BindKind BindKind
	Props    map[string][]*Ref
	Refs     []*Ref
}

func NewRef() *Ref {
	return &Ref{
		Node:     nil,
		Scope:    nil,
		Orig:     nil,
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

func (s *Scope) AddLocal(ref *Ref, checkDup bool) bool {
	cur := s
	name := ref.Node.Text()
	local := s.Local(name)

	// `try {} catch (foo) { let foo; }` is illegal
	if cur.IsKind(SPK_CATCH) && local != nil {
		return false
	}

	if ref.BindKind == BK_VAR {
		s = s.UpperFn()
	}

	if !checkDup {
		ref.Scope = s
		s.Refs[ref.Node.Text()] = ref
		return true
	}

	bindKind := ref.BindKind
	if local != nil && bindKind != BK_VAR {
		return false
	}

	ref.Scope = s
	s.Refs[ref.Node.Text()] = ref
	return true
}

func (s *Scope) RetainRef(ref *Ref) {
	n := ref.Node.Text()
	scope := s
	for scope != nil {
		if orig, ok := scope.Refs[n]; ok {
			ref.Scope = scope
			ref.Orig = orig
			orig.Refs = append(ref.Refs, ref)
			return
		}
		scope = scope.Up
	}
}

func (s *Scope) HasName(name string) bool {
	scope := s
	for scope != nil {
		if scope.Local(name) != nil {
			return true
		}
		scope = scope.Up
	}
	return false
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
	if s.Cur.IsKind(SPK_LOOP_DIRECT) {
		scope.Kind |= SPK_LOOP_INDIRECT
	}
	if s.Cur.IsKind(SPK_FUNC) {
		scope.Kind |= SPK_FUNC_INDIRECT
	}
	if s.Cur.IsKind(SPK_STRICT) {
		scope.Kind |= SPK_STRICT
	}
	if s.Cur.IsKind(SPK_CLASS_HAS_SUPER) {
		scope.Kind |= SPK_CLASS_HAS_SUPER
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
