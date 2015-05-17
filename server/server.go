// Package server implements a publisher/consumer ring-buffer experiment.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	// Allow dynamic profiling.
	_ "net/http/pprof"

	"github.com/composer22/ringoexp/logger"
	"github.com/composer22/ringoexp/ringbuffer"
	"golang.org/x/net/netutil"
	"golang.org/x/net/websocket"
)

// Server is the main structure that represents a server instance.
type Server struct {
	mu         sync.RWMutex        // For locking access to server attributes.
	running    bool                // Is the server running?
	info       *Info               // Basic server information used to run the server.
	opts       *Options            // Original options used to create the server.
	stats      *Stats              // Server statistics since it started.
	srvr       *http.Server        // HTTP/Socket server.
	ringbuffer []int               // A struct for the work
	rm         *ringbuffer.Manager // Manager wraps trackers for the work added to the ringbuffer
	quit       chan bool           // A channel to signal to web sockets and workers to close.
	log        *RingoExpLogger     // Log instance for recording error and other messages.
	wg         sync.WaitGroup      // Wait group to sync socket going down.
}

// New is a factory function that returns a new server instance.
func New(ops *Options) *Server {
	s := &Server{
		running: false,
		info: InfoNew(func(i *Info) {
			i.Name = ops.Name
			i.Hostname = ops.Hostname
			i.Port = ops.Port
			i.MaxConns = ops.MaxConns
			i.IsPublisher = ops.IsPublisher
			i.RingSize = ops.RingSize
			i.ConsumerHostname = ops.ConsumerHostname
			i.ConsumerPort = ops.ConsumerPort
			i.MaxWorkers = ops.MaxWorkers
			i.ProfPort = ops.ProfPort
			i.Debug = ops.Debug
		}),
		opts:       ops,
		stats:      StatsNew(),
		ringbuffer: make([]int, ops.RingSize),
		rm:         ringbuffer.ManagerNew(int64(ops.RingSize)),
		quit:       make(chan bool),
		log:        RingoExpLoggerNew(),
	}

	if s.info.Debug {
		s.log.SetLogLevel(logger.Debug)
	}

	// Setup the routes.
	http.Handle(wsRouteV1Ingest, websocket.Handler(s.ingestHandler))
	http.HandleFunc(httpRouteV1Alive, s.aliveHandler)
	http.HandleFunc(httpRouteV1Stats, s.statsHandler)
	s.srvr = &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.info.Hostname, s.info.Port),
	}

	s.handleSignals()
	return s
}

// PrintVersionAndExit prints the version of the server then exits.
func PrintVersionAndExit() {
	fmt.Printf("chattypantz version %s\n", version)
	os.Exit(0)
}

// Start spins up the server to accept incoming connections.
func (s *Server) Start() error {
	if s.isRunning() {
		return errors.New("Server already started.")
	}

	s.log.Infof("Starting ringoexp version %s\n", version)

	// Construct listener
	ln, err := net.Listen("tcp", s.srvr.Addr)
	if err != nil {
		s.log.Errorf("Cannot create net.listener: %s", err.Error())
		return err
	}
	// If we want to limit connections, created a special listener with a throttle.
	if s.info.MaxConns > 0 {
		ln = netutil.LimitListener(ln, s.info.MaxConns)
	}

	s.mu.Lock()

	// Pprof http endpoint for the profiler.
	if s.info.ProfPort > 0 {
		s.StartProfiler()
	}

	s.stats.Start = time.Now()
	s.running = true
	s.mu.Unlock()
	err = s.srvr.Serve(ln)

	// Done.
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	if err != nil {
		s.log.Emergencyf("Listen and Server Error: %s", err.Error())
	}
	return nil
}

// StartProfiler is called to enable dynamic profiling.
func (s *Server) StartProfiler() {
	s.log.Infof("Starting profiling on http port %d", s.opts.ProfPort)
	hp := fmt.Sprintf("%s:%d", s.info.Hostname, s.info.ProfPort)
	go func() {
		err := http.ListenAndServe(hp, nil)
		if err != nil {
			s.log.Emergencyf("Error starting profile monitoring service: %s", err)
		}
	}()
}

// Shutdown takes down the server gracefully back to an initialize state.
func (s *Server) Shutdown() {
	if !s.isRunning() {
		return
	}
	s.log.Infof("BEGIN server service stop.")
	s.log.Infof("Shutting down sockets...")
	close(s.quit)
	s.wg.Wait()
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	s.log.Infof("END server service stop.")
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			s.log.Infof("Server received signal: %v\n", sig)
			s.Shutdown()
			s.log.Infof("Server exiting.")
			os.Exit(0)
		}
	}()
}

// ingestHandler is the main entry point to handle chat connections to the client.
func (s *Server) ingestHandler(ws *websocket.Conn) {
	var ingester Ingester
	s.log.LogConnect(ws.Request())
	if s.opts.IsPublisher {
		ingester = IngestPublisherNew(ws, s.quit, s.ringbuffer, s.rm, s.log, &s.wg)
	} else {
		ingester = IngestConsumerNew(ws, s.quit, s.log, &s.wg)
	}
	ingester.Run()
}

// aliveHandler handles a client http:// "is the server alive?" request.
func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	s.log.LogConnect(r)
	s.initResponseHeader(w)
}

// statsHandler handles a client request for server information and statistics.
func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.log.LogConnect(r)
	s.initResponseHeader(w)
	s.mu.Lock()
	defer s.mu.Unlock()
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)
	b, _ := json.Marshal(
		&struct {
			Info    *Info             `json:"info"`
			Options *Options          `json:"options"`
			Stats   *Stats            `json:"stats"`
			Memory  *runtime.MemStats `json:"memStats"`
		}{
			Info:    s.info,
			Options: s.opts,
			Stats:   s.stats,
			Memory:  mStats,
		})
	w.Write(b)
}

// initResponseHeader sets up the common http response headers for the return of all json calls.
func (s *Server) initResponseHeader(w http.ResponseWriter) {
	h := w.Header()
	h.Add("Content-Type", "application/json;charset=utf-8")
	h.Add("Date", time.Now().UTC().Format(time.RFC1123Z))
	if s.info.Name != "" {
		h.Add("Server", s.info.Name)
	}
	h.Add("X-Request-ID", createV4UUID())
}

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
