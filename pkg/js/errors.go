package js

import "fmt"

type SourceError struct {
	file string
	line int
	col  int
}

func NewSourceError(file string, line, col int) *SourceError {
	return &SourceError{
		file: file,
		line: line,
		col:  col,
	}
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("unexpected rune at %sL%d:%d\n", e.file, e.line, e.col)
}

type LexerError struct {
	file string
	line int
	col  int
}

func NewLexerError(file string, line, col int) *LexerError {
	return &LexerError{
		file: file,
		line: line,
		col:  col,
	}
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("unexpected token at %sL%d:%c\n", e.file, e.line, e.col)
}
