package server

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	testInfoExpectedJSONResult = `{"version":"9.8.7","UUID":"ABCDEFGHIJKLMNOPQRSTUVWXYZ",` +
		`"name":"Test Server","hostname":"1.2.3.4","port":9999,"maxConns":9998,"isPublisher":` +
		`true,"ringBuffSize":9997,"consumerHostname":"4.5.6.7","consumerPort":9996,` +
		`"maxWorkers":9995,"profPort":9994,"debugEnabled":true}`
)

func TestInfoNew(t *testing.T) {
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Name = "Test Server"
		i.Hostname = "1.2.3.4"
		i.Port = 9999
		i.MaxConns = 9998
		i.IsPublisher = true
		i.RingBuffSize = 9997
		i.ConsumerHostname = "4.5.6.7"
		i.ConsumerPort = 9996
		i.MaxWorkers = 9995
		i.ProfPort = 9994
		i.Debug = true
	})
	tp := reflect.TypeOf(info)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Info not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Info not created as a struct.")
	}
	if tp.Name() != "Info" {
		t.Fatalf("Info struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Info struct is empty.")
	}
}

func TestInfoString(t *testing.T) {
	t.Parallel()
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Name = "Test Server"
		i.Hostname = "1.2.3.4"
		i.Port = 9999
		i.MaxConns = 9998
		i.IsPublisher = true
		i.RingBuffSize = 9997
		i.ConsumerHostname = "4.5.6.7"
		i.ConsumerPort = 9996
		i.MaxWorkers = 9995
		i.ProfPort = 9994
		i.Debug = true
	})
	actual := fmt.Sprint(info)
	if actual != testInfoExpectedJSONResult {
		t.Errorf("Info not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testInfoExpectedJSONResult, actual)
	}
}
