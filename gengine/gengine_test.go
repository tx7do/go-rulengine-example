package gengine

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/google/martian/log"
	"github.com/stretchr/testify/assert"
	rulengine "go-rulengine-example"
	"go-rulengine-example/model"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestGEngine(t *testing.T) {
	/**
	'@id',获取
	'@name',获取规则名
	'@desc',获取规则说明
	*/
	const dsl = `
rule "测试规则名称1" "rule desc"
begin
  aName = @name
  aId = @id
  aDesc = @desc
  Println("@id: ", aId)
  Println("@name: ", aName)
  Println("@desc: ", aDesc)
end

rule "测试规则名称2" "rule desc"
begin
  aName = @name
  aId = @id
  aDesc = @desc
  Println("@id: ", aId)
  Println("@name: ", aName)
  Println("@desc: ", aDesc)
end
`

	start1 := time.Now().UnixNano()

	// 创建上下文
	dataContext := context.NewDataContext()
	dataContext.Add("Println", fmt.Println)
	// 初始化规则构造器
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	// 从字符串中解析规则
	err := ruleBuilder.BuildRuleFromString(dsl)

	end1 := time.Now().UnixNano()

	fmt.Printf("rules num:%d, load rules cost time:%d ns\n", len(ruleBuilder.Kc.RuleEntities), end1-start1)

	if err != nil {
		log.Errorf("err:%s ", err)
	} else {
		eng := engine.NewGengine()

		start := time.Now().UnixNano()
		// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
		err := eng.Execute(ruleBuilder, true)
		end := time.Now().UnixNano()
		if err != nil {
			log.Errorf("execute rule error: %v", err)
		}
		log.Infof("execute rule cost %d ns", end-start)
	}
}

func TestSortModel(t *testing.T) {
	// 按排序的顺序执行
	// 没有定义salience,按照规则定义的先后顺序执行.
	// salience	定义规则优先级的整数，数值越大，优先级越高
	const dsl = `
rule "rule_1" "highest priority"
begin
	println("rule_1-->")
	println("cal.Data-->", cal.Data)
	println("cal.Name-->", cal.Name) 
end

rule "rule_2" "mid priority"
begin
	println("rule_2-->")
	cal.Name = "hello world"
end

rule "rule_3" "lowest priority"
begin
	println("rule_3-->")
	cal.Data = 5
end
`

	calculate := &model.Calculate{Data: 0}

	properties := make(rulengine.PropertiesMap)
	properties["cal"] = calculate

	functions := make(rulengine.FunctionsMap)
	functions["println"] = fmt.Println

	nodeName := "SortModel"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, dsl, functions, properties, rulengine.SortModel, false)
	assert.Nil(t, err)

	err = eng.Execute(nodeName)
	assert.Nil(t, err)
}

func TestMixModel(t *testing.T) {
	// 找出优先级最高的优先第一个执行,其他的并发.
	const dsl = `
rule "rule_1" "highest priority 第一个"  salience 1000
begin
	println("rule_1-->")
	println("cal.Data-->", cal.Data)
	println("cal.Name-->", cal.Name) 
end

rule "rule_2" "mid priority 并发" salience 5
begin
	println("rule_2-->")
	cal.Name = "hello world"
end

rule "rule_3" "most priority 并发" salience 2
begin
	println("rule_3-->")
	cal.Data = 5
end
`

	calculate := &model.Calculate{Data: 0}

	properties := make(rulengine.PropertiesMap)
	properties["cal"] = calculate

	functions := make(rulengine.FunctionsMap)
	functions["println"] = fmt.Println

	nodeName := "MixModel"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, dsl, functions, properties, rulengine.MixModel, false)
	assert.Nil(t, err)

	err = eng.Execute(nodeName)
	assert.Nil(t, err)
}

func TestInverseMixModel(t *testing.T) {
	// 先并发执行优先级最低以外的其他规则,最后才执行优先级最低的那一个.
	const dsl = `
rule "lowest_priority" "lowest priority 并发" salience 996
begin
	println("996-->")
	println("cal.Data-->", cal.Data)
	println("cal.Name-->", cal.Name) 
end

rule "lower_priority" "lower priority 并发" salience 998
begin
	println("998-->")
	cal.Name = "hello world"
end

rule "most_priority" "most priority 最后一个" salience 1000
begin
	println("1000-->")
	cal.Data = 5
end
`

	calculate := &model.Calculate{Data: 0}

	properties := make(rulengine.PropertiesMap)
	properties["cal"] = calculate

	functions := make(rulengine.FunctionsMap)
	functions["println"] = fmt.Println

	nodeName := "InverseMixModel"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, dsl, functions, properties, rulengine.InverseMixModel, false)
	assert.Nil(t, err)

	err = eng.Execute(nodeName)
	assert.Nil(t, err)
}

func TestConcurrentModel(t *testing.T) {
	dsl := `
rule "TemperatureRule" "温度事件计算规则"
begin
   println("/***************** 温度事件计算规则 ***************/")
   tempState = 0
   if Temperature.Value < 0 {
      tempState = 1
   } else if Temperature.Value > 80 {
      tempState = 2
   }
   if Temperature.State != tempState {
      if tempState == 0 {
         Temperature.Event = "温度正常"
      } else if tempState == 1 {
         Temperature.Event = "低温报警"
      } else {
         Temperature.Event = "高温报警"
      }
   } else {
      Temperature.Event = ""
   }
   Temperature.State = tempState
end
 
rule "WaterRule" "水浸事件计算规则"
begin
   println("/***************** 水浸事件计算规则 ***************/")
   tempState = 0
   if Water.Value != 0 {
      tempState = 1
   }
   if Water.State != tempState {
      if tempState == 0 {
         Water.Event = "水浸正常"
      } else {
         Water.Event = "水浸异常"
      }
   } else {
      Water.Event = ""
   }
   Water.State = tempState
end
 
rule "SmokeRule" "烟雾事件计算规则"
begin
   println("/***************** 烟雾事件计算规则 ***************/")
   tempState = 0
   if Smoke.Value != 0 {
      tempState = 1
   }
   if Smoke.State != tempState {
      if tempState == 0 {
         Smoke.Event = "烟雾正常"
      } else {
         Smoke.Event = "烟雾报警"
      }
   } else {
      Smoke.Event = ""
   }
   Smoke.State = tempState
end
`

	temperature := &model.Rule{
		Tag:   "temperature",
		Value: 90,
		State: 0,
		Event: "",
	}
	water := &model.Rule{
		Tag:   "water",
		Value: 0,
		State: 0,
		Event: "",
	}
	smoke := &model.Rule{
		Tag:   "smoke",
		Value: 1,
		State: 0,
		Event: "",
	}

	properties := make(rulengine.PropertiesMap)
	properties["Temperature"] = temperature
	properties["Water"] = water
	properties["Smoke"] = smoke

	functions := make(rulengine.FunctionsMap)
	functions["println"] = fmt.Println

	nodeName := "Station"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, dsl, functions, properties, rulengine.ConcurrentModel, false)
	assert.Nil(t, err)

	err = eng.Execute(nodeName)
	assert.Nil(t, err)

	fmt.Printf("temperature Event=%s\n", temperature.Event)
	fmt.Printf("water Event=%s\n", water.Event)
	fmt.Printf("smoke Event=%s\n", smoke.Event)
	for i := 0; i < 10; i++ {
		smoke.Value = float64(i % 3)
		err = eng.Execute(nodeName)
		assert.Nil(t, err)
		fmt.Printf("smoke Event=%s\n", smoke.Event)
	}
}

func TestModifyRule(t *testing.T) {
	ruleInit := `
rule "ruleScore" "rule-des" salience 10
begin
   if Student.Score > 60 {
      println(Student.Name, FormatInt(Student.Score, 10), "及格")
   } else {
      println(Student.Name, FormatInt(Student.Score, 10), "不及格")
   }
end
`

	ruleUpdate := `
rule "ruleScore" "rule-des" salience 10
begin
   if Student.Score > 80 {
      println(Student.Name, FormatInt(Student.Score, 10), "及格")
   } else {
      println(Student.Name, FormatInt(Student.Score, 10), "不及格")
   }
end
`
	ruleAdd := `
rule "ruleTeach " "rule-des" salience 10
begin
   if Student.Score < 70 {
	  println(Student.Name, FormatInt(Student.Score, 10), "需要补课")
   }
end
`

	student := &model.Student{
		Name:  "Phinx",
		Score: 100,
	}

	properties := make(rulengine.PropertiesMap)
	properties["Student"] = student

	functions := make(rulengine.FunctionsMap)
	functions["println"] = fmt.Println
	functions["FormatInt"] = strconv.FormatInt

	nodeName := "Student"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, ruleInit, functions, properties, rulengine.SortModel, false)
	assert.Nil(t, err)

	go func() {
		for {
			student.Score = rand.Int63n(50) + 50
			err2 := eng.Execute(nodeName)
			assert.Nil(t, err2)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		time.Sleep(3 * time.Second)

		err2 := eng.UpdateRule(nodeName, ruleUpdate)
		assert.Nil(t, err2)

		time.Sleep(3 * time.Second)

		err3 := eng.UpdateRule(nodeName, ruleAdd)
		assert.Nil(t, err3)
	}()

	time.Sleep(20 * time.Second)
}
