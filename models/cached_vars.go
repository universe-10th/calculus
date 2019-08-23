package models

import "github.com/universe-10th/calculus/expressions"


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
	for outputVar, _ := range cached.input {
		output[outputVar] = true
	}
	return output
}


// Merges several given instances of cachedVars into one.
func mergeCachedVars(elements ...cachedVars) cachedVars {
	input := expressions.Variables{}
	output := expressions.Variables{}
	for _, element := range elements {
		for inputVar, _ := range element.input {
			input[inputVar] = true
		}
		for outputVar := range element.output {
			output[outputVar] = true
		}
	}
	return cachedVars{
		input,
		output,
	}
}