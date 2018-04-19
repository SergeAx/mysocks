package monitoring

import (
	"sync/atomic"
)

type counter struct {
	metric
	value int64
}

func NewCounter(name string, tags map[string]string) *counter {
	return &counter{NewMetric(name, tags, "count"), 0}
}

func (m *counter) Increase() {
	m.metric.Record(atomic.AddInt64(&m.value, 1))
}

func (m *counter) Decrease() {
	m.metric.Record(atomic.AddInt64(&m.value, -1))
}
