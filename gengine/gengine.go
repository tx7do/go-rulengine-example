package gengine

import (
	"errors"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	rulengine "go-rulengine-example"
)

var (
	DefaultPoolMinLen int64 = 10
	DefaultPoolMaxLen int64 = 20
)

type Node struct {
	eng     *engine.Gengine
	builder *builder.RuleBuilder

	pool       *engine.GenginePool
	properties rulengine.PropertiesMap

	em rulengine.ExecuteModel
}

func (node *Node) createBuilder(functions rulengine.FunctionsMap, properties rulengine.PropertiesMap) *builder.RuleBuilder {
	dataContext := context.NewDataContext()
	for k, v := range functions {
		dataContext.Add(k, v)
	}
	for k, v := range properties {
		dataContext.Add(k, v)
	}
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	return ruleBuilder
}

// createNormalEngine 创建普通的规则引擎
// @param [in] rules 规则的字符串
// @param [in/out] properties 要绑定到规则引擎的传出属性表
func (node *Node) createNormalEngine(em rulengine.ExecuteModel, rules string, functions rulengine.FunctionsMap, properties rulengine.PropertiesMap) error {
	node.eng = engine.NewGengine()
	if node.eng == nil {
		return errors.New("engine init failed")
	}

	ruleBuilder := node.createBuilder(functions, properties)
	if ruleBuilder == nil {
		return errors.New("create rule builder failed")
	}

	err := ruleBuilder.BuildRuleFromString(rules)
	if err != nil {
		return err
	}

	node.em = em
	node.properties = properties
	node.builder = ruleBuilder

	return nil
}

// createPoolEngine 创建池化规则引擎
// @param [in] poolMinLen 池最小长度
// @param [in] poolMaxLen 池最大长度
// @param [in] em 执行模型
// @param [in] rules 规则的字符串
// @param [in/out] properties 要绑定到规则引擎的传出属性表(这里最好仅注入一些无状态函数，方便应用中的状态管理)
func (node *Node) createPoolEngine(poolMinLen, poolMaxLen int64, em rulengine.ExecuteModel, rules string, functions rulengine.FunctionsMap, properties rulengine.PropertiesMap) error {
	pool, err := engine.NewGenginePool(poolMinLen, poolMaxLen, int(node.em), rules, functions)
	if err != nil {
		return errors.New("create gengine pool failed")
	}
	node.em = em
	node.properties = properties
	node.pool = pool
	return nil
}

// updateRule 全量更新规则
func (node *Node) updateRule(rules string) error {
	if node.builder != nil {
		return node.builder.BuildRuleFromString(rules)
	} else if node.pool != nil {
		return node.pool.UpdatePooledRules(rules)
	} else {
		return errors.New("update rule failed")
	}
}

func (node *Node) execute() error {
	if node.builder != nil {
		switch node.em {
		case rulengine.SortModel:
			return node.eng.Execute(node.builder, true)
		case rulengine.ConcurrentModel:
			return node.eng.ExecuteConcurrent(node.builder)
		case rulengine.MixModel:
			return node.eng.ExecuteMixModel(node.builder)
		case rulengine.InverseMixModel:
			return node.eng.ExecuteInverseMixModel(node.builder)
		case rulengine.BucketModel:
			return errors.New("not support bucket model")
		default:
			return errors.New("unknown execute model")
		}
	} else if node.pool != nil {
		err, _ := node.pool.Execute(node.properties, true)
		return err
	} else {
		return errors.New("execute failed")
	}
}

// stop 停止
func (node *Node) stop() error {
	if node.builder != nil {
		node.builder = nil
	} else if node.pool != nil {
		node.pool = nil
	}
	return nil
}

type NodeMap map[string]Node

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
func (r *RuleEngine) AddNode(nodeName, rules string, functions rulengine.FunctionsMap, properties rulengine.PropertiesMap, executeModel rulengine.ExecuteModel, pooled bool) error {
	_, ok := r.nodes[nodeName]
	if ok {
		return errors.New("rule node already exists")
	}

	node := Node{em: executeModel, properties: properties}

	if !pooled {
		if err := node.createNormalEngine(executeModel, rules, functions, properties); err != nil {
			return err
		}
	} else {
		if err := node.createPoolEngine(DefaultPoolMinLen, DefaultPoolMaxLen, executeModel, rules, functions, properties); err != nil {
			return err
		}
	}

	r.nodes[nodeName] = node

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
func (r *RuleEngine) Execute(nodeName string) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}
	return node.execute()
}

// UpdateRule 全量更新规则
func (r *RuleEngine) UpdateRule(nodeName, rules string) error {
	node, ok := r.nodes[nodeName]
	if !ok {
		return errors.New("rule node does not exist")
	}
	return node.updateRule(rules)
}
