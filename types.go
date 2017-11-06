package qlog

type FlatMapS struct {
	K []string
	V []string
}

type FlatMapI struct {
	K []string
	V []interface{}
}

func NewMapS() *FlatMapS {
	return &FlatMapS{
		K: make([]string, 0, 4),
		V: make([]string, 0, 4),
	}
}

func NewMapI() *FlatMapI {
	return &FlatMapI{
		K: make([]string, 0, 4),
		V: make([]interface{}, 0, 4),
	}
}

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

func (m *FlatMapS) Get(key string) string {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			return m.V[k]
		}
	}
	return ""
}

func (m *FlatMapI) Get(key string) interface{} {
	for k := len(m.K) - 1; k >= 0; k-- {
		if m.K[k] == key {
			return m.V[k]
		}
	}
	return nil
}
