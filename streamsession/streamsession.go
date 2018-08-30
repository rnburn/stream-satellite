package streamsession
import (
  "github.com/rnburn/stream-satellite/circlebuffer"
  "net"
)

type Session struct {
  connection net.Conn
  requireHeader bool
  requiredSize int
  buffer circlebuffer.CircleBuffer
}

func NewSession(connection net.Conn) (*Session, error) {

  return nil, nil
}

func (session *Session) ReadUntilNextMessage() error {
  return nil
}

func (session *Session) ConsumeMessage() bool {
  return false
}

func (session *Session) doRead() error {
  return nil
}

func (session *Session) consumeHeader() {
}

func (session *Session) checkForNextMessage() bool {
  return true
}
