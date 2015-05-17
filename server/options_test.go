package server

import (
	"fmt"
	"testing"
)

const (
	testOptionsExpectedJSONResult = `{"name":"Test Server","hostname":"1.2.3.4",` +
		`"port":9999,"maxConns":9998,"isPublisher":true,"ringBuffSize":9997,"consumerHostname":` +
		`"5.6.7.8","consumerPort":9996,"maxWorkers":9995,"maxProcs":9994,"profPort":9993,` +
		`"debugEnabled":true}`
)

func TestOptionsString(t *testing.T) {
	t.Parallel()
	opts := &Options{
		Name:             "Test Server",
		Hostname:         "1.2.3.4",
		Port:             9999,
		MaxConns:         9998,
		IsPublisher:      true,
		RingBuffSize:     9997,
		ConsumerHostname: "5.6.7.8",
		ConsumerPort:     9996,
		MaxWorkers:       9995,
		MaxProcs:         9994,
		ProfPort:         9993,
		Debug:            true,
	}
	actual := fmt.Sprint(opts)
	if actual != testOptionsExpectedJSONResult {
		t.Errorf("Options not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testOptionsExpectedJSONResult, actual)
	}
}
