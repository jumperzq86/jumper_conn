package jumper_conn

import (
	"net"

	"github.com/jumperzq86/jumper_conn/def"

	"github.com/gorilla/websocket"
	"github.com/jumperzq86/jumper_conn/impl/conn"
	"github.com/jumperzq86/jumper_conn/interf"
)

func NewwsConn(c *websocket.Conn, co *def.ConnOptions, handler interf.Handler) (interf.Conn, error) {
	wsConn, err := conn.CreatewsConn(c, co, handler)
	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func NewtcpConn(c net.Conn, co *def.ConnOptions, handler interf.Handler) (interf.Conn, error) {
	tcpConn, err := conn.CreatetcpConn(c, co, handler)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}
