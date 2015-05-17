// ringoexp is a simple server experiment using a ringbuffer for optimizing memory and processor performance.
package main

import (
	"flag"
	"runtime"
	"strings"

	"github.com/composer22/ringoexp/server"
)

var (
	log *server.RingoExpLogger = server.RingoExpLoggerNew()
)

// main is the main entry point for the application or server launch.
func main() {
	opts := server.Options{}
	var showVersion bool

	flag.StringVar(&opts.Name, "N", "", "Name of the server.")
	flag.StringVar(&opts.Name, "--name", "", "Name of the server.")
	flag.StringVar(&opts.Hostname, "H", server.DefaultHostname, "Hostname of the server.")
	flag.StringVar(&opts.Hostname, "--hostname", server.DefaultHostname, "Hostname of the server.")
	flag.IntVar(&opts.Port, "p", server.DefaultPort, "Port to listen on.")
	flag.IntVar(&opts.Port, "--port", server.DefaultPort, "Port to listen on.")
	flag.IntVar(&opts.MaxConns, "n", server.DefaultMaxConns, "Maximum incoming connections allowed (s + http).")
	flag.IntVar(&opts.MaxConns, "--connections", server.DefaultMaxConns, "Maximum incoming connections allowed (ws + http).")
	flag.BoolVar(&opts.IsPublisher, "I", server.DefaultIsPublisher, "Is the server a publisher (true) or a consumer?")
	flag.BoolVar(&opts.IsPublisher, "--is_publisher", server.DefaultIsPublisher, "Is the server a publisher (true) or a consumer?")
	flag.IntVar(&opts.RingSize, "r", server.DefaultRingSize, "Maximum ringbuffer size if publisher.")
	flag.IntVar(&opts.RingSize, "--ring_size", server.DefaultRingSize, "Maximum ringbuffer size if publisher.")
	flag.StringVar(&opts.ConsumerHostname, "U", server.DefaultConsumerHostname, "Hostname of the remote consumer server.")
	flag.StringVar(&opts.ConsumerHostname, "--consumer_hostname", server.DefaultConsumerHostname, "Hostname of the remote consumer server.")
	flag.IntVar(&opts.ConsumerPort, "T", server.DefaultConsumerPort, "Port of the remote consumer server.")
	flag.IntVar(&opts.ConsumerPort, "--consumer_port", server.DefaultConsumerPort, "Port of the remote consumer server.")
	flag.IntVar(&opts.MaxWorkers, "W", server.DefaultMaxWorkers, "Maximum outgoing worker connections allowed if publisher.")
	flag.IntVar(&opts.MaxWorkers, "--workers", server.DefaultMaxWorkers, "Maximum outgoing worker connections allowed if publisher.")
	flag.IntVar(&opts.MaxProcs, "X", server.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.IntVar(&opts.MaxProcs, "--procs", server.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.IntVar(&opts.ProfPort, "L", server.DefaultProfPort, "Profiler port to listen on.")
	flag.IntVar(&opts.ProfPort, "--profiler_port", server.DefaultProfPort, "Profiler port to listen on.")
	flag.BoolVar(&opts.Debug, "d", false, "Enable debugging output.")
	flag.BoolVar(&opts.Debug, "--debug", false, "Enable debugging output.")
	flag.BoolVar(&showVersion, "V", false, "Show version.")
	flag.BoolVar(&showVersion, "--version", false, "Show version.")
	flag.Usage = server.PrintUsageAndExit
	flag.Parse()

	// Version flag request?
	if showVersion {
		server.PrintVersionAndExit()
	}

	// Check additional params beyond the flags.
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "version":
			server.PrintVersionAndExit()
		case "help":
			server.PrintUsageAndExit()
		}
	}

	// Set thread and proc usage.
	if opts.MaxProcs > 0 {
		runtime.GOMAXPROCS(opts.MaxProcs)
	}
	log.Infof("NumCPU %d GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(-1))

	s := server.New(&opts)
	s.Start()
}
