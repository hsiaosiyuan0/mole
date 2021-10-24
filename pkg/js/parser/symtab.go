package parser

type ScopeKind int

const (
	SPK_NONE          ScopeKind = 0
	SPK_LOOP_DIRECT             = 1 << 0
	SPK_LOOP_INDIRECT           = 1 << 1
	SPK_SWITCH                  = 1 << 2
	SPK_STRICT                  = 1 << 3
	SPK_BLOCK                   = 1 << 4
	SPK_GLOBAL                  = 1 << 5
	SPK_FUNC                    = 1 << 6
	SPK_FUNC_INDIRECT           = 1 << 7
	SPK_ASYNC                   = 1 << 8
	SPK_GENERATOR               = 1 << 9
)

type Binding struct {
	name     *Token
	local    bool
	legal    bool
	refByCnt int
}

type Scope struct {
	// an auto-increment number which is generated according
	// the depth-first walk over the entire AST
	Id   uint
	Kind ScopeKind

	Up   *Scope
	Down []*Scope

	Labels   map[string]Node
	Bindings map[string]int
}

func NewScope() *Scope {
	scope := &Scope{
		Id:       0,
		Down:     make([]*Scope, 0),
		Labels:   make(map[string]Node),
		Bindings: make(map[string]int),
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

func (s *Scope) HasLocal(name string) bool {
	_, ok := s.Bindings[name]
	return ok
}

func (s *Scope) AddBinding(name string) {
	if s.HasLocal(name) {
		return
	}
	s.Bindings[name] = len(s.Bindings)
}

func (s *Scope) HasBinding(name string) bool {
	scope := s
	for scope != nil {
		if scope.HasLocal(name) {
			return true
		}
		scope = scope.Up
	}
	return false
}

func (s *Scope) LocalIdx(name string) int {
	if s.HasLocal(name) {
		return s.Bindings[name]
	}
	return -1
}

func (s *Scope) HasLabel(name string) bool {
	_, ok := s.Labels[name]
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

func (s *SymTab) EnterScope(fn bool) *Scope {
	scope := NewScope()
	scope.Id = s.Cur.Id + 1

	if fn {
		scope.Kind = SPK_FUNC
	} else {
		// inherit labels from parent scope
		for k, v := range s.Cur.Labels {
			scope.Labels[k] = v
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
