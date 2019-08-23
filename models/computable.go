package models

import "github.com/universe-10th/calculus/expressions"


// This interface is intended for flows and flow chains.
type Computable interface {
	Compute(arguments expressions.Arguments) (expressions.Arguments, error)
	Input() expressions.Variables
	Output() expressions.Variables
}
