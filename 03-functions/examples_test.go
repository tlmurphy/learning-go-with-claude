package functions

import "testing"

func TestDemoMultipleReturns(t *testing.T) {
	DemoMultipleReturns()
}

func TestDemoVariadic(t *testing.T) {
	DemoVariadic()
}

func TestDemoFunctionsAsValues(t *testing.T) {
	DemoFunctionsAsValues()
}

func TestDemoClosures(t *testing.T) {
	DemoClosures()
}

func TestDemoDefer(t *testing.T) {
	DemoDefer()
}

func TestDemoDeferArgEvaluation(t *testing.T) {
	DemoDeferArgEvaluation()
}

func TestDemoFunctionComposition(t *testing.T) {
	DemoFunctionComposition()
}
