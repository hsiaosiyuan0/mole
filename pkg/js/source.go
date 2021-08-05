package js

import (
	"io/ioutil"
	"path/filepath"
	"sync"
	"unicode/utf8"
)

// hold the basic functionalities to manipulate the source
// does not read the source from filesystem so the `code`
// should be prepared by the caller
type Source struct {
	id    int
	store *SourceStore

	// should be absolute path
	path   string
	code   string
	pos    int
	peeked []rune
	line   int
	col    int
}

func NewSource(path string, code string) *Source {
	return &Source{
		path:   path,
		code:   code,
		pos:    0,
		peeked: make([]rune, 0),
		line:   1,
		col:    0,
	}
}

const (
	EOF int32 = iota - 2
	EOL
)

func (s *Source) RuneAtPos(pos int) (rune, int) {
	r, size := utf8.DecodeRuneInString(s.code[pos:])
	if r == utf8.RuneError && size == 0 {
		r = EOF
	}
	return r, size
}

// read and push back a rune into `s.peaked` also advance `s.pos`
func (s *Source) PeekRune() rune {
	r, size := s.RuneAtPos(s.pos)
	s.peeked = append(s.peeked, r)
	s.pos += size
	return r
}

// firstly try to pop the front of the `s.peaked` otherwise read
// a rune and advance `s.pos`
func (s *Source) NextRune() rune {
	if len(s.peeked) > 0 {
		r, rest := s.peeked[0], s.peeked[1:]
		s.peeked = rest
		return r
	}

	r, size := s.RuneAtPos(s.pos)
	s.pos += size
	return r
}

func (s *Source) AheadIsCh(c rune) bool {
	return s.PeekRune() == c
}

func (s *Source) AheadIsChs2(c1 rune, c2 rune) bool {
	p1 := s.PeekRune()
	if p1 != c1 {
		return false
	}

	p2 := s.PeekRune()
	return p2 == c2
}

func (s *Source) AheadIsChOr(c1 rune, c2 rune) bool {
	p := s.PeekRune()
	return p == c1 || p == c2
}

func IsLineTerminator(c rune) bool {
	return c == 0x0a || c == 0x0d || c == 0x2028 || c == 0x2029
}

func (s *Source) ReadIfNextIs(c rune) rune {
	if s.PeekRune() == c {
		return s.NextRune()
	}
	return utf8.RuneError
}

func (s *Source) NextJoinCRLF() rune {
	c := s.NextRune()
	if IsLineTerminator(c) {
		if c == '\r' {
			s.ReadIfNextIs('\n')
		}
		return EOL
	}
	return c
}

func (s *Source) NextIsEOF() bool {
	return s.NextRune() == EOF
}

func (s *Source) Line() int {
	return s.line
}

func (s *Source) Pos() int {
	return s.pos - len(s.peeked)
}

func (s *Source) NewOpenRange() *SourceRange {
	return &SourceRange{
		src: s,
		lo:  s.Pos(),
	}
}

func (s *Source) next(loose bool) (rune, error) {
	c := s.NextJoinCRLF()

	if c == utf8.RuneError && !loose {
		return 0, NewLexerError(s.path, s.line, s.Pos()-1)
	}

	if c == EOL {
		s.line += 1
		s.col = 0
	} else {
		s.col += 1
	}
	return c, nil
}

func (s *Source) NextStrict() (rune, error) {
	return s.next(false)
}

func (s *Source) Next() rune {
	r, _ := s.next(true)
	return r
}

// TODO: description
type SourceStore struct {
	basepath string
	sources  map[int]*Source
	counter  int
	mu       sync.Mutex
}

func NewSourceStore(basepath string) *SourceStore {
	return &SourceStore{
		basepath: basepath,
		sources:  make(map[int]*Source),
	}
}

func (s *SourceStore) AddSource(file string) (int, error) {
	path := file
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Rel(s.basepath, path)
		if err != nil {
			return 0, err
		}
	}
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	id := s.AddSourceFromString(path, code)
	return id, nil
}

func (s *SourceStore) AddSourceFromString(file string, code []byte) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	src := NewSource(file, string(code))
	src.store = s
	src.id = s.counter
	s.sources[s.counter] = src
	s.counter += 1
	return s.counter
}

// TODO: description
type SourceRange struct {
	src *Source
	lo  int
	hi  int
}
