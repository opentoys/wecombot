package wecombot

import "time"

// WebSocket endpoint for WeCom AI Bot long-connection.
const DefaultWSSURL = "wss://openws.work.weixin.qq.com"

// Command types defined by the WeCom AI Bot protocol.
const (
	CmdSubscribe           = "aibot_subscribe"
	CmdMsgCallback         = "aibot_msg_callback"
	CmdEventCallback       = "aibot_event_callback"
	CmdRespondWelcome      = "aibot_respond_welcome_msg"
	CmdRespondMsg          = "aibot_respond_msg"
	CmdRespondUpdateMsg    = "aibot_respond_update_msg"
	CmdSendMsg             = "aibot_send_msg"
	CmdPing                = "ping"
	CmdUploadMediaInit     = "aibot_upload_media_init"
	CmdUploadMediaChunk    = "aibot_upload_media_chunk"
	CmdUploadMediaFinish   = "aibot_upload_media_finish"
)

// Message types.
const (
	MsgTypeText    = "text"
	MsgTypeImage   = "image"
	MsgTypeVoice   = "voice"
	MsgTypeVideo   = "video"
	MsgTypeFile    = "file"
	MsgTypeMixed   = "mixed"
	MsgTypeEvent   = "event"
	MsgTypeStream  = "stream"
	MsgTypeMarkdown    = "markdown"
	MsgTypeTemplateCard = "template_card"
)

// Chat type constants for aibot_send_msg.
const (
	ChatTypeSingle = 1 // 单聊
	ChatTypeGroup  = 2 // 群聊
	ChatTypeAuto   = 0 // 兼容单聊/群聊
)

// Event types.
const (
	EventEnterChat        = "enter_chat"
	EventTemplateCard     = "template_card_event"
	EventFeedback         = "feedback_event"
	EventDisconnected     = "disconnected_event"
)

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

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(botID, secret string) *Config {
	return &Config{
		BotID:               botID,
		Secret:              secret,
		WSSURL:              DefaultWSSURL,
		HeartbeatInterval:   30 * time.Second,
		ReconnectMaxAttempts: 0, // infinite
		ReconnectWait:       3 * time.Second,
	}
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
	BotID   string `json:"bot_id"`
	Secret  string `json:"secret"`
}

// ---- Message Callback ----

// MsgCallbackBody is the body of aibot_msg_callback.
type MsgCallbackBody struct {
	MsgID   string        `json:"msgid"`
	AIBotID string        `json:"aibotid"`
	ChatID  string        `json:"chatid"`
	ChatType string       `json:"chattype"`
	From    FromInfo      `json:"from"`
	MsgType string        `json:"msgtype"`
	Text    *TextContent  `json:"text,omitempty"`
	Image   *ImageContent `json:"image,omitempty"`
	Voice   *VoiceContent `json:"voice,omitempty"`
	File    *FileContent  `json:"file,omitempty"`
	Video   *VideoContent `json:"video,omitempty"`
	Mixed   *MixedContent `json:"mixed,omitempty"`
}

// FromInfo identifies the message sender.
type FromInfo struct {
	UserID string `json:"userid"`
}

// TextContent represents a text message.
type TextContent struct {
	Content string `json:"content"`
}

// ImageContent represents an image message (with decryption key in long-connection mode).
type ImageContent struct {
	URL    string `json:"url"`
	AESKey string `json:"aeskey,omitempty"`
}

// VoiceContent represents a voice message.
type VoiceContent struct {
	URL    string `json:"url"`
	AESKey string `json:"aeskey,omitempty"`
}

// FileContent represents a file message.
type FileContent struct {
	URL    string `json:"url"`
	AESKey string `json:"aeskey,omitempty"`
}

// VideoContent represents a video message.
type VideoContent struct {
	URL    string `json:"url"`
	AESKey string `json:"aeskey,omitempty"`
}

// MixedContent represents mixed text + image content.
type MixedContent struct {
	Content string `json:"content"`
}

// ---- Event Callback ----

// EventCallbackBody is the body of aibot_event_callback.
type EventCallbackBody struct {
	MsgID     string    `json:"msgid"`
	CreateTime int64     `json:"create_time"`
	AIBotID   string    `json:"aibotid"`
	ChatID    string    `json:"chatid,omitempty"`
	ChatType  string    `json:"chattype,omitempty"`
	From      FromInfo  `json:"from"`
	MsgType   string    `json:"msgtype"`
	Event     EventInfo `json:"event"`
}

// EventInfo contains event-specific data.
type EventInfo struct {
	EventType string `json:"eventtype"`

	// For template_card_event
	TaskID   string            `json:"task_id,omitempty"`
	Response map[string]string `json:"response,omitempty"` // e.g. {"key":"confirm"}

	// For feedback_event
	FeedbackID string `json:"feedback_id,omitempty"`
	FeedbackUser string `json:"feedback_user,omitempty"`
	FeedbackContent string `json:"feedback_content,omitempty"`
}

// ---- Respond Messages ----

// WelcomeMsgBody is the body for aibot_respond_welcome_msg.
type WelcomeMsgBody struct {
	MsgType string        `json:"msgtype"`
	Text    *TextContent  `json:"text,omitempty"`
	Markdown *MarkdownContent `json:"markdown,omitempty"`
	Image   *ImageMedia   `json:"image,omitempty"`
	File    *FileMedia    `json:"file,omitempty"`
	Voice   *VoiceMedia   `json:"voice,omitempty"`
	Video   *VideoMedia   `json:"video,omitempty"`
}

// RespondMsgBody is the body for aibot_respond_msg.
type RespondMsgBody struct {
	MsgType       string          `json:"msgtype"`
	Stream        *StreamContent  `json:"stream,omitempty"`
	Text          *TextContent    `json:"text,omitempty"`
	Markdown      *MarkdownContent `json:"markdown,omitempty"`
	TemplateCard  *TemplateCard   `json:"template_card,omitempty"`
	File          *FileMedia      `json:"file,omitempty"`
	Image         *ImageMedia     `json:"image,omitempty"`
	Voice         *VoiceMedia     `json:"voice,omitempty"`
	Video         *VideoMedia     `json:"video,omitempty"`
}

// StreamContent is used for streaming responses.
type StreamContent struct {
	ID       string       `json:"id"`
	Finish   bool         `json:"finish"`
	Content  string       `json:"content"`
	Feedback *FeedbackRef `json:"feedback,omitempty"`
}

// FeedbackRef references a feedback callback ID.
type FeedbackRef struct {
	ID string `json:"id"`
}

// MarkdownContent for markdown messages.
type MarkdownContent struct {
	Content  string       `json:"content"`
	Feedback *FeedbackRef `json:"feedback,omitempty"`
}

// TemplateCard for template card messages.
type TemplateCard struct {
	CardType   string                 `json:"card_type"`
	MainTitle  *CardTitle             `json:"main_title,omitempty"`
	SubTitle   string                 `json:"sub_title,omitempty"`
	Source     *CardSource            `json:"source,omitempty"`
	CardAction *CardAction            `json:"card_action,omitempty"`
	ButtonList []CardButton           `json:"button_list,omitempty"`
	TaskID     string                 `json:"task_id,omitempty"`
	CardImage  string                 `json:"card_image,omitempty"`
	Feedback   *FeedbackRef           `json:"feedback,omitempty"`
	// Additional fields for different card types go here...
}

type CardTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type CardSource struct {
	Desc   string `json:"desc"`
	DescURL string `json:"desc_url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type CardAction struct {
	Type     int    `json:"type"`
	URL      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

type CardButton struct {
	Text    string `json:"text"`
	Style   int    `json:"style"`
	Key     string `json:"key"`
}

// UpdateCardBody is the body for aibot_respond_update_msg.
type UpdateCardBody struct {
	ResponseType string       `json:"response_type"`
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
}

// Media content for file/image/voice/video (uses media_id).
type FileMedia struct {
	MediaID string `json:"media_id"`
}

type ImageMedia struct {
	MediaID string `json:"media_id"`
}

type VoiceMedia struct {
	MediaID string `json:"media_id"`
}

type VideoMedia struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// ---- Active Send ----

// SendMsgBody is the body for aibot_send_msg.
type SendMsgBody struct {
	ChatID      string          `json:"chatid"`
	ChatType    uint32          `json:"chat_type,omitempty"`
	MsgType     string          `json:"msgtype"`
	Markdown    *MarkdownContent `json:"markdown,omitempty"`
	TemplateCard *TemplateCard  `json:"template_card,omitempty"`
	File        *FileMedia      `json:"file,omitempty"`
	Image       *ImageMedia     `json:"image,omitempty"`
	Voice       *VoiceMedia     `json:"voice,omitempty"`
	Video       *VideoMedia     `json:"video,omitempty"`
}

// ---- Upload Media ----

// UploadInitBody is the body for aibot_upload_media_init.
type UploadInitBody struct {
	Type        string `json:"type"`         // file/image/voice/video
	Filename    string `json:"filename"`
	TotalSize   int64  `json:"total_size"`
	TotalChunks int    `json:"total_chunks"`
	MD5         string `json:"md5,omitempty"`
}

// UploadInitResponse is the response to upload init.
type UploadInitResponse struct {
	Header  Header `json:"headers"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Body    struct {
		UploadID string `json:"upload_id"`
	} `json:"body,omitempty"`
}

// UploadChunkBody is the body for aibot_upload_media_chunk.
type UploadChunkBody struct {
	UploadID   string `json:"upload_id"`
	ChunkIndex int    `json:"chunk_index"`
	Base64Data string `json:"base64_data"`
}

// UploadFinishBody is the body for aibot_upload_media_finish.
type UploadFinishBody struct {
	UploadID string `json:"upload_id"`
}

// UploadFinishResponse is the response to upload finish.
type UploadFinishResponse struct {
	Header  Header `json:"headers"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Body    struct {
		Type      string `json:"type"`
		MediaID   string `json:"media_id"`
		CreatedAt int64  `json:"created_at"`
	} `json:"body,omitempty"`
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
