package support

import "github.com/universe-10th/calculus/expressions"


// Creates a polynomial expressions given an involved variable and their coefficients.
// Coefficients are specified from main to constant. Example: Polynomial(X, 2, 3, 5, 7)
// would generate 2*X^3 + 3*X^2 + 5*X + 7. Please note: while the coefficients parameter
// accepts anything, you must only use primitive or *big.(Int, Rat, Float) numbers or
// your expression will explode on evaluation.
func PolynomialExpression(variable expressions.Variable, coefficients ...interface{}) expressions.Expression {
	if len(coefficients) == 0 {
		return nil
	} else {
		length := len(coefficients)
		terms := make([]expressions.Expression, length)
		for index, coefficient := range coefficients[0:length - 1] {
			terms[index] = expressions.Mul(
				expressions.Num(coefficient),
				expressions.Pow(variable, expressions.Num(length - 1 - index)),
			)
		}
		terms[length-1] = expressions.Num(coefficients[length-1])
		addition, _ := expressions.Add(terms...).Simplify()
		return addition
	}
}
