package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
)


type AddExpr struct {
	terms []Expression
}


func (add AddExpr) Evaluate(args Arguments) (sets.Number, error) {
	terms := make([]sets.Number, len(add.terms))
	for index, term := range add.terms {
		if evaluated, err := term.Evaluate(args); err != nil {
			return nil, err
		} else {
			terms[index] = evaluated
		}
	}
	return ops.Add(terms...), nil
}


func (add AddExpr) Derivative(wrt Variable) (Expression, error) {
	derivedTerms := make([]Expression, len(add.terms))
	for index, term := range add.terms {
		if derivatedTerm, err := term.Derivative(wrt); err != nil {
			return nil, err
		} else {
			derivedTerms[index] = derivatedTerm
		}
	}
	return Add(derivedTerms...).Simplify()
}


func (add AddExpr) CollectVariables(variables Variables) {
	for _, term := range add.terms {
		term.CollectVariables(variables)
	}
}


func (add AddExpr) IsConstant(wrt Variable) bool {
	for _, term := range add.terms {
		if !term.IsConstant(wrt) {
			return false
		}
	}
	return true
}


func (add AddExpr) Simplify() (Expression, error) {
	simplifiedTerms := []sets.Number{}
	nonSimplifiedTerms := []Expression{}

	for _, term := range add.terms {
		if simplified, err := term.Simplify(); err != nil {
			return nil, err
		} else if num, ok := simplified.(Constant); ok {
			simplifiedTerms = append(simplifiedTerms, num.number)
		} else {
			nonSimplifiedTerms = append(nonSimplifiedTerms, num)
		}
	}

	simplifiedSummary := ops.Add(simplifiedTerms...)
	if len(nonSimplifiedTerms) != 0 {
		nonSimplifiedTerms = append(nonSimplifiedTerms, Constant{simplifiedSummary})
		return Add(nonSimplifiedTerms...), nil
	} else {
		return Constant{simplifiedSummary}, nil
	}
}


func (add AddExpr) String() string {
	// TODO
	return ""
}


func flattenTerms(terms []Expression) []Expression {
	flattenedTerms := make([]Expression, 2)
	for _, term := range terms {
		if addExpr, ok := term.(AddExpr); ok {
			for _, term := range flattenTerms(addExpr.terms) {
				flattenedTerms = append(flattenedTerms, term)
			}
		} else {
			flattenedTerms = append(flattenedTerms, term)
		}
	}
	return flattenedTerms
}


func Add(terms ...Expression) Expression {
	return AddExpr{flattenTerms(terms)}
}