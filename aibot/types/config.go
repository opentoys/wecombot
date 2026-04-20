package types

import "time"

// Config holds configuration for the WeCom Bot client.
type Config struct {
	// BotID is the unique identifier of the AI bot.
	BotID string

	// Secret is the long-connection specific secret key.
	Secret string

	// WSSURL is the WebSocket server URL. Defaults to DefaultWSSURL.
	WSSURL string

	// HeartbeatInterval is the duration between heartbeats. Default 30s.
	HeartbeatInterval time.Duration

	// ReconnectMaxAttempts is max reconnection attempts. 0 means infinite.
	ReconnectMaxAttempts int

	// ReconnectWait is initial wait before first reconnect attempt.
	ReconnectWait time.Duration

	// Debug enables verbose logging of all messages.
	Debug bool
}
