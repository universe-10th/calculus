package main

import (
	. "github.com/universe-10th/calculus/expressions"
	"fmt"
)


func Sample() Expression {
	return Mul(Add(
		X, Y, Negated(Z),
		Mul(X, Inverse(Y), Add(Y, Z), Inverse(Add(
			X, Y, Inverse(Z),
		))),
		Negated(Mul(Y, Z)),
	), Inverse(Mul(X, Y)), Pow(Add(Mul(X, Negated(Y)), Z), Add(X, Negated(Mul(Exp(Add(Y, Negated(Mul(Z, Num(34))))), Inverse(Z))))))
}


func main() {
	value := Sample()
	fmt.Printf("Currying expression %s with values (X=1, Y=2)...\n", value)
	result, err := value.Curry(Arguments{X: 1, Y: 2})
	fmt.Printf("Result: %s, %v\n", result, err)
}
