package server

import "encoding/json"

// Info provides basic config information to/about the running server.
type Info struct {
	Version          string `json:"version"`          // Version of the server.
	UUID             string `json:"UUID"`             // Unique ID of the server.
	Name             string `json:"name"`             // The name of the server.
	Hostname         string `json:"hostname"`         // The hostname of the server.
	Port             int    `json:"port"`             // Port the server is listening on.
	MaxConns         int    `json:"maxConns"`         // The maximum concurrent clients accepted.
	IsPublisher      bool   `json:"isPublisher"`      // Is the server a publisher (true) or a consumer (false)?
	RingSize         int    `json:"ringSize"`         // The ring buffer size in slots, if publisher else ignored.
	ConsumerHostname string `json:"consumerHostname"` // The hostname of the consumer server if this is a publisher.
	ConsumerPort     int    `json:"consumerPort"`     // The port of the consumer server if this is a publisher.
	MaxWorkers       int    `json:"maxWorkers"`       // The maximum outgoing workers allowed if publisher.
	ProfPort         int    `json:"profPort"`         // Profiler port the server is listening on.
	Debug            bool   `json:"debugEnabled"`     // Is debugging enabled on the server.
}

// InfoNew is a factory function that returns a new instance of Info.
// opts is an optional list of functions that initialize the structure
func InfoNew(opts ...func(*Info)) *Info {
	inf := &Info{
		Version: version,
		UUID:    createV4UUID(),
	}
	for _, f := range opts {
		f(inf)
	}
	return inf
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (i *Info) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
