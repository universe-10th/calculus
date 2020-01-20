package expressions

import (
	"github.com/universe-10th/calculus/v2/std/sets"
	"github.com/universe-10th/calculus/v2/std/ops"
	"strings"
)


// AddExpr is an addition expression. Represents both addition and subtraction.
// Actually, subtractions are instantiated as additions of negations.
type AddExpr struct {
	terms []Expression
}


// Curry will try currying each term independently, and then generate a new addition.
// The terms may be converted to constant terms. Considering this, Curry will try
// simplifying the expression before returning it.
func (add AddExpr) Curry(args Arguments) (Expression, error) {
	newTerms := make([]Expression, len(add.terms))
	for index, value := range add.terms {
		if curried, err := value.Curry(args); err != nil {
			return nil, err
		} else {
			newTerms[index] = curried
		}
	}
	return Add(newTerms...).Simplify()
}


// Evaluate computes the added evaluted values of the addition terms (recursively).
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


// Derivative applies the addition rule of derivatives.
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


// CollectVariables digs recursively in the addition terms.
func (add AddExpr) CollectVariables(variables Variables) {
	for _, term := range add.terms {
		term.CollectVariables(variables)
	}
}


// IsConstant is true only if all the terms are constant with respect to the given variable.
func (add AddExpr) IsConstant(wrt Variable) bool {
	for _, term := range add.terms {
		if !term.IsConstant(wrt) {
			return false
		}
	}
	return true
}


// Simplify compresses all the constant terms into one single constant term.
// The result is returned as a new expressions rather than modifying the current one.
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
		if len(nonSimplifiedTerms) == 1 {
			return nonSimplifiedTerms[0], nil
		} else {
			return Add(nonSimplifiedTerms...), nil
		}
	} else {
		if simplifiedSummary != nil {
			return Constant{simplifiedSummary}, nil
		} else {
			return Constant{ops.Zero(sets.N0)}, nil
		}
	}
}


// String represents the addition expression appropriately.
// This means: it adds parentheses appropriately and also replaces + with - for negated terms or negative constants.
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
		case Constant:
			if ops.IsNegative(v.number) {
				builder.WriteString(" - ")
				builder.WriteString(Negated(v).String())
			} else {
				builder.WriteString(" + ")
				builder.WriteString(expression.String())
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


// Add constructs a new addition node given their terms.
func Add(terms ...Expression) Expression {
	return AddExpr{flattenTerms(terms)}
}


// Sub constructs a new subtracting addition node given their terms.
func Sub(minuend Expression, subtrahends ...Expression) Expression {
	return Add(minuend, Negated(Add(subtrahends...)))
}