package webot

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultBaseURL is the WeCom webhook base URL.
	DefaultBaseURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

	// DefaultUploadURL is the media upload URL template.
	DefaultUploadURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media"

	// MaxMessagesPerMin is the API rate limit: 20 messages/minute per webhook.
	MaxMessagesPerMin = 20
)

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

// Bot is a WeCom Group Robot client that sends messages via webhook.
type Bot struct {
	webhookURL string // full send URL with key embedded
	uploadURL  string // full upload URL with key embedded
	httpClient *http.Client

	// rate limiter (optional, enabled via WithRateLimit)
	rateMu *sync.Mutex
	rateCh chan struct{}
	ctx    context.Context
}

// New creates a Bot with the given webhook key.
func New(key string) *Bot {
	return NewWithClient(key, defaultHTTPClient)
}

// NewWithURL creates a Bot with a full custom webhook URL.
func NewWithURL(webhookURL string) *Bot {
	return NewWithClientURL(webhookURL, defaultHTTPClient)
}

// NewWithClient creates a Bot with a custom HTTP client.
func NewWithClient(key string, hc *http.Client) *Bot {
	return &Bot{
		webhookURL: fmt.Sprintf("%s?key=%s", DefaultBaseURL, key),
		uploadURL:  fmt.Sprintf("%s?key=%s", DefaultUploadURL, key),
		httpClient: hc,
		rateMu:     &sync.Mutex{},
		ctx:        context.Background(),
	}
}

// NewWithClientURL creates a Bot with a full custom webhook URL and HTTP client.
func NewWithClientURL(webhookURL string, hc *http.Client) *Bot {
	return &Bot{
		webhookURL: webhookURL,
		uploadURL:  toUploadURL(webhookURL),
		httpClient: hc,
		rateMu:     &sync.Mutex{},
		ctx:        context.Background(),
	}
}

// WithRateLimit enables built-in token-bucket rate limiting (20 msgs/min).
// When the limit is reached, Send* calls block until a token is available.
func (b *Bot) WithRateLimit() *Bot {
	b.rateCh = make(chan struct{}, MaxMessagesPerMin)
	go b.refillLoop()
	return b
}

func (b *Bot) WithContext(ctx context.Context) *Bot {
	return &Bot{
		webhookURL: b.webhookURL,
		uploadURL:  b.uploadURL,
		httpClient: b.httpClient,
		rateMu:     b.rateMu,
		rateCh:     b.rateCh,
		ctx:        ctx,
	}
}

// ---- Text ----

// SendText sends a text message.
func (b *Bot) SendText(content string, opts ...TextOption) (*WebhookResponse, error) {
	payload := &TextPayload{Content: content}
	for _, opt := range opts {
		opt(payload)
	}
	req := WebhookRequest{
		MsgType: MsgTypeText,
		Text:    payload,
	}
	return b.send(&req)
}

// ---- Markdown ----

// SendMarkdown sends a markdown message.
func (b *Bot) SendMarkdown(content string) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType:  MsgTypeMarkdown,
		Markdown: &MarkdownPayload{Content: content},
	}
	return b.send(&req)
}

// SendMarkdownV2 sends a markdown_v2 message.
func (b *Bot) SendMarkdownV2(content string) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType:    MsgTypeMarkdownV2,
		MarkdownV2: &MarkdownPayload{Content: content},
	}
	return b.send(&req)
}

// ---- Image ----

// SendImage sends an image message using base64-encoded data.
// The md5 parameter should be the MD5 hash of the raw image bytes (before base64 encoding).
func (b *Bot) SendImage(base64Data, md5 string) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType: MsgTypeImage,
		Image: &ImagePayload{
			Base64: base64Data,
			MD5:    md5,
		},
	}
	return b.send(&req)
}

// SendImageFile reads a file, base64-encodes it and sends as an image.
func (b *Bot) SendImageFile(path string) (*WebhookResponse, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("webot: read file %s: %w", path, err)
	}
	b64, hash := encodeBase64MD5(data)
	return b.SendImage(b64, hash)
}

// ---- News (图文) ----

// SendNews sends a news (article) message.
func (b *Bot) SendNews(articles []Article) (*WebhookResponse, error) {
	if len(articles) == 0 {
		return nil, fmt.Errorf("webot: news requires at least one article")
	}
	if len(articles) > 8 {
		return nil, fmt.Errorf("webot: news supports max 8 articles, got %d", len(articles))
	}
	req := WebhookRequest{
		MsgType: MsgTypeNews,
		News:    &NewsPayload{Articles: articles},
	}
	return b.send(&req)
}

// SendSingleArticle sends a single-article news message (convenience wrapper).
func (b *Bot) SendSingleArticle(title, description, url, picURL string) (*WebhookResponse, error) {
	return b.SendNews([]Article{
		{Title: title, Description: description, URL: url, PicURL: picURL},
	})
}

// ---- File / Voice (require media_id from UploadMedia) ----

// SendFile sends a file message using a previously uploaded media_id.
func (b *Bot) SendFile(mediaID string) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType: MsgTypeFile,
		File:    &FilePayload{MediaID: mediaID},
	}
	return b.send(&req)
}

// SendVoice sends a voice message using a previously uploaded media_id.
func (b *Bot) SendVoice(mediaID string) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType: MsgTypeVoice,
		Voice:   &VoicePayload{MediaID: mediaID},
	}
	return b.send(&req)
}

// SendFileFromFile uploads a local file and then sends it as a file message.
func (b *Bot) SendFileFromFile(path string) (*WebhookResponse, error) {
	resp, err := b.UploadFile(path)
	if err != nil {
		return nil, err
	}
	return b.SendFile(resp.MediaID)
}

// ---- Template Card ----

// SendTemplateCard sends a template card message.
// card should be either *TextNoticeCard or *NewsNoticeCard.
func (b *Bot) SendTemplateCard(card interface{}) (*WebhookResponse, error) {
	req := WebhookRequest{
		MsgType:      MsgTypeTemplateCard,
		TemplateCard: card,
	}
	return b.send(&req)
}

// ---- Media Upload ----

const (
	MediaTypeFile  = "file"
	MediaTypeVoice = "voice"
)

// UploadMedia uploads a file or voice reader to obtain a media_id.
// mediaType must be MediaTypeFile or MediaTypeVoice.
func (b *Bot) UploadMedia(r io.Reader, filename, mediaType string) (*UploadMediaResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("media", filename)
	if err != nil {
		return nil, fmt.Errorf("webot: create form field: %w", err)
	}
	if _, err := io.Copy(fw, r); err != nil {
		return nil, fmt.Errorf("webot: copy file data: %w", err)
	}
	w.Close()

	url := fmt.Sprintf("%s&type=%s", b.uploadURL, mediaType)
	req, err := http.NewRequestWithContext(b.ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("webot: create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webot: upload request failed: %w", err)
	}
	defer resp.Body.Close()

	var result UploadMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("webot: decode response: %w", err)
	}
	if !result.IsOK() {
		return &result, fmt.Errorf("webot: upload error: code=%d msg=%s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}

// UploadFile uploads a regular file (max 20MB, >5 bytes). Returns media_id valid for 3 days.
func (b *Bot) UploadFile(path string) (*UploadMediaResponse, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("webot: open file %s: %w", path, err)
	}
	defer f.Close()
	return b.UploadMedia(f, filepath.Base(path), MediaTypeFile)
}

// UploadVoice uploads an AMR voice file (max 2MB, duration <= 60s). Returns media_id valid for 3 days.
func (b *Bot) UploadVoice(path string) (*UploadMediaResponse, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("webot: open voice %s: %w", path, err)
	}
	defer f.Close()
	return b.UploadMedia(f, filepath.Base(path), MediaTypeVoice)
}

// ---- Internal: HTTP send ----

func (b *Bot) send(req *WebhookRequest) (*WebhookResponse, error) {
	if b.rateCh != nil {
		<-b.rateCh // block until rate limiter allows
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("webot: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, b.webhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("webot: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("webot: send request failed: %w", err)
	}
	defer resp.Body.Close()

	var result WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("webot: decode response: %w", err)
	}
	return &result, nil
}

// ---- Internal: rate limiter ----

func (b *Bot) refillLoop() {
	ticker := time.NewTimer(time.Minute / MaxMessagesPerMin)
	defer ticker.Stop()
	for range ticker.C {
		b.rateMu.Lock()
		select {
		case b.rateCh <- struct{}{}:
		default:
		}
		b.rateMu.Unlock()
		ticker.Reset(time.Minute / MaxMessagesPerMin)
	}
}

// ---- Internal: helpers ----

// toUploadURL converts a send webhook URL to an upload URL by replacing the path segment.
func toUploadURL(sendURL string) string {
	// e.g., https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
	//   -> https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=xxx
	return strings.Replace(sendURL, "/webhook/send", "/webhook/upload_media", 1)
}

// encodeBase64MD5 returns (base64-encoded data, hex MD5 of original bytes).
func encodeBase64MD5(data []byte) (string, string) {
	h := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(data), fmt.Sprintf("%x", h)
}
