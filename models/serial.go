package models

import (
	"errors"
	"github.com/universe-10th/calculus/expressions"
)


// This computation must be considered a serial computation of
// several different model flows, one by one. The computation
// input must satisfy the first element's input requirements.
// The computation output of the nth step must satisfy the
// input requirements of the (n+1)th element. The output of
// the last element will become the output of the whole
// computation.
type SerialComputable struct {
	cachedVars
	elements []Computable
}


var ErrNoSerialComputablesGiven = errors.New("no serial computables given")
var ErrOutputToInputChainNotSatisfied = errors.New("there is at least one step where the output from it does not satisfy the input of the next")


// Creates a chained computable given all the elements.
func NewSerialComputable(elements ...Computable) (*SerialComputable, error) {
	length := len(elements)
	if length == 0 {
		return nil, ErrNoSerialComputablesGiven
	}
	if elements[0] == nil {
		return nil, ErrComputableIsNil
	}
	for index, element := range elements[0:length - 1] {
		nextElement := elements[index + 1]
		if nextElement == nil {
			return nil, ErrComputableIsNil
		} else {
			currentOutput := element.Output()
			nextInput := nextElement.Input()
			for inputVar, _ := range nextInput {
				if _, ok := currentOutput[inputVar]; !ok {
					return nil, ErrOutputToInputChainNotSatisfied
				}
			}
		}
	}
	return &SerialComputable{
		cachedVars{
			input: elements[0].Input(),
			output: elements[length - 1].Output(),
		},
		elements,
	}, nil
}


// Computes the models in serial order, and returns the results of the last step.
func (serialComputable *SerialComputable) Compute(arguments expressions.Arguments) (expressions.Arguments, error) {
	result := arguments
	var err error
	for _, element := range serialComputable.elements {
		if result, err = element.Compute(result); err != nil {
			return nil, err
		}
	}
	return result, nil
}
