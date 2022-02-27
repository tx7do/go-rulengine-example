package grule

import (
	"errors"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	rulengine "go-rulengine-example"
)

type Node struct {
	dc ast.IDataContext
	kb *ast.KnowledgeBase
}

type NodeMap map[string]*Node

type RuleEngine struct {
	eng   *engine.GruleEngine
	nodes NodeMap
	ver   string
}

// NewRuleEngine 创建一个新的规则引擎实例
func NewRuleEngine() *RuleEngine {
	re := &RuleEngine{eng: nil, nodes: NodeMap{}, ver: "0.0.1"}
	return re
}

// Start 启动规则引擎
func (r *RuleEngine) Start() error {
	r.eng = engine.NewGruleEngine()
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

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	dataContext := ast.NewDataContext()
	for k, v := range apis {
		err := dataContext.Add(k, v)
		if err != nil {
			return err
		}
	}

	byteRules := pkg.NewBytesResource([]byte(rules))
	err := ruleBuilder.BuildRuleFromResource(nodeName, r.ver, byteRules)
	if err != nil {
		return err
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance(nodeName, r.ver)
	if knowledgeBase == nil {
		return errors.New("no knowledge base")
	}

	var node Node
	node.dc = dataContext
	node.kb = knowledgeBase
	r.nodes[nodeName] = &node

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
func (r *RuleEngine) Execute(nodeName string, _ bool) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}
	if node.dc == nil || node.kb == nil {
		return errors.New("rule node does not exist")
	}

	return r.eng.Execute(node.dc, node.kb)
}

// ModifyRule 修改规则
func (r *RuleEngine) ModifyRule(_, _ string) error {
	return errors.New("NOT IMPL")
}
