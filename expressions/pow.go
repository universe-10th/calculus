package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
	"fmt"
)


// PowExpr stands both for roots and for powers. It will have a base and an exponent inner terms.
type PowExpr struct {
	base     Expression
	exponent Expression
}


func (pow PowExpr) wrappedPow(base, exponent sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.InvalidPowerOperation
		}
	}()
	result = ops.Pow(base, exponent)
	return
}


// Curry tries currying the underlying base and exponent expressions, and then generating a final expression.
// A simplification will also be attempted here if the curried expressions result to be constant.
func (pow PowExpr) Curry(args Arguments) (Expression, error) {
	var curriedBase, curriedExponent Expression
	var err error
	if curriedBase, err = pow.base.Curry(args); err != nil {
		return nil, err
	}
	if curriedExponent, err = pow.exponent.Curry(args); err != nil {
		return nil, err
	}
	return Pow(curriedBase, curriedExponent).Simplify()
}


// Evaluate computes the power of the expression considering base and exponent.
// Base expression and exponent expression are first evaluated.
func (pow PowExpr) Evaluate(args Arguments) (sets.Number, error) {
	if base, err := pow.base.Evaluate(args); err != nil {
		return nil, err
	} else if exponent, err := pow.exponent.Evaluate(args); err != nil {
		return nil, err
	} else {
		return pow.wrappedPow(base, exponent)
	}
}


func (pow PowExpr) derivativeBySpecialCases(
	simplifiedBase, simplifiedExponent Expression, wrt Variable,
) (Expression, error) {
	var err error
	var baseDerivative Expression
	var exponentDerivative Expression
	var isSBConstant bool
	var isSEConstant bool
	isSBConstant = simplifiedBase.IsConstant(wrt)
	isSEConstant = simplifiedExponent.IsConstant(wrt)

	if isSBConstant {
		if isSEConstant {
			return Num(0), nil
		} else {
			if exponentDerivative, err = pow.exponent.Derivative(wrt); err != nil {
				return nil, err
			}
			baseFactor, _ := Ln(simplifiedBase).Simplify()
			return Mul(Pow(simplifiedBase, simplifiedExponent), baseFactor, exponentDerivative), nil
			// a^[f(x)]*ln(a)*[df(x)/dx]
		}
	} else {
		if isSEConstant {
			if baseDerivative, err = pow.base.Derivative(wrt); err != nil {
				return nil, err
			}
			newExponent, _ := Add(simplifiedExponent, Num(-1)).Simplify()
			return Mul(
				simplifiedExponent,
				Pow(simplifiedBase, newExponent),
				baseDerivative,
			), nil
		} else {
			if baseDerivative, err = pow.base.Derivative(wrt); err != nil {
				return nil, err
			}
			if exponentDerivative, err = pow.exponent.Derivative(wrt); err != nil {
				return nil, err
			}
			first := Pow(pow.base, Add(pow.exponent, Num(-1)))
			second := Add(Mul(pow.exponent, baseDerivative), Mul(pow.base, Ln(pow.base), exponentDerivative))
			return Mul(first, second).Simplify()
		}
	}
}


// Derivative uses one of the power rules of the derivative: f(x)^a, a^f(x), f(x)^g(x).
// It also applies the chain rule appropriately over both functions.
func (pow PowExpr) Derivative(wrt Variable) (Expression, error) {
	// Say we have f(x), g(x)
	// d[f(x)^g(x)]/dx = f(x)^(g(x)-1) * [g(x)*df(x)/dx + f(x)*ln(f(x))*dg(x)/dx]
	var err error
	var simplifiedBase Expression
	var simplifiedExponent Expression
	if simplifiedBase, err = pow.base.Simplify(); err != nil {
		return nil, err
	}
	if simplifiedExponent, err = pow.exponent.Simplify(); err != nil {
		return nil, err
	}
	return pow.derivativeBySpecialCases(simplifiedBase, simplifiedExponent, wrt)
}


// CollectVariables digs into base and exponent expressions.
func (pow PowExpr) CollectVariables(variables Variables) {
	pow.base.CollectVariables(variables)
	pow.exponent.CollectVariables(variables)
}


// IsConstant returns whether both base and exponent are constant with respect to the given value.
func (pow PowExpr) IsConstant(wrt Variable) bool {
	return pow.base.IsConstant(wrt) && pow.exponent.IsConstant(wrt)
}


// Simplify attempts a constant simplification over both base and exponent.
// If both operands are constant, then it calculates the power into a new constant.
func (pow PowExpr) Simplify() (Expression, error) {
	if simplifiedBase, err := pow.base.Simplify(); err != nil {
		return nil, err
	} else if simplifiedExponent, err := pow.exponent.Simplify(); err != nil {
		return nil, err
	} else {
		// if both are constants, calculate.
		// otherwise, make new expression.
		simplifiedBaseNum, okBase := simplifiedBase.(Constant)
		simplifiedExponentNum, okPower := simplifiedExponent.(Constant)
		if okBase && okPower {
			if result, err := pow.wrappedPow(simplifiedBaseNum.number, simplifiedExponentNum.number); err != nil {
				return nil, err
			} else {
				return Constant{result}, nil
			}
		} else {
			return Pow(simplifiedBase, simplifiedExponent), nil
		}
	}
}


// String represents the power appropriately: (X)^(Y), perhaps removing parentheses appropriately.
func (pow PowExpr) String() string {
	baseStr := pow.base.String()
	if _, ok := pow.base.(SelfContained); !ok {
		baseStr = "(" + baseStr + ")"
	}
	exponentStr := pow.exponent.String()
	if _, ok := pow.exponent.(SelfContained); !ok {
		if _, ok := pow.exponent.(NegatedExpr); !ok {
			exponentStr = "(" + exponentStr + ")"
		}
	}
	return fmt.Sprintf("%s^%s", baseStr, exponentStr)
}


// LnExpr stands for a natural logarithm expression.
type LnExpr struct {
	FunctionExpr
	arg   Expression
}


func (ln LnExpr) wrappedLn(power sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.LogarithmOfNegative
		}
	}()
	result = ops.Ln(power)
	return
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (ln LnExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := ln.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Ln(curried).Simplify()
	}
}


// Evaluate returns the natural logarithm of the evaluated inner expression's value.
// It returns an error if the inner value is negative.
func (ln LnExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := ln.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ln.wrappedLn(result)
	}
}


// Derivative uses the rule of the natural logarithm and also applies chain rule.
func (ln LnExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := ln.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Inverse(ln.arg), derivative).Simplify()
	}
}


// Arguments returns a list of expression just being the inner argument.
func (ln LnExpr) Arguments() []Expression {
	return []Expression{ ln.arg }
}


// CollectVariables digs into the inner expression.
func (ln LnExpr) CollectVariables(variables Variables) {
	ln.arg.CollectVariables(variables)
}


// IsConstant returns whether the inner expression is constant with respect to the given variable.
func (ln LnExpr) IsConstant(wrt Variable) bool {
	return ln.arg.IsConstant(wrt)
}


// Simplify attempts a constant simplification of the inner value.
// If the simplified inner expression is constant, it attempts a calculation of its natural logarithm.
// It will be an error if the inner simplified value is negative.
func (ln LnExpr) Simplify() (Expression, error) {
	if simplified, err := ln.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := ln.wrappedLn(num.number); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return Ln(simplified), nil
	}
}


// String represents the natural logarithm as ln(X).
func (ln LnExpr) String() string {
	return FunctionDisplay(ln)
}


// LogExpr stands for a logarithm of any base.
type LogExpr struct {
	FunctionExpr
	power   Expression
	base    Expression
}


func (log LogExpr) wrappedLn(power, base sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.LogarithmOfNegative
		}
	}()
	result = ops.Log(power, base)
	return
}


// Curry tries currying the underlying base and power expressions, and then generating a final expression.
// A simplification will also be attempted here if the curried expressions result to be constant.
func (log LogExpr) Curry(args Arguments) (Expression, error) {
	var curriedBase, curriedPower Expression
	var err error
	if curriedBase, err = log.base.Curry(args); err != nil {
		return nil, err
	}
	if curriedPower, err = log.power.Curry(args); err != nil {
		return nil, err
	}
	return Log(curriedBase, curriedPower).Simplify()
}


// Evaluate computes the logarithm after evaluating both the base and the power expressions.
// It is an error if the power (or the base) evaluates a negative value.
func (log LogExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := log.power.Evaluate(args); err != nil {
		return nil, err
	} else if result2, err2 := log.base.Evaluate(args); err2 != nil {
		return nil, err
	} else {
		return log.wrappedLn(result, result2)
	}
}


// Derivative uses the generic logarithm rule of the derivative.
// It first converts the Log(a, b) into Ln(b)/Ln(a) and then computes the derivative.
func (log LogExpr) Derivative(wrt Variable) (Expression, error) {
	return Mul(Ln(log.power), Inverse(Ln(log.base))).Derivative(wrt)
}


// Arguments returns (base, exponent) as the Log function arguments.
func (log LogExpr) Arguments() []Expression {
	return []Expression{ log.base, log.power }
}


// CollectVariables digs into the power and the base.
func (log LogExpr) CollectVariables(variables Variables) {
	log.power.CollectVariables(variables)
	log.base.CollectVariables(variables)
}


// IsConstant returns whether both power and base are constant with respect to a given value.
func (log LogExpr) IsConstant(wrt Variable) bool {
	return log.base.IsConstant(wrt) && log.power.IsConstant(wrt)
}


// Simplify attempts simplifying both base and power expressions, and then simplifies the logarithm.
// If both base and exponent simplify into constants, then it attempts a logarithm calculation, which
// may fail if either of the constants is negative.
func (log LogExpr) Simplify() (Expression, error) {
	if simplifiedBase, err := log.base.Simplify(); err != nil {
		return nil, err
	} else if simplifiedPower, err := log.power.Simplify(); err != nil {
		return nil, err
	} else {
		// if both are constants, calculate.
		// otherwise, make new expression.
		simplifiedBaseNum, okBase := simplifiedBase.(Constant)
		simplifiedPowerNum, okPower := simplifiedPower.(Constant)
		if okBase && okPower {
			if result, err := log.wrappedLn(simplifiedPowerNum.number, simplifiedBaseNum.number); err != nil {
				return nil, err
			} else {
				return Constant{result}, nil
			}
		} else {
			return Log(simplifiedBase, simplifiedPower), nil
		}
	}
}


// String represents the logarithm as log(a, b).
func (log LogExpr) String() string {
	return FunctionDisplay(log)
}


// ExpExpr stands for an expression like e^X.
type ExpExpr struct {
	FunctionExpr
	exponent Expression
}


// Derivative computes the rule of e^X, also applying chain rule.
func (exp ExpExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := exp.exponent.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Exp(exp.exponent), derivative).Simplify()
	}
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (exp ExpExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := exp.exponent.Curry(args); err != nil {
		return nil, err
	} else {
		return Exp(curried).Simplify()
	}
}


// Evaluate computes the value of the inner expression (the exponent) and then computes e^(that value).
func (exp ExpExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := exp.exponent.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Exp(result), nil
	}
}


// Arguments returns a list of expressions only including the exponent.
func (exp ExpExpr) Arguments() []Expression {
	return []Expression{ exp.exponent }
}


// CollectVariables dig into the exponent expression.
func (exp ExpExpr) CollectVariables(variables Variables) {
	 exp.exponent.CollectVariables(variables)
}


// IsConstant returns whether the exponent expression is constant with respect to the given variable.
func (exp ExpExpr) IsConstant(wrt Variable) bool {
	return exp.exponent.IsConstant(wrt)
}


// Simplify attempts a constant evaluation of the exponent, and then of e^(exponent).
func (exp ExpExpr) Simplify() (Expression, error) {
	if simplified, err := exp.exponent.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Exp(num.number)}, nil
	} else {
		return Exp(simplified), nil
	}
}


// String represents this expression as exp(X).
func (exp ExpExpr) String() string {
	return FunctionDisplay(exp)
}


// Pow constructs a power expression with base and exponent.
func Pow(base, exponent Expression) Expression {
	return PowExpr{base, exponent}
}


// Ln constructs a natural logarithm operation expression.
func Ln(power Expression) Expression {
	return LnExpr{FunctionExpr{"ln"}, power}
}


// Log constructs a logarithm expression with base and power.
func Log(base, power Expression) Expression {
	return LogExpr{FunctionExpr{"log"}, power, base}
}


// Exp constructs an e^X expression.
func Exp(exponent Expression) Expression {
	return ExpExpr{FunctionExpr{"exp"}, exponent}
}