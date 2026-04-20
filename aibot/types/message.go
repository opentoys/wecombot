package types

// MarkdownContent for markdown messages.
type MarkdownContent struct {
	Content  string       `json:"content"`
	Feedback *FeedbackRef `json:"feedback,omitempty"`
}

// TemplateCard for template card messages.
type TemplateCard struct {
	CardType   string       `json:"card_type"`
	MainTitle  *CardTitle   `json:"main_title,omitempty"`
	SubTitle   string       `json:"sub_title,omitempty"`
	Source     *CardSource  `json:"source,omitempty"`
	CardAction *CardAction  `json:"card_action,omitempty"`
	ButtonList []CardButton `json:"button_list,omitempty"`
	TaskID     string       `json:"task_id,omitempty"`
	CardImage  string       `json:"card_image,omitempty"`
	Feedback   *FeedbackRef `json:"feedback,omitempty"`
	// Additional fields for different card types go here...
}

type CardTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type CardSource struct {
	Desc    string `json:"desc"`
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
	Text  string `json:"text"`
	Style int    `json:"style"`
	Key   string `json:"key"`
}

// UpdateCardBody is the body for aibot_respond_update_msg.
type UpdateCardBody struct {
	ResponseType string        `json:"response_type"`
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
	ChatID       string           `json:"chatid"`
	ChatType     uint32           `json:"chat_type,omitempty"`
	MsgType      string           `json:"msgtype"`
	Markdown     *MarkdownContent `json:"markdown,omitempty"`
	TemplateCard *TemplateCard    `json:"template_card,omitempty"`
	File         *FileMedia       `json:"file,omitempty"`
	Image        *ImageMedia      `json:"image,omitempty"`
	Voice        *VoiceMedia      `json:"voice,omitempty"`
	Video        *VideoMedia      `json:"video,omitempty"`
}

// ---- Upload Media ----

// UploadInitBody is the body for aibot_upload_media_init.
type UploadInitBody struct {
	Type        string `json:"type"` // file/image/voice/video
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

// ---- Respond Messages ----

// WelcomeMsgBody is the body for aibot_respond_welcome_msg.
type WelcomeMsgBody struct {
	MsgType  string           `json:"msgtype"`
	Text     *TextContent     `json:"text,omitempty"`
	Markdown *MarkdownContent `json:"markdown,omitempty"`
	Image    *ImageMedia      `json:"image,omitempty"`
	File     *FileMedia       `json:"file,omitempty"`
	Voice    *VoiceMedia      `json:"voice,omitempty"`
	Video    *VideoMedia      `json:"video,omitempty"`
}

// RespondMsgBody is the body for aibot_respond_msg.
type RespondMsgBody struct {
	MsgType      string           `json:"msgtype"`
	Stream       *StreamContent   `json:"stream,omitempty"`
	Text         *TextContent     `json:"text,omitempty"`
	Markdown     *MarkdownContent `json:"markdown,omitempty"`
	TemplateCard *TemplateCard    `json:"template_card,omitempty"`
	File         *FileMedia       `json:"file,omitempty"`
	Image        *ImageMedia      `json:"image,omitempty"`
	Voice        *VoiceMedia      `json:"voice,omitempty"`
	Video        *VideoMedia      `json:"video,omitempty"`
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

// ---- Message Callback ----

// MsgCallbackBody is the body of aibot_msg_callback.
type MsgCallbackBody struct {
	MsgID    string        `json:"msgid"`
	AIBotID  string        `json:"aibotid"`
	ChatID   string        `json:"chatid"`
	ChatType string        `json:"chattype"`
	From     FromInfo      `json:"from"`
	MsgType  string        `json:"msgtype"`
	Text     *TextContent  `json:"text,omitempty"`
	Image    *ImageContent `json:"image,omitempty"`
	Voice    *VoiceContent `json:"voice,omitempty"`
	File     *FileContent  `json:"file,omitempty"`
	Video    *VideoContent `json:"video,omitempty"`
	Mixed    *MixedContent `json:"mixed,omitempty"`
}
