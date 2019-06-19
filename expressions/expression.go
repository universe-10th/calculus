package expressions

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
	IsConstant(wrt Variable) bool
	Simplify() (Expression, error)
	fmt.Stringer
}


type SelfContained interface {
	Expression
	IsSelfContained() bool
}


func (variable Variable) CollectVariables(variables Variables) {
	variables[variable] = true
}


func (variable Variable) IsConstant(wrt Variable) bool {
	return variable != wrt
}


func (variable Variable) Simplify() (Expression, error) {
	return variable, nil
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


func (variable Variable) IsSelfContained() bool {
	return true
}


func (constant Constant) CollectVariables(Variables) {
	// Does nothing
}


func (constant Constant) IsConstant(wrt Variable) bool {
	return true
}


func (constant Constant) Simplify() (Expression, error) {
	return constant, nil
}


func (constant Constant) Evaluate(args Arguments) (sets.Number, error) {
	return constant.number, nil
}


func (constant Constant) Derivative(wrt Variable) (Expression, error) {
	return Constant{zero}, nil
}


func (constant Constant) String() string {
	switch c := constant.number.(type) {
	case *big.Float:
		return c.Text('f', -1)
	case *big.Int:
		return fmt.Sprintf("%d", c)
	case *big.Rat:
		return fmt.Sprintf("%s", c)
	}
	return "<?>"
}


func (constant Constant) IsSelfContained() bool {
	return true
}


func (constant Constant) Number() sets.Number {
	return constant.number
}


func Var(name string) Variable {
	return Variable{name}
}


func Num(n sets.Number) Constant {
	wrapped, _ := sets.Wrap(sets.Clone(n))
	return Constant{wrapped}
}


func (arguments Arguments) Wrap() Arguments {
	for key, value := range arguments {
		wrapped, _ := sets.Wrap(value)
		arguments[key] = wrapped
	}
	return arguments
}


var W = Variable{"W"}
var X = Variable{"X"}
var Y = Variable{"Y"}
var Z = Variable{"Z"}
