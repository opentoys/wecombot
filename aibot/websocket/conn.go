// Package websocket provides an abstract WebSocket connection interface
// so that different implementations can be swapped in (e.g., gorilla, nhooyr.io, golang.org/x/net).
package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opentoys/wecombot/types"
)

// ---- gorilla implementation ----

// GorillaDialer wraps [websocket.DefaultDialer] to implement Dialer.
type GorillaDialer struct {
	HandshakeTimeout time.Duration
	ReadBufferSize   int
	WriteBufferSize  int
}

func (d *GorillaDialer) dialer() *websocket.Dialer {
	ws := websocket.DefaultDialer
	if d.HandshakeTimeout != 0 {
		ws.HandshakeTimeout = d.HandshakeTimeout
	}
	if d.ReadBufferSize != 0 {
		ws.ReadBufferSize = d.ReadBufferSize
	}
	if d.WriteBufferSize != 0 {
		ws.WriteBufferSize = d.WriteBufferSize
	}
	return ws
}

// Dial opens a WebSocket connection using the gorilla implementation.
func (d *GorillaDialer) Dial(urlStr string, requestHeader http.Header) (types.Conn, error) {
	conn, _, err := d.dialer().Dial(urlStr, requestHeader)
	if err != nil {
		return nil, err
	}
	return &GorillaConn{Conn: conn}, nil
}

// GorillaConn wraps [*websocket.Conn] to implement Conn.
type GorillaConn struct {
	*websocket.Conn
}

// DefaultDialer is the default Dialer backed by gorilla/websocket.
var DefaultDialer = &GorillaDialer{}
