package types

import "net/http"

// Conn is the abstract WebSocket connection interface.
type Conn interface {
	// WriteJSON writes a JSON-encoded message.
	WriteJSON(v interface{}) error

	// ReadJSON reads a JSON-encoded message into v.
	ReadJSON(v interface{}) error

	// Close closes the underlying connection.
	Close() error
}

// Dialer creates new WebSocket connections.
type Dialer interface {
	Dial(uri string, header http.Header) (Conn, error)
}
