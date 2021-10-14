package parser

import "fmt"

type SourceError struct {
	msg  string
	file string
	line int
	col  int
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("%s at %sL%d:%d\n", e.msg, e.file, e.line, e.col)
}

func NewSourceError(msg, file string, line, col int) *SourceError {
	return &SourceError{
		file: file,
		line: line,
		col:  col,
	}
}

type LexerError struct {
	msg  string
	file string
	line int
	col  int
}

func NewLexerError(msg, file string, line, col int) *LexerError {
	return &LexerError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}

type ParserError struct {
	msg  string
	file string
	line int
	col  int
}

func NewParserError(msg, file string, line, col int) *ParserError {
	return &ParserError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}
