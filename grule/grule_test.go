package grule

import (
	"fmt"
	rulengine "go-rulengine-example"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyFact struct {
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}

func TestTutorial(t *testing.T) {
	const drl = `
rule CheckValues "Check the default values" salience 10 {
    when 
        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
		Retract("CheckValues");
}
`

	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	nodeName := "rule1"

	apis := make(rulengine.ExportMap)
	apis["MF"] = myFact

	err = eng.AddNode(nodeName, drl, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	assert.Equal(t, "Let say \"Hello Grule\"", myFact.WhatToSay)
	println(myFact.WhatToSay)
}

type RuleConditionContext struct {
	NetAmount float32
	Distance  int32
	Duration  int32
	Result    bool
}

func TestUseStruct(t *testing.T) {
	// DRL的规则
	const drl = `
rule  DuplicateRule1 "Duplicate Rule 1" salience 5 {
	when
		(RuleConditionContext.Distance > 5000 && RuleConditionContext.Duration > 120) && (RuleConditionContext.Result == false)
	Then
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

	apis := make(rulengine.ExportMap)
	apis["RuleConditionContext"] = ruleCondition

	err = eng.AddNode(nodeName, drl, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	fmt.Println(ruleCondition.Result)
}

type MyPoGo struct {
	Name string
}

func (p *MyPoGo) GetStringLength(sarg string) int {
	return len(sarg)
}

func (p *MyPoGo) AppendString(aString, subString string) string {
	return fmt.Sprintf("%s%s", aString, subString)
}

func TestCycleCallRule(t *testing.T) {
	drl := `
rule AgeNameCheck "test" {
	when
		Pogo.GetStringLength(Pogo.Name) < 100
	then
		Pogo.Name = Pogo.AppendString(Pogo.Name, "Groooling");
}
`
	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	nodeName := "Pogo"

	pogo := &MyPoGo{Name: "bobo"}

	apis := make(rulengine.ExportMap)
	apis["Pogo"] = pogo

	err = eng.AddNode(nodeName, drl, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	fmt.Println(pogo.Name)
}
