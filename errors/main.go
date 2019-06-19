package errors

import "errors"

var UndefinedValue           = errors.New("one or more parameters are undefined")
var LogarithmOfNegative      = errors.New("attempted to calculate logarithm of negative number")
var InvalidPowerOperation    = errors.New("attempted to calculate an invalid power (e.g. 0^0, x^y where x < 0)")
var DivisionByZero           = errors.New("attempted to calculate a division by 0")
var TangentOfVertical        = errors.New("attempted to calculate the tangent of a vertical angle")
var InvalidFactorialArgument = errors.New("argument to factorial does not fit in a 64bit unsigned integer")
// For derivative
var NotDerivableExpression   = errors.New("this expression is not derivable")
// For approximation methods
var IterationsExhausted      = errors.New("convergence did not occur at the specified number of iterations")