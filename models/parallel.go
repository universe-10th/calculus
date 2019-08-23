package models

import (
	"errors"
	"github.com/universe-10th/calculus/expressions"
)


// This computation must be considered a parallel computation
// of several different model flows, atomically. The computation
// input must satisfy the first element's input requirements.
// The computations must share input variables at most, but must
// not share any output variable.
type ParallelComputable struct {
	cachedVars
	elements []Computable
}


var ErrNoParallelComputablesGiven = errors.New("no parallel computables given")
var ErrSomeModelOutputsBelongToOtherModelInputs = errors.New("some output variables in one computable belong to input variables of other computable")
var ErrOutputVariableMergedTwice = errors.New("an output variable from one input computable clashes with an output variable in other computable(s)")

// Creates a parallel computable given all the elements.
func NewParallelComputable(elements ...Computable) (*ParallelComputable, error) {
	if len(elements) == 0 {
		return nil, ErrNoParallelComputablesGiven
	}
	cachedVars := cachedVars{}
	for _, element := range elements {
		if element == nil {
			return nil, ErrComputableIsNil
		}
		for outputVar, _ := range element.Output() {
			if cachedVars.DefinesOutput(outputVar) {
				return nil, ErrOutputVariableMergedTwice
			}
		}
		elementCachedVars := element.cachedVars()
		cachedVars.merge(elementCachedVars)
	}
	if !cachedVars.hasConsistentDomain() {
		return nil, ErrSomeModelOutputsBelongToOtherModelInputs
	}
	return &ParallelComputable{
		cachedVars,
		elements,
	}, nil
}


// Computes the models in parallel, and returns all the results.
func (parallelComputable *ParallelComputable) Compute(arguments expressions.Arguments) (expressions.Arguments, error) {
	result := expressions.Arguments{}
	for _, element := range parallelComputable.elements {
		if subResult, err := element.Compute(arguments); err != nil {
			return nil, err
		} else {
			for key, value := range subResult {
				result[key] = value
			}
		}
	}
	return result, nil
}