package span

import (
	"unicode"
	"unicode/utf8"
)

const RUNES_BUF_LEN = 4

type RunesBuf struct {
	buf      [RUNES_BUF_LEN]rune
	len      int // len of the alive cells in `buf`
	byte_len int // byte length occupied by all the alive cells
	r        int // the index in `buf` to perform next read
	w        int // the index in `buf` to perform next write
}

func (b *RunesBuf) incW() int {
	w := b.w + 1
	if w == RUNES_BUF_LEN {
		return 0
	}
	return w
}

func (b *RunesBuf) incR() int {
	r := b.r + 1
	if r == RUNES_BUF_LEN {
		return 0
	}
	return r
}

func (b *RunesBuf) writable() bool {
	return b.len < RUNES_BUF_LEN
}

func (b *RunesBuf) readable() bool {
	return b.len > 0
}

func (b *RunesBuf) write(r rune, len int) {
	if !b.writable() {
		panic("no place in runes buf")
	}
	b.buf[b.w] = r
	b.byte_len += len
	b.len += 1
	b.w = b.incW()
}

func (b *RunesBuf) read() rune {
	if !b.readable() {
		panic("runes buf is empty")
	}
	r := b.buf[b.r]
	b.byte_len -= utf8.RuneLen(r)
	b.len -= 1
	b.r = b.incR()
	return r
}

func (rb *RunesBuf) cur() rune {
	return rb.buf[rb.r]
}

type SourceState struct {
	rb          RunesBuf // runes buf to avoid decode rune twice
	ofst        int      // ofst base on byte
	pos         int      // pos base on rune
	line        int      // 1-based line number
	col         int      // 0-based column number
	metLineTerm bool     // whether line terminator appeared
}

func newSourceState() SourceState {
	s := SourceState{}
	s.line = 1
	return s
}

// process the utf8 encoded source file, it will panic if:
// - underlying runes buf is out of bounds
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
		state: newSourceState(),
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
func (s *Source) advance(cache bool) rune {
	ofst := s.state.ofst
	r, size := utf8.DecodeRuneInString(s.code[ofst:])
	if r == utf8.RuneError {
		if size == 0 {
			r = EOF
		}
	}

	if r == EOF {
		return EOF
	}

	s.state.ofst += size
	s.state.pos += 1
	if cache {
		s.state.rb.write(r, size)
	}
	return r
}

// ensure the peeked buffer have 2 chars and return the 2nd
func (s *Source) Ahead2nd() rune {
	if s.state.rb.len < 2 {
		s.advance(true)
	}
	if s.state.rb.len < 2 {
		s.advance(true)
	}
	if s.state.rb.len >= 2 {
		return s.state.rb.buf[s.state.rb.incR()]
	}
	return EOF
}

func (s *Source) AheadIsCh(c rune) bool {
	return s.Peek() == c
}

func (s *Source) AheadIsEOF() bool {
	return s.state.ofst == len(s.code) && s.state.rb.len == 0
}

// read from the runes buf at first, otherwise advance the
// underlying position
func (s *Source) readRune() rune {
	if s.state.rb.readable() {
		return s.state.rb.read()
	}
	return s.advance(false)
}

func (s *Source) Peek() rune {
	if !s.state.rb.readable() {
		return s.advance(true)
	}
	return s.state.rb.cur()
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
	r := s.readRune()
	if IsLineTerminator(r) {
		if r == '\r' {
			if s.Peek() == '\n' {
				s.readRune()
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
	a2 := s.Ahead2nd()
	return s.state.rb.cur() == c1 && a2 == c2
}

func (s *Source) AheadIsChOr(c1 rune, c2 rune) bool {
	r := s.Peek()
	return r == c1 || r == c2
}

func IsLineTerminator(c rune) bool {
	return c == 0x0a || c == 0x0d || c == 0x2028 || c == 0x2029
}

// ofst base on byte
func (s *Source) Ofst() int {
	return s.state.ofst - s.state.rb.byte_len
}

// pos base on rune
func (s *Source) Pos() int {
	return s.state.pos - s.state.rb.len
}

func (s *Source) Line() int {
	return s.state.line
}

func (s *Source) Col() int {
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
func (s *Source) Text(start, end int) string {
	return s.code[start:end]
}

type Range struct {
	Src *Source
	Lo  int
	Hi  int
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
