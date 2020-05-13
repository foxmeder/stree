package stree

// SortItem 可排序数据的interface
type Item interface {
	SetLevel(level int64) // 写level值
	GetMeta() ItemMeta    // 获取数据排序需要的数据
}

// ItemMeta 排序需要的元数据
type ItemMeta struct {
	ID    int64       //数据ID
	PID   int64       //数据父ID
	Level int64       //获取level值，非必须
	Sort  int64       //获取排序值
	Name  string      //展示名称
	val   interface{} //原始数据
}
