package stream

import (
	"fmt"

	"github.com/ddkwork/golibrary/mylog"
)

func TopologicalSort[T comparable](m *OrderedMap[T, []T], allowCyclicDependency bool) (sorted []T) {
	var zero T                            // 说白了就是树形结构转为去重+处理优先级的线性结构
	visited := NewOrderedMap(zero, false) // 用于检查孩子节点是否在容器节点中
	temp := NewOrderedMap(zero, false)    // 用于检测循环依赖
	var visitAll func(T)

	visitAll = func(id T) {
		if temp.Has(id) {
			if allowCyclicDependency {
				return // sln can be cyclic
			}
			mylog.Check(fmt.Errorf("cyclic dependency detected involving project %v", id))
		}
		if !visited.Has(id) { // 递归处理node及其children,最终得到一个拓扑排序sorted
			temp.Set(id, true)
			deps, ok := m.Get(id)
			if ok {
				for _, depID := range deps {
					visitAll(depID)
				}
			}
			temp.Set(id, false)
			visited.Set(id, true)
			// 排除每个容器及其孩子节点中存在的重复孩子节点,通过递归深度遍历得到了去重的线性节点切片
			sorted = append(sorted, id)
		}
	}

	for _, p := range m.List() {
		if !visited.Has(p.Key) {
			visitAll(p.Key)
		}
	}
	return
}
