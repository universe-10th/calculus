// Expressions package contains elements that are understood as mathematical expressions.
// Arithmetic, power/logarithm, and trigonometric operators/functions are supported. They
// are composed as nodes in trees.
package expressions

import (
	"fmt"
	"github.com/universe-10th/calculus/sets"
	"math/big"
	"github.com/universe-10th/calculus/errors"
)


var zero = big.NewFloat(0)
var one  = big.NewFloat(1)


// Constant stands for a simple number node, e.g. 3, 4/5, or 3.1415926535.
type Constant struct {
	number sets.Number
}


// Variable stands for a variable node, like X, Y, Z, W or whatever you like.
type Variable struct {
	name string
}


// Variables stands for a map of variable => bool. Intended to work like a set.
type Variables map[Variable]bool


// Arguments are, actually, an interpretation. They are a map of variables to their values.
type Arguments map[Variable]sets.Number


// Expression is the interface behind each node.
type Expression interface {
	// CollectVariables enumerates all the variables involved recursively in the node into the given set.
	CollectVariables(Variables)
	// Curry generates a new expression by evaluating/freezing the variables given some arguments.
	// E.g. Add(X, Y).Curry(Arguments{Y: 3}) would yield a new expression: Add(X, Num(3)).
	// In the bound cases, passing no relevant argument will return a cloned expression, and
	// passing all the relevant arguments will return a constant, simplified, expression.
	Curry(arguments Arguments) (Expression, error)
	// Evaluate tries to evaluate the expression, recursively, given some arguments.
	Evaluate(arguments Arguments) (sets.Number, error)
	// Derivative generates a new expression being the derivative of the current one.
	Derivative(wrt Variable) (Expression, error)
	// IsConstant tells whether this expression is constant with respect to a variable.
	// It asks recursively and fails if it finds at least one node being == variable.
	// E.g. (X + 3) is not a constant expression per se, but it is with respect to Y.
	IsConstant(wrt Variable) bool
	// Simplify generates a simplified expression of this one.
	// Currently, this only converts to constants every expression involving only
	// constants, and does not any heuristic regarding (x+1)/(x+1) or actual
	// simplification processes.
	Simplify() (Expression, error)
	fmt.Stringer
}


// SelfContained tells whether an expression is self-contained in representation.
// For example: negations are self-contained, as well as numbers or functions like
// Ln(x), Log(x), Sin(x), Cos(x), Tan(x). In human words: no extra parentheses would
// be ever needed around them. Never.
type SelfContained interface {
	Expression
	// A marker function distinguishing this interface.
	IsSelfContained() bool
}


// CollectVariables adds the current variable to the set.
func (variable Variable) CollectVariables(variables Variables) {
	variables[variable] = true
}


// IsConstant checks whether the current variable is the given one.
func (variable Variable) IsConstant(wrt Variable) bool {
	return variable != wrt
}


// Simplify returns the same variable. There is nothing to simplify here.
func (variable Variable) Simplify() (Expression, error) {
	return variable, nil
}


// Curry tries to partially evaluate the expression.
// In the case of variables, it will return the same variable (if absent from
// the arguments) or return a constant with the argument value (if present in
// the arguments).
func (variable Variable) Curry(args Arguments) (Expression, error) {
	if value, ok := args[variable]; ok {
		return Num(value), nil
	} else {
		return variable, nil
	}
}


// Evaluate just drags the appropriate value from the given arguments.
// It returns an error if a value for the current variable is not present.
func (variable Variable) Evaluate(args Arguments) (sets.Number, error) {
	if value, ok := args[variable]; !ok {
		return nil, errors.UndefinedValue
	} else {
		return value, nil
	}
}


// Derivative returns a 0 or 1 constant expression.
// The 0 will be in the case the variables are not the same.
func (variable Variable) Derivative(wrt Variable) (Expression, error) {
	if variable == wrt {
		return Constant{one}, nil
	} else {
		return Constant{zero}, nil
	}
}


// String returns the variable name.
func (variable Variable) String() string {
	return variable.name
}


func (variable Variable) IsSelfContained() bool {
	return true
}


// CollectVariables does nothing in this type.
func (constant Constant) CollectVariables(Variables) {
	// Does nothing
}


// IsConstant is always true regarding constant expressions.
func (constant Constant) IsConstant(wrt Variable) bool {
	return true
}


// Simplify returns the same constant. There is nothing to simplify here.
func (constant Constant) Simplify() (Expression, error) {
	return constant, nil
}


// Curry tries to partially evaluate the expression.
// This has no special meaning in constants: the same constant will be returned.
func (constant Constant) Curry(args Arguments) (Expression, error) {
	return constant, nil
}


// Evaluate ignores any argument and returns the constant's value.
func (constant Constant) Evaluate(args Arguments) (sets.Number, error) {
	return constant.number, nil
}


// Derivative is a 0 expression for any constant.
func (constant Constant) Derivative(wrt Variable) (Expression, error) {
	return Constant{zero}, nil
}


// String returns the constant's string representation, according to its inner type.
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


// Number returns the constant's inner value.
// Usually needed in simplification process.
func (constant Constant) Number() sets.Number {
	return constant.number
}


// Var constructs a new Variable node.
func Var(name string) Variable {
	return Variable{name}
}


// Num constructs a new Constant node.
func Num(n interface{}) Constant {
	wrapped, _ := sets.Wrap(sets.Clone(n))
	return Constant{wrapped}
}


// MakeExpression takes an arbitrary value and creates an expression out of it.
// If the value was already an expression, it returns it as-is.
// Otherwise, it makes a constant expression out of it.
func MakeExpression(n interface{}) Expression {
	if exp, ok := n.(Expression); ok {
		return exp
	} else {
		return Num(n.(sets.Number))
	}
}


// Wrap wraps all the values as arguments, and returns the same map.
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
