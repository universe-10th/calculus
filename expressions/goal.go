package expressions

import (
	"github.com/universe-10th/calculus/sets"
	"fmt"
	"github.com/universe-10th/calculus/errors"
)

// Goal-seeking algorithm engines will have only one method
// to care about: evaluating an expression to find their
// root (only one root is supported here - algorithms should
// choose wisely or allow the user to choose the root). The
// expression is already curried, with only one variable
// left.
type GoalSeekingAlgorithm interface {
	FindRoot(goalBasedExpression Expression) (sets.Number, error)
}


// Engine factories take some arguments (related to the
// intended evaluation and variables) and return an instance
// of goal-seeking algorithm.
type GoalSeekingAlgorithmFactory func(
	arguments Arguments, inverted Variable, fullDomain Variables,
) (GoalSeekingAlgorithm, error)


// Goal-seeks are a special kind of expressions involving two
// steps that are essentially distinct in nature:
// 1. Evaluate the underlying "goal" expression. Get a value.
// 2. Take the target expression, and normalize it. This will
//    create a new expression, being the subtraction of a
//    curried version of the target expression minus the
//    just-evaluated "goal".
// 4. Pass that expression to the underlying algorithm.
type GoalSeekExpr struct {
	// The goal expression is the one we will fully evaluate
	// by the usual Evaluate mechanism. Also the variables
	// involved in the goal will be cached
	goal         Expression
	targetDomain Variables
	// The target expression is the one we will curry and then
	// normalize to "zero" by subtracting the computed goal.
	// After currying, the new goal is to find the appropriate
	// value for the "inverted" variable making this expression
	// evaluate to 0.
	target Expression
	// The inverted variable. It must be present among the domain
	// variables in the target expression. It is the variable we
	// actually want to return.
	inverted Variable
	// The algorithm factory that we'll use to instantiate the
	// engine that will be used as goal-seeker.
	factory GoalSeekingAlgorithmFactory
}


// Derivative cannot be calculated
func (goalSeekExpr GoalSeekExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, errors.ErrNotDerivableExpression
}


// Simplifying simplifies the expression and, if only depending on constants, evaluates it.
func (goalSeekExpr GoalSeekExpr) Simplify() (Expression, error) {
	if simplifiedGoal, err := goalSeekExpr.goal.Simplify(); err != nil {
		return nil, err
	} else if simplifiedTarget, err := goalSeekExpr.target.Simplify(); err != nil {
		return nil, err
	} else {
		_, okGoal := simplifiedGoal.(Constant)
		_, okTarget := simplifiedTarget.(Constant)
		simplifiedExpr := GoalSeek(simplifiedGoal, simplifiedTarget, goalSeekExpr.inverted, goalSeekExpr.factory)
		if okTarget && okGoal {
			// Create a dummy one, and evaluate with no arguments.
			simplifiedExpr = GoalSeek(simplifiedGoal, simplifiedTarget, goalSeekExpr.inverted, goalSeekExpr.factory)
			if result, err := simplifiedExpr.Evaluate(Arguments{}); err != nil {
				return nil, err
			} else {
				return Num(result), nil
			}
		} else {
			return simplifiedExpr, nil
		}
	}
}


// CollectVariables digs into the target's variables (except for the inverted one),
// and also the goal's variables.
func (goalSeekExpr GoalSeekExpr) CollectVariables(variables Variables) {
	goalSeekExpr.target.CollectVariables(variables)
	delete(variables, goalSeekExpr.inverted)
	goalSeekExpr.goal.CollectVariables(variables)
}


// IsConstant will return whether the inner expression is constant with respect to the given variable.
// This check involves two requirements:
// - The goal is constant with respect to the variable.
// - The target is constant with respect to the variable, or the variable is the inverted.
func (goalSeekExpr GoalSeekExpr) IsConstant(wrt Variable) bool {
	return goalSeekExpr.goal.IsConstant(wrt) && (goalSeekExpr.inverted == wrt || goalSeekExpr.target.IsConstant(wrt))
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (goalSeekExpr GoalSeekExpr) Curry(args Arguments) (Expression, error) {
	if curriedGoal, err := goalSeekExpr.goal.Curry(args); err != nil {
		return nil, err
	} else 	if curriedTarget, err := goalSeekExpr.goal.Curry(goalSeekExpr.getNonInvertedArguments(args)); err != nil {
		return nil, err
	} else {
		return GoalSeek(curriedGoal, curriedTarget, goalSeekExpr.inverted, goalSeekExpr.factory).Simplify()
	}
}


// Gets all the arguments, but the inverted one.
func (goalSeekExpr GoalSeekExpr) getNonInvertedArguments(arguments Arguments) Arguments {
	argumentsCopy := Arguments{}
	for key, value := range arguments {
		if key != goalSeekExpr.inverted {
			argumentsCopy[key] = value
		}
	}
	return argumentsCopy
}


// Evaluate computes the factorial over the evaluated inner argument's value.
// It will be an error if the inner value does not evaluate into N0.
func (goalSeekExpr GoalSeekExpr) Evaluate(args Arguments) (sets.Number, error) {
	fmt.Print("Target domain:", goalSeekExpr.targetDomain, "inverted:", goalSeekExpr.inverted)
	if _, ok := goalSeekExpr.targetDomain[goalSeekExpr.inverted]; !ok {
		return nil, errors.ErrInvertedVariableNotInDomain
	} else if goal, err := goalSeekExpr.goal.Evaluate(args); err != nil {
		return nil, err
	} else if engine, err := goalSeekExpr.factory(args, goalSeekExpr.inverted, goalSeekExpr.targetDomain); err != nil {
		return nil, err
	} else if curried, err := goalSeekExpr.target.Curry(goalSeekExpr.getNonInvertedArguments(args)); err != nil {
		return nil, err
	} else {
		return engine.FindRoot(Sub(curried, Num(goal)))
	}
}


func (goalSeekExpr GoalSeekExpr) String() string {
	return fmt.Sprintf("GoalSeek({%s / %s == %s})", goalSeekExpr.inverted, goalSeekExpr.target, goalSeekExpr.goal)
}


func (goalSeekExpr GoalSeekExpr) IsSelfContained() bool {
	return true
}


func GoalSeek(goal, target Expression, inverted Variable, factory GoalSeekingAlgorithmFactory) Expression {
	targetDomain := Variables{}
	target.CollectVariables(targetDomain)
	return GoalSeekExpr{goal, targetDomain,target, inverted, factory}
}