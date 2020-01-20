package main

import (
	. "github.com/universe-10th/calculus/v2/std/expressions"
	"fmt"
)


func Sample() Expression {
	return Mul(Pow(X, Num(3)), Ln(Mul(X, Num(4))), Tan(X))
}


func main() {
	value := Sample()
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
}