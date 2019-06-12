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
	// TODO
	// Say we have f(x), g(x)
	// d[f(x)^g(x)]/dx = f(x)^(g(x)-1) * [g(x)*df(x)/dx + f(x)*ln(f(x))*dg(x)/dx]
	// TODO please note: we could use CollectVariables to tell whether
	// TODO wrt is included in either the base or exponent expressions and
	// TODO generate the special cases like f(x)^n, or a^f(x).
	return nil, nil
}


func (pow PowExpr) CollectVariables(variables Variables) {
	pow.base.CollectVariables(variables)
	pow.exponent.CollectVariables(variables)
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
	// TODO
	return nil, nil
}


func (ln LnExpr) Arguments() []Expression {
	return []Expression{ ln.arg }
}


func (ln LnExpr) CollectVariables(variables Variables) {
	ln.arg.CollectVariables(variables)
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
	// TODO
	return nil, nil
}


func (log LogExpr) Arguments() []Expression {
	return []Expression{ log.base, log.power }
}


func (log LogExpr) CollectVariables(variables Variables) {
	log.power.CollectVariables(variables)
	log.base.CollectVariables(variables)
}


type ExpExpr struct {
	FunctionExpr
	exponent Expression
}


func (exp ExpExpr) Derivative(wrt Variable) (Expression, error) {
	// TODO
	return nil, nil
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