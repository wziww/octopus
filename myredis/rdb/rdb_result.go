package rdb

import "container/heap"

const (
	DEFAULT_SIZE = 10
)

// Node ...
type Node struct {
	Key string
	Val int64
}

// LogHeap ...
type LogHeap []Node

// Len ...
func (l *LogHeap) Len() int {
	return len(*l)
}

// Less ...
func (l *LogHeap) Less(i, j int) bool {
	return (*l)[i].Val < (*l)[j].Val
}

// Swap ...
func (l *LogHeap) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

// Push ...
func (l *LogHeap) Push(x interface{}) {
	*l = append(*l, x.(Node))
}

// Pop ...
func (l *LogHeap) Pop() interface{} {
	res := (*l)[len(*l)-1]
	*l = (*l)[:len(*l)-1]
	return res
}

// Result ...
type Result struct {
	TotalNums      int64    // 键值总数
	Expires        int64    // 过期键总数
	AlreadyExpired int64    // 当前时间下已过期键数量
	LuaNums        int64    // lua 脚本缓存数量
	OffSetSize     int64    // 偏移量统计数值
	OffSetLog      *LogHeap // 偏移量统计
	OffSetCount    int64    // 当前偏移量日志数量
	ChildSize      int64    // 成员数量统计数值
	ChildLog       *LogHeap // 成员数量统计
	ChildCount     int64    // 当前成员日志数量
	Count          int64    // 日志记录条数
}

// CreateResult ...
// offsetSize 偏移量统计临界值
// childSize 成员数量统计数值
// count     日志记录条数
func CreateResult(offsetSize int64, childSize int64, count int64) *Result {
	r := &Result{
		OffSetSize: offsetSize,
		ChildSize:  childSize,
		Count:      count,
		// OffSetLog:  make([]Node, 0, DEFAULT_SIZE),
		// ChildLog:   make([]Node, 0, DEFAULT_SIZE),
	}
	var offsetlog, childlog LogHeap
	offsetlog = make([]Node, 0, DEFAULT_SIZE)
	childlog = make([]Node, 0, DEFAULT_SIZE)
	r.OffSetLog = &offsetlog
	r.ChildLog = &childlog
	heap.Init(r.OffSetLog)
	heap.Init(r.ChildLog)
	return r
}
