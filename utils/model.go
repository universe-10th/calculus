package utils

import (
	. "github.com/universe-10th/calculus/expressions"
	"errors"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/solvers"
)


var ErrNilModelMainExpression = errors.New("nil model's main expression")
var ErrMainVariableInsideMainExpression = errors.New("model's main variable is inside model's main expression")
var ErrMainExpressionIsConstant = errors.New("model's main expression is constant")
var ErrNoMissingVariable = errors.New("model cannot be evaluated: no missing variable from arguments")
var ErrMultipleMissingVariables = errors.New("model cannot be evaluated: multiple missing variables from arguments (perhaps a partial model was intended?)")
var ErrSolverNotGivenWhileNoCorollary = errors.New("model cannot be evaluated: attempted to evaluate a variable with no corollary, and no solver was given")
var ErrModelDoesNotReferVariable = errors.New("cannot change corollary: the given variable is not referred by the model")
var ErrCorollaryDoesNotReferExpectedVariables = errors.New("cannot change corollary: the given corollary does not refer the corresponding variables")


// Model involves an expression and a variable which is the
// concept behind the given expression. A model can also involve
// inverse expressions ("corollaries") for the involved variables
// in the main expression. By default, if those expressions are
// not given, the corollaries are implemented as a newton-raphson
// inverse check.
type Model struct {
	// The main involved variable.
	mainVariable   Variable
	// If this value is not nil, it is the main involved variable's
	// frozen value. This is because the main variable cannot be
	// actually partially evaluated, as other variables can.
	mainValue      sets.Number
	// The main involved expression.
	mainExpression Expression
	// The corollaries: corresponding expressions to clear out the
	// other involved variables. A nil expression will actually
	// mean a newton-raphson check over the main expression.
	corollaries    map[Variable]Expression
}


// NewModel creates a model given a variable and its expression.
// The given variable must not be inside the expression, and the
// expression must not be constant (this implies: the expression
// must imply at least one variable that is not the main one).
func NewModel(variable Variable, expression Expression) (*Model, error) {
	// Empty expressions are not allowed.
	if expression == nil {
		return nil, ErrNilModelMainExpression
	}

	// Collecting the variables.
	variables := Variables{}
	expression.CollectVariables(variables)
	// Constant expressions are neither allowed.
	if len(variables) == 0 {
		return nil, ErrMainExpressionIsConstant
	}

	// The main variable must not be inside the expression.
	if _, ok := variables[variable]; ok {
		return nil, ErrMainVariableInsideMainExpression
	}

	// Creating the corollaries map with empty values.
	corollaries := map[Variable]Expression{}
	for key := range variables {
		corollaries[key] = nil
	}

	// If everything is satisfied, we're ok.
	return &Model{variable, nil,expression, corollaries}, nil
}


// SetCorollary sets a corollary for the given variable in the model.
// The corollary must use the model's main variable and only those variables that are part of the model.
// Example: If the model is A=B*C + D, a corollary for D would be: D=A-B*C.
// A valid corollary is the nil value. It will use a solver instead when evaluating.
func (model *Model) SetCorollary(variable Variable, expression Expression) error {
	if _, ok := model.corollaries[variable]; !ok {
		return ErrModelDoesNotReferVariable
	} else {
		if expression == nil {
			// remove the corollary - a solver will be used.
			model.corollaries[variable] = nil
		} else {
			// The expression must include all the model variables EXCEPT the given one.
			// Also, the expression must contain the main variable.
			// Finally, the expression must not contain any other variable.
			variables := Variables{}
			expression.CollectVariables(variables)
			// The expression must refer the same number of variables the model's main
			// expression refers.
			if len(variables) != len(model.corollaries) {
				return ErrCorollaryDoesNotReferExpectedVariables
			}
			// The expression must refer the model's main variable.
			if _, ok := variables[model.mainVariable]; !ok {
				return ErrCorollaryDoesNotReferExpectedVariables
			}
			// The expression must refer all the other variables (which are keys in the
			// corollaries map) that are not the given one.
			for innerVariable := range variables {
				if innerVariable != model.mainVariable {
					if _, ok := model.corollaries[innerVariable]; !ok {
						return ErrCorollaryDoesNotReferExpectedVariables
					}
				}
			}
			model.corollaries[variable] = expression
		}
		return nil
	}
}


// Partial creates a partial model based on this one. It will partially evaluate
// the expression and also partially evaluate the main variable, if given.
func (model *Model) Partial(arguments Arguments) (*Model, error) {
	mainValue          := arguments[model.mainVariable]
	partialCorollaries := map[Variable]Expression{}
	var mainPartialExpression Expression
	var err                   error
	if mainPartialExpression, err = model.mainExpression.Curry(arguments); err != nil {
		return nil, err
	}
	for variable, corollary := range model.corollaries {
		if _, ok := arguments[variable]; !ok {
			// If the variable is not frozen, we must add the corollary but partially evaluated.
			if corollary != nil {
				if partialCorollary, err := corollary.Curry(arguments); err != nil {
					return nil, err
				} else {
					partialCorollaries[variable] = partialCorollary
				}
			} else {
				partialCorollaries[variable] = nil
			}
		}
	}
	newModel := &Model{
		mainVariable: model.mainVariable, mainValue: mainValue,
		mainExpression: mainPartialExpression, corollaries: partialCorollaries,
	}
	return newModel, nil
}


// Gets the only missing variable, if exactly one is missing from the arguments.
// An error will be returned if there are no exactly ONE missing variable, because
// a model needs to know which variable must it calculate.
func (model *Model) getMissingVariable(args Arguments) (Variable, error) {
	missingInArguments := false
	missingVariable    := Variable{}
	for variable := range model.corollaries {
		if _, ok := args[variable]; !ok {
			if missingInArguments {
				return Variable{}, ErrMultipleMissingVariables
			} else {
				missingInArguments = true
				missingVariable    = variable
			}
		}
	}
	if _, ok := args[model.mainVariable]; ok {
		if missingInArguments {
			return missingVariable, nil
		} else {
			return Variable{}, ErrNoMissingVariable
		}
	} else {
		if missingInArguments {
			return Variable{}, ErrMultipleMissingVariables
		} else {
			return model.mainVariable, nil
		}
	}
}


// Evaluate actually tries an evaluation of the model with the given arguments.
// Evaluating a model implying (N+1) variables will involve specifying N arguments
// which correspond to the variables and leaving exactly one variable out. The
// returned value will correspond to that missing variable.
func (model *Model) Evaluate(args Arguments, solver solvers.Solver) (sets.Number, error) {
	if variable, err := model.getMissingVariable(args); err != nil {
		// WRONG: Exactly one variable must be missing.
		return nil, err
	} else if variable == model.mainVariable {
		// Evaluate the main expression to compute the value for the main variable.
		return model.mainExpression.Evaluate(args)
	} else if corollary := model.corollaries[variable]; corollary == nil {
		// There is no corollary to evaluate. We will take the main expression
		// and partially evaluate it with the given arguments (the main variable
		// among the arguments will be ignored).
		// After it is curried, a variable will be free (the variable that we
		// want to calculate its value) so we need to convert our curried
		// expression which is f(variable)==V to g(variable)-V==0, where V is
		// the value of the main variable (either present in the current model or
		// present among the arguments).
		// After getting the goal function, a solver must be present and it will
		// be run to compute the result value.
		if solver == nil {
			return nil, ErrSolverNotGivenWhileNoCorollary
		} else if curried, err := model.mainExpression.Curry(args); err != nil {
			return nil, err
		} else if model.mainValue != nil {
			return solver(GoalBasedExpression(curried, model.mainValue))
		} else {
			return solver(GoalBasedExpression(curried, args[model.mainVariable]))
		}
		return nil, nil
	} else {
		return corollary.Evaluate(args)
	}
}