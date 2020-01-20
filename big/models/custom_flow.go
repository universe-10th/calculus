package models

import (
	"github.com/universe-10th/calculus/v2/big/expressions"
	"github.com/universe-10th/calculus/v2/big/models/errors"
)


// A custom model will only define the involved variables
// but its behaviour must be implemented. Typically used
// as an inheritance component of an intended descendant
// (composed) type.
type CustomModelFlow struct {
	cachedVars
}


// Creates a base instance for a custom model flow.
// The result of this function will be used as a component
// for a derived-class instance.
func NewCustomModelFlow(inputVars, outputVars expressions.Variables) (*CustomModelFlow, error) {
	cachedVars := cachedVars{inputVars, outputVars}
	if !cachedVars.HasConsistentDomain() {
		return nil, errors.ErrOutputVariableInsideFlowExpression
	}
	return &CustomModelFlow{cachedVars}, nil
}


func (customModelFlow *CustomModelFlow) Evaluate(arguments expressions.Arguments) (expressions.Arguments, error) {
	panic("this behaviour must be implemented")
}