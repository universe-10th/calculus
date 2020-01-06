package models

import (
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/models/errors"
)


// A model flow involving one single model expression over some input
// variables, and to some other output variable (output variable must
// not be also present among input variables.
type SingleOutputModelFlow struct {
	cachedVars
	output expressions.Variable
	expression expressions.Expression
}


// Creates the flow given all the expression.
func NewFlow(output expressions.Variable, expression expressions.Expression) (*SingleOutputModelFlow, error) {
	inputVars := expressions.Variables{}
	outputVars := expressions.Variables{output: true}

	if expression == nil {
		return nil, errors.ErrExpressionIsNil
	} else {
		expression.CollectVariables(inputVars)
	}

	cachedVars := cachedVars{
		input: inputVars,
		output: outputVars,
	}

	if !cachedVars.hasConsistentDomain() {
		return nil, errors.ErrOutputVariableInsideFlowExpression
	}

	return &SingleOutputModelFlow{
		cachedVars: cachedVars,
		expression: expression,
	}, nil
}


// Given the arguments, tries to resolve all the involved expressions.
//
// If at least one of the required arguments is not present, the flow will fail.
// It at least one of the flow expressions returns an error, the whole flow will fail.
func (flow *SingleOutputModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	for key, _ := range flow.input {
		if _, ok := arguments[key]; !ok {
			// Required input is not present.
			return nil, errors.ErrInsufficientArguments
		}
	}

	result := expressions.Arguments{}
	if value, err := flow.expression.Evaluate(arguments); err != nil {
		return nil, err
	} else {
		result[flow.output] = value
	}
	return result, nil
}
