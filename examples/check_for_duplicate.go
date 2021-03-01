package examples

import (
	"fmt"
	compare "go-compare-expressions"
)

func ExampleToCheckForDuplicate() {
	expr1 := "(foo == 1 && bar == 1) || baz == 0 "
	expr2 := "baz == 0 || bar == 1 && foo == 1"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}
