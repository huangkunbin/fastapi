package fastapi

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/funny/link"
	proto "github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"time"
)

func (app *App) newClientCodec(rw io.ReadWriter) (link.Codec, error) {
	return app.newCodec(rw, app.newResponse), nil
}

func (app *App) newServerCodec(rw io.ReadWriter) (link.Codec, error) {
	return app.newCodec(rw, app.newRequest), nil
}

func (app *App) newCodec(rw io.ReadWriter, newMessage func(byte, byte) (Message, error)) link.Codec {
	c := &codec{
		app:        app,
		conn:       rw.(net.Conn),
		reader:     bufio.NewReaderSize(rw, app.ReadBufSize),
		newMessage: newMessage,
	}
	c.headBuf = c.headData[:]
	return c
}

func (app *App) newRequest(serviceID, messageID byte) (Message, error) {
	if service := app.services[serviceID]; service != nil {
		if msg := service.(Service).NewRequest(messageID); msg != nil {
			return msg, nil
		}
		return nil, DecodeError{fmt.Sprintf("Unsupported Message Type: [%d, %d]", serviceID, messageID)}
	}
	return nil, DecodeError{fmt.Sprintf("Unsupported Service: [%d, %d]", serviceID, messageID)}
}

func (app *App) newResponse(serviceID, messageID byte) (Message, error) {
	if service := app.services[serviceID]; service != nil {
		if msg := service.(Service).NewResponse(messageID); msg != nil {
			return msg, nil
		}
		return nil, DecodeError{fmt.Sprintf("Unsupported Message Type: [%d, %d]", serviceID, messageID)}
	}
	return nil, DecodeError{fmt.Sprintf("Unsupported Service: [%d, %d]", serviceID, messageID)}
}

const packetHeadSize = 4 + 2

type codec struct {
	app        *App
	headBuf    []byte
	headData   [packetHeadSize]byte
	conn       net.Conn
	reader     *bufio.Reader
	newMessage func(byte, byte) (Message, error)
}

func (c *codec) Conn() net.Conn {
	return c.conn
}

type AddReq struct {
	A int32 `protobuf:"varint,1,opt,name=A,json=a" json:"A,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=B,json=b" json:"B,omitempty"`
}

func (c *codec) Receive() (msg interface{}, err error) {
	if c.app.RecvTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.app.RecvTimeout))
		defer c.conn.SetReadDeadline(time.Time{})
	}

	if _, err = io.ReadFull(c.reader, c.headBuf); err != nil {
		return
	}

	packetSize := int(binary.LittleEndian.Uint32(c.headBuf))

	if packetSize > c.app.MaxRecvSize {
		return nil, DecodeError{fmt.Sprintf("Too Large Receive Packet Size: %d", packetSize)}
	}

	packet := c.app.Pool.Alloc(packetSize)

	if _, err = io.ReadFull(c.reader, packet); err == nil {
		msg1, err1 := c.newMessage(c.headData[4], c.headData[5])
		if err1 == nil {
			func() {
				defer func() {
					if panicErr := recover(); panicErr != nil {
						err = DecodeError{panicErr}
					}
				}()
				err = proto.Unmarshal(packet, msg1)
				// msg1.UnmarshalPacket(packet)
			}()
			msg = msg1
		} else {
			err = err1
		}
	}

	c.app.Pool.Free(packet)
	return
}

func (c *codec) Send(m interface{}) (err error) {
	msg := m.(Message)

	// packetSize := msg.BinarySize()
	packetSize := proto.Size(msg)

	if packetSize > c.app.MaxSendSize {
		panic(EncodeError{fmt.Sprintf("Too Large Send Packet Size: %d", packetSize)})
	}

	packet := c.app.Pool.Alloc(packetHeadSize + packetSize)
	binary.LittleEndian.PutUint32(packet, uint32(packetSize))
	packet[4] = msg.ServiceID()
	packet[5] = msg.MessageID()

	func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				err = EncodeError{panicErr}
			}
		}()
		// msg.MarshalPacket(packet[packetHeadSize:])
		// proto.Marshal()
		pb := []byte{}
		pb, err = proto.Marshal(msg)
		copy(packet[packetHeadSize:], pb)

	}()

	if c.app.SendTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.app.SendTimeout))
		defer c.conn.SetWriteDeadline(time.Time{})
	}

	_, err = c.conn.Write(packet)
	c.app.Pool.Free(packet)
	return
}

func (c *codec) Close() error {
	return c.conn.Close()
}

type msgFormat struct {
	newMessage func(byte, byte) (Message, error)
}

func (f *msgFormat) EncodeMessage(msg interface{}) (buf []byte, err error) {
	msg2 := msg.(Message)
	defer func() {
		if panicErr := recover(); panicErr != nil {
			buf = nil
			err = EncodeError{panicErr}
		}
	}()
	buf = make([]byte, 2+proto.Size(msg2))
	buf[0] = msg2.ServiceID()
	buf[1] = msg2.MessageID()
	// msg2.MarshalPacket(buf[2:])
	pb, err := proto.Marshal(msg2)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	copy(buf[2:], pb)
	return
}

func (f *msgFormat) DecodeMessage(buf []byte) (msg interface{}, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = DecodeError{panicErr}
		}
	}()
	var msg2 Message
	msg2, err = f.newMessage(buf[0], buf[1])
	if err == nil {
		// msg2.UnmarshalPacket(buf[2:])
		err := proto.Unmarshal(buf[2:], msg2)
		if err != nil {
			log.Fatal(" send unmarshaling error: ", err)
		}
		msg = msg2
	}
	return
}
