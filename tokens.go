package main

import "fmt"

type TokenType int

const (
	//Literals
	ID TokenType = iota
	STRING
	NUM

	// Single Character Tokens
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
	MOD
	SLASH
	HASHTAG

	//One or two character tokens
	NOT
	NE
	EQ
	EQEQ
	GT
	GE
	LT
	LE

	//Keyword
	AND
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	TRUE
	WHILE
	PRINTLN
	CLASS
	THIS
	SUPER
	VAR

	EOF
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

func NewToken(tt TokenType, lexeme string, line int) Token {
	return Token{
		Type:   tt,
		Lexeme: lexeme,
		Line:   line,
	}
}

func (token *Token) String() string {
	return fmt.Sprintf("%v\t%v\t%v", token.Type.String(), token.Lexeme, token.Line)
}

func (t TokenType) String() string {
	switch t {
	case ID:
		return "ID"
	case STRING:
		return "STRING"
	case NUM:
		return "NUM"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case STAR:
		return "STAR"
	case MOD:
		return "MOD"
	case SLASH:
		return "SLASH"
	case HASHTAG:
		return "HASHTAG"
	case NOT:
		return "NOT"
	case NE:
		return "NE"
	case EQ:
		return "EQ"
	case EQEQ:
		return "EQEQ"
	case GT:
		return "GT"
	case GE:
		return "GE"
	case LT:
		return "LT"
	case LE:
		return "LE"
	case AND:
		return "AND"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case TRUE:
		return "TRUE"
	case WHILE:
		return "WHILE"
	case PRINTLN:
		return "PRINTLN"
	case CLASS:
		return "CLASS"
	case THIS:
		return "THIS"
	case SUPER:
		return "SUPER"
	case VAR:
		return "VAR"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
