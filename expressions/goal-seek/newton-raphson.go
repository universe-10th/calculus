package goal_seek

import (
	"fmt"
	"math/big"
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/goals"
	"errors"
)

// Newton-Raphson algorithms take a notion of tolerance
// (differential error that allows us to consider to values
// as equal), an initial guess, and two iteration settings.
type NRGoalSeekingAlgorithm struct {
	initialGuess, epsilon *big.Float
	maxIterations, maxArgCorrectionsPerIteration uint32
	inverted expressions.Variable
}


// Executes the well-known Newton-Raphson method.
func (nrGoalSeekingAlgorithm NRGoalSeekingAlgorithm) FindRoot(goalBasedExpression expressions.Expression) (sets.Number, error) {
	inverted := nrGoalSeekingAlgorithm.inverted
	if goalBasedDerivativeExpression, err := goalBasedExpression.Derivative(inverted); err != nil {
		return nil, err
	} else {
		fmt.Println("Goal:", goalBasedExpression, "derivative:", goalBasedDerivativeExpression)
		goalBasedExpressionFunction := func(current *big.Float) (sets.Number, error) {
			return goalBasedExpression.Evaluate(expressions.Arguments{inverted: current})
		}
		goalBasedDerivativeExpressionFunction := func(current *big.Float) (sets.Number, error) {
			return goalBasedDerivativeExpression.Evaluate(expressions.Arguments{inverted: current})
		}
		return goals.NewtonRaphson(
			goalBasedExpressionFunction, goalBasedDerivativeExpressionFunction, nrGoalSeekingAlgorithm.initialGuess,
			nrGoalSeekingAlgorithm.epsilon, nrGoalSeekingAlgorithm.maxIterations, nrGoalSeekingAlgorithm.maxArgCorrectionsPerIteration,
		)
	}
}


var ErrNRGoalSeekingAlgorithmBadParams = errors.New("epsilon, max iterations, and max corrections must be all positive")


// Function type that provides the relevant arguments for the
// instantiation of NRGoalSeekingAlgorithms.
type NRGoalSeekAlgorithmArgsProvider func(arguments expressions.Arguments) (
	initialGuess, epsilon *big.Float,
	maxIterations, maxArgCorrectionsPerIteration uint32,
)


// Builds a new Newton-Raphson algorithm.
func NRGoalSeek(
	goal, target expressions.Expression, invertedVariable expressions.Variable,
	argsProvider NRGoalSeekAlgorithmArgsProvider,
) expressions.Expression {
	return expressions.GoalSeek(goal, target, invertedVariable, func(
		arguments expressions.Arguments, inverted expressions.Variable, fullDomain expressions.Variables,
	) (expressions.GoalSeekingAlgorithm, error) {
		initial, epsilon, maxIterations, maxCorrections := argsProvider(arguments)
		if epsilon.Sign() < 1 || maxIterations == 0 || maxCorrections == 0 {
			return nil, ErrNRGoalSeekingAlgorithmBadParams
		} else {
			return NRGoalSeekingAlgorithm{
				initial, epsilon,
				maxIterations, maxCorrections,
				inverted,
			}, nil
		}
	})
}
