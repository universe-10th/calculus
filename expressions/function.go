package expressions

import (
	"strings"
	"fmt"
)


type Function interface {
	fmt.Stringer
	StandardName() string
	Arguments()    []Expression
}


type FunctionExpr struct {
	standardName string
}


func (function FunctionExpr) IsSelfContained() bool {
	return true
}


func (function FunctionExpr) StandardName() string {
	return function.standardName
}


func (function FunctionExpr) Arguments() []Expression {
	panic("not implemented")
}


func (function FunctionExpr) String() string {
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