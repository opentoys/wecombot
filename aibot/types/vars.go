package types

// WebSocket endpoint for WeCom AI Bot long-connection.
const DefaultWSSURL = "wss://openws.work.weixin.qq.com"

// Command types defined by the WeCom AI Bot protocol.
const (
	CmdSubscribe         = "aibot_subscribe"
	CmdMsgCallback       = "aibot_msg_callback"
	CmdEventCallback     = "aibot_event_callback"
	CmdRespondWelcome    = "aibot_respond_welcome_msg"
	CmdRespondMsg        = "aibot_respond_msg"
	CmdRespondUpdateMsg  = "aibot_respond_update_msg"
	CmdSendMsg           = "aibot_send_msg"
	CmdPing              = "ping"
	CmdUploadMediaInit   = "aibot_upload_media_init"
	CmdUploadMediaChunk  = "aibot_upload_media_chunk"
	CmdUploadMediaFinish = "aibot_upload_media_finish"
)

// Message types.
const (
	MsgTypeText         = "text"
	MsgTypeImage        = "image"
	MsgTypeVoice        = "voice"
	MsgTypeVideo        = "video"
	MsgTypeFile         = "file"
	MsgTypeMixed        = "mixed"
	MsgTypeEvent        = "event"
	MsgTypeStream       = "stream"
	MsgTypeMarkdown     = "markdown"
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
	EventEnterChat    = "enter_chat"
	EventTemplateCard = "template_card_event"
	EventFeedback     = "feedback_event"
	EventDisconnected = "disconnected_event"
)
