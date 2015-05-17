package server

import "encoding/json"

// Options represents parameters that are passed to the application to be used in constructing
// the server.
type Options struct {
	Name             string `json:"name"`             // The name of the server.
	Hostname         string `json:"hostname"`         // The hostname of the server.
	Port             int    `json:"port"`             // The default port of the server.
	MaxConns         int    `json:"maxConns"`         // The maximum incoming connections allowed.
	IsPublisher      bool   `json:"isPublisher"`      // Is the server a publisher (true) or a consumer (false)?
	RingSize         int    `json:"ringSize"`         // The ring buffer size in slots, if publisher else ignored.
	ConsumerHostname string `json:"consumerHostname"` // The hostname of the consumer server if this is a publisher.
	ConsumerPort     int    `json:"consumerPort"`     // The port of the consumer server if this is a publisher.
	MaxWorkers       int    `json:"maxWorkers"`       // The maximum outgoing workers allowed if publisher.
	MaxProcs         int    `json:"maxProcs"`         // The maximum number of processor cores available.
	ProfPort         int    `json:"profPort"`         // The profiler port of the server.
	Debug            bool   `json:"debugEnabled"`     // Is debugging enabled in the application or server.
}

// String is an implentation of the Stringer interface so the structure is returned as a string
// to fmt.Print() etc.
func (o *Options) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
