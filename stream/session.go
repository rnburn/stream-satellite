package stream

import (
  "github.com/rnburn/stream-satellite/circlebuffer"
  "net"
)

type Session struct {
  connection net.Conn
  requireHeader bool
  requiredSize int
  buffer *circlebuffer.CircleBuffer
}

func NewSession(connection net.Conn) *Session {
  result := &Session{
    connection: connection,
    requireHeader: true,
    buffer: circlebuffer.NewCircleBuffer(1024*10),
  }
  return result
}

func (session *Session) ReadUntilNextMessage() error {
  return nil
}

func (session *Session) ConsumeHeader() {
}

func (session *Session) ConsumeMessage() bool {
  return false
}

func (session *Session) doRead() error {
  numRead, err := session.connection.Read(session.buffer.FreeSpaceAsSlice())
  if (err != nil) {
    return err
  }
  session.buffer.Grow(int64(numRead))
  return nil
}

func (session *Session) consumeHeader() {
}

func (session *Session) checkForNextMessage() bool {
  if session.requireHeader {
    if int(session.buffer.Size()) < session.requiredSize {
      return false
    }
    session.ConsumeHeader()
  }
  if int(session.buffer.Size()) < session.requiredSize {
    return false
  }
  return true
}
