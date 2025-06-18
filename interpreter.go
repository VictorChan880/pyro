package main

import (
	"fmt"
	"math"
	"strconv"
)

type Interpreter struct {
	Environment *Environment
	Globals     *Environment
}

func (a *Interpreter) interpret(statements []Stmt) error {
	for _, statement := range statements {
		err := a.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Interpreter) execute(stmt Stmt) error {
	return stmt.Accept(a)
}
func stringify(value interface{}) string {
	switch v := value.(type) {
	case float64:
		str := fmt.Sprintf("%g", v) // compact format (e.g., avoids trailing .0 by default)
		return str
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (a *Interpreter) evalute(expr Expr) (interface{}, error) {
	return expr.Accept(a)
}

func (a *Interpreter) VisitCallExpr(expr Call) (interface{}, error) {
	callee, err := a.evalute(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, 0)
	for _, argument := range expr.Arguments {
		eval, err := a.evalute(argument)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, eval)
	}

	function, isCallable := callee.(Callable)
	if !isCallable {
		rtErr := NewRunTimeError(expr.Paren, "Can only call functions and classes.")
		report(rtErr.Err)
		return nil, rtErr
	}
	if len(arguments) != function.Arity() {
		rtErr := NewRunTimeError(expr.Paren, "Expected "+strconv.Itoa(function.Arity())+" arguments but got "+strconv.Itoa(len(arguments)))
		report(rtErr.Err)
		return nil, rtErr
	}

	return function.Call(a, arguments)
}

func (a *Interpreter) VisitFunctionStmt(stmt Function) error {
	function := NewPyroFunction(stmt)
	a.Environment.define(function.Declaration.Name.Lexeme, function)
	return nil
}

func (a *Interpreter) VisitWhileStmt(expr While) error {
	cond, err := a.evalute(expr.Condition)
	if err != nil {
		return err
	}
	for isTruthy(cond) {
		err = a.execute(expr.Body)
		if err != nil {
			return err
		}

		cond, err = a.evalute(expr.Condition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Interpreter) VisitLogicalExpr(expr Logical) (interface{}, error) {
	left, err := a.evalute(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return a.evalute(expr.Right)
}

func (a *Interpreter) VisitIfStmt(stmt If) error {
	cond, err := a.evalute(stmt.Condition)
	if err != nil {
		return err
	}

	if isTruthy(cond) {
		err = a.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		err = a.execute(*stmt.ElseBranch)
	}

	return err
}

func (a *Interpreter) VisitBlockStmt(stmt Block) error {
	err := a.executeBlock(stmt.Statements, NewEnclosedEnvironment(a.Environment))
	return err
}

func (a *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := a.Environment
	a.Environment = environment

	for _, statement := range statements {
		err := a.execute(statement)
		if err != nil {
			a.Environment = previous
			return err
		}
	}
	a.Environment = previous
	return nil
}

func (a *Interpreter) VisitAssignExpr(expr Assign) (interface{}, error) {
	value, err := a.evalute(expr.Value)
	if err != nil {
		return nil, err
	}

	err = a.Environment.assign(expr.Name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (a *Interpreter) VisitVarStmt(stmt Var) error {
	var value interface{}
	var err error
	if stmt.Initalizer != nil {
		value, err = a.evalute(*stmt.Initalizer)
		if err != nil {
			return err
		}
	}
	a.Environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (a *Interpreter) VisitVariableExpr(expr Variable) (interface{}, error) {
	value, err := a.Environment.get(expr.Name)
	return value, err
}

func (a *Interpreter) VisitExpressionStmt(stmt Expression) error {
	_, err := a.evalute(stmt.Expression)
	return err
}

func (a *Interpreter) VisitPrintStmt(stmt Print) error {
	value, err := a.evalute(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func (a *Interpreter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return expr.Value, nil
}

func (a *Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return a.evalute(expr.Expression)
}

func (a *Interpreter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	right, err := a.evalute(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case NE:
		return !isTruthy(right), nil
	case MINUS:
		err := checkNumOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}

	//unreachable
	return nil, nil
}

func (a *Interpreter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	left, err := a.evalute(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := a.evalute(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case GT:
		{
			err := checkNumOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) > right.(float64), nil
		}
	case GE:
		{
			err := checkNumOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) >= right.(float64), nil
		}
	case LT:
		{
			err := checkNumOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) < right.(float64), nil
		}
	case LE:
		{
			err := checkNumOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) <= right.(float64), nil
		}
	case MINUS:
		{
			err := checkNumOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) - right.(float64), nil
		}
	case NE:
		return !isEqual(left, right), nil
	case EQEQ:
		return isEqual(left, right), nil
	case PLUS:
		lStr, lIsStr := left.(string)
		rStr, rIsStr := right.(string)

		if lIsStr && rIsStr {
			return lStr + rStr, nil
		}

		lNum, lIsNum := left.(float64)
		rNum, rIsNum := right.(float64)

		if lIsNum && rIsNum {
			return lNum + rNum, nil
		}

		err := NewRunTimeError(expr.Operator, "Operands must be two nums or two strings")
		report(err.Err)
		return nil, err

	case SLASH:
		err := checkNumOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		err := checkNumOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case MOD:
		err := checkNumOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return math.Mod(left.(float64), right.(float64)), nil

	}

	//unreachable
	return nil, nil
}

func checkNumOperands(operator Token, left interface{}, right interface{}) error {
	_, lIsNum := left.(float64)
	_, rIsNum := right.(float64)

	if lIsNum && rIsNum {
		return nil
	}
	err := NewRunTimeError(operator, "Operands must be a number")
	report(err.Err)
	return err

}

func checkNumOperand(operator Token, operand interface{}) error {
	if _, isFloat := operand.(float64); isFloat {
		return nil
	}
	err := NewRunTimeError(operator, "Operand must be a number")
	report(err.Err)
	return err

}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if b, isBool := obj.(bool); isBool {
		return b
	}
	return true

}
