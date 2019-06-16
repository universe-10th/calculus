package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"strings"
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
		if derivedTerm, err := term.Derivative(wrt); err != nil {
			return nil, err
		} else {
			derivedTerms[index] = derivedTerm
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
			nonSimplifiedTerms = append(nonSimplifiedTerms, simplified)
		}
	}

	simplifiedSummary := ops.Add(simplifiedTerms...)
	if len(nonSimplifiedTerms) != 0 {
		if simplifiedSummary != nil && !ops.IsZero(simplifiedSummary) {
			nonSimplifiedTerms = append(nonSimplifiedTerms, Constant{simplifiedSummary})
		}
		return Add(nonSimplifiedTerms...), nil
	} else {
		if simplifiedSummary != nil {
			return Constant{simplifiedSummary}, nil
		} else {
			return Constant{ops.Zero(sets.N0)}, nil
		}
	}
}


func (add AddExpr) String() string {
	builder := strings.Builder{}
	if len(add.terms) == 0 {
		return ""
	} else {
		builder.WriteString(add.terms[0].String())
	}

	for _, expression := range add.terms[1:] {
		// Unary operators: -E will be - V or - (Add) or - Expression
		switch v := expression.(type) {
		case NegatedExpr:
			inner := v.arg
			builder.WriteString(" - ")
			switch inner.(type) {
			case AddExpr:
				builder.WriteString("(")
				builder.WriteString(inner.String())
				builder.WriteString(")")
			default:
				builder.WriteString(inner.String())
			}
		default:
			builder.WriteString(" + ")
			builder.WriteString(expression.String())
		}
	}

	return builder.String()
}


func flattenTerms(terms []Expression) []Expression {
	flattenedTerms := make([]Expression, 0)
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