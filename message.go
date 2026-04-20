package wecombot

// ---- Welcome Message (respond to enter_chat event) ----

// RespondWelcome sends a welcome text message.
// Must be called within 5 seconds of receiving the enter_chat event.
func (c *Client) RespondWelcome(reqID, content string) error {
	return c.sendRequest(CmdRespondWelcome, reqID, &WelcomeMsgBody{
		MsgType: MsgTypeText,
		Text:    &TextContent{Content: content},
	})
}

// RespondWelcomeMarkdown sends a welcome markdown message.
func (c *Client) RespondWelcomeMarkdown(reqID, content string) error {
	return c.sendRequest(CmdRespondWelcome, reqID, &WelcomeMsgBody{
		MsgType:  MsgTypeMarkdown,
		Markdown: &MarkdownContent{Content: content},
	})
}

// ---- Message Response (respond to aibot_msg_callback) ----

// RespondText replies with a text message to the given callback req_id.
func (c *Client) RespondText(reqID, content string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeText,
		Text:    &TextContent{Content: content},
	})
}

// RespondMarkdown replies with a markdown message.
func (c *Client) RespondMarkdown(reqID, content string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType:  MsgTypeMarkdown,
		Markdown: &MarkdownContent{Content: content},
	})
}

// StreamStart begins a streaming response. Returns the stream ID.
// Call RespondStreamUpdate or RespondStreamFinish with the same stream ID
// and the original reqID to update the message.
// Must complete within 10 minutes of first send.
func (c *Client) StreamStart(reqID string) (streamID string, err error) {
	streamID = genReqID()
	err = c.streamSend(reqID, streamID, false, "")
	return
}

// StreamStartWithContent begins a streaming response with initial content.
func (c *Client) StreamStartWithContent(reqID, content string) (streamID string, err error) {
	streamID = genReqID()
	err = c.streamSend(reqID, streamID, false, content)
	return
}

// StreamUpdate updates an existing streaming message with new content.
func (c *Client) StreamUpdate(reqID, streamID, content string) error {
	return c.streamSend(reqID, streamID, false, content)
}

// StreamFinish completes the streaming message with final content.
// After this call, the message can no longer be updated.
func (c *Client) StreamFinish(reqID, streamID, finalContent string) error {
	return c.streamSend(reqID, streamID, true, finalContent)
}

func (c *Client) streamSend(reqID, streamID string, finish bool, content string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeStream,
		Stream: &StreamContent{
			ID:      streamID,
			Finish:  finish,
			Content: content,
		},
	})
}

// RespondTemplateCard replies with a template card message.
func (c *Client) RespondTemplateCard(reqID string, card *TemplateCard) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType:      MsgTypeTemplateCard,
		TemplateCard: card,
	})
}

// RespondFile replies with a file message using media_id.
func (c *Client) RespondFile(reqID, mediaID string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeFile,
		File:    &FileMedia{MediaID: mediaID},
	})
}

// RespondImage replies with an image using media_id.
func (c *Client) RespondImage(reqID, mediaID string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeImage,
		Image:   &ImageMedia{MediaID: mediaID},
	})
}

// RespondVoice replies with a voice message using media_id.
func (c *Client) RespondVoice(reqID, mediaID string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeVoice,
		Voice:   &VoiceMedia{MediaID: mediaID},
	})
}

// RespondVideo replies with a video using media_id.
func (c *Client) RespondVideo(reqID, mediaID string) error {
	return c.sendRequest(CmdRespondMsg, reqID, &RespondMsgBody{
		MsgType: MsgTypeVideo,
		Video:   &VideoMedia{MediaID: mediaID},
	})
}

// ---- Update Template Card (respond to template_card_event) ----
// Must be called within 5 seconds of receiving the event.

// UpdateTemplateCard updates a template card after user clicks its button.
func (c *Client) UpdateTemplateCard(reqID string, card *TemplateCard) error {
	return c.sendRequest(CmdRespondUpdateMsg, reqID, &UpdateCardBody{
		ResponseType: "update_template_card",
		TemplateCard: card,
	})
}

// ---- Active Push (aibot_send_msg) ----
// Requires prior interaction in the chat session.

// SendText actively pushes a text message to a single chat or group.
func (c *Client) SendText(chatID string, chatType uint32, content string) error {
	reqID := genReqID()
	return c.sendRequest(CmdSendMsg, reqID, &SendMsgBody{
		ChatID:   chatID,
		ChatType: chatType,
		MsgType:  MsgTypeText,
		Markdown: &MarkdownContent{Content: content}, // use markdown content field
	})
}

// SendMarkdown actively pushes a markdown-formatted message.
func (c *Client) SendMarkdown(chatID string, chatType uint32, content string) error {
	reqID := genReqID()
	return c.sendRequest(CmdSendMsg, reqID, &SendMsgBody{
		ChatID:   chatID,
		ChatType: chatType,
		MsgType:  MsgTypeMarkdown,
		Markdown: &MarkdownContent{Content: content},
	})
}

// SendTemplateCard actively pushes a template card message.
func (c *Client) SendTemplateCard(chatID string, chatType uint32, card *TemplateCard) error {
	reqID := genReqID()
	return c.sendRequest(CmdSendMsg, reqID, &SendMsgBody{
		ChatID:       chatID,
		ChatType:     chatType,
		MsgType:      MsgTypeTemplateCard,
		TemplateCard: card,
	})
}
