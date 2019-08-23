package models

import (
	"github.com/universe-10th/calculus/expressions"
	"errors"
)


// This interface is intended for flows and flow chains.
type Computable interface {
	Compute(arguments expressions.Arguments) (expressions.Arguments, error)
	cachedVars() cachedVars
	Input() expressions.Variables
	Output() expressions.Variables
}


var ErrComputableIsNil = errors.New("a given computable is nil")
