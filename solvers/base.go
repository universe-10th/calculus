package solvers

import (
	"github.com/universe-10th/calculus/expressions"
	"math/big"
)

// Solver is just a type of function that evaluates on an expression.
// Different algorithms (e.g. Brent, Newton-Raphson) will provide means to
// build the parametrized solvers on their own.
//
// Solvers always take an expression and try to make them evaluate to zero.
// Typically, you'll provide expressions with the intention to evaluate them
// to zero in order to find their root(s).
//
// Solvers always work with float values, since they are approximate methods,
// so they will always return a *big.Float and an error.
type Solver func(expressions.Expression) (result *big.Float, exception error)
