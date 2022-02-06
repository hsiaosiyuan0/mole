package span

import (
	"unicode"
	"unicode/utf8"
)

type Runes struct {
	raw string
}

// read and return a rune from the beginning of
// the underlying `bytes` as well as advance the head
// of bytes to the beginning of the next rune
//
// returned rune may be either `EOF` or `utf8.RuneError`
// in both case with size zero
func (s *Runes) advance() (rune, int) {
	r, size := utf8.DecodeRuneInString(s.raw)
	if r == utf8.RuneError {
		if size == 0 {
			r = EOF
		}
		size = 0
	} else {
		s.raw = s.raw[size:]
	}
	return r, size
}

type SourceState struct {
	runes       Runes
	ofst        uint32 // ofst base on byte
	pos         uint32 // pos base on rune
	line        uint32 // 1-based line number
	col         uint32 // 0-based column number
	metLineTerm bool   // whether line terminator appeared
}

func newSourceState(code string) SourceState {
	s := SourceState{}
	s.runes.raw = code
	s.line = 1
	return s
}

// process the utf8 encoded source file, it will panic if:
//
// the returned rune will be `utf8.RuneError` if the position
// is not well encoded in utf8
//
// `\r`, `\r\n` will be unified to `\n` which has an alias `span.EOL`
type Source struct {
	Path string
	code string

	// state is encapsulated in `SourceState` for push/pop easily
	state SourceState

	// a stack to mark/remark the parsing position to deal with some
	// complex syntaxes such as the `<` in typescript which can be
	// either the beginning of type params(`a<b>()`) or operator in
	// binary expression(`a < b`)
	ss []SourceState
}

func NewSource(path string, code string) *Source {
	return &Source{
		Path:  path,
		code:  code,
		state: newSourceState(code),
		ss:    make([]SourceState, 0),
	}
}

const (
	// used to unify `\r`, `\r\n` to `\n`, despite other line-terminators in unicode
	EOL = rune('\n')
	EOF = rune(-1)
)

// push current state into the state-stack
func (s *Source) PushState() {
	s.ss = append(s.ss, s.state)
}

// pop the state-stack and discard the popped state
func (s *Source) DiscardState() {
	last := len(s.ss) - 1
	s.ss = s.ss[:last]
}

// pop the state-stack and apply the popped state
func (s *Source) PopState() {
	last := len(s.ss) - 1
	rest, state := s.ss[:last], s.ss[last]
	s.ss = rest
	s.state = state
}

// advance the underlying chars and cache the advanced char
func (s *Source) advance() rune {
	r, size := s.state.runes.advance()
	if r == EOF {
		return EOF
	}

	s.state.ofst += uint32(size)
	s.state.pos += 1
	return r
}

func (s *Source) Peek() rune {
	runes := s.state.runes
	r, _ := runes.advance()
	return r
}

func (s *Source) Ahead2nd() rune {
	runes := s.state.runes
	r, size := runes.advance()
	if size == 0 {
		return r
	}
	r, _ = runes.advance()
	return r
}

func (s *Source) AheadIsCh(c rune) bool {
	return s.Peek() == c
}

func (s *Source) AheadIsEOF() bool {
	return s.state.ofst == uint32(len(s.code))
}

func (s *Source) AdvanceIf(c rune) bool {
	if s.Peek() == c {
		s.Read()
		return true
	}
	return false
}

// read rune as well as join the CR/LF to record the line/column info
func (s *Source) Read() rune {
	r := s.advance()
	if IsLineTerminator(r) {
		if r == '\r' {
			if s.Peek() == '\n' {
				s.advance()
			}
		}
	}
	if r == '\r' || r == '\n' {
		r = EOL
		s.state.line += 1
		s.state.col = 0
	} else {
		s.state.col += 1
	}
	return r
}

func (s *Source) AheadIsChs2(c1 rune, c2 rune) bool {
	r1 := utf8.RuneError
	r2 := utf8.RuneError
	size := 0
	runes := s.state.runes
	r1, size = runes.advance()
	if size != 0 {
		r2, _ = runes.advance()
	}
	return r1 == c1 && r2 == c2
}

func (s *Source) AheadIsChOr(c1 rune, c2 rune) bool {
	r := s.Peek()
	return r == c1 || r == c2
}

func IsLineTerminator(c rune) bool {
	return c == 0x0a || c == 0x0d || c == 0x2028 || c == 0x2029
}

// ofst base on byte
func (s *Source) Ofst() uint32 {
	return s.state.ofst
}

// pos base on rune
func (s *Source) Pos() uint32 {
	return s.state.pos
}

func (s *Source) Line() uint32 {
	return s.state.line
}

func (s *Source) Col() uint32 {
	return s.state.col
}

func (s *Source) MetLineTerm() bool {
	return s.state.metLineTerm
}

func (s *Source) NewOpenRange() *Range {
	return &Range{
		Src: s,
		Lo:  s.Ofst(),
		Hi:  s.Ofst(),
	}
}

func (s *Source) OpenRange(rng *Range) *Range {
	rng.Src = s
	rng.Lo = s.Ofst()
	rng.Hi = s.Ofst()
	return rng
}

func (s *Source) SkipSpace() *Source {
	s.state.metLineTerm = false
	for {
		c := s.Peek()
		if !unicode.IsSpace(c) || c == EOF {
			break
		}
		if IsLineTerminator(c) {
			s.state.metLineTerm = true
		}
		s.Read()
	}
	return s
}

// return the string in the span `[start,end)`
func (s *Source) Text(start, end uint32) string {
	return s.code[start:end]
}

type Range struct {
	Src *Source
	Lo  uint32
	Hi  uint32
}

func (r *Range) Text() string {
	return r.Src.code[r.Lo:r.Hi]
}

func (r *Range) Clone() *Range {
	return &Range{
		Src: r.Src,
		Lo:  r.Lo,
		Hi:  r.Hi,
	}
}
