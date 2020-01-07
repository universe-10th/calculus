package models

import (
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/models/errors"
)


// Goal-seeker algorithm engines will have only one method
// to care about: evaluating an expression to find their
// root (only one root is supported here - algorithms should
// choose wisely or allow the user to choose the root). The
// expression is already curried, with only one variable
// left.
type GoalSeekerAlgorithm interface {
	FindRoot(curriedExpression expressions.Expression, variable expressions.Variable) (sets.Number, error)
}


// Engine factories take some arguments (related to the
// intended evaluation and variables) and return an instance
// of goal-seeking algorithm.
type GoalSeekerAlgorithmFactory func(
	arguments expressions.Arguments, coDomain,
	inverted expressions.Variable, fullDomain expressions.Variables,
) (GoalSeekerAlgorithm, error)


// Goal-seeker models are quite like Single Output
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
type GoalSeekerModelFlow struct {
	CustomModelFlow
	expression expressions.Expression
	coDomainVariable expressions.Variable
	invertedVariable expressions.Variable
	factory GoalSeekerAlgorithmFactory
}


// Given an expression to evaluate, the co-domain variable for
// the expression, and the variable to invert (root-find their
// value) creates an abstract Goal-seeker model flow. Additionally,
// the algorithm factory function is also needed. This said, this
// function should be used as part of construction of other
// (specific, totally or partially implemented) goal-seeker model
// flows, but it can safely be used directly (although more
// verbose than using custom constructors).
func NewGoalSeekerModelFlow(expression expressions.Expression, coDomain, invertedVariable expressions.Variable,
	                            factory GoalSeekerAlgorithmFactory) (*GoalSeekerModelFlow, error) {
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
		return &GoalSeekerModelFlow{
			*custom, expression,
			coDomain, invertedVariable,
			factory,
		}, nil
	}
}


// Gets all the arguments, but the inverted one.
func (goalSeekerModelFlow GoalSeekerModelFlow) getNonInvertedArguments(arguments expressions.Arguments) expressions.Arguments {
	argumentsCopy := expressions.Arguments{}
	for key, value := range arguments {
		if key != goalSeekerModelFlow.coDomainVariable {
			argumentsCopy[key] = value
		}
	}
	return argumentsCopy
}


// Evaluates the model. First, it creates a goal-based expression from the given
// one (and the co-domain value). Then, it curries such expression with all the
// input values except for the inverted value. Finally, with only one variable
// left (the inverted one), it tries the root-finding algorithm.
func (goalSeekerModelFlow *GoalSeekerModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	fullDomain := goalSeekerModelFlow.Input()
	coDomain := goalSeekerModelFlow.coDomainVariable
	inverted := goalSeekerModelFlow.invertedVariable
	if engine, err := goalSeekerModelFlow.factory(arguments, coDomain, inverted, fullDomain); err != nil {
		return nil, err
	} else {
		// 1. Keep only non-inverted arguments.
		arguments = goalSeekerModelFlow.getNonInvertedArguments(arguments)
		// 2. Curry the expression given the arguments. The remaining
		// expression must have exactly one free variable, corresponding
		// to the inverted variable.
		if curried, err := goalSeekerModelFlow.expression.Curry(arguments); err != nil {
			return nil, err
		} else {
			// Check the only free variable is the inverted.
			variables := expressions.Variables{}
			curried.CollectVariables(variables)
			if _, ok := variables[inverted]; !ok || len(variables) != 1 {
				return nil, errors.ErrInsufficientArguments
			}

			// 3. Now, tell the engine to run the algorithm given the curried
			// expression and the inverted variable
			if result, err := engine.FindRoot(curried, inverted); err != nil {
				return nil, err
			} else {
				return expressions.Arguments{inverted: result}, nil
			}
		}
	}
}