package stream

// import (
//   "fmt"
// )

func min(x, y int64) int64 {
    if x < y {
        return x
    }
    return y
  }

type CircleBuffer struct {
  data []byte
  head int64
  tail int64
}

func NewCircleBuffer(capacity int64) *CircleBuffer {
  // if capacity <= 0 {
  //   return nil, fmt.Errorf("Capacity must be positive")
  // }
  result := &CircleBuffer {
    data: make([]byte, capacity),
  }
  return result
}

func (circleBuffer *CircleBuffer) Capacity() int64 {
  return int64(len(circleBuffer.data))
}

func (circleBuffer *CircleBuffer) freeSpace() int64 {
  if circleBuffer.head >= circleBuffer.tail {
    return int64(len(circleBuffer.data)) - 1 - (circleBuffer.head - circleBuffer.tail)
  }
  return circleBuffer.tail - circleBuffer.head - 1
}

func (circleBuffer *CircleBuffer) FreeSpaceAsSlice() []byte {
  n := min(circleBuffer.freeSpace(), circleBuffer.Capacity() - circleBuffer.head)
  return circleBuffer.data[circleBuffer.head:circleBuffer.head+n]
}

func (circleBuffer *CircleBuffer) Grow(numBytes int64) {
  circleBuffer.head = (circleBuffer.head + numBytes) % circleBuffer.Capacity()
}

func (circleBuffer *CircleBuffer) Size() int64 {
  if circleBuffer.head >= circleBuffer.tail {
    return circleBuffer.head - circleBuffer.tail
  }
  return circleBuffer.Capacity() - (circleBuffer.tail - circleBuffer.head)
}

func (circleBuffer *CircleBuffer) Peek(numBytes int64) []byte {
  size1 := min(circleBuffer.Capacity() - circleBuffer.tail, numBytes)
  slice1 := circleBuffer.data[circleBuffer.tail:circleBuffer.tail + size1]
  if size1 == numBytes {
    return slice1
  }
  size2 := numBytes - size1
  slice2 := circleBuffer.data[0:size2]
  result := make([]byte, size1 + size2)
  copy(slice1, result[0:])
  copy(slice2, result[size1:])
  return result
}

func (circleBuffer *CircleBuffer) Consume(numBytes int64) {
  // TODO: Check for errors here
  circleBuffer.tail = (circleBuffer.tail + numBytes) % circleBuffer.Capacity()
}
