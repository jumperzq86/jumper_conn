package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/jumperzq86/jumper_conn"
	"github.com/jumperzq86/jumper_conn/def"
	"github.com/jumperzq86/jumper_conn/interf"
	"github.com/jumperzq86/jumper_conn/util"
	"github.com/jumperzq86/jumper_transform"
	jtd "github.com/jumperzq86/jumper_transform/def"
	jti "github.com/jumperzq86/jumper_transform/interf"
)

const addr = "localhost:8801"

func main() {
	var conn net.Conn
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial failed, err: %s\n", err)
		return
	}

	//note: transform 可以只定义一个，他本身是线程安全对
	ts := jumper_transform.Newtransform()
	ts.AddOp(jtd.PacketBinary, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(c net.Conn) {
		defer func() {
			fmt.Printf("============")
			wg.Done()
		}()

		var h Handler

		tcpOp := def.ConnOptions{
			MaxMsgSize:     def.MaxMsgSize,
			ReadTimeout:    def.ReadTimeout,
			WriteTimeout:   def.WriteTimeout,
			AsyncWriteSize: def.AsyncWriteSize,
			Side:           def.ClientSide,
		}

		jconn, err := jumper_conn.NewtcpConn(c, &tcpOp, &h)
		if err != nil {
			fmt.Printf("new tcp conn failed. err: %s\n", err)
			return
		}

		h.Init(jconn, ts)

		fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

		//send hello
		state := fmt.Sprintf("this is tcp_client %s, hello", jconn.LocalAddr())

		msg := &jti.Message{
			Type:    1,
			Content: []byte(state),
		}

		var output []byte
		err = h.Execute(jtd.Forward, msg, &output)
		if err != nil {
			fmt.Printf("transform failed, err: %s\n", err)
			return
		}

		length := len(output)
		head := make([]byte, 4)
		binary.BigEndian.PutUint32(head, uint32(length))

		sendMsg := make([]byte, 0, def.TcpHeadSize+length)
		sendMsg = append(sendMsg, head...)
		sendMsg = append(sendMsg, output...)

		h.Write(sendMsg)
		if err != nil {
			fmt.Printf("write failed, err: %s\n", err)
			return
		}

	}(conn)

	wg.Wait()
}

type Handler struct {
	interf.Conn
	jti.Transform
}

func (this *Handler) Init(conn interf.Conn, ts jti.Transform) {
	this.Conn = conn
	this.Transform = ts
	this.Run()
}

func (this *Handler) OnMessage(data []byte) error {
	util.TraceLog("handler.OnMessage")
	fmt.Printf("handler get data: %v\n", data)
	return nil
}

func (this *Handler) OnClose(err error) {
	util.TraceLog("handler.OnClose")
}
