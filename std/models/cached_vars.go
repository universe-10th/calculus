package models

import "github.com/universe-10th/calculus/v2/std/expressions"


// Holds a reference to cached variables. Can return them later.
type cachedVars struct {
	input       expressions.Variables
	output      expressions.Variables
}


// Returns a copy of the input vars set.
func (cached *cachedVars) Input() expressions.Variables {
	input := expressions.Variables{}
	for inputVar, _ := range cached.input {
		input[inputVar] = true
	}
	return input
}


// Returns a copy of the output vars set.
func (cached *cachedVars) Output() expressions.Variables {
	output := expressions.Variables{}
	for outputVar, _ := range cached.output {
		output[outputVar] = true
	}
	return output
}


// Tells whether a variable is contained among the input(s).
func (cached *cachedVars) DefinesInput(variable expressions.Variable) bool {
	return cached.input[variable]
}


// Tells whether a variable is contained among the output(s).
func (cached *cachedVars) DefinesOutput(variable expressions.Variable) bool {
	return cached.output[variable]
}


// Tells whether domain variables are present among the
// co-domain as well. That case is a model inconsistency.
func (cached *cachedVars) HasConsistentDomain() bool {
	for outputVar, _ := range cached.output {
		if _, ok := cached.input[outputVar]; ok {
			return false
		}
	}
	return true
}


// Returns the same object.
func (cached *cachedVars) CachedVars() cachedVars {
	return *cached
}


// merges a cachedVars set into another one.
func (cached *cachedVars) Merge(source cachedVars) {
	for inputVar, _ := range source.input {
		cached.input[inputVar] = true
	}
	for outputVar := range source.output {
		cached.output[outputVar] = true
	}
}
