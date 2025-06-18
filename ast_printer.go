package main

// import "fmt"

// type AstPrinter struct{}

// func (a AstPrinter) Print(expr Expr) (interface{}, error) {
// 	_, err := expr.Accept(a)
// 	return nil, err
// }


// func (a AstPrinter) VisitVariableExpr(expr Variable) (interface{}, error) {
// 	fmt.Print("Var ")
// 	fmt.Print(expr.Name.Lexeme)
// 	return nil, nil
// }

// func (a AstPrinter) VisitBinaryExpr(expr Binary) (interface{}, error) {
// 	fmt.Print("(")
// 	a.Print(expr.Left)
// 	fmt.Print(expr.Operator.Lexeme)
// 	a.Print(expr.Right)
// 	fmt.Print(")")
// 	return nil, nil
// }

// func (a AstPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
// 	fmt.Print("(")
// 	fmt.Print(expr.Operator.Lexeme)
// 	a.Print(expr.Right)
// 	fmt.Print(")")
// 	return nil, nil
// }

// func (a AstPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
// 	fmt.Print(expr.Value)
// 	return nil, nil
// }

// func (a AstPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
// 	fmt.Print("(")
// 	a.Print(expr.Expression)
// 	fmt.Print(")")
// 	return nil, nil
// }
