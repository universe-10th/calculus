package models

import (
	"github.com/universe-10th/calculus/expressions"
	"errors"
)


type ModelFlowExpressions map[expressions.Variable]expressions.Expression


// A model flow involves one or many simultaneous operations over some input
// variables, and to some other output variables (outputs must not be also
// present as input variables.
//
// Expressions will run in an atomic-like fashion (resolution must be
// considered atomic for the sake of this system).
type ModelFlow struct {
	input       expressions.Variables
	output      expressions.Variables
	expressions ModelFlowExpressions
}


var ErrNoFlowExpressionsGiven = errors.New("no flow expressions given (at least one is mandatory)")
var ErrExpressionIsNil = errors.New("one or more expressions for this flow are nil")
var ErrOutputVariableInsideFlowExpressions = errors.New("model flow's output variable is inside flow's expression")
var ErrInsufficientArguments = errors.New("insufficient arguments for the model flow")


// Creates the flow given all the expression.
func NewFlow(flowExpressions ModelFlowExpressions) (*ModelFlow, error) {
	if flowExpressions == nil {
		return nil, ErrNoFlowExpressionsGiven
	}

	inputVars := expressions.Variables{}
	outputVars := expressions.Variables{}

	for outputVar, expression := range flowExpressions {
		if expression == nil {
			return nil, ErrExpressionIsNil
		} else {
			outputVars[outputVar] = true
			expression.CollectVariables(inputVars)
		}
	}

	for outputVar, _ := range outputVars {
		if _, ok := inputVars[outputVar]; ok {
			return nil, ErrOutputVariableInsideFlowExpressions
		}
	}

	return &ModelFlow{
		input: inputVars,
		expressions: flowExpressions,
		output: outputVars,
	}, nil
}


// Given the arguments, tries to resolve all the involved expressions.
//
// If at least one of the required arguments is not present, the flow will fail.
// It at least one of the flow expressions returns an error, the whole flow will fail.
func (flow *ModelFlow) Compute(arguments expressions.Arguments) (expressions.Arguments, error) {
	for key, _ := range flow.input {
		if _, ok := arguments[key]; !ok {
			// Required input is not present.
			return nil, ErrInsufficientArguments
		}
	}

	result := expressions.Arguments{}
	for outputVar, expression := range flow.expressions {
		if value, err := expression.Evaluate(arguments); err != nil {
			return nil, err
		} else {
			result[outputVar] = value
		}
	}
	return result, nil
}


// Returns a copy of the input vars set.
func (flow *ModelFlow) Input() expressions.Variables {
	input := expressions.Variables{}
	for inputVar, _ := range flow.input {
		input[inputVar] = true
	}
	return input
}


// Returns a copy of the output vars set.
func (flow *ModelFlow) Output() expressions.Variables {
	output := expressions.Variables{}
	for outputVar, _ := range flow.input {
		output[outputVar] = true
	}
	return output
}
