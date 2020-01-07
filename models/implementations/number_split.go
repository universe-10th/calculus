package implementations

import (
	"github.com/universe-10th/calculus/models"
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/ops"
	"github.com/universe-10th/calculus/errors"
)


type NumberSplitModelFlow struct {
	models.CustomModelFlow
	input expressions.Variable
	intOutput expressions.Variable
	fracOutput expressions.Variable
}


func NewNumberSplitModelFlow(input, intOutput, fracOutput expressions.Variable) (*NumberSplitModelFlow, error) {
	if custom, err := models.NewCustomModelFlow(
		expressions.Variables{input: true},
		expressions.Variables{intOutput: true, fracOutput: true},
	); err != nil {
		return nil, err
	} else {
		return &NumberSplitModelFlow{*custom, input, intOutput, fracOutput}, nil
	}
}


func (numberSplitModelFlow *NumberSplitModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	if intPart, fracPart := ops.Split(arguments[numberSplitModelFlow.input]); intPart == nil {
		return nil, errors.InfiniteCannotBeRounded
	} else {
		return expressions.Arguments{
			numberSplitModelFlow.intOutput: intPart,
			numberSplitModelFlow.fracOutput: fracPart,
		}, nil
	}
}