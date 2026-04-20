// Package webot provides the WeCom (企业微信) Group Robot Webhook HTTP SDK.
//
// It implements the message-push (消息推送) API described at:
// https://developer.work.weixin.qq.com/document/path/99110
//
// Usage:
//
//	bot := webot.New("YOUR_WEBHOOK_KEY")
//	err := bot.SendText("hello world")
//
// Supported message types:
//   - text (with @mention support)
//   - markdown / markdown_v2
//   - image (base64)
//   - news (articles)
//   - file / voice (via media_id upload)
//   - template_card (text_notice / news_notice)
package webot
