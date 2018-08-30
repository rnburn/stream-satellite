package main

import (
  "github.com/rnburn/stream-satellite/stream"
)

func main() {
  server := stream.NewServer()
  server.Stop()
}
