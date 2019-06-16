package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"strings"
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


func (mul MulExpr) CollectVariables(variables Variables) {
	for _, term := range mul.factors {
		term.CollectVariables(variables)
	}
}


func (mul MulExpr) IsConstant(wrt Variable) bool {
	for _, factor := range mul.factors {
		if !factor.IsConstant(wrt) {
			return false
		}
	}
	return true
}


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
	if len(nonSimplifiedTerms) != 0 {
		if simplifiedSummary != nil {
			nonSimplifiedTerms = append(nonSimplifiedTerms, Constant{simplifiedSummary})
		}
		return Add(nonSimplifiedTerms...), nil
	} else {
		if simplifiedSummary != nil {
			return Constant{simplifiedSummary}, nil
		} else {
			return Constant{ops.One(sets.N0)}, nil
		}
	}
}


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


func Mul(factors ...Expression) Expression {
	return MulExpr{flattenFactors(factors)}
}