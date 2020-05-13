package stree

// CateTree layui treeselect需要的json数据结构
type CateTree struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Open     bool        `json:"open"`
	Checked  bool        `json:"checked"`
	Children []*CateTree `json:"children,omitempty"`
}
