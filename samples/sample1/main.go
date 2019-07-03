package main

import (
	. "github.com/universe-10th-calculus/expressions"
	"fmt"
	"math/big"
	"github.com/universe-10th-calculus/solvers"
	"github.com/universe-10th-calculus/solvers/support"
	"github.com/universe-10th-calculus/utils"
)


func Sample1() Expression {
	return Mul(Add(
		X, Y, Negated(Z),
		Mul(X, Inverse(Y), Add(Y, Z), Inverse(Add(
			X, Y, Inverse(Z),
		))),
		Negated(Mul(Y, Z)),
	), Inverse(Mul(X, Y)), Pow(Add(Mul(X, Negated(Y)), Z), Add(X, Negated(Mul(Exp(Add(Y, Negated(Mul(Z, Num(34))))), Inverse(Z))))))
}


func Sample2() Expression {
	return Mul(Pow(Z, Add(Y, X)), Num(3))
}


func Sample3() Expression {
	return Mul(Pow(X, Num(3)), Ln(Mul(X, Num(4))), Tan(X))
}


func SampleNR() Expression {
	return Add(
		Num(2),
		Mul(Num(-3), X),
		Mul(Num(5), Pow(X, Num(2))),
		Mul(Num(-7), Pow(X, Num(3))),
		Mul(Num(11), Pow(X, Num(4))),
		Mul(Num(-13), Pow(X, Num(5))),
	)
}


func main() {
	value := Sample3()
	fmt.Println("Display: ", value.String())
	result, err := value.Evaluate(Arguments{
		X: 1, Y: 2, Z: 4,
	}.Wrap())
	fmt.Println("Evaluating with (X=1, Y=2, Z=4): ", result, err)
	if derivative, err := value.Derivative(X); err != nil {
		fmt.Printf("Error when deriving")
	} else {
		fmt.Println("Derivatie display: ", derivative.String())
		result, err := derivative.Evaluate(Arguments{
			X: 1, Y: 2, Z: 4,
		}.Wrap())
		fmt.Println("Evaluating derivative with (X=1, Y=2, Z=4): ", result, err)
	}
	nrSample := SampleNR()
	guess := big.NewFloat(1000000.0)
	epsilon := support.Epsilon(10)
	fmt.Printf("Testing Newton-Raphson sample: %s with initial guess: %f\n", nrSample, guess)
	if nrResult, err := solvers.NewtonRaphson(nrSample, guess, epsilon, 10000, 10); err != nil {
		fmt.Printf("Error in Newton-Raphson: %s\n", err)
	} else {
		fmt.Printf("Returned value: %f\n", nrResult)
	}
	value = Sample1()
	fmt.Printf("Currying expression %s with values (X=1, Y=2)...\n", value)
	result, err = value.Curry(Arguments{X: 1, Y: 2})
	fmt.Printf("Result: %s, %v\n", result, err)
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2))
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2, -3))
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2, -3, 5))
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2, -3, 5, -7))
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2, -3, 5, -7, 11))
	fmt.Printf("A polynomial expression: %s\n", utils.PolynomialExpression(X, 2, -3, 5, -7, 11, -13))
}