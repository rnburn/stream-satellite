package stream

import (
  // "fmt"
  "io"
  "net"
  "time"
  "sync"
  "errors"
  "github.com/golang/protobuf/proto"
  "github.com/rnburn/stream-satellite/egresspb"
)

type Server struct {
  mutex sync.Mutex
  listeners map[net.Listener]bool
  connections map[io.Closer]bool
}

func NewServer() *Server {
  result := &Server {
    listeners: make(map[net.Listener]bool),
    connections: make(map[io.Closer]bool),
  }
  return result
}

func (server *Server) Serve(listener net.Listener) error {
  server.mutex.Lock()
  if server.listeners == nil {
    server.mutex.Unlock()
    listener.Close()
    return errors.New("Server stopped")
  }
  server.listeners[listener] = true
  defer func() {
    server.mutex.Lock()
    if server.listeners != nil && server.listeners[listener] {
      listener.Close()
      delete(server.listeners, listener)
    }
    server.mutex.Unlock()
  }()

  for {
    connection, err := listener.Accept()
    if err != nil {
      // TODO: Add retry logic?
      return err
    }
    go handleConnection(connection)
  }
}

func readMessages(connection net.Conn, messageChan chan proto.Message) {
  session := NewSession(connection)
  for {
    session.ReadMessages(messageChan)
    // TODO: check and log error messages. Perhaps continue reading.
    break
  }
  close(messageChan)
}

func processMessage(message proto.Message) {
}

func flushSpans() {
}

func makeReport(initMessage *egresspb.StreamInitialization) *egresspb.ReportRequest {
  return &egresspb.ReportRequest {
    Reporter: &egresspb.Reporter {
      ReporterId: initMessage.ReporterId,
      Tags: initMessage.Tags,
    },
  }
}

func handleConnection(connection net.Conn) {
  messageChan := make(chan proto.Message)
  go readMessages(connection, messageChan)
  initMessage, ok := (<-messageChan).(*egresspb.StreamInitialization)
  if !ok {
    // TODO: error invalid initialization
    return
  }
  if initMessage == nil {
    // TODO: error initializing
    return
  }
  tickerChan := time.Tick(500 * time.Millisecond)
  report := makeReport(initMessage)
  for {
    select {
    case message := <-messageChan:
      span, ok := message.(*egresspb.Span)
      if !ok {
        // TODO: unexpected message
        continue
      }
      report.Spans = append(report.Spans, span)
    case <-tickerChan:
      // flush report
      if len(report.Spans) == 0 {
        continue
      }
      report = makeReport(initMessage)
    }
  }
}

func (server *Server) Stop() {
  server.mutex.Lock()
  defer server.mutex.Unlock()
  listeners := server.listeners
  server.listeners = nil
  connections := server.connections
  server.connections = nil

  for listener := range listeners {
    listener.Close()
  }

  for connection := range connections {
    connection.Close()
  }
}
