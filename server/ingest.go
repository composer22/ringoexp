package server

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Ingester interface {
	Run()
	receive()
	signalTrap()
	shutDown()
}

// Ingest is a wrapper around an incoming connection to a publishing/consuming server.
type Ingest struct {
	start time.Time       // The start time of the connection.
	ws    *websocket.Conn // The socket to the remote client.
	quit  chan bool       // Channel to signal service should disconnect and close down from server.
	done  chan bool       // Channel to tell signalTrap() should close go routine.
	log   *RingoExpLogger // Log file out.
	swg   *sync.WaitGroup // Server synchronization of server close.
	wg    sync.WaitGroup  // Synchronization of channel close.
}

// IngestNew is a factory function that returns a new Ingest instance
func IngestNew(w *websocket.Conn, q chan bool, l *RingoExpLogger, swg *sync.WaitGroup) *Ingest {
	return &Ingest{
		ws:   w,
		quit: q,
		done: make(chan bool),
		log:  l,
		swg:  swg,
	}
}

// Run starts the event loop that manages the receiving of information from the remote client.
func (i *Ingest) Run() {
	i.start = time.Now()
	i.swg.Add(1)      // We let the big boss know so it can micromanage us on server close.
	i.wg.Add(1)       //   but we also have our own signal to signalTrap().
	go i.signalTrap() // Spawn a background task to check for close requests.
	i.receive()       // Then wait on incoming requests.
}

// receive polls and handles any commands or information sent from the remote client.
func (i *Ingest) receive() {
	defer i.swg.Done()
	remoteAddr := i.ws.Request().RemoteAddr
	var req []byte
	var err error
	for {
		// Receive data.

		if err = websocket.Message.Receive(i.ws, &req); err != nil {
			switch {
			case err.Error() == "EOF":
				i.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
				i.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"): // cntl-c safety.
				i.shutDown()
			default:
				i.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
				i.shutDown()
			}
			return
		}

		// Implement your own version of this and perform work

		// ACK back we received.
		if err = websocket.Message.Send(i.ws, 'a'); err != nil {
			switch {
			case err.Error() == "EOF":
				i.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
				i.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"): // cntl-c safety.
				i.shutDown()
			default:
				i.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
				i.shutDown()
			}
			return
		}
	}
}

// signalTrap is a go routine used to poll quit requests from the server.
func (i *Ingest) signalTrap() {
	defer i.wg.Done()
	for {
		select {
		case <-i.quit: // Server or receiver shutdown signal.
			i.ws.Close() // Close the socket sends signal to receive()
			return
		case <-i.done:
			i.ws.Close() // Close the socket sends signal to receive()
			return
		default:
			runtime.Gosched()
		}
	}
}

// shutDown shuts down sending/receiving.
func (i *Ingest) shutDown() {
	close(i.done) // Signal to signalTrap()
	i.wg.Wait()   // Wait for signalTrap()
}
