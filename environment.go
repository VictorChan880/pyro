package main

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func (e Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func NewEnvironment() *Environment {
	return &Environment{
		Enclosing: nil,
		Values:    make(map[string]interface{}),
	}
}
func NewEnclosedEnvironment(enclosed *Environment) *Environment {
	return &Environment{
		Enclosing: enclosed,
		Values:    make(map[string]interface{}),
	}
}

func (e Environment) get(name Token) (interface{}, error) {
	if _, exists := e.Values[name.Lexeme]; exists {
		return e.Values[name.Lexeme], nil
	}
	if e.Enclosing != nil {
		return e.Enclosing.get(name)
	}
	err := NewRunTimeError(name, " Undefined variable "+name.Lexeme)
	report(err.Err)
	return nil, err
}

func (e Environment) assign(name Token, value interface{}) error {
	if _, exists := e.Values[name.Lexeme]; exists {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}

	err := NewRunTimeError(name, " Undefined variable "+name.Lexeme)
	report(err.Err)

	return err
}
