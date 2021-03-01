package examples

import (
	"fmt"
	compare "go-compare-expressions"
)

/**
Working example showing how to use this library
 */
func ExampleToCheckForDuplicate() {
	expr1 := "(foo == 1 && bar == 1) || baz == 0 "
	expr2 := "baz == 0 || (bar == 1 && foo == 1)"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}

/**
Since && has precendence over ||, note in the example that brackets are not necessary
*/
func ExampleToCompareExpressionWithNoBrackets() {
	expr1 := "foo == 1 && bar == 1 || baz == 0 "
	expr2 := "baz == 0 || bar == 1 && foo == 1"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}

/**
We can provide nested bracket as below
*/
func ExampleToCompareExpressionWithNestedBrackets() {
	expr1 := "(foo == 1 && (bar == 1 || baz == 0) || boo == 0"
	expr2 := "((baz == 0 || bar == 1) && foo == 1) || boo == 0"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}

/**
The below method throws error since the parameters in expr1 (foo) does not match with parameters in expr2 ( foo, boo)
*/
func ExampleToCompareExpressionWithMismatchParams() {
	expr1 := "foo == 1"
	expr2 := "foo == 1 || boo == 0"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}

/**
The below method throws error since the expression contains unsupported operators like and, or
*/
func ExampleToCompareExpressionWithUnsupportedOperators() {
	expr1 := "boo == 1 and foo == 0"
	expr2 := "foo == 1 or boo == 0"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}


/**
The below method throws error since the expression contains unsupported comparators like > , <
*/
func ExampleToCompareExpressionWithUnsupportedComparators() {
	expr1 := "boo > 1 || foo == 0"
	expr2 := "foo > 1 || boo == 0"
	result, err := compare.CheckIfDuplicateExpressions(expr1, expr2)
	if err != nil {
		fmt.Printf("Received error= %v", err.Error())
	} else {
		fmt.Printf("Expressions expr1=%s, expr2=%s are %v", expr1, expr2, result)
	}
}
