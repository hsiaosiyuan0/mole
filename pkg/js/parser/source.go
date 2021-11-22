package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Source struct {
	path string
	code string

	ofst int // ofst base on byte
	pos  int // pos base on codepint
	line int
	col  int

	peeked [peekBufLen]rune
	prl    [peekBufLen]int // byte len of each rune in `peeked`
	pl     int             // len of `peeked`
	pbl    int             // bytes length of peeked runes
	pr     int             // offset in buffer for reading
	pw     int             // offset in buffer for writing

	metLineTerminator bool
}

const peekBufLen = 4

func NewSource(path string, code string) *Source {
	return &Source{
		path: path,
		code: code,
		line: 1,
	}
}

const (
	EOF = rune(-1)
	EOL = rune(0x0a)
)

func (s *Source) RuneAtOfst(ofst int) (rune, int) {
	r, size := utf8.DecodeRuneInString(s.code[ofst:])
	if r == utf8.RuneError && size == 0 {
		r = EOF
	}
	return r, size
}

// read and push back a rune into `s.peaked` as well as advance `s.ofst`,
// return `utf8.RuneError` if the rune is deformed
//
// be careful with the calling times of this method since it will panic
// if its internal buffer for caching peeked rune is full
func (s *Source) peekGrow() rune {
	if s.pl == peekBufLen {
		panic(s.error(fmt.Sprintf("peek buffer of source is full, max len is %d\n", s.pl)))
	}

	r, size := s.RuneAtOfst(s.ofst)
	if r == EOF {
		return EOF
	}

	s.peeked[s.pw] = r
	s.prl[s.pw] = size

	s.pw = s.pwInc()
	s.pl += 1
	s.pbl += size

	s.ofst += size
	s.pos += 1
	return r
}

func (s *Source) Peek() rune {
	if s.pl > 0 {
		return s.peeked[s.pr]
	}
	return s.peekGrow()
}

func (s *Source) pwInc() int {
	w := s.pw + 1
	if w == peekBufLen {
		return 0
	}
	return w
}

func (s *Source) prInc() int {
	r := s.pr + 1
	if r == peekBufLen {
		return 0
	}
	return r
}

// try to pop the front of the `s.peaked` otherwise read
// a rune and advance `s.ofst`
func (s *Source) NextRune() rune {
	if s.pl > 0 {
		pr := s.pr
		r := s.peeked[pr]
		s.pr = s.prInc()
		s.pl -= 1
		s.pbl -= s.prl[pr]
		return r
	}

	r, size := s.RuneAtOfst(s.ofst)
	s.ofst += size
	s.pos += 1
	return r
}

func (s *Source) AheadIsCh(c rune) bool {
	return s.Peek() == c
}

func (s *Source) AheadIsEOF() bool {
	return s.ofst == len(s.code) && s.pl == 0
}

func (s *Source) Ahead2() rune {
	if s.pl < 2 {
		s.peekGrow()
	}
	if s.pl < 2 {
		s.peekGrow()
	}
	if s.pl < 2 {
		return utf8.RuneError
	}
	return s.peeked[s.prInc()]
}

func (s *Source) AheadIsChs2(c1 rune, c2 rune) bool {
	a2 := s.Ahead2()
	return s.peeked[s.pr] == c1 && a2 == c2
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
		s.Read()
		return true
	}
	return false
}

// join CR，LF, returns `utf8.RuneError` if the rune is deformed
func (s *Source) Read() rune {
	c := s.NextRune()
	r := c
	if IsLineTerminator(c) {
		if c == '\r' {
			if s.Peek() == '\n' {
				s.NextRune()
			}
		}
		r = EOL
	}
	if c == '\r' || c == '\n' {
		s.line += 1
		s.col = 0
	} else {
		s.col += 1
	}
	return r
}

// ofst base on byte
func (s *Source) Ofst() int {
	return s.ofst - s.pbl
}

// pos base on codepint
func (s *Source) Pos() int {
	return s.pos - s.pl
}

func (s *Source) NewOpenRange() *SourceRange {
	return &SourceRange{
		src: s,
		lo:  s.Ofst(),
		hi:  s.Ofst(),
	}
}

func (s *Source) openRange(rng *SourceRange) *SourceRange {
	rng.src = s
	rng.lo = s.Ofst()
	rng.hi = s.Ofst()
	return rng
}

func (s *Source) SkipSpace() *Source {
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
	return s
}

func (s *Source) error(msg string) *SourceError {
	return NewSourceError(msg, s.path, s.line, s.Ofst()-1)
}

type SourceRange struct {
	src *Source
	lo  int
	hi  int
}

func (r *SourceRange) Text() string {
	return r.src.code[r.lo:r.hi]
}
