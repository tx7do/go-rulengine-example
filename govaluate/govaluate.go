package govaluate

import (
	"errors"
	"github.com/dop251/goja"
	rulengine "go-rulengine-example"
)

type Node struct {
}

type NodeMap map[string]*Node

type RuleEngine struct {
	nodes NodeMap
}

// NewRuleEngine 创建一个新的规则引擎实例
func NewRuleEngine() *RuleEngine {
	re := &RuleEngine{nodes: NodeMap{}}
	return re
}

// Start 启动规则引擎
func (r *RuleEngine) Start() error {

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

	vm := goja.New()
	if vm == nil {
		return errors.New("create js vm failed")
	}

	for k, v := range apis {
		err := vm.Set(k, v)
		if err != nil {
			return err
		}
	}

	var node Node
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
	_, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}

	//_, err := node.vm.RunString(node.js)
	//return err

	return nil
}

// ModifyRule 修改规则
func (r *RuleEngine) ModifyRule(nodeName, rules string) error {
	_, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}

	return nil
}
