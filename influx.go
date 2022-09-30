package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/influxdata/influxdb-client-go/api/write"
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

func Insert(id string, values map[string]interface{}) {
	writeApi.WritePoint(write.NewPoint(config.Influxdb.Measurement, map[string]string{"id": id}, values, time.Now()))
}

func Query(id, field, start, end, window, fn string) ([]Point, error) {
	//metric := fmt.Sprintf("%d", id)
	bucket := "zgwit"

	flux := "from(bucket: \"" + bucket + "\")\n"
	flux += "|> range(start: " + start + ", stop: " + end + ")\n"
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
