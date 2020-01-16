package models

import (
	"github.com/universe-10th/calculus/v2/expressions"
	"github.com/universe-10th/calculus/v2/models/errors"
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


// Creates a parallel model flow given all the elements.
func NewParallelModelFlow(elements ...ModelFlow) (*ParallelModelFlow, error) {
	if len(elements) == 0 {
		return nil, errors.ErrNoParallelModelFlowsGiven
	}
	cachedVars := cachedVars{
		expressions.Variables{},
		expressions.Variables{},
	}
	for _, element := range elements {
		if element == nil {
			return nil, errors.ErrModelFlowIsNil
		}
		for outputVar, _ := range element.Output() {
			if cachedVars.DefinesOutput(outputVar) {
				return nil, errors.ErrOutputVariableMergedTwice
			}
		}
		elementCachedVars := element.CachedVars()
		cachedVars.Merge(elementCachedVars)
	}
	if !cachedVars.HasConsistentDomain() {
		return nil, errors.ErrSomeModelOutputsBelongToOtherModelInputs
	}
	return &ParallelModelFlow{
		cachedVars,
		elements,
	}, nil
}


// Computes the models in parallel, and returns all the results.
func (parallelModelFlow *ParallelModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	result := expressions.Arguments{}
	for _, element := range parallelModelFlow.elements {
		if subResult, err := element.Evaluate(arguments); err != nil {
			return nil, err
		} else {
			for key, value := range subResult {
				result[key] = value
			}
		}
	}
	return result, nil
}