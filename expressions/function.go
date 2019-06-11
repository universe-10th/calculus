package diff


type Function interface {
	StandardName() string
}


type FunctionExpr struct {
	standardName string
}


func (function FunctionExpr) StandardName() string {
	return function.standardName
}