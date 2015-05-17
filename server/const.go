package server

const (
	version                 = "0.1.0"     // Application and server version.
	DefaultHostname         = "localhost" // The hostname of the server.
	DefaultPort             = 6660        // Port to receive requests: see IANA Port Numbers.
	DefaultProfPort         = 0           // Profiler port to receive requests. *
	DefaultConsumerHostname = "localhost" // The hostname of the remote consumer server.
	DefaultConsumerPort     = 6660        // The port of the remote consumer server.
	DefaultIsPublisher      = true        // Is the server a publisher? true = pub; false = consumer.
	DefaultMaxConns         = 0           // Maximum number of incoming connections allowed (ws and/or web). *
	DefaultMaxWorkers       = 1024        // Maximum number of outgoing worker connections allowed ( to consumer).
	DefaultRingSize         = 4096        // Ring buffer size. Note this should be a power of 2. Ignored if consumer.
	DefaultMaxProcs         = 0           // Maximum number of computer processors to utilize. *

	// * zeros = no change or no limitation or not enabled.

	// http and ws routes for servers.
	wsRouteV1Ingest  = "/v1.0/ingest" // For the publisher or subscriber, this is the external endpoint.
	httpRouteV1Alive = "/v1.0/alive"
	httpRouteV1Stats = "/v1.0/stats"
)
