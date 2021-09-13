package parser

type LoopKind int

const (
	LOOP_NONE LoopKind = iota
	LOOP_DIRECT
	LOOP_INDIRECT
)

type Binding struct {
	at   *Token
	refs []Pos
}

type Scope struct {
	// an auto-increment number which is generated according
	// the depth-first walk over the entire AST
	Id       uint
	LoopKind LoopKind

	Parent *Scope
	Subs   []*Scope

	Params   map[string]int
	Bindings map[string]int
}

func NewScope() *Scope {
	scope := &Scope{
		Id:       0,
		Subs:     make([]*Scope, 0),
		Params:   make(map[string]int),
		Bindings: make(map[string]int),
	}
	return scope
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
		scope = scope.Parent
	}
	return false
}

func (s *Scope) AddParam(name string) {
	s.Params[name] = len(s.Params)
}

func (s *Scope) HasParam(name string) bool {
	_, ok := s.Params[name]
	return ok
}

func (s *Scope) LocalIdx(name string) int {
	if s.HasLocal(name) {
		return s.Bindings[name]
	}
	return -1
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

func (s *SymTab) EnterScope() {
	scope := NewScope()
	scope.Id = s.Cur.Id + 1
	s.Scopes[scope.Id] = scope

	scope.Parent = s.Cur
	s.Cur.Subs = append(s.Cur.Subs, scope)

	s.Cur = scope
}

func (s *SymTab) LeaveScope() {
	s.Cur = s.Cur.Parent
}

func (s *SymTab) HasExternal(name string) bool {
	for _, ext := range s.Externals {
		if ext == name {
			return true
		}
	}
	return false
}
