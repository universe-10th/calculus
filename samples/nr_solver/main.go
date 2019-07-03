package main

import (
	"math/big"
	. "github.com/universe-10th/calculus/expressions"
	"github.com/universe-10th/calculus/solvers/support"
	"github.com/universe-10th/calculus/solvers"
	"fmt"
)


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
	nrSample := SampleNR()
	guess := big.NewFloat(1000000.0)
	epsilon := support.Epsilon(10)
	fmt.Printf("Testing Newton-Raphson sample: %s with initial guess: %f\n", nrSample, guess)
	if nrResult, err := solvers.NewtonRaphson(nrSample, guess, epsilon, 10000, 10); err != nil {
		fmt.Printf("Error in Newton-Raphson: %s\n", err)
	} else {
		fmt.Printf("Returned value: %f\n", nrResult)
	}
}
