package server

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

// IngestConsumer is a wrapper around an incoming connection to a publishing server.
type IngestConsumer struct {
	*Ingest
	// add database pointer here for receive()
}

// IngestConsumerNew is a factory function that returns a new IngestConsumer instance
func IngestConsumerNew(w *websocket.Conn, q chan bool, l *RingoExpLogger, swg *sync.WaitGroup) *IngestConsumer {
	return &IngestConsumer{
		IngestNew(w, q, l, swg),
	}
}

// receive polls and handles any commands or information sent from the remote client.
func (i *IngestConsumer) receive() {
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

		// Store value to Database

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
