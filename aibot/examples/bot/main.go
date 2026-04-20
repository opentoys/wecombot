package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentoys/wecombot"
	"github.com/opentoys/wecombot/types"
	"github.com/opentoys/wecombot/websocket"
)

func main() {
	cfg := wecombot.DefaultConfig(
		os.Getenv("WECOM_BOT_ID"),
		os.Getenv("WECOM_BOT_SECRET"),
	)
	cfg.Debug = true

	client, err := wecombot.New(cfg, websocket.DefaultDialer)
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	// Register handlers
	client.OnConnected(func() {
		log.Println("[bot] connected to WeCom AI Bot")
	})

	client.OnDisconnected(func(err error) {
		log.Printf("[bot] disconnected: %v", err)
	})

	client.OnReconnecting(func(attempt int) {
		log.Printf("[bot] reconnecting... attempt %d", attempt)
	})

	client.OnEvent(func(reqID string, event *types.EventCallbackBody) {
		switch event.Event.EventType {
		case types.EventEnterChat:
			log.Printf("[bot] user %s entered chat", event.From.UserID)
			err := client.RespondWelcome(reqID, "您好！我是智能助手，有什么可以帮您的？")
			if err != nil {
				log.Printf("[bot] respond welcome error: %v", err)
			}

		case types.EventTemplateCard:
			log.Printf("[bot] template card event: task=%s response=%+v",
				event.Event.TaskID, event.Event.Response)

		case types.EventFeedback:
			log.Printf("[bot] feedback from %s: content=%s",
				event.Event.FeedbackUser, event.Event.FeedbackContent)
		}
	})

	client.OnMessage(func(reqID string, msg *types.MsgCallbackBody) {
		log.Printf("[bot] message from %s in %s: type=%s",
			msg.From.UserID, msg.ChatType, msg.MsgType)

		switch msg.MsgType {
		case types.MsgTypeText:
			content := msg.Text.Content
			log.Printf("[bot] text: %s", content)

			// Example: streaming response
			streamID, err := client.StreamStartWithContent(reqID, "正在思考...")
			if err != nil {
				log.Printf("[bot] stream error: %v", err)
				return
			}

			// Simulate processing...
			reply := fmt.Sprintf("收到您的消息: %s\n\n已为您记录。", content)
			if err := client.StreamFinish(reqID, streamID, reply); err != nil {
				log.Printf("[bot] stream finish error: %v", err)
			}

		default:
			log.Printf("[bot] unsupported message type: %s", msg.MsgType)
		}
	})

	// Start connection (blocks until context cancelled or Close called)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGTERM / SIGINT for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("[bot] shutting down...")
		cancel()
		os.Exit(0)
	}()

	log.Println("[bot] connecting...")
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("[bot] connect error: %v", err)
	}
}
