package js

import "fmt"

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
	return fmt.Sprintln("un")
}

func TestLexer() {
	fmt.Println("1")
}
