package models

import (
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/models/errors"
)


// This computation must be considered a serial computation of
// several different model flows, one by one. The computation
// input must satisfy the first element's input requirements.
// The computation output of the nth step must satisfy the
// input requirements of the (n+1)th element. The output of
// the last element will become the output of the whole
// computation.
type SerialModelFlow struct {
	cachedVars
	elements []ModelFlow
}


// Creates a chained model flow given all the elements.
func NewSerialModelFlow(elements ...ModelFlow) (*SerialModelFlow, error) {
	length := len(elements)
	if length == 0 {
		return nil, errors.ErrNoSerialModelFlowsGiven
	}
	if elements[0] == nil {
		return nil, errors.ErrModelFlowIsNil
	}
	for index, element := range elements[0:length - 1] {
		nextElement := elements[index + 1]
		if nextElement == nil {
			return nil, errors.ErrModelFlowIsNil
		} else {
			currentOutput := element.Output()
			nextInput := nextElement.Input()
			for inputVar, _ := range nextInput {
				if _, ok := currentOutput[inputVar]; !ok {
					return nil, errors.ErrOutputToInputChainNotSatisfied
				}
			}
		}
	}
	return &SerialModelFlow{
		cachedVars{
			input: elements[0].Input(),
			output: elements[length - 1].Output(),
		},
		elements,
	}, nil
}


// Computes the models in serial order, and returns the results of the last step.
func (serialModelFlow *SerialModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	result := arguments
	var err error
	for _, element := range serialModelFlow.elements {
		if result, err = element.Evaluate(result); err != nil {
			return nil, err
		}
	}
	return result, nil
}
