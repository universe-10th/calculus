package goals

import (
	"time"
	"errors"
	"math/big"
	"math/rand"
	"github.com/universe-10th/calculus/sets"
	diffUtils "github.com/universe-10th/calculus/core/utils/diff"
	goalErrors "github.com/universe-10th/calculus/core/goals/errors"
)


var ErrMaxNewtonRaphsonArgCorrections = errors.New("could not correct the argument to avoid a zero derivative")


func nextNewtonRaphsonStep(expression, derivative func(*big.Float) (sets.Number, error),
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
		if currentDerRes, err = derivative(current); err != nil {
			return nil, err
		}
		// This function may panic, but we're covered in the parent call.
		currentDerImg = sets.UpCastOneTo(currentDerRes, sets.R).(*big.Float)
		if currentDerImg.Sign() != 0 {
			// evaluate f(currentArg) into currentImg if told to recalculate.
			if recalculate {
				if currentRes, err = expression(current); err != nil {
					return nil, err
				}
				// This function may panic, but we're covered in the parent call.
				currentImg = sets.UpCastOneTo(currentRes, sets.R).(*big.Float)
			}
			// Compute x = x - f(x)/f'(x).
			current.Sub(current, quot.Quo(currentImg, currentDerImg))
			return current, nil
		} else {
			recalculate = true
			random.SetFloat64(rand.Float64() - 0.5)
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
func NewtonRaphson(expression, derivative func(*big.Float) (sets.Number, error), initialGuess, epsilon *big.Float,
	               maxIterations, maxArgCorrectionsPerIteration uint32) (result *big.Float, exception error) {
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
	currentArg = initialGuess
	if maxIterations == 0 {
		maxIterations = 100
	}
	var iteration uint32
	defer func() {
		if recovered := recover(); recovered != nil {
			result = nil
			exception = errors.New(recovered.(string))
		}
	}()
	for iteration = 0; iteration < maxIterations; iteration++ {
		// evaluate f(currentArg) into currentImg.
		if currentRes, err = expression(currentArg); err != nil {
			return nil, err
		}
		// This function may panic, but we're covered here.
		currentImg = sets.UpCastOneTo(currentRes, sets.R).(*big.Float)
		// If it is close to zero, then return.
		if diffUtils.CloseTo(currentImg, zero, diff, dist, epsilon) {
			return currentArg, nil
		}
		// Otherwise, process a newton-raphson step.
		if currentArg, err = nextNewtonRaphsonStep(
			expression, derivative, currentRes, currentDerRes, currentImg, currentDerImg, quot,
			currentArg, epsilon, delta, random, maxArgCorrectionsPerIteration,
		); err != nil {
			return nil, err
		}
	}
	return nil, goalErrors.ErrIterationsExhausted
}