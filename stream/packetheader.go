package stream

import (
  "encoding/binary"
)

const PacketHeaderSize = 5

type PacketHeader struct {
  Version uint8
  BodySize uint32
}

func DeserializePacketHeader(data []byte) PacketHeader {
  // if len(data) != PacketHeaderSize {
    // TODO: error
  // }
  return PacketHeader {
    Version: uint8(data[0]),
    BodySize: binary.LittleEndian.Uint32(data[1:]),
  }
}
