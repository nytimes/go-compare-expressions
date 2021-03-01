package compare_expressions

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateFormat(t *testing.T) {
	tables := []struct {
		name   string
		expr   string
		result []string
		err    error
	}{
		{"success_when_simple_expression", "a == 1", []string{"a"}, nil},
		{"success_when_2_params", "a == 1 && b == 1", []string{"a", "b"}, nil},
		{"success_when_3_params", "a == 1 && b == 1 || c == 1", []string{"a", "b", "c"}, nil},
		{"success_when_2_params_with_brackets", "(a == 1 && b == 1 )", []string{"a", "b"}, nil},
		{"success_when_3_params_with_brackets", "(a == 1) && (b == 1 || c == 1)", []string{"a", "b", "c"}, nil},


		{"error_when_invalid_format_equals", "a === 1", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},
		{"error_when_invalid_format_ones", "a == 11", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},
		{"error_when_invalid_format_zeroes", "a == 00", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},

		{"error_when_invalid_format_equals", "a == 1 && b ==== 1", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},
		{"error_when_invalid_format_ones", "a == 1 && b == 11", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},
		{"error_when_invalid_format_zeroes", "a == 00 && b == 0", nil, errors.New(fmt.Sprintf("Invalid expression, Required Format 'variable == 1 or 0'"))},

		{"error_when_invalid_format_ands", "a == 0 &&& b == 1", nil, errors.New(fmt.Sprintf("Invalid expression, Allowed combinators && or ||"))},
		{"error_when_invalid_format_ors", "a == 0 |||| b == 1", nil, errors.New(fmt.Sprintf("Invalid expression, Allowed combinators && or ||"))},
	}

	for _, table := range tables {
		result, err := ValidateFormat(table.expr)
		if table.err != nil && table.err.Error() != err.Error() {
			t.Errorf("validateFormat error failed in test %s, got: %v, want: %v.", table.name, err, table.err)
		} else if table.err == nil && !reflect.DeepEqual(result, table.result) {
			t.Errorf("validateFormat failed in test %s, got: %v, want: %v.", table.name, result, table.result)
		}
	}

}

func TestListContains(t *testing.T) {
	tables := []struct {
		name    string
		params1 []string
		params2 []string
		err     error
	}{
		{"success_when_one_params", []string{"a"}, []string{"a"}, nil},
		{"success_when_2_params", []string{"a", "b"}, []string{"a", "b"}, nil},
		{"success_when_3_params", []string{"a", "b", "c"}, []string{"a", "b", "c"}, nil},

		{"error_when_nil", []string{"a"}, nil, errors.New(fmt.Sprintf("expressions have different number of parameters, params1:%v , params2:%v", []string{"a"}, []string{}))},
		{"error_when_nil", nil, []string{"a", "b", "c"}, errors.New(fmt.Sprintf("expressions have different number of parameters, params1:%v , params2:%v", []string{}, []string{"a", "b", "c"}))},
		{"error_when_nil", []string{"a", "b", "d"}, []string{"a", "b", "c"}, errors.New(fmt.Sprintf("expressions have different parameters, params1:%v , params2:%v", []string{"a", "b", "d"}, []string{"a", "b", "c"}))},
	}

	for _, table := range tables {
		err := ListContains(table.params1, table.params2)
		if table.err != nil && table.err.Error() != err.Error() {
			t.Errorf("ListContains failed in test %s, got: %v, want: %v.", table.name, err, table.err)
		}
	}

}

func TestFilterDuplicates(t *testing.T) {
	tables := []struct {
		name   string
		input  []string
		result []string
	}{
		{"filter_duplicate_all_same", []string{"a", "a", "a"}, []string{"a"}},
		{"filter_duplicate_all_diff", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"filter_duplicate_same_in_begin", []string{"a", "a", "c"}, []string{"a", "c"}},
		{"filter_duplicate_same_in_end", []string{"a", "c", "c"}, []string{"a", "c"}},
	}

	for _, table := range tables {
		result := FilterDuplicates(table.input)
		if !reflect.DeepEqual(result, table.result) {
			t.Errorf("FilterDuplicates failed in test %s, got: %v, want: %v.", table.name, result, table.result)
		}
	}

}

func TestCheckIfDuplicateExpressions(t *testing.T) {
	tables := []struct {
		name   string
		expr1  string
		expr2  string
		result bool
		err    error
	}{
		{"success_when_one_params", "a == 1", "a == 1", true, nil},
		{"success_when_two_params", "a == 1 && b == 1", "b == 1 && a == 1", true, nil},
		{"success_when_three_params", "a == 1 && b == 1 && c == 0", "b == 1 && a == 1 && c == 0", true, nil},
		{"success_when_three_params", "a == 1 || b == 1 || c == 0", "b == 1 || a == 1 || c == 0", true, nil},
		{"success_when_three_params", "a == 1 || b == 1 && c == 0", "c == 0 &&  b == 1 || a == 1", true, nil},
		{"success_when_three_params", "(a == 1 || b == 1) && c == 0", "c == 0 &&  (b == 1 || a == 1)", true, nil},
		{"success_when_three_params", "a == 1 || (b == 1 && c == 0)", "(c == 0 &&  b == 1) || a == 1", true, nil},
		{"success_when_four_params", "(a == 1 || b == 1 ) && (c == 0 || a == 1)", "(a == 1 || c == 0 ) && (b == 1 || a == 1)", true, nil},

		{"error_when_four_params", "(a == 1 || b == 1 ) && (c == 0 || a == 1)", "(a == 0 || c == 0 ) && (b == 1 || a == 1)", false, nil},

		{"error_when_diff_params", "(a == 1 || b == 1 ) && (c == 0 || a == 1)", "(a == 0 ) && (b == 1 || a == 1)", false, errors.New(fmt.Sprintf("expressions have different number of parameters, params1:%v , params2:%v", []string{"a", "b", "c"}, []string{"a", "b"}))},
	}

	for _, table := range tables {
		result, err := CheckIfDuplicateExpressions(table.expr1, table.expr2)
		if table.err != nil && table.err.Error() != err.Error() {
			t.Errorf("CheckIfDuplicateExpressions error failed in test %s, got: %v, want: %v.", table.name, err, table.err)
		} else if table.err == nil && !reflect.DeepEqual(result, table.result) {
			t.Errorf("CheckIfDuplicateExpressions failed in test %s, got: %v, want: %v.", table.name, result, table.result)
		}
	}

}
