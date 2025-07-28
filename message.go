package wayland

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
)

type Message struct {
	ObjectId uint32
	OpCode   uint16
	Size     uint16
	Body     []byte
	n        uint16
	mu       sync.Mutex
	Fds      []int
	nextFd   int
}

func (msg *Message) Bytes() []byte {
	if msg == nil {
		return nil
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, msg.ObjectId)
	binary.Write(&buf, binary.LittleEndian, msg.OpCode)
	binary.Write(&buf, binary.LittleEndian, uint16(0))
	binary.Write(&buf, binary.LittleEndian, msg.Body)

	result := buf.Bytes()
	binary.LittleEndian.PutUint16(result[6:8], uint16(buf.Len()))
	return result
}

func (msg *Message) String() string {
	if msg == nil {
		return ""
	}

	return fmt.Sprintf("objectId: %d, size: %d, opcode: %d, body: %x", msg.ObjectId, msg.Size, msg.OpCode, msg.Body)
}

func (msg *Message) ReadUint32() uint32 {
	var result uint32
	msg.mu.Lock()
	binary.Read(bytes.NewReader(msg.Body[msg.n:msg.n+uint16(4)]), binary.LittleEndian, &result)
	msg.n += 4
	msg.mu.Unlock()
	return result
}

func (msg *Message) ReadInt32() int32 {
	return int32(msg.ReadUint32())
}

func (msg *Message) ReadString() string {
	var length uint32
	msg.mu.Lock()
	binary.Read(bytes.NewReader(msg.Body[msg.n:msg.n+uint16(4)]), binary.LittleEndian, &length)
	msg.n += 4
	result := string(msg.Body[msg.n : msg.n+uint16(length)])
	msg.n += uint16(length)
	if length%4 != 0 {
		msg.n += uint16(4 - length%4)
	}
	msg.mu.Unlock()
	return strings.TrimSuffix(result, "\x00")
}

func (msg *Message) ReadFixed() Fixed {
	return Fixed(msg.ReadUint32())
}

func (msg *Message) ReadArray() []uint32 {
	// TODO: implement reading arrays
	return nil
}

func (msg *Message) ReadFd() int {
	if msg.nextFd < len(msg.Fds) {
		return msg.Fds[msg.nextFd]
	}
	return 0
}

func (msg *Message) WithFds(fd ...int) *Message {
	msg.Fds = fd
	return msg
}

func NewMessage(objectId uint32, opcode uint16, args ...any) *Message {
	result := &Message{
		ObjectId: objectId,
		OpCode:   opcode,
	}

	var buf bytes.Buffer
	for _, arg := range args {
		switch arg := arg.(type) {
		case int:
			binary.Write(&buf, binary.LittleEndian, uint32(arg))
		case int32:
			binary.Write(&buf, binary.LittleEndian, uint32(arg))
		case uint32:
			binary.Write(&buf, binary.LittleEndian, arg)
		case string:
			binary.Write(&buf, binary.LittleEndian, uint32(len(arg)+1))
			binary.Write(&buf, binary.LittleEndian, []byte(arg+"\x00"))
			if (len(arg)+1)%4 != 0 {
				binary.Write(&buf, binary.LittleEndian, make([]byte, 4-(len(arg)+1)%4))
			}
		case []uint32:
			// TODO: writing arrays
		default:
			return nil
		}
	}

	result.Body = buf.Bytes()
	result.Size = uint16(8 + len(result.Body))

	return result
}
