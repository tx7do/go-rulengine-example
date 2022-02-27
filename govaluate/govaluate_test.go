package govaluate

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvaluableExpression(t *testing.T) {
	// 支持多个逻辑表达式
	expr, err := govaluate.NewEvaluableExpression(
		"(10 > 0) && (2.1 == 2.1) && 'service is ok' == 'service is ok'" +
			" && 1 in (1,2) && 'code1' in ('code3','code2',1)")
	assert.Nil(t, err)

	result, err := expr.Evaluate(nil)
	assert.Nil(t, err)
	assert.False(t, result.(bool))

	fmt.Println(result)

	// 逻辑表达式包含变量
	expression, err := govaluate.NewEvaluableExpression("http_response_body == 'service is ok'")
	parameters := make(map[string]interface{}, 8)
	parameters["http_response_body"] = "service is ok"
	res, err := expression.Evaluate(parameters)
	assert.Nil(t, err)
	assert.True(t, res.(bool))
	fmt.Println(res)

	// 算数表达式包含变量
	expression1, err := govaluate.NewEvaluableExpression("requests_made * requests_succeeded / 100")
	assert.Nil(t, err)

	parameters1 := make(map[string]interface{}, 8)
	parameters1["requests_made"] = 100
	parameters1["requests_succeeded"] = 80
	result1, err := expression1.Evaluate(parameters1)
	assert.Nil(t, err)

	assert.Equal(t, result1, float64(80))
	fmt.Println(result1)
}

func TestNewEvaluableExpressionWithFunctions(t *testing.T) {
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return length, nil
		},
	}

	exprString := "strlen('teststring')"
	expr, _ := govaluate.NewEvaluableExpressionWithFunctions(exprString, functions)
	result, _ := expr.Evaluate(nil)
	assert.Equal(t, result, 10)
	fmt.Println(result)
}
