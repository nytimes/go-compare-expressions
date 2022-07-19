package compare_expressions

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

/**
This is the main function that is to be called to verify if 2 expressions are duplicate or not.

Usage:
parameters:  expr1 a == 1 && b == 1 ,   expr1 b == 1 && a == 1
return (true, nil)

Note that the expressions provided should be restricted to have values mentioned below
Currently provides support for
operators                           ----  && ||
comparators                         ----  ==
values or right side of expressions ----  to be binary only. So it can be 1 or 0
attribute or left side of expression ---- can be any valid string name
*/
func CheckIfDuplicateExpressions(expr1 string, expr2 string) (bool, error) {
	parameters, err := ValidateInput(expr1, expr2)
	if err != nil {
		fmt.Errorf("Error validating the input %v ", err.Error())
		return false, err
	}
	var count []bool
	combinedExpression := "(" + expr1 + ") == (" + expr2 + ")"
	parametersMap := make(map[string]interface{})
	err = GenerateTruthTable(combinedExpression, parameters, parametersMap, 0, &count)
	if err != nil {
		fmt.Errorf("Unable to generate truth table %v ", err.Error())
		return false, err
	}
	countOfMatches := 0
	for i := range count {
		if count[i] {
			countOfMatches++
		}
	}
	if countOfMatches == len(count) {
		return true, nil
	}
	return false, nil
}

/**
This function is to validate the input in the following way.
1. It checks if the given expressions have valid format
2. It checks if the given expressions have same number of parameters
Example:
parameters:  expr1 a == 1 && b == 1,expr1 b == 1 && a == 1   return: ["a","b"], nil
parameters:  expr1 a == 1 && b == 1,expr1 b == 1 && a == 1 and c == 1  "expressions have different number of parameters"
*/
func ValidateInput(expr1, expr2 string) ([]string, error) {
	params1, err := ValidateFormat(expr1)
	if err != nil {
		return nil, err
	}
	params2, err := ValidateFormat(expr2)
	if err != nil {
		return nil, err
	}
	filteredList1 := FilterDuplicates(params1)
	filteredList2 := FilterDuplicates(params2)
	if err := ListContains(filteredList1, filteredList2); err != nil {
		fmt.Printf("Error comparing lists")
		return nil, err
	}
	if err := ListContains(filteredList2, filteredList1); err != nil {
		fmt.Printf("Error comparing lists")
		return nil, err
	}

	var finalParams []string
	for _, p := range params1 {
		finalParams = append(finalParams, strings.ReplaceAll(p, ".", "_"))
	}

	return finalParams, nil
}

/**
This function filters the duplicates from given array of string
*/
func FilterDuplicates(params []string) []string {
	parametersMap := make(map[string]interface{})
	list := []string{}
	for _, entry := range params {
		if _, value := parametersMap[entry]; !value {
			parametersMap[entry] = true
			list = append(list, entry)
		}
	}
	return list

}

/**
This function validates the given input expression. It uses regular expressions to validate the format.
Expr should be something like "a == 1" and
Currently provides support for
operators                           ----  && ||
comparators                         ----  ==
values or right side of expressions ----  to be binary only. So it can be 1 or 0
attribute or left side of expression ---- can be any valid string name
*/
func ValidateFormat(expr string) ([]string, error) {
	regex := regexp.MustCompile(`\s*[!=><]{2}?\s*[\d]`)
	replaceExpr := regex.ReplaceAllString(expr, " ")
	result := strings.Fields(replaceExpr)

	invalidRegex := regexp.MustCompile(`\s+[!=><\d]{1}\s*`)
	invalidExpr := invalidRegex.MatchString(replaceExpr)
	if invalidExpr {
		return nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))
	}

	regex = regexp.MustCompile(`\s+&{2}?\s+`)
	replaceExpr = regex.ReplaceAllString(replaceExpr, " ")
	result = strings.Fields(replaceExpr)

	regex = regexp.MustCompile(`\s+\|{2}?\s+`)
	replaceExpr = regex.ReplaceAllString(replaceExpr, " ")
	result = strings.Fields(replaceExpr)

	andsFound := strings.Contains(replaceExpr, "&")
	orsFound := strings.Contains(replaceExpr, "|")
	if andsFound || orsFound {
		return nil, errors.New(fmt.Sprintf("Invalid expression, Allowed combinators && or ||"))
	}

	regex = regexp.MustCompile(`[(|)]`)
	replaceExpr = regex.ReplaceAllString(replaceExpr, " ")
	result = strings.Fields(replaceExpr)

	regex = regexp.MustCompile(`[\D+]`)
	emptyExpression := regex.ReplaceAllString(replaceExpr, " ")
	if emptyExpression != "" {
		errors.New(fmt.Sprintf("Invalid expression, invalid characters=%v ", emptyExpression))
	}
	return result, nil
}

/**
This function checks whether all strings present in params2 is contained in params1
*/
func ListContains(params1, params2 []string) error {
	if len(params1) != len(params2) {
		return errors.New(fmt.Sprintf("expressions have different number of parameters, params1:%v , params2:%v", params1, params2))
	}
	count := 0
	for _, v := range params1 {
		for _, w := range params2 {
			if v == w {
				count++
			}
		}
	}
	if count != len(params2) {
		return errors.New(fmt.Sprintf("expressions have different parameters, params1:%v , params2:%v", params1, params2))
	}

	return nil
}

/**
This function generates the truth table for the given expression and set of parameters present in the expression.
It used tail recursion to evaluate expression for all possible values to the set of parameters.
*/
func GenerateTruthTable(expr string, parameters []string, parametersMap map[string]interface{}, index int, count *[]bool) error {
	separatRegex := regexp.MustCompile(`[.]`)
	replexpr := separatRegex.ReplaceAllString(expr, "_")
	if index == len(parameters) {
		result, err := EvaluateExpression(replexpr, parametersMap)
		if err != nil {
			fmt.Errorf("Unable to  evaluate expression, error: %v", err.Error())
			return err
		} else {
			if result.(bool) {
				*count = append(*count, true)
			} else {
				*count = append(*count, false)
			}
			return nil
		}
	}
	parametersMap[parameters[index]] = 1
	if err := GenerateTruthTable(expr, parameters, parametersMap, index+1, count); err != nil {
		return err
	}
	parametersMap[parameters[index]] = 0
	if err := GenerateTruthTable(expr, parameters, parametersMap, index+1, count); err != nil {
		return err
	}
	return nil
}

/**
This function uses govaluate library to evaluate the provided boolean expression
*/
func EvaluateExpression(expr string, parameters map[string]interface{}) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return nil, errors.New("Unable to initialize the expression. error:" + err.Error())
	}
	result, err := expression.Evaluate(parameters)
	if err != nil {
		return nil, errors.New("Unable to evaluate the expression. error:" + err.Error())
	}
	return result, nil
}
