package main

import (
  "github.com/rnburn/stream-satellite/stream"
  "fmt"
  "os"
  "net"
)

func main() {
  server := stream.NewServer()

  listener, err := net.Listen("tcp", ":8081")
  if err != nil {
    fmt.Printf("Failed to listen")
    os.Exit(1)
  }
  server.Serve(listener)
}
