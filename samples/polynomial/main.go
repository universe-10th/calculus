package main

import (
	"fmt"
	"github.com/universe-10th/calculus/core/support"
	. "github.com/universe-10th/calculus/expressions"
)

func main() {
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2))
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2, -3))
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2, -3, 5))
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2, -3, 5, -7))
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2, -3, 5, -7, 11))
	fmt.Printf("A polynomial expression: %s\n", support.PolynomialExpression(X, 2, -3, 5, -7, 11, -13))
}
