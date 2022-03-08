package goja

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	rulengine "go-rulengine-example"
	"testing"
)

func TestGoJa(t *testing.T) {
	const SCRIPT = `
	var hasX = false;
	var hasY = false;
	for (var key in o) {
		switch (key) {
		case "x":
			if (hasX) {
				throw "Already have x";
			}
			hasX = true;
			delete o.y;
			break;
		case "y":
			if (hasY) {
				throw "Already have y";
			}
			hasY = true;
			delete o.x;
			break;
		default:
			throw "Unexpected property: " + key;
		}
	}
	
	hasX && !hasY || hasY && !hasX;
	`
	r := goja.New()
	_ = r.Set("o", map[string]interface{}{
		"x": 40,
		"y": 2,
	})
	v, err := r.RunString(SCRIPT)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
}

type RuleConditionContext struct {
	NetAmount float32
	Distance  int32
	Duration  int32
	Result    bool
}

func TestUseStruct(t *testing.T) {
	const SCRIPT = `
if ((RuleConditionContext.Distance > 5000 && RuleConditionContext.Duration > 120) && (RuleConditionContext.Result == false)) {
	   RuleConditionContext.Result = true;
}
`
	ruleCondition := &RuleConditionContext{
		Distance: 6000,
		Duration: 121,
		Result:   false,
	}

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	nodeName := "rule1"

	apis := make(rulengine.PropertiesMap)
	apis["RuleConditionContext"] = ruleCondition

	err = eng.AddNode(nodeName, SCRIPT, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	fmt.Println(ruleCondition.Result)
}

type MyPoGo struct {
	Name string
}

func (p *MyPoGo) GetNameLength() int {
	return len(p.Name)
}

func (p *MyPoGo) AppendString(subString string) string {
	p.Name += subString
	return p.Name
}

func TestCycleCallRule(t *testing.T) {
	drl := `
while (Pogo.GetNameLength() < 100) {
  Pogo.AppendString(" Groooling");
}
`
	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	nodeName := "Pogo"

	pogo := &MyPoGo{Name: "bobo"}

	apis := make(rulengine.PropertiesMap)
	apis["Pogo"] = pogo

	err = eng.AddNode(nodeName, drl, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	fmt.Println(pogo.Name)
}
