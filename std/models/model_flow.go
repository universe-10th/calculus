package models

import (
	"github.com/universe-10th/calculus/v2/std/expressions"
)


// This interface is intended for flows and flow chains.
type ModelFlow interface {
	Evaluate(arguments expressions.Arguments) (expressions.Arguments, error)
	CachedVars() cachedVars
	Input() expressions.Variables
	Output() expressions.Variables
	DefinesInput(variable expressions.Variable) bool
	DefinesOutput(variable expressions.Variable) bool
}
