/*

sort 提供了无限级分类排序方法

包内实现了两个简单易用的方法 GetList 和 GetCateTree

分别获取排序后的 slice 和 layui treeselect需要使用的json数据结构

使用此包进行排序，被排序的数据需要实现 Item

然后使用

	var list []Item
	err := sort.GetList(&list)

或

	var list []Item
	tree,err := sort.GetCateTree(list)

注意 GetList 方法的参数必须为slice的指针，GetCateTree 的参数为slice

*/
package stree
