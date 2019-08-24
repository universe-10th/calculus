package models

import (
	"errors"
	"github.com/universe-10th/calculus/expressions"
)


// This model flow must be considered a parallel computation
// of several different model flows, atomically. The computation
// input must satisfy the first element's input requirements.
// The computations must share input variables at most, but must
// not share any output variable.
type ParallelModelFlow struct {
	cachedVars
	elements []ModelFlow
}


var ErrNoParallelModelFlowsGiven = errors.New("no parallel model flows given")
var ErrSomeModelOutputsBelongToOtherModelInputs = errors.New("some output variables in one model flow belong to input variables of other model flow(s)")
var ErrOutputVariableMergedTwice = errors.New("an output variable from one input model flow clashes with an output variable in other model flow(s)")

// Creates a parallel model flow given all the elements.
func NewParallelModelFlow(elements ...ModelFlow) (*ParallelModelFlow, error) {
	if len(elements) == 0 {
		return nil, ErrNoParallelModelFlowsGiven
	}
	cachedVars := cachedVars{}
	for _, element := range elements {
		if element == nil {
			return nil, ErrModelFlowIsNil
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
	return &ParallelModelFlow{
		cachedVars,
		elements,
	}, nil
}


// Computes the models in parallel, and returns all the results.
func (parallelModelFlow *ParallelModelFlow) Compute(arguments expressions.Arguments) (expressions.Arguments, error) {
	result := expressions.Arguments{}
	for _, element := range parallelModelFlow.elements {
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