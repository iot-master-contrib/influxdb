package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"time"
)

var writeApi api.WriteAPI
var queryApi api.QueryAPI

type Point struct {
	Value interface{} `json:"value"`
	Time  time.Time   `json:"time"`
}

func OpenInfluxdb() {
	org := "zgwit"
	bucket := "zgwit"
	token := "K79bHLvUj0aXqqJH-QBnI-pPl_ZGKWbT6f5Jx1Ol1ZEY1LTJJUwt_aiWJ5vxUc1PcngV7BmlwkJ6ckqkQv5x2g=="
	//measurement := "data_up"

	client := influxdb2.NewClient("http://git.zgwit.com:8086", token)
	//defer client.Close()

	writeApi = client.WriteAPI(org, bucket)
	queryApi = client.QueryAPI(org)
}

func Query(id string, field string, start, end, window string) ([]Point, error) {
	//metric := fmt.Sprintf("%d", id)
	bucket := "zgwit"

	flux := "from(bucket: \"" + bucket + "\")\n"
	flux += "|> range(start: " + start + ", stop: " + end + ")\n"
	flux += "|> filter(fn: (r) => r[\"id\"] == \"" + id + "\")\n"
	flux += "|> filter(fn: (r) => r[\"_field\"] == \"" + field + "\")"
	flux += "|> aggregateWindow(every: " + window + ", fn: mean, createEmpty: false)\n"
	flux += "|> yield(name: \"mean\")"

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
