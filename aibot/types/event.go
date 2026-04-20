package types

// FromInfo identifies the message sender.
type FromInfo struct {
	UserID string `json:"userid"`
}

// ---- Handler Types ----

// OnMessageFunc handles incoming text/multimedia messages from users.
type OnMessageFunc func(reqID string, msg *MsgCallbackBody)

// OnEventFunc handles event callbacks (enter_chat, card_click, etc).
type OnEventFunc func(reqID string, event *EventCallbackBody)

// OnConnectedFunc is called after successful subscription.
type OnConnectedFunc func()

// OnDisconnectedFunc is called when connection is lost or kicked.
type OnDisconnectedFunc func(err error)

// OnReconnectingFunc is called before each reconnection attempt.
type OnReconnectingFunc func(attempt int)

// ---- Event Callback ----

// EventCallbackBody is the body of aibot_event_callback.
type EventCallbackBody struct {
	MsgID      string    `json:"msgid"`
	CreateTime int64     `json:"create_time"`
	AIBotID    string    `json:"aibotid"`
	ChatID     string    `json:"chatid,omitempty"`
	ChatType   string    `json:"chattype,omitempty"`
	From       FromInfo  `json:"from"`
	MsgType    string    `json:"msgtype"`
	Event      EventInfo `json:"event"`
}

// EventInfo contains event-specific data.
type EventInfo struct {
	EventType string `json:"eventtype"`

	// For template_card_event
	TaskID   string            `json:"task_id,omitempty"`
	Response map[string]string `json:"response,omitempty"` // e.g. {"key":"confirm"}

	// For feedback_event
	FeedbackID      string `json:"feedback_id,omitempty"`
	FeedbackUser    string `json:"feedback_user,omitempty"`
	FeedbackContent string `json:"feedback_content,omitempty"`
}

// ---- Protocol Envelope ----

// Header is the common header structure for requests and responses.
type Header struct {
	ReqID string `json:"req_id"`
}

// Request is the generic envelope for outgoing commands.
type Request struct {
	Cmd    string      `json:"cmd"`
	Header Header      `json:"headers"`
	Body   interface{} `json:"body,omitempty"`
}

// Response is the generic envelope for incoming responses.
type Response struct {
	Header  Header `json:"headers"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// IsOK returns true if errcode indicates success.
func (r *Response) IsOK() bool { return r.ErrCode == 0 }

// ---- Subscribe ----

// SubscribeBody is the body for aibot_subscribe.
type SubscribeBody struct {
	BotID  string `json:"bot_id"`
	Secret string `json:"secret"`
}
