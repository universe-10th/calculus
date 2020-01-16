package models

import (
	"github.com/universe-10th/calculus/v2/expressions"
	"github.com/universe-10th/calculus/v2/models/errors"
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
func NewSingleOutputModelFlow(output expressions.Variable, expression expressions.Expression) (*SingleOutputModelFlow, error) {
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

	if !cachedVars.HasConsistentDomain() {
		return nil, errors.ErrOutputVariableInsideFlowExpression
	}

	return &SingleOutputModelFlow{
		cachedVars: cachedVars,
		expression: expression,
		output: output,
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


type SingleOutputModelFlowSpec struct {
	Variable expressions.Variable
	Expression expressions.Expression
}


// Creates an array of single output model flows.
// This is a convenience function to wrap several
// calls to NewSingleOutputModelFlow by making a
// short-circuit traversal until an error occurs,
// and either returning all the created flows or
// returning the raised error.
func NewSingleOutputModelFlowSet(args ...SingleOutputModelFlowSpec) ([]ModelFlow, error) {
	result := make([]ModelFlow, len(args))
	var err error
	for index, value := range args {
		if result[index], err = NewSingleOutputModelFlow(value.Variable, value.Expression); err != nil {
			return nil, err
		}
	}
	return result, nil
}