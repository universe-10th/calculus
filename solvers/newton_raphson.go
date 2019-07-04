package solvers

import (
	"math/big"
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/solvers/support"
	"github.com/universe-10th/calculus/sets"
	"errors"
	"math/rand"
	"time"
)


var ErrMaxNewtonRaphsonArgCorrections = errors.New("could not correct the argument to avoid a zero derivative")


func nextNewtonRaphsonStep(expression, derivative expressions.Expression, variable expressions.Variable,
	                       currentRes, currentDerRes sets.Number, currentImg, currentDerImg, quot *big.Float,
	                       arg, epsilon, delta, random *big.Float, maxArgCorrections uint32) (*big.Float, error) {
	var err error
	var current = arg
	var recalculate = false
	// This will proceed as follows: it will try calculating the derivative a lof of times until the result is
	// not zero. Otherwise, successive steps will modify the current x by a very small random step. When the
	// derivative is found, the function value is recalculated (if failed at least once), the quotient is computed,
	// and the new value for the argument is corrected.
	for tryOut := uint32(0); tryOut < maxArgCorrections; tryOut++ {
		if currentDerRes, err = derivative.Evaluate(expressions.Arguments{variable: current}); err != nil {
			return nil, err
		}
		if currentDerImg, err = support.ForceFloat(currentDerRes); err != nil {
			return nil, err
		}
		if currentDerImg.Sign() != 0 {
			// evaluate f(currentArg) into currentImg if told to recalculate.
			if recalculate {
				if currentRes, err = expression.Evaluate(expressions.Arguments{variable: current}); err != nil {
					return nil, err
				}
				if currentImg, err = support.ForceFloat(currentRes); err != nil {
					return nil, err
				}
			}
			// Compute x = x - f(x)/f'(x).
			current.Sub(current, quot.Quo(currentImg, currentDerImg))
			return current, nil
		} else {
			recalculate = true
			random.SetFloat64(rand.Float64())
			current.Add(current, delta.Mul(epsilon, random))
		}
	}
	return nil, ErrMaxNewtonRaphsonArgCorrections
}


// NewtonRaphson tries a Newton-Raphson iteration to find the nearest approximate root of the given function.
// If you want to find x so that f(x) == target, just define g(x) = f(x) - target and you will be on your way.
// While the usual newton-raphson fails when f'(x) (or g'(x)) is zero, we replace such value by the same epsilon
// in the tolerance. The number of max iterations, and the number of max corrections of the argument when the
// derivative is zero must be provided as well.
func NewtonRaphson(expression expressions.Expression, initialGuess, epsilon *big.Float,
	               maxIterations, maxArgCorrectionsPerIteration uint32) (result *big.Float, exception error) {
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
	var delta = big.NewFloat(0)
	var random = big.NewFloat(0)
	rand.Seed(time.Now().UTC().UnixNano())
	if expression, err = expression.Simplify(); err != nil {
		return nil, err
	} else if variable, err = support.GetTheOnlyVariable(expression); err != nil {
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
			// Otherwise, process a newton-raphson step.
			if currentArg, err = nextNewtonRaphsonStep(
				expression, derivative, variable, currentRes, currentDerRes, currentImg, currentDerImg, quot,
				currentArg, epsilon, delta, random, maxArgCorrectionsPerIteration,
			); err != nil {
				return nil, err
			}
		}
		return nil, support.ErrIterationsExhausted
	}
}
