package log

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Write(Data) error
}

type Data struct {
	Method           string
	Endpoint         string
	Code             int
	Referer          string
	Duration         int64
	IP               string
	UserAgent        string
	UserAgentVersion string
	OS               string
	OSVersion        string
	Device           string
	DeviceType       string
	User             string
	Time             time.Time
}

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	return &LogrusLogger{
		logger: logrus.New(),
	}
}

func (l *LogrusLogger) Write(data Data) error {
	l.logger.WithFields(logrus.Fields{
		"Method":           data.Method,
		"Endpoint":         data.Endpoint,
		"Code":             data.Code,
		"Referer":          data.Referer,
		"Duration":         data.Duration,
		"IP":               data.IP,
		"UserAgent":        data.UserAgent,
		"UserAgentVersion": data.UserAgentVersion,
		"OS":               data.OS,
		"OSVersion":        data.OSVersion,
		"Device":           data.Device,
		"DeviceType":       data.DeviceType,
		"User":             data.User,
		"Time":             data.Time,
	}).Infof("request")
	return nil
}

type InfluxLogger struct {
	client influxdb2.Client
}

const (
	influxOrg    = "test"
	influxBucket = "webmetrics"
)

func NewInfluxLogger(client influxdb2.Client) *InfluxLogger {
	return &InfluxLogger{
		client: client,
	}
}

func (i *InfluxLogger) Write(data Data) error {

	writeAPI := i.client.WriteAPIBlocking(influxOrg, influxBucket)
	p := influxdb2.NewPoint(
		"request",
		map[string]string{
			"method":      data.Method,
			"endpoint":    data.Endpoint,
			"code":        fmt.Sprintf("%d", data.Code),
			"ip":          data.IP,
			"os":          data.OS,
			"decive_type": data.DeviceType,
			"user":        data.User,
		},
		map[string]interface{}{
			"os_version":        data.OSVersion,
			"device":            data.Device,
			"useragent":         data.UserAgent,
			"useragent_version": data.UserAgentVersion,
			"duration":          data.Duration,
			"referer":           data.Referer,
		},
		data.Time,
	)
	return writeAPI.WritePoint(context.Background(), p)
}
