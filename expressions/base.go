package expressions


func areAllConstants(expressions ...Expression) bool {
	variables := Variables{}
	for _, expression := range expressions {
		expression.CollectVariables(variables)
	}
	return len(variables) == 0
}
