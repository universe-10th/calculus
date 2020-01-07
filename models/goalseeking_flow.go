package models

import (
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/models/errors"
	"github.com/universe-10th/calculus/utils"
)


// Goal-seeking algorithm engines will have only one method
// to care about: evaluating an expression to find their
// root (only one root is supported here - algorithms should
// choose wisely or allow the user to choose the root). The
// expression is already curried, with only one variable
// left.
type GoalSeekingAlgorithm interface {
	FindRoot(goalBasedExpression expressions.Expression, variable expressions.Variable) (sets.Number, error)
}


// Engine factories take some arguments (related to the
// intended evaluation and variables) and return an instance
// of goal-seeking algorithm.
type GoalSeekingAlgorithmFactory func(
	arguments expressions.Arguments, coDomain,
	inverted expressions.Variable, fullDomain expressions.Variables,
) (GoalSeekingAlgorithm, error)


// Goal-seeking models are quite like Single Output
// flows, but instead of using an expression and
// being interested in the traditional result of
// it on evaluation (the direct evaluation output),
// it uses a goal-seek algorithm, which involves:
// specifying the output variable (which must not be
// present among the input variables in the expression),
// specifying a goal value (value for the output variable),
// specifying which variable among the expressions's input
// variables to ignore, and specifying the values for all
// the other input variables. The evaluation will try to
// root-find the ignored (inverted) variable and return
// it as output in this flow.
//
// There are many different algorithms for goal-seeks
// (root-finding algorithms) and each of them will
// initialize in different ways the underlying algorithm
// engine to use.
//
// This class is abstract: it must be inherited since the
// underlying algorithm engine instantiation must be
// implemented.
type GoalSeekingModelFlow struct {
	CustomModelFlow
	expression expressions.Expression
	coDomainVariable expressions.Variable
	invertedVariable expressions.Variable
	factory GoalSeekingAlgorithmFactory
}


// Given an expression to evaluate, the co-domain variable for
// the expression, and the variable to invert (root-find their
// value) creates an abstract Goal-seeking model flow. Additionally,
// the algorithm factory function is also needed. This said, this
// function should be used as part of construction of other
// (specific, totally or partially implemented) goal-seeking model
// flows, but it can safely be used directly (although more
// verbose than using custom constructors).
func NewGoalSeekingModelFlow(expression expressions.Expression, coDomain, invertedVariable expressions.Variable,
	                            factory GoalSeekingAlgorithmFactory) (*GoalSeekingModelFlow, error) {
	inputVars := expressions.Variables{}
	outputVars := expressions.Variables{invertedVariable: true}

	if expression == nil {
		return nil, errors.ErrExpressionIsNil
	} else {
		expression.CollectVariables(inputVars)
		if _, ok := inputVars[invertedVariable]; !ok {
			return nil, errors.ErrRootFindingVariableMissingFromExpression
		} else {
			delete(inputVars, invertedVariable)
			inputVars[coDomain] = true
		}
	}

	if custom, err := NewCustomModelFlow(inputVars, outputVars); err != nil {
		return nil, err
	} else {
		return &GoalSeekingModelFlow{
			*custom, expression,
			coDomain, invertedVariable,
			factory,
		}, nil
	}
}


// Gets all the arguments, but the inverted one.
func (goalSeekingModelFlow GoalSeekingModelFlow) getNonInvertedArguments(arguments expressions.Arguments) expressions.Arguments {
	argumentsCopy := expressions.Arguments{}
	for key, value := range arguments {
		if key != goalSeekingModelFlow.coDomainVariable {
			argumentsCopy[key] = value
		}
	}
	return argumentsCopy
}


// Evaluates the model. First, it creates a goal-based expression from the given
// one (and the co-domain value). Then, it curries such expression with all the
// input values except for the inverted value. Finally, with only one variable
// left (the inverted one), it tries the root-finding algorithm.
func (goalSeekingModelFlow *GoalSeekingModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	fullDomain := goalSeekingModelFlow.Input()
	coDomain := goalSeekingModelFlow.coDomainVariable
	inverted := goalSeekingModelFlow.invertedVariable
	if engine, err := goalSeekingModelFlow.factory(arguments, coDomain, inverted, fullDomain); err != nil {
		return nil, err
	} else {
		// 1. Keep only non-inverted arguments.
		arguments = goalSeekingModelFlow.getNonInvertedArguments(arguments)
		// 2. Curry the expression given the arguments. The remaining
		// expression must have exactly one free variable, corresponding
		// to the inverted variable.
		if curried, err := goalSeekingModelFlow.expression.Curry(arguments); err != nil {
			return nil, err
		} else {
			// 3. Turn it into a goal-based expression. For this to work,
			// a value must be set, among the arguments, for the co-domain
			// variable
			if goal, ok := arguments[coDomain]; !ok {
				return nil, errors.ErrInsufficientArguments
			} else {
				goalBased := utils.GoalBasedExpression(curried, goal)

				// Check the only free variable is the inverted.
				variables := expressions.Variables{}
				goalBased.CollectVariables(variables)
				if _, ok := variables[inverted]; !ok || len(variables) != 1 {
					return nil, errors.ErrInsufficientArguments
				}

				// 4. Now, tell the engine to run the algorithm given the goal-based
				// expression and the inverted variable.
				if result, err := engine.FindRoot(goalBased, inverted); err != nil {
					return nil, err
				} else {
					return expressions.Arguments{inverted: result}, nil
				}
			}
		}
	}
}