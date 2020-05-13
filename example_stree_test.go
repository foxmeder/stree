package stree_test

import (
	"encoding/json"
	"fmt"

	"github.com/foxmeder/stree"
)

type TreeItem struct {
	ID    int64
	PID   int64
	Level int64
	Sort  int64
	Name  string
}

func (ti *TreeItem) SetLevel(lvl int64) {
	ti.Level = lvl
}

func (ti *TreeItem) GetMeta() stree.ItemMeta {
	return stree.ItemMeta{
		ID:    ti.ID,
		PID:   ti.PID,
		Level: ti.Level,
		Sort:  ti.Sort,
		Name:  ti.Name,
	}
}

func Example() {
	list := []*TreeItem{
		{1, 0, 0, 2, "cate1"},
		{2, 0, 0, 1, "cate2"},
		{3, 1, 0, 0, "cate1-3"},
		{4, 2, 0, 5, "cate2-4"},
		{5, 2, 0, 1, "cate2-5"},
	}
	if len(list) != 0 {
		// 编译阶段接口实现判断
		var _ stree.Item = list[0]
		if err := stree.GetList(&list); err != nil {
			panic(err)
		} else {
			b, _ := json.Marshal(list)
			fmt.Printf("%s\n", b)
		}
		if tree, err := stree.GetCateTree(list, 0); err != nil {
			panic(err)
		} else {
			b, _ := json.Marshal(tree)
			fmt.Printf("%s\n", b)
		}
	}
	// Output:[{"ID":2,"PID":0,"Level":0,"Sort":1,"Name":"cate2"},{"ID":5,"PID":2,"Level":1,"Sort":1,"Name":"cate2-5"},{"ID":4,"PID":2,"Level":1,"Sort":5,"Name":"cate2-4"},{"ID":1,"PID":0,"Level":0,"Sort":2,"Name":"cate1"},{"ID":3,"PID":1,"Level":1,"Sort":0,"Name":"cate1-3"}]
	// [{"id":2,"name":"cate2","open":true,"checked":false,"children":[{"id":5,"name":"cate2-5","open":true,"checked":false},{"id":4,"name":"cate2-4","open":true,"checked":false}]},{"id":1,"name":"cate1","open":true,"checked":false,"children":[{"id":3,"name":"cate1-3","open":true,"checked":false}]}]
}
