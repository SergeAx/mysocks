package monitoring

import (
	"os"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/aspcartman/mysocks/env"
)

var log = env.Log.WithField("module", "monitoring")

var (
	client influxdb.Client
	batch  influxdb.BatchPoints
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

		batch, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Precision: "ms",
			Database:  db,
		})
		if err != nil {
			panic(err)
		}
	} else {
		client = nullclient{}
		batch = nullbatch{}
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := client.Write(batch); err != nil {
				log.WithError(err).Error("failed writing influxdb metrics")
			}
		}
	}()
}
