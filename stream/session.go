package stream

import (
  "net"
  "github.com/golang/protobuf/proto"
  "github.com/rnburn/stream-satellite/egresspb"
)

type Session struct {
  connection net.Conn
  requireHeader bool
  requiredSize int
  buffer *CircleBuffer
}

func NewSession(connection net.Conn) *Session {
  result := &Session{
    connection: connection,
    requireHeader: true,
    requiredSize: PacketHeaderSize,
    buffer: NewCircleBuffer(1024*10),
  }
  return result
}

func (session *Session) ReadUntilNextMessage() error {
  for {
    if session.checkForNextMessage() {
      return nil
    }
    err := session.doRead()
    if err != nil {
      return err
    }
  }
}

func (session *Session) ReadMessages(messageChan chan proto.Message) error {
  err := session.ReadUntilNextMessage()
  if err != nil {
    return err
  }
  initialization := &egresspb.StreamInitialization{}
  err = session.ConsumeMessage(initialization)
  if err != nil {
    return err
  }
  messageChan <- initialization
  for {
    err = session.ReadUntilNextMessage()
    if err != nil {
      return err
    }
    span := &egresspb.Span{}
    err = session.ConsumeMessage(span)
    if err != nil {
      return err
    }
    messageChan <- span
  }
}

func (session *Session) consumeHeader() {
  header := DeserializePacketHeader(session.buffer.Peek(PacketHeaderSize))
  session.buffer.Consume(PacketHeaderSize)
  session.requireHeader = false
  session.requiredSize = int(header.BodySize)
  // if session.requiredSize > session.buffer.Capacity() - 1 {
    // TODO: Error
  // }
}

func (session *Session) ConsumeMessage(message proto.Message) error {
  // if session.requireHeader || session.buffer.Size() < session.requiredSize {
    // error
  // }
  err := proto.Unmarshal(session.buffer.Peek(int64(session.requiredSize)), message)
  session.buffer.Consume(int64(session.requiredSize))
  session.requireHeader = true
  session.requiredSize = PacketHeaderSize
  return err
}

func (session *Session) doRead() error {
  numRead, err := session.connection.Read(session.buffer.FreeSpaceAsSlice())
  if (err != nil) {
    return err
  }
  session.buffer.Grow(int64(numRead))
  return nil
}

func (session *Session) checkForNextMessage() bool {
  if session.requireHeader {
    if int(session.buffer.Size()) < session.requiredSize {
      return false
    }
    session.consumeHeader()
  }
  if int(session.buffer.Size()) < session.requiredSize {
    return false
  }
  return true
}
