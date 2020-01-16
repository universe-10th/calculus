package models

import (
	"github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/errors"
	"github.com/universe-10th/calculus/sets"
)


// Just a mapping of variables and the
// model they are used into, as arguments.
type ModelFlowMapping struct {
	variables expressions.Variables
	modelFlow ModelFlow
}


// A model is a collection of model flows.
// Evaluating a model implies evaluating
// the corresponding involved flow, which
// is picked according to the involved
// variables among the arguments.
type Model struct {
	modelFlowMappings []ModelFlowMapping
}


// Creates an empty model. Flows must be added
// later via calls to newModel.AddFlow(...) by
// passing a single argument implementing the
// ModelFlow interface.
func NewModel() *Model {
	return &Model{
		[]ModelFlowMapping{},
	}
}


func isSubset(x, y expressions.Variables) bool {
	for vx := range x {
		if _, ok := y[vx]; !ok {
			return false
		}
	}
	return true
}


func satisfies(args expressions.Arguments, spec expressions.Variables) bool {
	for variable := range spec {
		if _, ok := args[variable]; !ok {
			return false
		}
	}
	return true
}


// Adds a flow to the model. Added flows cannot be removed later.
// A flow cannot be added if its input spec has a conflict with
// an already-added flow's input spec. This means: the input of
// the new flow is equal, subset, or superset of one or more
// input specs already registered.
func (model *Model) AddFlow(modelFlow ModelFlow) error {
	input := modelFlow.Input()
	for _, mapping := range model.modelFlowMappings {
		if isSubset(mapping.variables, input) || isSubset(input, mapping.variables) {
			return errors.ErrAmbiguousInputSpec
		}
	}
	model.modelFlowMappings = append(model.modelFlowMappings, ModelFlowMapping{
		input, modelFlow,
	})
	return nil
}


// Gets the model of the first mapping having a set
// of input variables being satisfied by the argument
// variables, or an error if no model input matches at
// least a subset of the specified variables as input.
func (model *Model) GetFlow(arguments expressions.Arguments) (ModelFlow, error) {
	for _, mapping := range model.modelFlowMappings {
		if satisfies(arguments, mapping.variables) {
			return mapping.modelFlow, nil
		}
	}
	return nil, errors.ErrNoRegisteredFlowFowInput
}


// Attempts evaluating the model, finding the
// appropriate flow, given a set of arguments.
func (model *Model) Evaluate(arguments expressions.Arguments) (sets.Number, error) {
	if flow, err := model.GetFlow(arguments); err != nil {
		return nil, err
	} else {
		return flow.Evaluate(arguments)
	}
}