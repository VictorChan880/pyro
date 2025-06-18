package main

type PyroFunction struct {
	Declaration Function
}

func NewPyroFunction (declaration Function) PyroFunction{
	return PyroFunction{
		Declaration: declaration,
	}
}

func (pf PyroFunction) Arity() int {
	return len(pf.Declaration.Params)
}

func (pf PyroFunction) toString() string {
	return "<fn " + pf.Declaration.Name.Lexeme + ">"
}


func (pf PyroFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error){
	environment := NewEnclosedEnvironment(interpreter.Environment)

	for i := 0; i < len(pf.Declaration.Params); i++ {
		environment.define(pf.Declaration.Params[i].Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(pf.Declaration.Body, environment)
	return nil, err
}