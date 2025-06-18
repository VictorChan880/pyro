package main

type Scanner struct {
	Source   string
	Tokens   []Token
	Start    int
	Current  int
	Line     int
	Keywords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	var keywords = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	return &Scanner{
		Source:   source,
		Start:    0,
		Current:  0,
		Line:     1,
		Keywords: keywords,
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {

		s.Start = s.Current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, NewToken(EOF, "", s.Line))
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(LPAREN)
	case ')':
		s.addToken(RPAREN)
	case '{':
		s.addToken(LBRACE)
	case '}':
		s.addToken(RBRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		tt := NOT
		if s.match('=') {
			tt = NE
		}
		s.addToken(tt)
	case '/':
		s.addToken(SLASH)
	case '%':
		s.addToken(MOD)
	case '=':
		tt := EQ
		if s.match('=') {
			tt = EQEQ
		}
		s.addToken(tt)
	case '>':
		tt := GT
		if s.match('=') {
			tt = GE
		}
		s.addToken(tt)
	case '<':
		tt := LT
		if s.match('=') {
			tt = LE
		}
		s.addToken(tt)
	case '#':
		for s.peek() != '\n' && !s.isAtEnd() {
			s.advance()
		}
	case '"':
		s.scanString()
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.Line++
	default:
		if isDigit(c) {
			s.scanNum()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			err := NewError(s.Line, "Unexpected character.", "")
			report(err)
		}
	}
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]
	tt, exists := s.Keywords[text]
	if !exists {
		tt = ID
	}
	s.addToken(tt)

}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isAlphaNumeric(c rune) bool {
	return (isDigit(c) || isAlpha(c))
}
func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		err := NewError(s.Line, "Unterminated string.", "")
		report(err)
		return
	}

	s.advance() //closing "

	s.addTokenString()

}

func (s *Scanner) scanNum() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}
	s.addToken(NUM)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return []rune(s.Source)[s.Current]
}

func (s *Scanner) peekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return ' '
	}
	return []rune(s.Source)[s.Current+1]
}

func (s *Scanner) match(c rune) bool {
	if s.isAtEnd() {
		return false
	}
	if []rune(s.Source)[s.Current] != c {
		return false
	}
	s.Current++
	return true
}
func (s *Scanner) advance() rune {
	s.Current += 1
	return []rune(s.Source)[s.Current-1]
}

func (s *Scanner) addToken(tt TokenType) {
	s.addTokenScanner(tt)
}

func (s *Scanner) addTokenScanner(tt TokenType) {
	value := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, NewToken(tt, value, s.Line))
}

func (s *Scanner) addTokenString() {
	value := s.Source[s.Start+1 : s.Current-1]
	s.Tokens = append(s.Tokens, NewToken(STRING, value, s.Line))

}
