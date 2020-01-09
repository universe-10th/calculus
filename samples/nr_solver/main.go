package main

import (
	"math/big"
	. "github.com/universe-10th/calculus/expressions"
	"fmt"
	"github.com/universe-10th/calculus/utils"
	"github.com/universe-10th/calculus/goals"
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
	epsilon := utils.Epsilon(10)
	fmt.Printf("Testing Newton-Raphson sample: %s with initial guess: %f\n", nrSample, guess)
	if nrResult, err := goals.NewtonRaphson(nrSample, guess, epsilon, 10000, 10); err != nil {
		fmt.Printf("Error in Newton-Raphson: %s\n", err)
	} else {
		fmt.Printf("Returned value: %f\n", nrResult)
	}
}
