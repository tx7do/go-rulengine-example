package rulengine

type ExecuteModel int

const (
	SortModel       ExecuteModel = 1 // 顺序模式
	ConcurrentModel ExecuteModel = 2 // 并发执行模式
	MixModel        ExecuteModel = 3 // 混合执行模式
	InverseMixModel ExecuteModel = 4 // 逆混合执行模式
	BucketModel     ExecuteModel = 5 // 桶模式
)

type PropertiesMap map[string]interface{}
type FunctionsMap map[string]interface{}

type RuleEngine interface {
	// Start 启动规则引擎
	Start() error

	// Stop 停止规则引擎
	Stop() error

	// Execute 执行规则
	// @param [in] concurrent 是否并发执行
	Execute(nodeName string, executeModel ExecuteModel) error

	// AddNode 添加一个规则节点
	// @param [in] nodeName 节点名
	// @param [in] properties 属性,方法的键值对
	// @param [in] rules 字符串规则
	AddNode(nodeName, rules string, properties PropertiesMap) error

	// RemoveNode 移除一个规则节点
	RemoveNode(nodeName string) error

	// NodeCount 规则节点个数
	NodeCount() int

	// ModifyRule 修改规则
	// @param [in] nodeName 节点名
	// @param [in] rules 字符串规则
	ModifyRule(nodeName, rules string) error
}
