package models

import (
	"github.com/universe-10th/calculus/expressions"
	"errors"
)


// This interface is intended for flows and flow chains.
type ModelFlow interface {
	Evaluate(arguments expressions.Arguments) (expressions.Arguments, error)
	cachedVars() cachedVars
	Input() expressions.Variables
	Output() expressions.Variables
	DefinesInput(variable expressions.Variable) bool
	DefinesOutput(variable expressions.Variable) bool
}


var ErrModelFlowIsNil = errors.New("a given model flow is nil")
