// Package wecombot provides a Go SDK for WeCom (企业微信) AI Bot WebSocket long-connection API.
//
// The SDK implements the full protocol including:
//   - WebSocket connection management with auto-reconnect
//   - Message callback handling (text, image, voice, file, video, mixed)
//   - Event callbacks (enter_chat, template_card_event, feedback_event, disconnected)
//   - Message responses (text, markdown, stream, template_card)
//   - Active message push
//   - Heartbeat keep-alive
//   - Temporary media upload (chunked upload)
package aibot
