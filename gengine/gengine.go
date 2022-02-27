package gengine

import (
	"errors"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	rulengine "go-rulengine-example"
)

type NodeMap map[string]*builder.RuleBuilder

type RuleEngine struct {
	eng   *engine.Gengine
	nodes NodeMap
}

// NewRuleEngine 创建一个新的规则引擎实例
func NewRuleEngine() *RuleEngine {
	re := &RuleEngine{eng: nil, nodes: NodeMap{}}
	return re
}

// Start 启动规则引擎
func (r *RuleEngine) Start() error {
	r.eng = engine.NewGengine()
	if r.eng == nil {
		return errors.New("engine init failed")
	}
	return nil
}

// Stop 停止规则引擎
func (r *RuleEngine) Stop() error {
	return nil
}

// AddNode 添加一个规则节点
func (r *RuleEngine) AddNode(nodeName, rules string, apis rulengine.ExportMap) error {
	_, ok := r.nodes[nodeName]
	if ok {
		return errors.New("rule node already exists")
	}

	ruleBuilder := r.createBuilder(apis)
	if ruleBuilder == nil {
		return errors.New("create rule builder failed")
	}

	err := ruleBuilder.BuildRuleFromString(rules)
	if err != nil {
		return err
	}

	r.nodes[nodeName] = ruleBuilder

	return nil
}

// RemoveNode 移除一个规则节点
func (r *RuleEngine) RemoveNode(nodeName string) error {
	delete(r.nodes, nodeName)
	return nil
}

// NodeCount 规则节点个数
func (r *RuleEngine) NodeCount() int {
	return len(r.nodes)
}

// Execute 执行规则
func (r *RuleEngine) Execute(nodeName string, concurrent bool) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}

	if concurrent {
		return r.eng.ExecuteConcurrent(node)
	} else {
		return r.eng.Execute(node, true)
	}
}

// ModifyRule 修改规则
func (r *RuleEngine) ModifyRule(nodeName, rules string) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}
	return node.BuildRuleWithIncremental(rules)
}

// RemoveRules 移除规则
// @param [in] nodeName 节点名
// @param [in] ruleNames 规则名列表
func (r *RuleEngine) RemoveRules(nodeName string, ruleNames []string) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}
	return node.RemoveRules(ruleNames)
}

func (r *RuleEngine) createBuilder(apis rulengine.ExportMap) *builder.RuleBuilder {
	dataContext := context.NewDataContext()
	for k, v := range apis {
		dataContext.Add(k, v)
	}
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	return ruleBuilder
}
