package stream

import (
  "fmt"
  "io"
  "net"
  "sync"
  "errors"
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
    go server.handleConnection(connection)
  }
}

func (server *Server) handleConnection(connection net.Conn) {
  session := NewSession(connection)
  err := session.ReadUntilNextMessage()
  if err != nil {
    // TODO: Check for EOF? Log other errors?
    return
  }
  initialization := &egresspb.StreamInitialization{}
  err = session.ConsumeMessage(initialization)
  if err != nil {
    return
  }
  for {
    err = session.ReadUntilNextMessage()
    if err != nil {
      // TODO: Check for EOF? Log other errors?
      return
    }
    span := &egresspb.Span{}
    err = session.ConsumeMessage(span)
    fmt.Printf("Received span %d\n", span.SpanContext.TraceId)
    if err != nil {
      return
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
