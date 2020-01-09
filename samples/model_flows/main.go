package main

import (
	. "github.com/universe-10th/calculus/expressions"
	. "github.com/universe-10th/calculus/expressions/goal-seek"
	"github.com/universe-10th/calculus/models"
	"fmt"
	"github.com/universe-10th/calculus/models/implementations"
	"math/big"
	diffUtils "github.com/universe-10th/calculus/core/support/diff"
)

func main() {
	A := Var("A")
	B := Var("B")
	C := Var("C")
	linearExpr := Add(B, Mul(X, A))
	quadraticExpr := Add(C, Mul(B, X), Mul(A, Pow(X, Num(2))))
	exponentialExpr := Exp(Y)
	invertedLinear := NRGoalSeek(Y, linearExpr, X, func(arguments Arguments) (initialGuess, epsilon *big.Float, maxIterations, maxArgCorrectionsPerIteration uint32) {
		initialGuess = big.NewFloat(0.0)
		epsilon = diffUtils.Epsilon(70)
		maxIterations = 100
		maxArgCorrectionsPerIteration = 100
		return
	})

	linearModel, _ := models.NewSingleOutputModelFlow(Y, linearExpr)
	quadraticModel, _ := models.NewSingleOutputModelFlow(Z, quadraticExpr)
	exponentialModel, _ := models.NewSingleOutputModelFlow(W, exponentialExpr)
	serialModel, _ := models.NewSerialModelFlow(linearModel, exponentialModel)
	parallelModel, _ := models.NewParallelModelFlow(linearModel, quadraticModel)
	invertedLinearModel, _ := models.NewSingleOutputModelFlow(X, invertedLinear)
	numberSplitModel, _ := implementations.NewNumberSplitModelFlow(Z, Y, X)

	fmt.Println("Expected input for linear model:", linearModel.Input())
	fmt.Println("Expected input for quadratic model:", quadraticModel.Input())
	fmt.Println("Expected input for exponential model:", exponentialModel.Input())
	fmt.Println("Expected input for serial model:", serialModel.Input())
	fmt.Println("Expected input for parallel model:", parallelModel.Input())
	fmt.Println("Expected input for inverted model:", invertedLinearModel.Input())
	fmt.Println("Expected input for split model:", numberSplitModel.Input())

	result, err := linearModel.Evaluate(Arguments{
		A: 2, B: 3, X: 5,
	}.Wrap())
	fmt.Println("Evaluating linear model:", result, err)

	result, err = quadraticModel.Evaluate(Arguments{
		A: 2, B: 3, C: 5, X: 7,
	}.Wrap())
	fmt.Println("Evaluating quadratic model:", result, err)

	result, err = exponentialModel.Evaluate(Arguments{
		Y: 2,
	}.Wrap())
	fmt.Println("Evaluating exponential model:", result, err)

	result, err = serialModel.Evaluate(Arguments{
		A: 2, B: 3, X: 5,
	}.Wrap())
	fmt.Println("Evaluating serial model:", result, err)

	result, err = parallelModel.Evaluate(Arguments{
		A: 2, B: 3, C: 5, X: 7,
	}.Wrap())
	fmt.Println("Evaluating parallel model:", result, err)

	result, err = invertedLinearModel.Evaluate(Arguments{
		A: 2, B: 3, Y: 5,
	}.Wrap())
	fmt.Println("Evaluating inverted linear model:", result, err)

	result, err = numberSplitModel.Evaluate(Arguments{
		Z: 13.78,
	}.Wrap())
	fmt.Println("Evaluating number split model:", result, err)
}
