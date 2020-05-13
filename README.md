# stree

[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/foxmeder/stree?tab=doc)

cate sort library by go

## Example

```go
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
```

