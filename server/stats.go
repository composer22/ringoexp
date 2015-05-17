package server

import (
	"encoding/json"
	"time"
)

// Stats contains runtime statistics for the server.
type Stats struct {
	Start time.Time `json:"startTime"` // The start time of the server.
}

// StatsNew is a factory function that returns a new instance of statistics.
// options is an optional list of functions that initialize the structure
func StatsNew(opts ...func(*Stats)) *Stats {
	s := &Stats{
		Start: time.Now(),
	}
	for _, f := range opts {
		f(s)
	}
	return s
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Stats) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
