package webot

// ---- Message Types ----

const (
	MsgTypeText          = "text"
	MsgTypeMarkdown      = "markdown"
	MsgTypeMarkdownV2    = "markdown_v2"
	MsgTypeImage         = "image"
	MsgTypeNews          = "news"
	MsgTypeFile          = "file"
	MsgTypeVoice         = "voice"
	MsgTypeTemplateCard  = "template_card"
)

// ---- Request envelope ----

// WebhookRequest is the outer JSON payload sent to the webhook URL.
type WebhookRequest struct {
	MsgType       string           `json:"msgtype"`
	Text          *TextPayload     `json:"text,omitempty"`
	Markdown      *MarkdownPayload `json:"markdown,omitempty"`
	MarkdownV2    *MarkdownPayload `json:"markdown_v2,omitempty"`
	Image         *ImagePayload    `json:"image,omitempty"`
	News          *NewsPayload     `json:"news,omitempty"`
	File          *FilePayload     `json:"file,omitempty"`
	Voice         *VoicePayload    `json:"voice,omitempty"`
	TemplateCard  interface{}      `json:"template_card,omitempty"`
}

// ---- Text ----

// TextPayload represents a text message.
type TextPayload struct {
	Content            string   `json:"content"`
	MentionedList      []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

// TextOption is a functional option for building text messages.
type TextOption func(*TextPayload)

// WithMentionList sets userid list to @mention.
func WithMentionList(userIDs ...string) TextOption {
	return func(t *TextPayload) { t.MentionedList = userIDs }
}

// WithMentionMobile sets phone number list to @mention.
func WithMentionMobile(mobiles ...string) TextOption {
	return func(t *TextPayload) { t.MentionedMobileList = mobiles }
}

// MentionAll adds @all to both mention lists.
func MentionAll() TextOption {
	return func(t *TextPayload) {
		t.MentionedList = append(t.MentionedList, "@all")
		t.MentionedMobileList = append(t.MentionedMobileList, "@all")
	}
}

// ---- Markdown ----

// MarkdownPayload represents a markdown or markdown_v2 message.
type MarkdownPayload struct {
	Content string `json:"content"`
}

// ---- Image ----

// ImagePayload represents an image message (base64-encoded).
type ImagePayload struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

// ---- News (图文) ----

// NewsPayload represents a news (article) message.
type NewsPayload struct {
	Articles []Article `json:"articles"`
}

// Article is a single article in a news message.
type Article struct {
	Title       string `json:"title"`        // max 128 bytes
	Description string `json:"description"`   // max 512 bytes
	URL         string `json:"url"`           // click-through link
	PicURL      string `json:"picurl"`        // optional, JPG/PNG
}

// ---- File / Voice (via media_id) ----

// FilePayload represents a file message using a media_id.
type FilePayload struct {
	MediaID string `json:"media_id"`
}

// VoicePayload represents a voice message using a media_id.
type VoicePayload struct {
	MediaID string `json:"media_id"`
}

// ---- Template Card (模板卡片) ----

const (
	CardTypeTextNotice = "text_notice" // 文本通知模版卡片
	CardTypeNewsNotice = "news_notice" // 图文展示模版卡片
)

// TemplateCardSource is the source style info for template cards.
type TemplateCardSource struct {
	IconURL    string `json:"icon_url,omitempty"`
	Desc       string `json:"desc,omitempty"`
	DescColor  int    `json:"desc_color,omitempty"` // 0=gray,1=black,2=red,3=green
}

// TemplateCardMainTitle is the main title area of a template card.
type TemplateCardMainTitle struct {
	Title string `json:"title,omitempty"` // max 26 chars
	Desc  string `json:"desc,omitempty"`  // max 30 chars
}

// TemplateCardEmphasisContent highlights key data.
type TemplateCardEmphasisContent struct {
	Title string `json:"title,omitempty"` // max 10 chars
	Desc  string `json:"desc,omitempty"`  // max 15 chars
}

// TemplateCardQuoteArea is the quote/citation area.
type TemplateCardQuoteArea struct {
	Type      int    `json:"type,omitempty"`      // 0=none, 1=url, 2=miniprogram
	URL       string `json:"url,omitempty"`       // when type==1
	AppID     string `json:"appid,omitempty"`      // when type==2
	PagePath  string `json:"pagepath,omitempty"`   // when type==2
	Title     string `json:"title,omitempty"`
	QuoteText string `json:"quote_text,omitempty"`
}

// TemplateCardHorizontalContent is a key-value row in horizontal_content_list.
type TemplateCardHorizontalContent struct {
	Keyname string `json:"keyname"`                // required, max 5 chars
	Value   string `json:"value,omitempty"`
	Type    int    `json:"type,omitempty"`          // 1=url, 2=file, 3=user profile
	URL     string `json:"url,omitempty"`           // when type==1
	MediaID string `json:"media_id,omitempty"`      // when type==2
	UserID  string `json:"userid,omitempty"`        // when type==3
}

// TemplateCardJumpItem is a jump guide entry in jump_list.
type TemplateCardJumpItem struct {
	Type     int    `json:"type,omitempty"`     // 0=none, 1=url, 2=miniprogram
	Title    string `json:"title,omitempty"`    // required, max 13 chars
	URL      string `json:"url,omitempty"`      // when type==1
	AppID    string `json:"appid,omitempty"`     // when type==2
	PagePath string `json:"pagepath,omitempty"`  // when type==2
}

// TemplateCardAction is the overall card click action.
type TemplateCardAction struct {
	Type     int    `json:"type"`               // 1=url, 2=miniprogram
	URL      string `json:"url,omitempty"`      // when type==1
	AppID    string `json:"appid,omitempty"`     // when type==2
	PagePath string `json:"pagepath,omitempty"`  // when type==2
}

// TextNoticeCard represents a text_notice template card.
type TextNoticeCard struct {
	CardType             string                           `json:"card_type"`
	Source               *TemplateCardSource              `json:"source,omitempty"`
	MainTitle            *TemplateCardMainTitle           `json:"main_title,omitempty"`
	EmphasisContent      *TemplateCardEmphasisContent     `json:"emphasis_content,omitempty"`
	QuoteArea            *TemplateCardQuoteArea           `json:"quote_area,omitempty"`
	SubTitleText         string                           `json:"sub_title_text,omitempty"`
	HorizontalContentList []*TemplateCardHorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList             []*TemplateCardJumpItem           `json:"jump_list,omitempty"`
	CardAction           *TemplateCardAction              `json:"card_action"`
}

// NewsNoticeCardImage is the card image for news_notice cards.
type NewsNoticeCardImage struct {
	URL         string  `json:"url"`
	AspectRatio float64 `json:"aspect_ratio,omitempty"` // 1.3~2.25, default 1.3
}

// NewsNoticeImageTextArea is the left-image-right-text area.
type NewsNoticeImageTextArea struct {
	Type     int    `json:"type,omitempty"`     // 0=none, 1=url, 2=miniprogram
	URL      string `json:"url,omitempty"`
	AppID    string `json:"appid,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
	Title    string `json:"title,omitempty"`
	Desc     string `json:"desc,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// NewsNoticeVerticalContent is a vertical content item.
type NewsNoticeVerticalContent struct {
	Title string `json:"title,omitempty"` // max 26 chars
	Desc  string `json:"desc,omitempty"`  // max 112 chars
}

// NewsNoticeCard represents a news_notice template card.
type NewsNoticeCard struct {
	CardType             string                            `json:"card_type"`
	Source               *TemplateCardSource               `json:"source,omitempty"`
	MainTitle            *TemplateCardMainTitle            `json:"main_title,omitempty"`
	CardImage            *NewsNoticeCardImage              `json:"card_image,omitempty"`
	ImageTextArea        *NewsNoticeImageTextArea          `json:"image_text_area,omitempty"`
	QuoteArea            *TemplateCardQuoteArea            `json:"quote_area,omitempty"`
	VerticalContentList  []*NewsNoticeVerticalContent      `json:"vertical_content_list,omitempty"`
	HorizontalContentList []*TemplateCardHorizontalContent  `json:"horizontal_content_list,omitempty"`
	JumpList             []*TemplateCardJumpItem            `json:"jump_list,omitempty"`
	CardAction           *TemplateCardAction               `json:"card_action"`
}

// ---- Response ----

// WebhookResponse is the JSON response from the webhook API.
type WebhookResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// IsOK returns true if the response indicates success.
func (r *WebhookResponse) IsOK() bool {
	return r.ErrCode == 0
}

// ---- Upload Media Response ----

// UploadMediaResponse is the JSON response from the media upload API.
type UploadMediaResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// IsOK returns true if the upload succeeded.
func (r *UploadMediaResponse) IsOK() bool {
	return r.ErrCode == 0
}
