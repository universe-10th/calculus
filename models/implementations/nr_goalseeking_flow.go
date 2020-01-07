package implementations

import (
	"math/big"
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/solvers"
	"github.com/universe-10th/calculus/models"
	"errors"
)


// Newton-Raphson algorithms take a notion of tolerance
// (differential error that allows us to consider to values
// as equal), an initial guess, and two iteration settings.
type NRGoalSeekingAlgorithm struct {
	initialGuess, epsilon *big.Float
	maxIterations, maxArgCorrectionsPerIteration uint32
}


// Executes the well-known Newton-Raphson method.
func (nrGoalSeekingAlgorithm NRGoalSeekingAlgorithm) FindRoot(goalBasedExpression expressions.Expression) (sets.Number, error) {
	return solvers.NewtonRaphson(
		goalBasedExpression, nrGoalSeekingAlgorithm.initialGuess,
		nrGoalSeekingAlgorithm.epsilon, nrGoalSeekingAlgorithm.maxIterations, nrGoalSeekingAlgorithm.maxArgCorrectionsPerIteration,
	)
}


var ErrNRGoalSeekingAlgorithmBadParams = errors.New("epsilon, max iterations, and max corrections must be all positive")


// Function type that provides the relevant arguments for the
// instantiation of NRGoalSeekingAlgorithms.
type NRGoalSeekAlgorithmArgsProvider func(arguments expressions.Arguments) (
	initialGuess, epsilon *big.Float,
	maxIterations, maxArgCorrectionsPerIteration uint32,
)


// Builds a new Newton-Raphson algorithm.
func NewNewtonRaphsonGoalSeekingFlow(
	expression expressions.Expression, coDomain, invertedVariable expressions.Variable,
	argsProvider NRGoalSeekAlgorithmArgsProvider,
) (*models.GoalSeekingModelFlow, error) {
	return models.NewGoalSeekingModelFlow(expression, coDomain, invertedVariable, func(
		arguments expressions.Arguments, coDomain, inverted expressions.Variable,
		fullDomain expressions.Variables,
	) (models.GoalSeekingAlgorithm, error) {
		initial, epsilon, maxIterations, maxCorrections := argsProvider(arguments)
		if epsilon.Sign() < 1 || maxIterations == 0 || maxCorrections == 0 {
			return nil, ErrNRGoalSeekingAlgorithmBadParams
		} else {
			return NRGoalSeekingAlgorithm{
				initial, epsilon, maxIterations, maxCorrections,
			}, nil
		}
	})
}