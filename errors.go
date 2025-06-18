package main

import "fmt"

type Error struct {
	Line    int
	Message string
	Where   string
	Token   *Token
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", e.Line, e.Where, e.Message)
}

type ParseError struct {
	Err Error
}

type RunTimeError struct {
	Err   Error
	Token Token
}

func (rte RunTimeError) Error() string {
	return rte.Err.Error()
}

func (pe ParseError) Error() string {
	return pe.Err.Error()
}

func NewRunTimeError(token Token, message string) RunTimeError {
	err := Error{
		Line:    token.Line,
		Message: message,
	}

	return RunTimeError{
		Err:   err,
		Token: token,
	}
}
func NewParseError(token Token, message string) ParseError {
	err := Error{
		Line:    token.Line,
		Message: message,
	}

	if token.Type == EOF {
		err.Where = " at end"
	} else {
		err.Where = " at '" + token.Lexeme + "'"
	}
	return ParseError{Err: err}
}

func NewError(line int, msg string, where string) Error {
	return Error{
		Line:    line,
		Message: msg,
		Where:   where,
	}
}

func report(err Error) {
	fmt.Println("[line ", err.Line, " Error"+err.Where+": "+err.Message + "]")
	hasError = true
}
