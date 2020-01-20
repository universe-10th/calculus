package errors

import "errors"


var ErrExpressionIsNil = errors.New("one or more expressions for this flow are nil")
var ErrOutputVariableInsideFlowExpression = errors.New("model flow's output variable is inside flow's expression")
var ErrInsufficientArguments = errors.New("insufficient arguments for the model flow")
var ErrNoSerialModelFlowsGiven = errors.New("no serial model flows given")
var ErrOutputToInputChainNotSatisfied = errors.New("there is at least one step model flow where the output from it does not satisfy the input of the next")
var ErrNoParallelModelFlowsGiven = errors.New("no parallel model flows given")
var ErrSomeModelOutputsBelongToOtherModelInputs = errors.New("some output variables in one model flow belong to input variables of other model flow(s)")
var ErrOutputVariableMergedTwice = errors.New("an output variable from one input model flow clashes with an output variable in other model flow(s)")
var ErrModelFlowIsNil = errors.New("a given model flow is nil")
var ErrRootFindingVariableMissingFromExpression = errors.New("Root-finding / inverted variable missing from expression domain")
