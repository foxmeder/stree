package stree

import (
	"reflect"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

// Sorter 实现了排序需要的方法
type Sorter struct {
	sync.Mutex
	data   map[int64][]ItemMeta
	len    int
	sorted bool
}

// Insert 将需要排序的数据放入 Sorter
// list必须为slice
func (s *Sorter) Insert(list interface{}) error {
	vlist := reflect.ValueOf(list)
	if vlist.Kind() != reflect.Slice {
		return errors.New("not a list")
	}
	if vlist.Len() == 0 {
		return nil
	}
	s.Lock()
	defer s.Unlock()
	if s.data == nil {
		s.data = make(map[int64][]ItemMeta)
	}
	for j := 0; j < vlist.Len(); j++ {
		vsi, ok := vlist.Index(j).Interface().(Item)
		if !ok {
			return errors.Errorf("index %d is not Item", j)
		}
		s.insert(vsi)
	}
	return nil
}

// 按数据pid分组
func (s *Sorter) insert(item ...Item) {
	if len(item) == 0 {
		return
	} else {
		s.sorted = false
		for _, v := range item {
			meta := v.GetMeta()
			meta.val = v
			//pid := v.GetPID()
			s.data[meta.PID] = append(s.data[meta.PID], meta)
		}
		s.len = s.len + len(item)
	}
}

// GetList 获取排序后的数据
// list必须是slice的指针，方法会根据list的数据类型尝试进行类型转换
func (s *Sorter) GetList(list interface{}) error {
	vlist := reflect.ValueOf(list)
	// 检查是否为指针
	if vlist.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	// 是否空指针
	if vlist.IsNil() {
		return errors.New("nil pointer")
	}
	// 指针指向是否为slice
	ret := reflect.Indirect(vlist)
	if ret.Kind() != reflect.Slice {
		return errors.New("not a list")
	}
	// 排序
	s.sort()
	// 根据list类型创建slice
	typo := ret.Type().Elem()
	sv := reflect.MakeSlice(ret.Type(), s.len, s.len)
	i := 0
	// 获取排序后的数据
	c := s.fillChan()
	for v := range c {
		//beego.Info(v.GetID(), v.GetPID())
		val := reflect.ValueOf(v.val)
		if val.Type().ConvertibleTo(typo) {
			sv.Index(i).Set(val.Convert(typo))
			i++
		} else {
			return errors.Errorf("can not convert to %v", typo)
		}
	}
	// 为list赋值
	ret.Set(sv.Slice(0, i))
	return nil
}

// GetTree 获取前端js layui treeseelct使用的分类树
// level为正数时 大于此level的数据不返回
func (s *Sorter) GetTree(level int) []*CateTree {
	ret := make([]*CateTree, 0)
	tmp := make(map[int64]*CateTree)
	c := s.fillChan()
	for v := range c {
		if level > 0 && v.Level > int64(level) {
			// 大于特定level的不要
			continue
		}
		cateTree := CateTree{
			ID:      v.ID,
			Name:    v.Name,
			Open:    true,
			Checked: false,
			//Children: make([]*CateTree, 0),
		}
		tmp[v.ID] = &cateTree
		if v.PID == 0 {
			// 一级分类
			ret = append(ret, &cateTree)
		} else if cate, ok := tmp[v.PID]; ok {
			// 增加子节点
			cate.Children = append(cate.Children, &cateTree)
		} else {
			// 其他可能有问题的数据别丢了
			ret = append(ret, &cateTree)
		}
	}
	return ret
}

// 创建channel按顺序放入排序后的数据
func (s *Sorter) fillChan() chan ItemMeta {
	c := make(chan ItemMeta, s.len)
	go func() {
		defer close(c)
		s.fillSlice(c, 0, 0)
	}()
	return c
}

// 递归向channel存入数据
func (s *Sorter) fillSlice(c chan ItemMeta, pid, level int64) {
	if data, ok := s.data[pid]; ok {
		for _, v := range data {
			if item, ok := v.val.(Item); ok {
				item.SetLevel(level)
			}
			v.Level = level
			//beego.Info(v.GetID(), v.GetPID())
			//v.val.(Item).SetLevel(level)
			c <- v
			s.fillSlice(c, v.ID, level+1)
		}
	}
}

// 排序数据
func (s *Sorter) sort() {
	s.Lock()
	defer s.Unlock()
	if s.sorted {
		return
	}
	for i, v := range s.data {
		s.data[i] = sortItem(v)
	}
	s.sorted = true
}

func sortItem(si []ItemMeta) []ItemMeta {
	sort.SliceStable(si, func(i, j int) bool {
		return si[i].Sort < si[j].Sort
	})
	return si
}

// GetCateTree 从slice获取前端js需要的结构
// list必须为slice
func GetCateTree(list interface{}, level int) ([]*CateTree, error) {
	sorter := &Sorter{}
	if err := sorter.Insert(list); err != nil {
		return make([]*CateTree, 0), err
	}
	return sorter.GetTree(level), nil
}

// GetList 获取排序过的列表
// list必须为slice的指针
func GetList(list interface{}) error {
	vlist := reflect.ValueOf(list)
	if vlist.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	sorter := &Sorter{}
	if err := sorter.Insert(vlist.Elem().Interface()); err != nil {
		return errors.Wrap(err, "insert list")
	}
	if err := sorter.GetList(list); err != nil {
		return errors.Wrap(err, "sort list")
	}
	return nil
}
