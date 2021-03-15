# go-compare-expressions

Provides support for comparing boolean/logical expressions. 
You can use this package to check if two given expressions are duplicate/equivalent

# Description
 Logically two expressions `P` and `Q`
 are considered duplicate/equivalent, if 
 - when P is true, Q is true as well and  when P is false, Q is false as well and vice versa
 
 This is done by generating truth tables for the given expression and comparing the truth tables. If both the generated truth tables are equal, we concur that given expressions are equal as well. 
## Installation
Download this package into your golang project.
 ```go
go get github.com/nytimes/go-compare-expressions
```

## Usage
 To compare any two boolean/logical expression, call `AreDuplicateExpressions` method by passing the two expression that you want to compare.
 ```go
    expr1 := "(foo == 1 && bar == 1) || baz = 0 " 
	expr2 := "baz == 1 || bar == 1 && foo == 1"
	result, err := AreDuplicateExpressions(expr1,expr2)
```

## Test
```go
go test .
```


## Examples 
Examples can be found [here](https://github.com/nytimes/go-compare-expressions/tree/master/examples) 

## Operations supported:
    
   - Comparators: `==`
   - Logical operators: `&&` `||`
   - Only Binary values allowed: `0` `1`
   - Parenthesis can be used to define precedence of operations `()`
  
## License
  This project is licensed under the MIT general use license.
   
   
   




