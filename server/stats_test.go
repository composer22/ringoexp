package server

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	testStatsExpectedJSONResult = `{"startTime":"2006-01-02T13:24:56Z"}`
)

func TestStatsNew(t *testing.T) {
	s := StatsNew()
	tp := reflect.TypeOf(s)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Stats not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Stats not created as a struct.")
	}
	if tp.Name() != "Stats" {
		t.Fatalf("Stats struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Stats struct is empty.")
	}
}

func TestStatString(t *testing.T) {
	t.Parallel()
	mockTime, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 13:24:56 -0000")
	s := StatsNew(func(sts *Stats) {
		sts.Start = mockTime
	})
	actual := fmt.Sprint(s)
	if actual != testStatsExpectedJSONResult {
		t.Errorf("Stats not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testStatsExpectedJSONResult, actual)
	}
}
