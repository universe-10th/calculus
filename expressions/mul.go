package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
)


type MulExpr struct {
	factors []Expression
}


func (mul MulExpr) Evaluate(args Arguments) (sets.Number, error) {
	factors := make([]sets.Number, len(mul.factors))
	for index, term := range mul.factors {
		if evaluated, err := term.Evaluate(args); err != nil {
			return nil, err
		} else {
			factors[index] = evaluated
		}
	}
	return ops.Mul(factors...), nil
}


func (mul MulExpr) Derivative(wrt Variable) (Expression, error) {
	factors := make([]Expression, len(mul.factors))
	for index, factor := range mul.factors {
		factors[index] = factor
	}
	derivatedFactors := make([]Expression, len(mul.factors))
	for index, factor := range factors {
		if derivatedFactor, err := factor.Derivative(wrt); err != nil {
			derivatedFactors[index] = derivatedFactor
		}
	}
	terms := make([]Expression, len(mul.factors))
	for index := range terms {
		termFactors := make([]Expression, len(mul.factors))
		for index2 := range terms {
			if index == index2 {
				termFactors[index2] = derivatedFactors[index2]
			} else {
				termFactors[index2] = factors[index2]
			}
		}
		terms[index] = Mul(termFactors...)
	}
	return Add(terms...), nil
}


func (mul MulExpr) CollectVariables(variables Variables) {
	for _, term := range mul.factors {
		term.CollectVariables(variables)
	}
}


func (mul MulExpr) String() string {
	// TODO
	return ""
}


func flattenFactors(factors []Expression) []Expression {
	flattenedFactors := make([]Expression, 2)
	for _, term := range factors {
		if mulExpr, ok := term.(MulExpr); ok {
			for _, term := range flattenFactors(mulExpr.factors) {
				flattenedFactors = append(flattenedFactors, term)
			}
		} else {
			flattenedFactors = append(flattenedFactors, term)
		}
	}
	return flattenedFactors
}


func Mul(factors ...Expression) Expression {
	return MulExpr{flattenFactors(factors)}
}