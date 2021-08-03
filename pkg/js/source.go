package js

import "unicode/utf8"

// hold the basic functionalities to manipulate the source
// does not read the source from filesystem so the `code`
// should be prepared by the caller
type Source struct {
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

func (s *Source) PeekRune() rune {
	r, size := s.RuneAtPos(s.pos)
	s.peeked = append(s.peeked, r)
	s.pos += size
	return r
}

// func (s *Source) PeekRune() rune {
// 	r, _ := s.PeekRuneSized()
// 	return r
// }

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

func (s *Source) next(loose bool) rune {
	c := s.NextJoinCRLF()

	if c == utf8.RuneError && !loose {
		// TODO: error report
		panic("deformed codepoint at")
	}

	if c == EOL {
		s.line += 1
		s.col = 0
	} else {
		s.col += 1
	}
	return c
}

func (s *Source) Next() rune {
	return s.next(false)
}

func (s *Source) NextLoose() rune {
	return s.next(true)
}
