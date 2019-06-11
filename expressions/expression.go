package diff

import (
	"fmt"
	"github.com/universe-10th-calculus/sets"
	"math/big"
	"github.com/universe-10th-calculus/errors"
)


var zero = big.NewFloat(0)
var one  = big.NewFloat(1)


type Constant struct {
	number sets.Number
}


type Variable struct {
	name string
}
type Variables map[Variable]bool


type Arguments map[Variable]sets.Number


type Expression interface {
	CollectVariables(Variables)
	Evaluate(arguments Arguments) (sets.Number, error)
	Derivative(wrt Variable) (Expression, error)
	fmt.Stringer
}


func (variable Variable) CollectVariables(variables Variables) {
	variables[variable] = true
}


func (variable Variable) Evaluate(args Arguments) (sets.Number, error) {
	if value, ok := args[variable]; !ok {
		return nil, errors.UndefinedValue
	} else {
		return value, nil
	}
}


func (variable Variable) Derivative(wrt Variable) (Expression, error) {
	if variable == wrt {
		return Constant{one}, nil
	} else {
		return Constant{zero}, nil
	}
}


func (variable Variable) String() string {
	return variable.name
}


func (constant Constant) CollectVariables(Variables) {
	// Does nothing
}


func (constant Constant) Evaluate(args Arguments) (sets.Number, error) {
	return constant.number, nil
}


func (constant Constant) Derivative(wrt Variable) (Expression, error) {
	return Constant{zero}, nil
}


func (constant Constant) String() string {
	return fmt.Sprintf("%s", constant.number)
}


func Var(name string) Variable {
	return Variable{name}
}


func Num(n sets.Number) Constant {
	return Constant{n}
}


var W = Variable{"W"}
var X = Variable{"X"}
var Y = Variable{"Y"}
var Z = Variable{"Z"}
