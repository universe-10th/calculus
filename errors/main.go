package errors

import "errors"

var ErrUndefinedValue = errors.New("one or more parameters are undefined")
var ErrLogarithmOfNegative = errors.New("attempted to calculate logarithm of negative number")
var ErrInvalidPowerOperation = errors.New("attempted to calculate an invalid power (e.g. 0^0, x^y where x < 0)")
var ErrDivisionByZero = errors.New("attempted to calculate a division by 0")
var ErrTangentOfVertical = errors.New("attempted to calculate the tangent of a vertical angle")
var ErrInvalidFactorialArgument = errors.New("argument to factorial does not fit in a 64bit unsigned integer")
// For derivative
var ErrNotDerivableExpression = errors.New("this expression is not derivable")
var ErrUndefinedOnInteger = errors.New("this expression is undefined on integer values")
var ErrInfiniteCannotBeRounded = errors.New("infinite numbers cannot be rounded")
