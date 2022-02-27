package model

type Station struct {
	Temperature  int64 //温度
	Humidity     int64 //湿度
	Water        int64 //水浸
	Smoke        int64 //烟雾
	Door1        int64 //门禁1
	Door2        int64 //门禁2
	StationState int64 //探测站状态:   0正常；1预警；2异常；3未知
}

type Rule struct {
	Tag   string  //标签点名称
	Value float64 //数据值
	State int64   //状态
	Event string  //报警事件
}
