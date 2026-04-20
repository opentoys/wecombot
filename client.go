package wecombot

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ErrNotConnected is returned when trying to send on a disconnected connection.
	ErrNotConnected = errors.New("wecombot: not connected")

	// ErrClosed is returned when the client has been explicitly closed.
	ErrClosed = errors.New("wecombot: client closed")
)

// Client is the WeCom AI Bot WebSocket long-connection client.
type Client struct {
	cfg *Config

	conn     *websocket.Conn
	connMu   sync.RWMutex

	handlers struct {
		OnMessage      OnMessageFunc
		OnEvent        OnEventFunc
		OnConnected    OnConnectedFunc
		OnDisconnected OnDisconnectedFunc
		OnReconnecting OnReconnectingFunc
	}

	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once

	wg sync.WaitGroup
}

// New creates a new WeCom Bot client with the given configuration.
func New(cfg *Config) (*Client, error) {
	if cfg.WSSURL == "" {
		cfg.WSSURL = DefaultWSSURL
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = 30 * time.Second
	}
	if cfg.ReconnectWait == 0 {
		cfg.ReconnectWait = 3 * time.Second
	}
	return &Client{cfg: cfg}, nil
}

// ---- Handler Registration ----

// OnMessage registers a handler for user message callbacks.
func (c *Client) OnMessage(fn OnMessageFunc) { c.handlers.OnMessage = fn }

// OnEvent registers a handler for event callbacks.
func (c *Client) OnEvent(fn OnEventFunc) { c.handlers.OnEvent = fn }

// OnConnected registers a callback invoked after successful subscription.
func (c *Client) OnConnected(fn OnConnectedFunc) { c.handlers.OnConnected = fn }

// OnDisconnected registers a callback when connection is lost or kicked.
func (c *Client) OnDisconnected(fn OnDisconnectedFunc) { c.handlers.OnDisconnected = fn }

// OnReconnecting registers a callback before each reconnect attempt.
func (c *Client) OnReconnecting(fn OnReconnectingFunc) { c.handlers.OnReconnecting = fn }

// ---- Connection Lifecycle ----

// Connect establishes the WebSocket connection, subscribes, starts heartbeat,
// and begins reading messages. It blocks until context is cancelled or Close is called.
func (c *Client) Connect(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	attempt := 0
	for {
		err := c.connectOnce()
		if err == ErrClosed {
			return err
		}

		if err != nil {
			select {
			case <-c.ctx.Done():
				return c.ctx.Err()
			default:
			}

			attempt++
			if c.cfg.ReconnectMaxAttempts > 0 && attempt > c.cfg.ReconnectMaxAttempts {
				return fmt.Errorf("wecombot: max reconnect attempts (%d) exceeded", c.cfg.ReconnectMaxAttempts)
			}

			wait := c.cfg.ReconnectWait * time.Duration(attempt)
			if c.handlers.OnReconnecting != nil {
				c.handlers.OnReconnecting(attempt)
			} else {
				log.Printf("[wecombot] reconnecting in %v (attempt %d)...", wait, attempt)
			}

			select {
			case <-time.After(wait):
			case <-c.ctx.Done():
				return c.ctx.Err()
			}
			continue
		}

		// Successful connect — reset attempt counter
		attempt = 0
	}
}

// connectOnce performs a single connection + subscribe + read loop cycle.
// Returns nil only if the server intentionally kicked us (disconnected_event)
// or if we were closed. Returns an error for transient failures that should trigger reconnect.
func (c *Client) connectOnce() error {
	// Dial WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(c.cfg.WSSURL, nil)
	if err != nil {
		return fmt.Errorf("wecombot: dial error: %w", err)
	}

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	defer func() {
		c.closeConn()
	}()

	// Send subscribe request
	subReqID := genReqID()
	if err := c.sendRequest(CmdSubscribe, subReqID, &SubscribeBody{
		BotID:  c.cfg.BotID,
		Secret: c.cfg.Secret,
	}); err != nil {
		return fmt.Errorf("wecombot: subscribe failed: %w", err)
	}

	// Read subscribe response
	var subResp Response
	if err := c.readJSON(&subResp); err != nil {
		return fmt.Errorf("wecombot: subscribe response read failed: %w", err)
	}
	if !subResp.IsOK() {
		return fmt.Errorf("wecombot: subscribe rejected: code=%d msg=%s", subResp.ErrCode, subResp.ErrMsg)
	}

	if c.cfg.Debug {
		log.Printf("[wecombot] subscribed successfully")
	}

	// Notify connected
	if c.handlers.OnConnected != nil {
		go c.handlers.OnConnected()
	}

	// Start heartbeat
	stopHB := make(chan struct{})
	go c.heartbeatLoop(conn, stopHB)

	// Read loop
	for {
		var raw json.RawMessage
		if err := c.readJSON(&raw); err != nil {
			close(stopHB)
			return err // trigger reconnect
		}

		cmd, err := extractCmd(raw)
		if err != nil {
			if c.cfg.Debug {
				log.Printf("[wecombot] unknown message format: %s", string(raw))
			}
			continue
		}

		switch cmd {
		case CmdPing:
			// Server may send ping; respond with pong-like behavior
			if c.cfg.Debug {
				log.Printf("[wecombot] received ping from server")
			}

		case CmdMsgCallback:
			msg := &CallbackEnvelope{}
			if err := json.Unmarshal(raw, msg); err == nil {
				if c.handlers.OnMessage != nil {
					c.handlers.OnMessage(msg.Header.ReqID, &msg.Body)
				}
			}

		case CmdEventCallback:
			ev := &EventEnvelope{}
			if err := json.Unmarshal(raw, ev); err == nil {
				// Disconnected event means we were kicked — return to allow reconnect
				if ev.Body.Event.EventType == EventDisconnected {
					close(stopHB)
					if c.cfg.Debug {
						log.Printf("[wecombot] received disconnected_event (kicked by new connection)")
					}
					return errors.New("wecombot: disconnected by new connection")
				}
				if c.handlers.OnEvent != nil {
					c.handlers.OnEvent(ev.Header.ReqID, &ev.Body)
				}
			}

		default:
			if c.cfg.Debug {
				log.Printf("[wecombot] unhandled cmd: %s", cmd)
			}
		}
	}
}

// Close gracefully shuts down the client.
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		if c.cancel != nil {
			c.cancel()
		}
		c.closeConn()
	})
	return nil
}

func (c *Client) closeConn() {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// Connected returns true if the WebSocket connection is active.
func (c *Client) Connected() bool {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.conn != nil
}

// ---- Internal Helpers ----

func (c *Client) sendRequest(cmd, reqID string, body interface{}) error {
	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return ErrNotConnected
	}

	req := Request{
		Cmd:    cmd,
		Header: Header{ReqID: reqID},
		Body:   body,
	}
	return conn.WriteJSON(req)
}

func (c *Client) readJSON(v interface{}) error {
	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return ErrNotConnected
	}
	return conn.ReadJSON(v)
}

func (c *Client) heartbeatLoop(conn *websocket.Conn, stop <-chan struct{}) {
	ticker := time.NewTicker(c.cfg.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reqID := genReqID()
			msg := Request{Cmd: CmdPing, Header: Header{ReqID: reqID}}
			c.connMu.RLock()
			currentConn := c.conn
			c.connMu.RUnlock()

			if currentConn != conn {
				// Connection changed (reconnected), stop this loop
				return
			}
			if err := currentConn.WriteJSON(msg); err != nil {
				if c.cfg.Debug {
					log.Printf("[wecombot] heartbeat write error: %v", err)
				}
				return
			}
			if c.cfg.Debug {
				log.Printf("[wecombot] ping sent (%s)", reqID)
			}

		case <-stop:
			return

		case <-c.ctx.Done():
			return
		}
	}
}

// ---- Envelope types for deserialization ----

// CallbackEnvelope wraps a message callback with its header.
type CallbackEnvelope struct {
	Header Header          `json:"headers"`
	Body   MsgCallbackBody `json:"body"`
}

// EventEnvelope wraps an event callback with its header.
type EventEnvelope struct {
	Header Header           `json:"headers"`
	Body   EventCallbackBody `json:"body"`
}

// extractCmd extracts the "cmd" field from raw JSON without full unmarshal.
func extractCmd(raw json.RawMessage) (string, error) {
	tmp := struct {
		Cmd string `json:"cmd"`
	}{}
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return "", err
	}
	return tmp.Cmd, nil
}

// genReqID generates a unique request ID.
func genReqID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return hex.EncodeToString(b)
}
