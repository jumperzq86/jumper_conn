package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jumperzq86/jumper_conn"
	"github.com/jumperzq86/jumper_conn/def"
	"github.com/jumperzq86/jumper_conn/interf"
	"github.com/jumperzq86/jumper_conn/util"
	"github.com/jumperzq86/jumper_transform"
	jtd "github.com/jumperzq86/jumper_transform/def"
	jti "github.com/jumperzq86/jumper_transform/interf"
)

const addr = "localhost:8802"

func main() {

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws_connect"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("dial failed, err: %s\n", err)
		return
	}
	//note: transform 可以只定义一个，他本身是线程安全对
	ts := jumper_transform.Newtransform()
	ts.AddOp(jtd.PacketBinary, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(c *websocket.Conn) {
		defer func() {
			fmt.Printf("============")
			wg.Done()
		}()

		var h Handler

		wsOp := def.ConnOptions{
			MaxMsgSize:     def.MaxMsgSize,
			ReadTimeout:    def.ReadTimeout,
			WriteTimeout:   def.WriteTimeout,
			AsyncWriteSize: def.AsyncWriteSize,
			Side:           def.ClientSide,

			PingPeriod:       def.PingPeriod,
			PongWait:         def.PongWait,
			CloseGracePeriod: def.CloseGracePeriod,
		}

		jconn, err := jumper_conn.NewwsConn(c, &wsOp, &h)
		if err != nil {
			fmt.Printf("new tcp conn failed. err: %s\n", err)
			return
		}
		h.Init(jconn, ts)

		fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

		//send hello
		str := fmt.Sprintf("this is tcp_client %s, hello", jconn.LocalAddr())
		msg := jti.Message{
			Type:    1,
			Content: []byte(str),
		}

		var output []byte
		err = h.Execute(jtd.Forward, &msg, &output)
		if err != nil {
			fmt.Printf("transform failed, err: %s\n", err)
			return
		}

		fmt.Printf("sendMsg: %v\n", output)

		err = h.Write(output)
		if err != nil {
			fmt.Printf("write failed, err: %s\n", err)
			return
		}

	}(c)

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
