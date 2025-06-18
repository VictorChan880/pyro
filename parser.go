package main

import (
	"strconv"
)

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)

	for !p.isAtEnd() {
		declar, err := p.declaration()
		if err != nil {
			if _, isParse := err.(ParseError); !isParse {
				return statements, err
			}
		} else {
			statements = append(statements, declar)
		}
	}

	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		declar, err := p.varDeclaration()

		if _, isParse := err.(ParseError); isParse {
			p.synchronize()
			return declar, nil
		}
		return declar, err
	}
	stmt, err := p.statement()
	if _, isParse := err.(ParseError); isParse {
		p.synchronize()
	}
	return stmt, err
}

func (p *Parser) varDeclaration() (Var, error) {
	name, err := p.consume(ID, "Expect variable name")
	if err != nil {
		return Var{}, err
	}

	var initalizer *Expr

	if p.match(EQ) {
		temp, err := p.expression()
		initalizer = &temp
		if err != nil {
			return Var{}, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after value")
	if err != nil {
		return Var{}, err
	}

	return NewVar(name, initalizer), nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		printStmt, err := p.printStatement()
		if err != nil {
			return nil, err
		}
		return printStmt, nil
	} else if p.match(LBRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(block), nil
	} else if p.match(IF) {
		ifStmt, err := p.ifStatement()
		return ifStmt, err
	} else if p.match(WHILE) {
		whileStmt, err := p.whileStatement()
		return whileStmt, err
	} else if p.match(FOR) {
		_, err := p.consume(LPAREN, "Expect '(' after for")
		if err != nil {
			return nil, err
		}

		var intializer Stmt
		if p.match(SEMICOLON) {
			intializer = nil
		} else if p.match(VAR) {
			intializer, err = p.varDeclaration()
		} else {
			intializer, err = p.expressionStatement()
		}
		if err != nil {
			return nil, err
		}

		var condition *Expr
		if !p.check(SEMICOLON) {
			temp, err := p.expression()
			condition = &temp
			if err != nil {
				return nil, err
			}
		}
		_, err = p.consume(SEMICOLON, "Expect ';' after loop condition")
		if err != nil {
			return nil, err
		}

		var increment *Expr
		if !p.check(RPAREN) {
			temp, err := p.expression()
			increment = &temp
			if err != nil {
				return nil, err
			}
		}
		_, err = p.consume(RPAREN, "Expect ')' after for clauses")
		if err != nil {
			return nil, err
		}

		body, err := p.statement()
		if err != nil {
			return nil, err
		}

		if increment != nil {
			statements := []Stmt{body, NewExpression(*increment)}
			body = NewBlock(statements)
		}

		if condition == nil {
			body = NewWhile(NewLiteral(true), body)
		} else {
			body = NewWhile(*condition, body)
		}

		if intializer != nil {
			statements := []Stmt{intializer, body}
			body = NewBlock(statements)
		}
		return body, nil

	} else if (p.match(FUN)) {
		function, err := p.function("function")
		if err != nil {
			return nil, err
		}
		return function, nil
	}

	expressionStmt, err := p.expressionStatement()
	return expressionStmt, err
}

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(ID, "Expect " + kind + " name")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LPAREN, "Expect '(' after " + kind + "name")
	if err != nil {
		return nil, err
	}

	parameters := make([]Token, 0) 

	if !p.check(RPAREN) {
		for {
			if len(parameters) >= 255 {
				report(NewParseError(p.peek(), "Can't have more than 255 parameters").Err)
			}
			parameter, err := p.consume(ID, "Expected parameter name")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameter)
			if !p.match(COMMA) {
				break
			}
		}
	}
	_, err = p.consume(RPAREN, "Expected ')' after paramters")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LBRACE, "Expect '{' before " + kind + " body")
	if err != nil {
		return nil, err
	}

	body, err  := p.block()
	if err != nil {
		return nil, err
	}

	function := NewFunction(name, parameters, body)
	return function, err 
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LPAREN, "Expected '(' after while")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RPAREN, "Expected ')' after while condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return NewWhile(condition, body), nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LPAREN, "Expected '(' after 'if'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RPAREN, "Expected ')' after 'if' condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch *Stmt
	if p.match(ELSE) {
		temp, err := p.statement()
		elseBranch = &temp
		if err != nil {
			return nil, err
		}
	}

	return NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(RBRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := p.consume(RBRACE, "Expect '}' after block")
	if err != nil {
		return nil, err
	}

	return statements, nil

}

func (p *Parser) printStatement() (Print, error) {
	value, err := p.expression()
	if err != nil {
		return Print{}, err
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after value")
	if err != nil {
		return Print{}, err
	}
	return NewPrint(value), nil
}

func (p *Parser) expressionStatement() (Expression, error) {
	expr, err := p.expression()
	if err != nil {
		return Expression{}, err
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after value")
	if err != nil {
		return Expression{}, err
	}

	return NewExpression(expr), nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQ) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, isVariable := expr.(Variable); isVariable {
			name := v.Name
			return NewAssign(name, value), nil
		}

		pErr := NewParseError(equals, "Invalid assignment target.")
		report(pErr.Err)

		return nil, pErr
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(NE, EQEQ) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(LT, LE, GT, GE) {
		operator := p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()

	if err != nil {
		return nil, err
	}
	for p.match(PLUS, MINUS, MOD) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(NE, MINUS) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return NewUnary(operator, right), nil
	}
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LPAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)

	if !p.check(RPAREN) {
		for {
			if len(arguments) >= 255 {
				argError := NewParseError(p.peek(), "Can't have more than 255 arguments")
				report(argError.Err)
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RPAREN, "Expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	return NewCall(callee, paren, arguments), nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(NIL) {
		return NewLiteral(nil), nil
	} else if p.match(TRUE) {
		return NewLiteral(true), nil
	} else if p.match(FALSE) {
		return NewLiteral(false), nil
	} else if p.match(NUM) {
		num, _ := strconv.ParseFloat(p.previous().Lexeme, 64)
		return NewLiteral(num), nil
	} else if p.match(STRING) {
		return NewLiteral(p.previous().Lexeme), nil
	} else if p.match(LPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RPAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}

		return expr, nil //changed grouping
	} else if p.match(ID) {
		return NewVariable(p.previous()), nil
	}
	err := NewParseError(p.peek(), "Expected expression")

	report(err.Err)
	return nil, err
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}
		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) consume(tt TokenType, errMsg string) (Token, error) {
	if p.check(tt) {
		return p.advance(), nil
	}
	parseError := NewParseError(p.peek(), errMsg)
	report(parseError.Err)
	return Token{}, parseError
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tt := range types {
		if p.check(tt) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tt
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}
