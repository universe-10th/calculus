package expressions

import (
	"strings"
	"fmt"
)


// Function interface requires an implementation of standard name and arguments methods.
// Standard name is, e.g., "Ln" while arguments is an array of the inner involved expressions.
type Function interface {
	fmt.Stringer
	StandardName() string
	Arguments()    []Expression
}


// FunctionExpr is a base structure for other function nodes. It is also marked as self-contained,
// in the same way of variables and constants.
type FunctionExpr struct {
	standardName string
}


func (function FunctionExpr) IsSelfContained() bool {
	return true
}


// The function's standard name, e.g. Ln, Sin, Tan.
func (function FunctionExpr) StandardName() string {
	return function.standardName
}


func (function FunctionExpr) Arguments() []Expression {
	panic("not implemented")
}


// FunctionDisplay appropriately represents a function object as a human-readable string.
func FunctionDisplay(function Function) string {
	arguments := function.Arguments()
	placeholders := make([]string, len(arguments))
	for index := range placeholders {
		placeholders[index] = "%s"
	}
	fmtString := "%s(" + strings.Join(placeholders, ", ") + ")"
	fmtArgs := make([]interface{}, len(arguments) + 1)
	fmtArgs[0] = function.StandardName()
	for index, argument := range arguments {
		fmtArgs[index + 1] = argument
	}
	return fmt.Sprintf(fmtString, fmtArgs...)
}