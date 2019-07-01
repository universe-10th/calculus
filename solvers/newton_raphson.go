package solvers

import (
	"math/big"
	"github.com/universe-10th-calculus/expressions"
	"github.com/universe-10th-calculus/solvers/support"
	"github.com/universe-10th-calculus/sets"
)


// NewtonRaphson tries a Newton-Raphson iteration to find the nearest approximate root of the given function.
// If you want to find x so that f(x) == target, just define g(x) = f(x) - target and you will be on your way.
// While the usual newton-raphson fails when f'(x) (or g'(x)) is zero, we replace such value by the same epsilon
// in the tolerance.
func NewtonRaphson(expression expressions.Expression, initialGuess, epsilon *big.Float,
	               maxIterations uint32) (result *big.Float, exception error) {
	var derivative expressions.Expression
	var variable expressions.Variable
	var err error
	var currentArg *big.Float
	var currentRes sets.Number
	var currentDerRes sets.Number
	var currentImg *big.Float
	var currentDerImg *big.Float
	var zero = big.NewFloat(0)
	var diff = big.NewFloat(0)
	var dist = big.NewFloat(0)
	var quot = big.NewFloat(0)
	if variable, err = support.GetTheOnlyVariable(expression); err != nil {
		return nil, err
	} else if derivative, err = expression.Derivative(variable); err != nil {
		return nil, err
	} else {
		currentArg = initialGuess
		if maxIterations == 0 {
			maxIterations = 100
		}
		var iteration uint32
		for iteration = 0; iteration < maxIterations; iteration++ {
			// evaluate f(currentArg) into currentImg.
			if currentRes, err = expression.Evaluate(expressions.Arguments{variable: currentArg}); err != nil {
				return nil, err
			}
			if currentImg, err = support.ForceFloat(currentRes); err != nil {
				return nil, err
			}
			// If it is close to zero, then return.
			if support.CloseTo(currentImg, zero, diff, dist, epsilon) {
				return currentArg, nil
			}
			// Calculate (and perhaps correct) the derivative at that value.
			if currentDerRes, err = derivative.Evaluate(expressions.Arguments{variable: currentArg}); err != nil {
				return nil, err
			}
			if currentDerImg, err = support.ForceFloat(currentDerRes); err != nil {
				return nil, err
			}
			// Correct the derivative if zero. Use epsilon instead.
			// TODO This is wrong. If the derivative is zero, one has to add the epsilon to X instead,
			// TODO   and then attempt again a calculation of f(x) and g(x), repeating this process until
			// TODO   a derivative != 0 is found.
			if currentDerImg.Sign() == 0 {
				currentDerImg = epsilon
			}
			// Compute nextArg = currentArg - func/der.
			currentArg.Sub(currentArg, quot.Quo(currentImg, currentDerImg))
		}
		return nil, support.ErrIterationsExhausted
	}
}
