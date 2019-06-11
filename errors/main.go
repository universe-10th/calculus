package errors

import "errors"

var LogarithmOfNegative = errors.New("attempted to calculate logarithm of negative number")
var PowerOfNonPositive  = errors.New("attempted to calculate an invalid power (e.g. 0^0, x^y where x < 0)")
var DivisionByZero      = errors.New("attempted to calculate a division by 0")
var TangentOfVertical   = errors.New("attempted to calculate the tangent of a vertical angle")
