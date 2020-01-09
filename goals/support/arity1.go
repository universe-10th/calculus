package support

import (
	"github.com/universe-10th/calculus/expressions"
	"errors"
)


var ErrMultipleVariables = errors.New("the given function has not exactly one parameter")


// GetTheOnlyVariable takes an expression and, if it involves just a single variable, returns that variable.
// If the expression involves no variables or involves several variables, it returns an error.
func GetTheOnlyVariable(expression expressions.Expression) (expressions.Variable, error) {
	variables := make(expressions.Variables)
	expression.CollectVariables(variables)
	if len(variables) != 1 {
		return expressions.Variable{}, ErrMultipleVariables
	} else {
		// Return the only existing variable.
		key := expressions.Variable{}
		for key, _ = range variables {
			break
		}
		// This line is actually never reached.
		return key, nil
	}
}
