package diff

import (
	"fmt"
	"github.com/universe-10th-calculus/sets"
	"math/big"
)


var zero = big.NewFloat(0)
var one  = big.NewFloat(1)


type Constant struct {
	number sets.Number
}


type Variable struct {
	name string
}


type Arguments map[Variable]sets.Number


type Expression interface {
	Evaluate(arguments Arguments) (sets.Number, error)
	Derivative(wrt Variable) Expression
	fmt.Stringer
}


func (variable Variable) Evaluate(args Arguments) (sets.Number, error) {
	return args[variable], nil
}


func (variable Variable) Derivative(wrt Variable) Expression {
	if variable == wrt {
		return Constant{one}
	} else {
		return Constant{zero}
	}
}


func (variable Variable) String() string {
	return variable.name
}


func (constant Constant) Evaluate(args Arguments) (sets.Number, error) {
	return constant.number, nil
}


func (constant Constant) Derivative(wrt Variable) Expression {
	return Constant{zero}
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
