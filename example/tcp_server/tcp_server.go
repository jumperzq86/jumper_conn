package main

import (
	"fmt"
	"net"
	"time"

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
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen failed, err: %s\n", err)
		return
	}

	//note: transform 可以只定义一个，他本身是线程安全对
	ts := jumper_transform.Newtransform()
	ts.AddOp(jtd.PacketBinary, nil)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed, err: %s\n", err)
			return
		}

		go func(c net.Conn) {
			var h Handler

			tcpOp := def.ConnOptions{
				MaxMsgSize:     def.MaxMsgSize,
				ReadTimeout:    def.ReadTimeout,
				WriteTimeout:   def.WriteTimeout,
				AsyncWriteSize: def.AsyncWriteSize,
				Side:           def.ServerSide,
			}

			jconn, err := jumper_conn.NewtcpConn(c, &tcpOp, &h)
			if err != nil {
				fmt.Printf("new tcp conn failed. err: %s\n", err)
				return
			}

			h.Init(jconn, ts)
			//h.Set("random_num", 10000)

			fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

		}(conn)
	}
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
	defer util.TraceLog("handler.OnMessage")()
	fmt.Printf("handler get data: %v\n", data)
	var msg jti.Message
	err := this.Execute(jtd.Backward, data, &msg)
	if err != nil {
		fmt.Printf("transform failed, err: %s\n", err)
		return err
	}
	//num := this.Get("random_num").(int)
	//fmt.Printf("num: %d\n", num)
	//this.Del("random_num")
	//if n := this.Get("random_num"); n != nil {
	//	fmt.Printf("delete failed")
	//	return errors.New("delete failed.")
	//}

	fmt.Printf("type: %d, content: %s\n", msg.Type, msg.Content)

	time.Sleep(1 * time.Second)
	this.Close()

	return nil
}

func (this *Handler) OnClose(err error) {
	defer util.TraceLog("handler.OnClose")()
}
