package influxdb

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/zgwit/iot-master/v4/history"
	"time"
)

type Historian struct {
	Url, Org, Bucket, Token string

	client influxdb2.Client
	writer api.WriteAPI
	reader api.QueryAPI
}

func (h *Historian) Open() {
	h.client = influxdb2.NewClient(h.Url, h.Token)
	h.writer = h.client.WriteAPI(h.Org, h.Bucket)
	h.reader = h.client.QueryAPI(h.Org)
}

func (h *Historian) Close() {
	h.client.Close()
}

func (h *Historian) Write(table, id string, timestamp int64, values map[string]any) error {
	h.writer.WritePoint(write.NewPoint(table, map[string]string{"id": id}, values, time.UnixMilli(timestamp)))
	return nil
}

func (h *Historian) Query(table, id, name, start, end, window, method string) ([]history.Point, error) {
	flux := "from(Bucket: \"" + h.Bucket + "\")\n"
	flux += "|> range(start: " + start + ", stop: " + end + ")\n"
	flux += "|> filter(fn: (r) => r[\"_measurement\"] == \"" + table + "\")\n"
	flux += "|> filter(fn: (r) => r[\"id\"] == \"" + id + "\")\n"
	flux += "|> filter(fn: (r) => r[\"_field\"] == \"" + name + "\")"
	flux += "|> aggregateWindow(every: " + window + ", fn: " + method + ", createEmpty: false)\n"
	flux += "|> yield(name: \"" + method + "\")"

	result, err := h.reader.Query(context.Background(), flux)
	if err != nil {
		return nil, err
	}

	var records []history.Point
	for result.Next() {
		//result.TableChanged() 查询多个数值的情况？？？
		records = append(records, history.Point{
			Value: result.Record().Value(),
			Time:  result.Record().Time().UnixMilli(),
		})
	}
	return records, result.Err()
}
