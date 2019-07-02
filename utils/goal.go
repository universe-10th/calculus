package utils

import (
	"github.com/universe-10th-calculus/expressions"
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
)


// Takes an existing expression -say- f(...) and generates its g(...) = f(...) - goal.
func GoalBasedExpression(expression expressions.Expression, goal interface{}) expressions.Expression {
	wrapped, _ := sets.Wrap(goal)
	return expressions.Add(expression, expressions.Num(ops.Neg(wrapped)))
}