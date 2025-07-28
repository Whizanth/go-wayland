package wayland

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sys/unix"
)

type Client struct {
	conn      net.Conn
	objectId  uint32
	rmu       sync.Mutex
	listeners map[int]map[int]func(message *Message)
}

// NewObjectId returns a new object ID that hasn't been used yet
func (client *Client) NewObjectId() uint32 {
	client.objectId++
	return client.objectId
}

// Read waits for and reads the next message from the compositor
func (client *Client) Read() *Message {
	// TODO: read fds

	client.rmu.Lock()

	header := make([]byte, 8)
	if _, err := io.ReadFull(client.conn, header); err != nil {
		client.rmu.Unlock()
		fmt.Printf("error reading from socket: %v\n", err)
		return nil
	}

	size := binary.LittleEndian.Uint16(header[6:8])

	body := make([]byte, size-8)
	if _, err := io.ReadFull(client.conn, body); err != nil {
		client.rmu.Unlock()
		fmt.Printf("error reading message body: %v\n", err)
		return nil
	}

	client.rmu.Unlock()

	result := &Message{
		ObjectId: binary.LittleEndian.Uint32(header[0:4]),
		Size:     size,
		OpCode:   binary.LittleEndian.Uint16(header[4:6]),
		Body:     body,
	}
	//fmt.Println("<<< " + result.String())
	return result
}

// Write sends a message to the compositor, optionally passing through any file descriptors
func (client *Client) Write(msg *Message) error {
	//fmt.Println(">>> " + msg.String())
	if len(msg.Fds) == 0 {
		_, err := client.conn.Write(msg.Bytes())
		return err
	} else {
		_, _, err := client.conn.(*net.UnixConn).WriteMsgUnix(msg.Bytes(), unix.UnixRights(msg.Fds...), nil)
		return err
	}
}

// Request composes a message and sends it as request to the compositor
func (client *Client) Request(objectId uint32, opcode uint16, args ...any) {
	client.Write(NewMessage(objectId, opcode, args...))
}

// On calls listener if the client receives an event with the specified objectId and opcode
func (client *Client) On(objectId uint32, opcode uint16, listener func(message *Message)) chan struct{} {
	wait := make(chan struct{})

	listeners, ok := client.listeners[int(objectId)]
	if !ok {
		client.listeners[int(objectId)] = make(map[int]func(message *Message))
		listeners = client.listeners[int(objectId)]
	}
	listeners[int(opcode)] = func(msg *Message) {
		listener(msg)

		if len(wait) == 0 {
			go func() {
				wait <- struct{}{}
			}()
		}
	}

	return wait
}

// Listen reads and delivering messages to the appropriate listener if registered
func (client *Client) Listen() {
	for {
		msg := client.Read()
		if listeners, ok := client.listeners[int(msg.ObjectId)]; ok {
			if listener, ok := listeners[int(msg.OpCode)]; ok {
				listener(msg)
			}
		}
	}
}

// Close disconnects the client
func (client *Client) Close() {
	client.conn.Close()
}

// NewClient creates a new client and tries to connect to the compositor
func NewClient() (*Client, error) {
	address := os.Getenv("WAYLAND_SOCKET")
	if address == "" {
		waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
		xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
		if waylandDisplay != "" && xdgRuntimeDir != "" {
			address = filepath.Join(xdgRuntimeDir, waylandDisplay)
		} else if xdgRuntimeDir != "" {
			address = filepath.Join(xdgRuntimeDir, "wayland-0")
		} else {
			return nil, errors.New("no Wayland compositor detected")
		}
	}

	result := &Client{
		listeners: make(map[int]map[int]func(message *Message)),
	}

	var err error
	result.conn, err = net.Dial("unix", address)
	if err != nil {
		return nil, err
	}

	return result, nil
}
