package implementations

import (
	"github.com/universe-10th/calculus/v2/big/models"
	merrors "github.com/universe-10th/calculus/v2/big/models/errors"
	"github.com/universe-10th/calculus/v2/big/expressions"
	"github.com/universe-10th/calculus/v2/big/ops"
	"github.com/universe-10th/calculus/v2/big/errors"
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
	if number, ok := arguments[numberSplitModelFlow.input]; !ok {
		return nil, merrors.ErrInsufficientArguments
	} else if intPart, fracPart := ops.Split(number); intPart == nil {
		return nil, errors.ErrInfiniteCannotBeRounded
	} else {
		return expressions.Arguments{
			numberSplitModelFlow.intOutput: intPart,
			numberSplitModelFlow.fracOutput: fracPart,
		}, nil
	}
}