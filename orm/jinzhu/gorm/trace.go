package gorm

import (
	"fmt"
	"runtime/debug"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type TTraceItem struct {
	// GormDB    *DB
	Tx        string
	TraceTime string
	Sql       string
	Stack     string
}

var traceMap = cmap.New()

func GetTraceMap() cmap.ConcurrentMap {
	return traceMap
}

func (s *DB) TxTrace(enable bool) *DB {
	if enable {
		s.txTrace = true
	} else {
		s.txTrace = false
	}
	return s
}

func (s *DB) beginTxTrace() {

	key := fmt.Sprintf("%p", s.db)

	item := &TTraceItem{
		Tx:        key,
		TraceTime: time.Now().Format("2006-01-02 15:04:05.000"),
		Sql:       "TX BEGIN",
		Stack:     string(debug.Stack()),
	}

	itemsobj, ok := traceMap.Get(key)
	if ok {
		panic("tx already exists...")
		items := itemsobj.([]*TTraceItem)
		items = append(items, item)
		traceMap.Set(key, items)
	} else {
		var items []*TTraceItem
		items = append(items, item)
		traceMap.Set(key, items)
	}
}

func (s *DB) addTxTrace(sql string) {

	key := fmt.Sprintf("%p", s.db)

	item := &TTraceItem{
		Tx:        key,
		TraceTime: time.Now().Format("2006-01-02 15:04:05.000"),
		Sql:       sql,
		Stack:     string(debug.Stack()),
	}

	itemsobj, ok := traceMap.Get(key)
	if ok {
		items := itemsobj.([]*TTraceItem)
		items = append(items, item)
		traceMap.Set(key, items)
	} else { // 不存在则不记录
		// panic("tx not exists...")
		// var items []*TTraceItem
		// items = append(items, item)
		// traceMap.Set(key, items)
	}
}

func (s *DB) clearTxTrace() {
	key := fmt.Sprintf("%p", s.db)
	traceMap.Remove(key)
}
