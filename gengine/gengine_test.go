package gengine

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/google/martian/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	rule_engine "go-rulengine-example"
	"go-rulengine-example/model"
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
	const grl = `
rule "测试规则名称1" "rule desc"
begin
  aName = @name
  aId = @id
  aDesc = @desc
  PrintName("@id: ", aId)
  PrintName("@name: ", aName)
  PrintName("@desc: ", aDesc)
end

rule "rule name" "rule desc"
begin
  aName = @name
  aId = @id
  aDesc = @desc
  PrintName("@id: ", aId)
  PrintName("@name: ", aName)
  PrintName("@desc: ", aDesc)
end
`

	start1 := time.Now().UnixNano()

	// 创建上下文
	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", fmt.Println)
	// 初始化规则构造器
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	// 从字符串中解析规则
	err := ruleBuilder.BuildRuleFromString(grl)

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

func TestConcurrent(t *testing.T) {
	grl := `
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

	apis := make(rule_engine.ExportMap)
	apis["Temperature"] = temperature
	apis["Water"] = water
	apis["Smoke"] = smoke
	apis["println"] = fmt.Println

	nodeName := "Station"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, grl, apis)
	assert.Nil(t, err)

	err = eng.Execute(nodeName, true)
	assert.Nil(t, err)

	fmt.Printf("temperature Event=%s\n", temperature.Event)
	fmt.Printf("water Event=%s\n", water.Event)
	fmt.Printf("smoke Event=%s\n", smoke.Event)
	for i := 0; i < 10; i++ {
		smoke.Value = float64(i % 3)
		err = eng.Execute(nodeName, true)
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
		Name:  "菲乐",
		Score: 100,
	}

	apis := make(rule_engine.ExportMap)
	apis["FormatInt"] = strconv.FormatInt
	apis["println"] = fmt.Println
	apis["Student"] = student

	nodeName := "Student"

	eng := NewRuleEngine()
	err := eng.Start()
	assert.Nil(t, err)

	err = eng.AddNode(nodeName, ruleInit, apis)
	assert.Nil(t, err)

	go func() {
		for {
			student.Score = rand.Int63n(50) + 50
			err2 := eng.Execute(nodeName, false)
			assert.Nil(t, err2)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		time.Sleep(3 * time.Second)

		err2 := eng.ModifyRule(nodeName, ruleUpdate)
		assert.Nil(t, err2)

		time.Sleep(3 * time.Second)

		err3 := eng.ModifyRule(nodeName, ruleAdd)
		assert.Nil(t, err3)
	}()

	time.Sleep(20 * time.Second)
}
