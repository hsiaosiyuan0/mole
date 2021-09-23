package parser

type LoopKind int

const (
	LOOP_NONE LoopKind = iota
	LOOP_DIRECT
	LOOP_INDIRECT
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
	Id       uint
	LoopKind LoopKind
	Strict   bool

	Up   *Scope
	Down []*Scope

	Bindings map[string]int
}

func NewScope() *Scope {
	scope := &Scope{
		Id:       0,
		Down:     make([]*Scope, 0),
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

func (s *SymTab) EnterScope() *Scope {
	scope := NewScope()
	scope.Id = s.Cur.Id + 1
	scope.Strict = s.Cur.Strict
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
