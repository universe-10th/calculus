package implementations

import (
	"github.com/universe-10th/calculus/models"
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/ops"
)


type NumberSplitModelFlow struct {
	models.CustomModelFlow
	input expressions.Variable
	intOutput expressions.Variable
	fracOutput expressions.Variable
}


func NewNumberSplitModelFlow(input, intOutput, fracOutput expressions.Variable) NumberSplitModelFlow {
	return NumberSplitModelFlow{
		models.NewCustomModelFlow(
			expressions.Variables{input: true},
			expressions.Variables{intOutput: true, fracOutput: true},
		),
		input, intOutput, fracOutput,
	}
}


func (numberSplitModelFlow NumberSplitModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	intPart, fracPart := ops.Split(arguments[numberSplitModelFlow.input])
	return expressions.Arguments{
		numberSplitModelFlow.intOutput: intPart,
		numberSplitModelFlow.fracOutput: fracPart,
	}, nil
}