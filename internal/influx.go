package internal

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

var writeApi api.WriteAPI
var queryApi api.QueryAPI

type Point struct {
	Value interface{} `json:"value"`
	Time  time.Time   `json:"time"`
}

func OpenInfluxdb() {
	client := influxdb2.NewClient(config.Influxdb.Url, config.Influxdb.Token)
	writeApi = client.WriteAPI(config.Influxdb.Org, config.Influxdb.Bucket)
	queryApi = client.QueryAPI(config.Influxdb.Org)
}

func Insert(measurement, id string, fields map[string]interface{}, ts time.Time) {
	writeApi.WritePoint(write.NewPoint(measurement, map[string]string{"id": id}, fields, ts))
}

func Query(measurement, id, field, start, end, window, fn string) ([]Point, error) {
	bucket := config.Influxdb.Bucket

	flux := "from(bucket: \"" + bucket + "\")\n"
	flux += "|> range(start: " + start + ", stop: " + end + ")\n"
	flux += "|> filter(fn: (r) => r[\"_measurement\"] == \"" + measurement + "\")\n"
	flux += "|> filter(fn: (r) => r[\"id\"] == \"" + id + "\")\n"
	flux += "|> filter(fn: (r) => r[\"_field\"] == \"" + field + "\")"
	flux += "|> aggregateWindow(every: " + window + ", fn: " + fn + ", createEmpty: false)\n"
	flux += "|> yield(name: \"" + fn + "\")"

	result, err := queryApi.Query(context.Background(), flux)
	if err != nil {
		return nil, err
	}

	records := make([]Point, 0)
	for result.Next() {
		//result.TableChanged()
		records = append(records, Point{
			Value: result.Record().Value(),
			Time:  result.Record().Time(),
		})
	}
	return records, result.Err()
}
