package main

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitVarStmt(stmt Var) error
	VisitPrintStmt(stmt Print) error
	VisitExpressionStmt(stmt Expression) error
	VisitBlockStmt(stmt Block) error
	VisitIfStmt(stmt If) error
	VisitWhileStmt(stmt While) error
	VisitFunctionStmt(stmt Function) error
}

type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewFunction(name Token, params []Token, body []Stmt) Function {
	return Function{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (f Function) Accept(visitor StmtVisitor) error {
	return visitor.VisitFunctionStmt(f)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) While {
	return While{
		Condition: condition,
		Body:      body,
	}
}

func (w While) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(w)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch *Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch *Stmt) If {
	return If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (i If) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(i)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) Block {
	return Block{
		Statements: statements,
	}
}

func (b Block) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(b)
}

type Var struct {
	Name       Token
	Initalizer *Expr
}

func (v Var) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(v)
}

func NewVar(name Token, intializer *Expr) Var {
	return Var{
		Name:       name,
		Initalizer: intializer,
	}
}

type Print struct {
	Expression Expr
}

func (p Print) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(p)
}

func NewPrint(expr Expr) Print {
	return Print{
		Expression: expr,
	}
}

type Expression struct {
	Expression Expr
}

func (e Expression) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(e)
}

func NewExpression(expr Expr) Expression {
	return Expression{
		Expression: expr,
	}
}
