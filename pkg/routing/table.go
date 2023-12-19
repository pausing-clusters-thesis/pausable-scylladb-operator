package routing

import (
	"sync"
	"sync/atomic"
)

// The routing table implementation uses the ReadMostly pattern: https://pkg.go.dev/sync/atomic#example-Value-ReadMostly.

type Table struct {
	v atomic.Value
	m sync.Mutex
}

func NewTable() *Table {
	t := &Table{}
	t.v.Store(map[string]string{})

	return t
}

func (t *Table) Get(host string) (string, bool) {
	backendHost, ok := t.v.Load().(map[string]string)[host]
	return backendHost, ok
}

func (t *Table) Store(routingTable map[string]string) {
	t.m.Lock()
	defer t.m.Unlock()
	t.v.Store(routingTable)
}
