package monitoring

import (
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

type nullclient struct{}

func (nullclient) Ping(timeout time.Duration) (time.Duration, string, error) { return 0, "", nil }

func (nullclient) Write(bp influxdb.BatchPoints) error { return nil }

func (nullclient) Query(q influxdb.Query) (*influxdb.Response, error) { return nil, nil }

func (nullclient) Close() error { return nil }

type nullbatch struct{}

func (nullbatch) AddPoint(p *influxdb.Point) {}

func (nullbatch) AddPoints(ps []*influxdb.Point) {}

func (nullbatch) Database() string { return "" }

func (nullbatch) Points() []*influxdb.Point { return nil }

func (nullbatch) Precision() string { return "" }

func (nullbatch) RetentionPolicy() string { return "" }

func (nullbatch) SetDatabase(s string) {}

func (nullbatch) SetPrecision(s string) error { return nil }

func (nullbatch) SetRetentionPolicy(s string) {}

func (nullbatch) SetWriteConsistency(s string) {}

func (nullbatch) WriteConsistency() string { return "" }
