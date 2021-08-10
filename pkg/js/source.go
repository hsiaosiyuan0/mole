package js

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"unicode"
	"unicode/utf8"
)

// hold the basic functionalities to manipulate the source
// does not read the source from filesystem so the `code`
// should be prepared by the caller
type Source struct {
	id    int
	store *SourceStore

	// should be absolute path
	path string
	code string
	pos  int
	line int
	col  int

	peeked    [sizeOfPeekedRune]rune
	peekedLen int
	peekedR   int
	peekedW   int

	metLineTerminator bool
}

const sizeOfPeekedRune = 4

func NewSource(path string, code string) *Source {
	return &Source{
		path: path,
		code: code,
		line: 1,
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
// returns `utf8.RuneError` if the rune is deformed
//
// be careful with the calling times of this method since it will panic
// if its internal buffer for caching peeked rune is full
func (s *Source) peek() rune {
	if s.peekedLen == sizeOfPeekedRune {
		panic(s.error(fmt.Sprintf("peek buffer of source is full, max len is %d\n", s.peekedLen)))
	}

	r, size := s.RuneAtPos(s.pos)
	if r == EOF {
		return EOF
	}

	s.peeked[s.peekedW] = r
	s.peekedW += 1
	s.peekedLen += 1
	if s.peekedW == sizeOfPeekedRune {
		s.peekedW = 0
	}
	s.pos += size
	return r
}

func (s *Source) Peek() rune {
	if s.peekedLen > 0 {
		return s.peeked[s.peekedR]
	}
	return s.peek()
}

func (s *Source) peekedRInc() int {
	r := s.peekedR + 1
	if r == sizeOfPeekedRune {
		return 0
	}
	return r
}

// firstly try to pop the front of the `s.peaked` otherwise read
// a rune and advance `s.pos`
func (s *Source) NextRune() rune {
	if s.peekedLen > 0 {
		r := s.peeked[s.peekedR]
		s.peekedR = s.peekedRInc()
		s.peekedLen -= 1
		return r
	}

	r, size := s.RuneAtPos(s.pos)
	s.pos += size
	return r
}

func (s *Source) AheadIsCh(c rune) bool {
	return s.Peek() == c
}

func (s *Source) AheadIsEof() bool {
	return s.pos == len(s.code)
}

func (s *Source) AheadIsChThenConsume(c rune) bool {
	if s.Peek() == c {
		s.Read()
		return true
	}
	return false
}

func (s *Source) AheadIsChs2(c1 rune, c2 rune) bool {
	if s.peekedLen < 2 {
		s.peek()
	}
	if s.peekedLen < 2 {
		s.peek()
	}
	return s.peeked[s.peekedR] == c1 && s.peeked[s.peekedRInc()] == c2
}

func (s *Source) AheadIsChOr(c1 rune, c2 rune) bool {
	c := s.Peek()
	return c == c1 || c == c2
}

func IsLineTerminator(c rune) bool {
	return c == 0x0a || c == 0x0d || c == 0x2028 || c == 0x2029
}

func (s *Source) ReadIfNextIs(c rune) bool {
	if s.Peek() == c {
		s.NextRune()
		return true
	}
	return false
}

func (s *Source) readJoinCRLF() rune {
	c := s.NextRune()
	if IsLineTerminator(c) {
		if c == '\r' {
			s.ReadIfNextIs('\n')
		}
		return EOL
	}
	return c
}

func (s *Source) Line() int {
	return s.line
}

func (s *Source) Pos() int {
	return s.pos - s.peekedLen
}

func (s *Source) NewOpenRange() *SourceRange {
	pos := s.Pos()
	return &SourceRange{
		src: s,
		lo:  pos,
		hi:  pos,
	}
}

// returns `utf8.RuneError` if the rune is deformed
func (s *Source) Read() rune {
	c := s.readJoinCRLF()
	if c == EOL {
		s.line += 1
		s.col = 0
	} else {
		s.col += 1
	}
	return c
}

// skip spaces except line terminator
func (s *Source) SkipSpace() {
	s.metLineTerminator = false
	for {
		c := s.Peek()
		if unicode.IsSpace(c) {
			if IsLineTerminator(c) {
				s.metLineTerminator = true
			}
			s.Read()
		} else {
			break
		}
	}
}

func (s *Source) error(msg string) *SourceError {
	return NewSourceError(msg, s.path, s.line, s.Pos()-1)
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

func (s *SourceRange) Text() string {
	return s.src.code[s.lo:s.hi]
}
