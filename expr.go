package main

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(expr Binary) (interface{}, error)
	VisitUnaryExpr(expr Unary) (interface{}, error)
	VisitLiteralExpr(expr Literal) (interface{}, error)
	VisitGroupingExpr(expr Grouping) (interface{}, error)
	VisitVariableExpr(expr Variable) (interface{}, error)
	VisitAssignExpr(expr Assign) (interface{}, error)
	VisitLogicalExpr(expr Logical) (interface{}, error)
	VisitCallExpr(expr Call) (interface{}, error)

}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewBinary(left Expr, op Token, right Expr) Binary {
	return Binary{
		Left:     left,
		Operator: op,
		Right:    right,
	}
}

type Unary struct {
	Operator Token
	Right    Expr
}

func NewUnary(op Token, right Expr) Unary {
	return Unary{
		Operator: op,
		Right:    right,
	}
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Literal {
	return Literal{
		Value: value,
	}
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) Grouping {
	return Grouping{
		Expression: expr,
	}
}

type Variable struct {
	Name Token
}

func NewVariable(name Token) Variable {
	return Variable{
		Name: name,
	}
}

type Assign struct {
	Name  Token
	Value Expr
}

func NewAssign(name Token, value Expr) Assign {
	return Assign{
		Name:  name,
		Value: value,
	}
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogical(left Expr, operator Token, right Expr) Logical {
	return Logical{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

type Call struct {
	Callee Expr
	Paren Token
	Arguments []Expr
}

func NewCall(callee Expr, paren Token, arguments []Expr) Call {
	return Call{
		Callee: callee,
		Paren: paren,
		Arguments: arguments,
	}
}

func (b Binary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

func (u Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

func (u Literal) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(u)
}

func (g Grouping) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

func (v Variable) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(v)
}
func (a Assign) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(a)
}

func (l Logical) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(l)
}

func (c Call) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(c)
}




