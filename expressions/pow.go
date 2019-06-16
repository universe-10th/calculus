package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
	"fmt"
)


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


func (pow PowExpr) Evaluate(args Arguments) (sets.Number, error) {
	if base, err := pow.base.Evaluate(args); err != nil {
		return nil, err
	} else if exponent, err := pow.exponent.Evaluate(args); err != nil {
		return nil, err
	} else {
		return pow.wrappedPow(base, exponent)
	}
}


func (pow PowExpr) Derivative(wrt Variable) (Expression, error) {
	// Say we have f(x), g(x)
	// d[f(x)^g(x)]/dx = f(x)^(g(x)-1) * [g(x)*df(x)/dx + f(x)*ln(f(x))*dg(x)/dx]
	var err error
	var baseDerivative Expression
	var exponentDerivative Expression
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


func (pow PowExpr) CollectVariables(variables Variables) {
	pow.base.CollectVariables(variables)
	pow.exponent.CollectVariables(variables)
}


func (pow PowExpr) IsConstant(wrt Variable) bool {
	return pow.base.IsConstant(wrt) && pow.exponent.IsConstant(wrt)
}


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


func (ln LnExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := ln.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ln.wrappedLn(result)
	}
}


func (ln LnExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := ln.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Inverse(ln.arg), derivative).Simplify()
	}
}


func (ln LnExpr) Arguments() []Expression {
	return []Expression{ ln.arg }
}


func (ln LnExpr) CollectVariables(variables Variables) {
	ln.arg.CollectVariables(variables)
}


func (ln LnExpr) IsConstant(wrt Variable) bool {
	return ln.arg.IsConstant(wrt)
}


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


func (ln LnExpr) String() string {
	return FunctionDisplay(ln)
}


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


func (log LogExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := log.power.Evaluate(args); err != nil {
		return nil, err
	} else if result2, err2 := log.base.Evaluate(args); err2 != nil {
		return nil, err
	} else {
		return log.wrappedLn(result, result2)
	}
}


func (log LogExpr) Derivative(wrt Variable) (Expression, error) {
	return Mul(Ln(log.power), Inverse(Ln(log.base))).Derivative(wrt)
}


func (log LogExpr) Arguments() []Expression {
	return []Expression{ log.base, log.power }
}


func (log LogExpr) CollectVariables(variables Variables) {
	log.power.CollectVariables(variables)
	log.base.CollectVariables(variables)
}


func (log LogExpr) IsConstant(wrt Variable) bool {
	return log.base.IsConstant(wrt) && log.power.IsConstant(wrt)
}


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


func (log LogExpr) String() string {
	return FunctionDisplay(log)
}


type ExpExpr struct {
	FunctionExpr
	exponent Expression
}


func (exp ExpExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := exp.exponent.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Exp(exp.exponent), derivative).Simplify()
	}
}


func (exp ExpExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := exp.exponent.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Exp(result), nil
	}
}


func (exp ExpExpr) Arguments() []Expression {
	return []Expression{ exp.exponent }
}


func (exp ExpExpr) CollectVariables(variables Variables) {
	 exp.exponent.CollectVariables(variables)
}


func (exp ExpExpr) IsConstant(wrt Variable) bool {
	return exp.exponent.IsConstant(wrt)
}


func (exp ExpExpr) Simplify() (Expression, error) {
	if simplified, err := exp.exponent.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Exp(num.number)}, nil
	} else {
		return Exp(simplified), nil
	}
}


func (exp ExpExpr) String() string {
	return FunctionDisplay(exp)
}


func Pow(base, exponent Expression) Expression {
	return PowExpr{base, exponent}
}


func Ln(power Expression) Expression {
	return LnExpr{FunctionExpr{"ln"}, power}
}


func Log(base, power Expression) Expression {
	return LogExpr{FunctionExpr{"log"}, power, base}
}


func Exp(exponent Expression) Expression {
	return ExpExpr{FunctionExpr{"exp"}, exponent}
}