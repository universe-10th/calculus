package goals

import (
	"time"
	"errors"
	"math/rand"
	"github.com/universe-10th/calculus/v2/std/sets"
	diffUtils "github.com/universe-10th/calculus/v2/std/core/support/diff"
	goalErrors "github.com/universe-10th/calculus/v2/std/core/goals/errors"
)


var ErrMaxNewtonRaphsonArgCorrections = errors.New("could not correct the argument to avoid a zero derivative")


func nextNewtonRaphsonStep(expression, derivative func(float64) (sets.Number, error),
	                       currentRes, currentDerRes sets.Number, currentImg, currentDerImg, quot float64,
	                       arg, epsilon, delta, random float64, maxArgCorrections uint32) (float64, error) {
	var err error
	var current = arg
	var recalculate = false
	// This will proceed as follows: it will try calculating the derivative a lof of times until the result is
	// not zero. Otherwise, successive steps will modify the current x by a very small random step. When the
	// derivative is found, the function value is recalculated (if failed at least once), the quotient is computed,
	// and the new value for the argument is corrected.
	for tryOut := uint32(0); tryOut < maxArgCorrections; tryOut++ {
		if currentDerRes, err = derivative(current); err != nil {
			return 0, err
		}
		// This function may panic, but we're covered in the parent call.
		currentDerImg = sets.UpCastOneTo(currentDerRes, sets.R).(float64)
		if currentDerImg != 0 {
			// evaluate f(currentArg) into currentImg if told to recalculate.
			if recalculate {
				if currentRes, err = expression(current); err != nil {
					return 0, err
				}
				// This function may panic, but we're covered in the parent call.
				currentImg = sets.UpCastOneTo(currentRes, sets.R).(float64)
			}
			// Compute x = x - f(x)/f'(x).
			return current - (currentImg / currentDerImg), nil
		} else {
			recalculate = true
			random = rand.Float64() - 0.5
			current = current + epsilon * random
		}
	}
	return 0, ErrMaxNewtonRaphsonArgCorrections
}


// NewtonRaphson tries a Newton-Raphson iteration to find the nearest approximate root of the given function.
// If you want to find x so that f(x) == target, just define g(x) = f(x) - target and you will be on your way.
// While the usual newton-raphson fails when f'(x) (or g'(x)) is zero, we replace such value by the same epsilon
// in the tolerance. The number of max iterations, and the number of max corrections of the argument when the
// derivative is zero must be provided as well.
func NewtonRaphson(expression, derivative func(float64) (sets.Number, error), initialGuess, epsilon float64,
	               maxIterations, maxArgCorrectionsPerIteration uint32) (result float64, exception error) {
	var err error
	var currentArg float64
	var currentRes sets.Number
	var currentDerRes sets.Number
	var currentImg float64
	var currentDerImg float64
	var zero = float64(0)
	var quot = float64(0)
	var delta = float64(0)
	var random = float64(0)
	rand.Seed(time.Now().UTC().UnixNano())
	currentArg = initialGuess
	if maxIterations == 0 {
		maxIterations = 100
	}
	var iteration uint32
	defer func() {
		if recovered := recover(); recovered != nil {
			result = 0
			exception = errors.New(recovered.(string))
		}
	}()
	for iteration = 0; iteration < maxIterations; iteration++ {
		// evaluate f(currentArg) into currentImg.
		if currentRes, err = expression(currentArg); err != nil {
			return 0, err
		}
		// This function may panic, but we're covered here.
		currentImg = sets.UpCastOneTo(currentRes, sets.R).(float64)
		// If it is close to zero, then return.
		if diffUtils.CloseTo(currentImg, zero, epsilon) {
			return currentArg, nil
		}
		// Otherwise, process a newton-raphson step.
		if currentArg, err = nextNewtonRaphsonStep(
			expression, derivative, currentRes, currentDerRes, currentImg, currentDerImg, quot,
			currentArg, epsilon, delta, random, maxArgCorrectionsPerIteration,
		); err != nil {
			return 0, err
		}
	}
	return 0, goalErrors.ErrIterationsExhausted
}