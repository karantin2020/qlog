package qlog

import (
	"io"
)

type FormatWriter interface {
	io.Writer
	Format(*Entry)
}

type Formatter func(*Entry)
type Hook func(*Entry)
type Output func(*Entry)

// FlatMapS structure to store map[string]string-like data
type FlatMapS struct {
	K []string
	V []string
}

// FlatMapI structure to store map[string]interface{}-like data
type FlatMapI struct {
	K []string
	V []interface{}
}

// NewMapS returns pointer to new FlatMapS structure instance
func NewMapS() *FlatMapS {
	return &FlatMapS{
		K: make([]string, 0, 4),
		V: make([]string, 0, 4),
	}
}

// NewMapI returns pointer to new FlatMapI structure instance
func NewMapI() *FlatMapI {
	return &FlatMapI{
		K: make([]string, 0, 4),
		V: make([]interface{}, 0, 4),
	}
}

// Add adds new key-value pair to struct
func (m *FlatMapS) Add(key string, val string) {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			m.V[k] = val
			return
		}
	}
	m.K = append(m.K, key)
	m.V = append(m.V, val)
}

// Add adds new key-value pair to struct
func (m *FlatMapI) Add(key string, val interface{}) {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			m.V[k] = val
			return
		}
	}
	m.K = append(m.K, key)
	m.V = append(m.V, val)
}

// Get returns value for given key. If no key exists in structure then return empty string
func (m *FlatMapS) Get(key string) string {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			return m.V[k]
		}
	}
	return ""
}

// Get returns value for given key. If no key exists in structure then return nil
func (m *FlatMapI) Get(key string) interface{} {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			return m.V[k]
		}
	}
	return nil
}
