# go-rulengine-example

本项目是测试golang下面的规则引擎

## 什么是规则引擎

规则引擎是一个逻辑或条件，例如“当某些条件被评估为真时，然后执行某些任务。”规则引擎可以被视为一个复杂的 if/then 语句解释器。被解释的 if/then 语句称为规则。将规则引擎想象成一个将数据和规则作为输入的系统。它将这些规则应用于数据，并根据规则定义为我们提供输出。

![什么是规则引擎](.\docs\rule_egine_1_.svg)

## GO规则引擎列表

| 名称	                                                         | 规则描述语言	    | 使用场景	     | 使用复杂性 |
|-------------------------------------------------------------|------------|-----------|-------|
| [govaluate](https://github.com/Knetic/govaluate)            | 类Golang	   | 表达式解析	    | 低     |
| [YQL](https://github.com/caibirdme/yql)                     | 类SQL	      | 表达式解析	    | 低     |
| [Gval](https://github.com/PaesslerAG/gval)                  | 类Golang	   | 表达式解析	    | 低     |
| [Grule](https://github.com/hyperjumptech/grule-rule-engine) | 自定义DSL     | 规则执行	     | 中     |
| [Gengine](https://github.com/bilibili/gengine)              | 自定义DSL     | 规则执行	     | 中     |
| [goja](https://github.com/dop251/goja)                      | JavaScript | 规则解析	     | 中     |
| [cel-go](https://github.com/google/cel-go)                  | 类C	        | 通用表达式语言		 | 中     |
| [GopherLua](https://github.com/yuin/gopher-lua)             | lua	       | 规则解析	     | 高     |

## 规则执行模式

通过对各种业务场景的分析提炼，一个规则引擎至少应该满足3种执行模式。但实际上，规则执行模式至少有5种，具体执行模式，如下图所示：

![规则执行模式](.\docs\rule-model.png)

### 1. **顺序模式(Sort Model)**

![顺序模式](.\docs\sort_model.png)

规则优先级高越高的越先执行，规则优先级低的越后执行。这也是drools支持的模式。此模式的缺点很明显：随着规则链越来越长，执行规则返回的速度也越来越慢。

### 2. **并发执行模式(Concurrent Model)**

![并发执行模式](.\docs\concurrent_model.png)

在此执行模式下，多个规则执行时，不考虑规则之间的优先级，规则与规则之间并发执行。规则执行的返回的速度等于所有规则中的执行时间最长的那个规则的速度（**逆木桶原理**）。执行性能优异，但无法满足规则优先级。

### 3. **混合执行模式（Mix Model）**

![混合执行模式](.\docs\mix_model.png)

规则引擎选择一个优先级最高规则的最先执行，剩下的规则并发执行。规则执行返回耗时 = 最高优先级的那个规则执行时间 + 并发执行中执行时间最长的那个规则耗时；此模式兼顾优先级和性能，适合于有豁免规则(或前置规则)的场景。

### 4. **逆混合执行模式(Inverse Mix Model)**

![逆混合执行模式](.\docs\inverse_mix_model.png)

优先级最高的n-1个规则并发执行，执行完毕之后，再执行剩下的一个优先级最低的规则。这种模式适用于有很多前导判断规则的场景。其特性与混合模式类似，兼顾性能和优先级。

### 5. **桶模式（Bucket Model）**

![桶模式](.\docs\bucket_model.png)

规则引擎基于规则优先级进行分桶，优先级相同的规则置于同一个桶中，桶内的规则并发执行，桶间的规则基于规则优先级顺序执行。

## 参考资料

- [RulesEngine - Martin Fowler](https://martinfowler.com/bliki/RulesEngine.html)
- [Business rules engine](https://en.wikipedia.org/wiki/Business_rules_engine)
- [What is Rule-Engine?](https://medium.com/@er.rameshkatiyar/what-is-rule-engine-86ea759ad97d)
- [What is a rule engine ?](http://www.mastertheboss.com/bpm/drools/what-is-a-rule-engine/)
- [B站新一代golang规则引擎的设计与实现](https://cloud.tencent.com/developer/news/667806)
- [规则引擎在哔哩哔哩的应用](https://www.biaodianfu.com/bilibili-gengine.html)