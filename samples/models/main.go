package main

import (
	. "github.com/universe-10th/calculus/v2/expressions"
	"fmt"
	"github.com/universe-10th/calculus/v2/models"
)


func SampleModel() *models.Model {
	// Z = X/Y + W, and corollaries
	model := models.NewModel()
	flowZ, _ := models.NewSingleOutputModelFlow(Z, Add(W, Div(X, Y)))
	model.AddFlow(flowZ)
	flowW, _ := models.NewSingleOutputModelFlow(W, Sub(Z, Div(X, Y)))
	model.AddFlow(flowW)
	flowX, _ := models.NewSingleOutputModelFlow(X, Mul(Y, Sub(Z, W)))
	model.AddFlow(flowX)
	flowY, _ := models.NewSingleOutputModelFlow(Y, Div(X, Sub(Z, W)))
	model.AddFlow(flowY)
	return model
}


func TestModel(model *models.Model, arguments Arguments) {
	result, err := model.Evaluate(arguments)
	fmt.Println("Testing model with arguments:", arguments, "returns", result, err)
}


func main() {
	model := SampleModel()
	TestModel(model, Arguments{X: 6, Y: 3, W: 2}.Wrap())
	TestModel(model, Arguments{X: 6, Y: 3, Z: 4}.Wrap())
	TestModel(model, Arguments{X: 6, Z: 4, W: 2}.Wrap())
	TestModel(model, Arguments{Z: 4, Y: 3, W: 2}.Wrap())
}