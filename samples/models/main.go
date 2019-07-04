package main

import (
	"github.com/universe-10th/calculus/utils"
	. "github.com/universe-10th/calculus/expressions"
	"fmt"
	"github.com/universe-10th/calculus/solvers"
	"math/big"
	"github.com/universe-10th/calculus/solvers/support"
)


func SampleModel() *utils.Model {
	// Z = X/Y + W
	model, _ := utils.NewModel(Z, Add(Mul(X, Inverse(Y)), W))
	// W = Z - X/Y
	fmt.Println(model.SetCorollary(W, Add(Z, Negated(Mul(X, Inverse(Y))))))
	// X (by solver) will be: Y*(Z - W)
	// Y (by solver) will be: X/(Z - W)
	return model
}


func main() {
	model := SampleModel()
	var result interface{}
	var err    error
	result, err = model.Evaluate(Arguments{
		X: 1, Y: 2, W: 3,
	}.Wrap(), nil)
	fmt.Printf("Return after evaluating Z=Model(X: 1, Y: 2, W: 3): %v %v\n", result, err)
	result, err = model.Evaluate(Arguments{
		X: 1, Y: 2, Z: 3,
	}.Wrap(), nil)
	fmt.Printf("Return after evaluating W=Model^-1(X: 1, Y: 2, Z: 3): %v %v\n", result, err)
	nrSolver := solvers.MakeNewtonRaphsonSolver(big.NewFloat(0).SetInt64(0), support.Epsilon(10), 100, 100)
	result, err = model.Evaluate(Arguments{
		W: 1, Y: 2, Z: 3,
	}.Wrap(), nrSolver)
	fmt.Printf("Return after evaluating W=Model^-1(W: 1, Y: 2, Z: 3): %v %v\n", result, err)
	result, err = model.Evaluate(Arguments{
		X: 1, W: 2, Z: 3,
	}.Wrap(), nrSolver)
	fmt.Printf("Return after evaluating W=Model^-1(X: 1, W: 2, Z: 3): %v %v\n", result, err)
}