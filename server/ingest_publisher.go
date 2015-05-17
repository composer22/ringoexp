package server

import (
	"bytes"
	"encoding/binary"
	"strings"
	"sync"

	"github.com/composer22/ringoexp/ringbuffer"
	"golang.org/x/net/websocket"
)

// IngestPublisher is a wrapper around an incoming connection to a publishing server.
type IngestPublisher struct {
	*Ingest
	rb []int               // Ringbuffer for the data.
	rm *ringbuffer.Manager // Synchronizer for work.
}

// IngestPublisherrNew is a factory function that returns a new IngestPublisher instance
func IngestPublisherNew(w *websocket.Conn, q chan bool, r []int, m *ringbuffer.Manager,
	l *RingoExpLogger, swg *sync.WaitGroup) *IngestPublisher {
	return &IngestPublisher{
		Ingest: IngestNew(w, q, l, swg),
		rb:     r,
		rm:     m,
	}
}

// receive polls and handles any commands or information sent from the remote client.
func (i *IngestPublisher) receive() {
	defer i.swg.Done()
	remoteAddr := i.ws.Request().RemoteAddr
	var req []byte
	var err error
	var mask = i.rm.Leader.Mask()
	for {
		// Receive data.

		if err = websocket.Message.Receive(i.ws, &req); err != nil {
			switch {
			case err.Error() == "EOF":
				i.log.Infof("Client %s disconnected.", remoteAddr)
				i.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"): // cntl-c safety.
				i.shutDown()
			default:
				i.log.Infof("Couldn't receive from %s. Error: %s", remoteAddr, err.Error())
				i.shutDown()
			}
			return
		}

		// Store into ring.
		indx := i.rm.Leader.Reserve(1)
		j, _ := binary.ReadVarint(bytes.NewBuffer(req))
		i.rb[indx&mask] = int(j)
		i.rm.Leader.Commit(indx, indx)

		// ACK back we received.
		if err = websocket.Message.Send(i.ws, 'a'); err != nil {
			switch {
			case err.Error() == "EOF":
				i.log.Infof("Client %s disconnected.", remoteAddr)
				i.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"): // cntl-c safety.
				i.shutDown()
			default:
				i.log.Infof("Couldn't receive from %s. Error: %s", remoteAddr, err.Error())
				i.shutDown()
			}
			return
		}
	}
}
