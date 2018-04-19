package monitoring

import (
	"os"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/aspcartman/mysocks/env"
	"sync"
)

var log = env.Log.WithField("module", "monitoring")

var (
	client influxdb.Client
	batch  influxdb.BatchPoints
	mtx    sync.Mutex
)

func init() {
	addr := os.Getenv("INFLUX_ADDR")
	db := os.Getenv("INFLUX_DB")

	if len(addr) > 0 {
		var err error
		client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
			Addr: addr,
		})
		if err != nil {
			panic(err)
		}

		_, err = client.Query(influxdb.Query{
			Command:  "CREATE DATABASE " + db,
			Database: db,
		})
		if err != nil {
			panic(err)
		}

		batch = newbatch(db)
	} else {
		client = nullclient{}
		batch = nullbatch{}
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)

			b := newbatch(db)
			mtx.Lock()
			b, batch = batch, b
			mtx.Unlock()

			if err := client.Write(b); err != nil {
				log.WithError(err).Error("failed writing influxdb metrics")
			}
		}
	}()
}

func newbatch(db string) influxdb.BatchPoints {
	b, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Precision: "ms",
		Database:  db,
	})
	if err != nil {
		panic(err)
	}
	return b
}
