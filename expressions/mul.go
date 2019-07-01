package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"strings"
)


// MulExpr is a product expression. Represents both multiplication and division.
// Actually, divisions are instantiated as products of inverses.
type MulExpr struct {
	factors []Expression
}


// Evaluate evaluates the product of the evaluated factor's values.
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


// Derivative applies the multiplication rule of derivatives.
// While we are used to the simple case of (fg)' = f'g + g'f,
// the general case involves n terms of the n factors being
// multiplied, and inside each different term, a different
// function is being derived.
func (mul MulExpr) Derivative(wrt Variable) (Expression, error) {
	factors := make([]Expression, len(mul.factors))
	for index, factor := range mul.factors {
		factors[index] = factor
	}
	derivedFactors := make([]Expression, len(mul.factors))
	for index, factor := range factors {
		if derivedFactor, err := factor.Derivative(wrt); err == nil {
			derivedFactors[index] = derivedFactor
		} else {
			return nil, err
		}
	}
	terms := make([]Expression, len(mul.factors))
	for index := range mul.factors {
		termFactors := make([]Expression, len(mul.factors))
		for index2 := range mul.factors {
			if index == index2 {
				termFactors[index2] = derivedFactors[index2]
			} else {
				termFactors[index2] = factors[index2]
			}
		}
		terms[index] = Mul(termFactors...)
	}
	return Add(terms...).Simplify()
}


// CollectVariables digs into the inner factors.
func (mul MulExpr) CollectVariables(variables Variables) {
	for _, term := range mul.factors {
		term.CollectVariables(variables)
	}
}


// IsConstant returns whether all the factors are constant with respect to the given variable.
func (mul MulExpr) IsConstant(wrt Variable) bool {
	for _, factor := range mul.factors {
		if !factor.IsConstant(wrt) {
			return false
		}
	}
	return true
}


// Simplify compresses all the constant factors into one single constant factor.
// The result is returned as a new expressions rather than modifying the current one.
// Note: this method is optimized: if at least one factor results in constant 0, the
// constant 0 expression will be returned.
func (mul MulExpr) Simplify() (Expression, error) {
	simplifiedTerms := []sets.Number{}
	nonSimplifiedTerms := []Expression{}

	for _, factor := range mul.factors {
		if simplified, err := factor.Simplify(); err != nil {
			return nil, err
		} else if num, ok := simplified.(Constant); ok {
			simplifiedTerms = append(simplifiedTerms, num.number)
		} else {
			nonSimplifiedTerms = append(nonSimplifiedTerms, simplified)
		}
	}

	simplifiedSummary := ops.Mul(simplifiedTerms...)
	if simplifiedSummary != nil && ops.IsZero(simplifiedSummary) {
		return Num(0), nil
	}

	if len(nonSimplifiedTerms) != 0 {
		if simplifiedSummary != nil && !ops.IsOne(simplifiedSummary) {
			nonSimplifiedTerms = append(nonSimplifiedTerms, Constant{simplifiedSummary})
		}
		return Mul(nonSimplifiedTerms...), nil
	} else {
		if simplifiedSummary != nil {
			return Constant{simplifiedSummary}, nil
		} else {
			return Constant{ops.One(sets.N0)}, nil
		}
	}
}


// String represents the multiplication expression appropriately.
// This means: it adds parentheses appropriately and also replaces * with / for inverted terms.
func (mul MulExpr) String() string {
	builder := strings.Builder{}
	if len(mul.factors) == 0 {
		return ""
	} else {
		// Since addition has the lowest precedence,
		// it must be wrapped
		expression := mul.factors[0]
		switch expression.(type) {
		case AddExpr:
			builder.WriteString("(")
			builder.WriteString(expression.String())
			builder.WriteString(")")
		default:
			builder.WriteString(expression.String())
		}
	}

	for _, expression := range mul.factors[1:] {
		// Inv operator will be added: / (E) for inner mul or addition, / E for any other inner.
		// Add operation will be added: * (E)
		// Other operators: * E
		switch v := expression.(type) {
		case InverseExpr:
			inner := v.arg
			builder.WriteString(" / ")
			switch inner.(type) {
			case MulExpr, AddExpr:
				builder.WriteString("(")
				builder.WriteString(inner.String())
				builder.WriteString(")")
			default:
				builder.WriteString(inner.String())
			}
		case AddExpr:
			builder.WriteString(" * (")
			builder.WriteString(expression.String())
			builder.WriteString(")")
		default:
			builder.WriteString(" * ")
			builder.WriteString(expression.String())
		}
	}

	return builder.String()
}


func flattenFactors(factors []Expression) []Expression {
	flattenedFactors := make([]Expression, 0)
	for _, factor := range factors {
		if mulExpr, ok := factor.(MulExpr); ok {
			for _, term := range flattenFactors(mulExpr.factors) {
				flattenedFactors = append(flattenedFactors, term)
			}
		} else {
			flattenedFactors = append(flattenedFactors, factor)
		}
	}
	return flattenedFactors
}


// Mul constructs a new multiplication node given their factors.
func Mul(factors ...Expression) Expression {
	return MulExpr{flattenFactors(factors)}
}