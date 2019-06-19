package newton_raphson

import "errors"


var CannotIterateOverZeroDerivative = errors.New("cannot iterate when the derivative is zero")
