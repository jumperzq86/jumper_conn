package interf

import (
	"net"
)

type Conn interface {
	GetConn() net.Conn
	Close()
	IsClosed() bool

	Write(data []byte) error
	AsyncWrite(data []byte) error

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Set(string, interface{})
	Get(string) interface{}
	Del(string)

	Run()
}
