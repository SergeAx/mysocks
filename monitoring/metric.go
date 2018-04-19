package monitoring

import (
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

type metric struct {
	name   string
	tags   map[string]string
	fields []string
}

func NewMetric(name string, tags map[string]string, fields ...string) metric {
	return metric{name: name, tags: tags, fields: fields}
}

func (m *metric) Record(values ...interface{}) {
	m.record(values, time.Now())
}

func (m *metric) record(values []interface{}, t time.Time) {
	if len(values) != len(m.fields) {
		panic("wrong number of records")
	}

	rec := make(map[string]interface{}, len(m.fields))
	for i, f := range m.fields {
		switch v := values[i].(type) {
		case time.Duration:
			rec[f] = int64(v)
		default:
			rec[f] = v
		}
	}

	p, err := influxdb.NewPoint(m.name, m.tags, rec, t)
	if err != nil {
		panic(err)
	}

	batch.AddPoint(p)
}
