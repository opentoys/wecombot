# Wecom Bot (企微机器人长连接 SDK)

> Go SDK for 企业微信 智能机器人 WebSocket 长连接 API

基于 [官方文档](https://developer.work.weixin.qq.com/document/path/101463) 实现的完整 Go 版本。

## 特性

- **WebSocket 长连接管理** — 自动订阅、心跳保活、断线重连（指数退避）
- **消息接收** — text/image/voice/file/video/mixed 全类型支持
- **事件回调** — enter_chat、template_card_event、feedback_event、disconnected
- **消息回复** — 文本/Markdown/流式输出/模板卡片/多媒体文件
- **主动推送** — 无需用户触发，主动向会话发送消息（需先有交互历史）
- **临时素材上传** — 分片上传，支持 file/image/voice/video

## 快速开始

```go
package main

import (
    "context"
    "log"
    "os"

    wecom "github.com/opentoys/wecombot"
)

func main() {
    cfg := wecom.DefaultConfig(
        os.Getenv("WECOM_BOT_ID"),   // 企微机器人 BotID
        os.Getenv("WECOM_BOT_SECRET"), // 长连接专用 Secret
    )

    client, err := wecom.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 注册事件处理器
    client.OnMessage(func(reqID string, msg *wecom.MsgCallbackBody) {
        // 处理用户消息
        client.RespondText(reqID, "收到您的消息！")
    })

    client.OnEvent(func(reqID string, ev *wecom.EventCallbackBody) {
        switch ev.Body.Event.EventType {
        case wecom.EventEnterChat:
            client.RespondWelcome(reqID, "您好！有什么可以帮您？")
        }
    })

    // 启动连接（阻塞运行）
    if err := client.Connect(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## API 参考

### Client 创建与生命周期

| 方法 | 签名 | 说明 |
|------|------|------|
| `New` | `New(cfg *Config) (*Client, error)` | 创建客户端实例 |
| `Connect` | `(c *Client) Connect(ctx context.Context) error` | 建立连接、订阅、启动心跳和读循环 |
| `Close` | `(c *Client) Close() error` | 关闭连接 |
| `Connected` | `(c *Client) Connected() bool` | 是否处于连接状态 |

### Config 配置项

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `BotID` | string | - | **必填** 机器人唯一标识 |
| `Secret` | string | - | **必填** 长连接专用密钥 |
| `WSSURL` | string | `wss://openws.work.weixin.qq.com` | WebSocket 地址 |
| `HeartbeatInterval` | time.Duration | 30s | 心跳间隔 |
| `ReconnectMaxAttempts` | int | 0 (无限) | 最大重连次数 |
| `ReconnectWait` | time.Duration | 3s | 首次重连等待时间 |
| `Debug` | bool | false | 调试日志开关 |

### 事件回调注册

| 方法 | 触发场景 |
|------|----------|
| `OnMessage(fn)` | 用户发来文本/图片/语音/文件/视频等消息 |
| `OnEvent(fn)` | 进入会话/模板卡片点击/反馈/连接被踢 |
| `OnConnected(fn)` | 订阅成功，连接就绪 |
| `OnDisconnected(fn)` | 连接断开或被踢 |
| `OnReconnecting(fn)` | 每次重连前 |

### 回复消息（响应 aibot_msg_callback）

| 方法 | 签名 | 说明 |
|------|------|------|
| `RespondText` | `(reqID, content string) error` | 回复纯文本 |
| `RespondMarkdown` | `(reqID, content string) error` | 回复 Markdown 格式（支持标题/列表/表格/代码块等） |
| `StreamStart` | `(reqID string) (streamID, err)` | 开始流式回复，返回 stream ID |
| `StreamStartWithContent` | `(reqID, content string) (streamID, err)` | 带初始内容的流式回复 |
| `StreamUpdate` | `(reqID, streamID, content string) error` | 更新流式内容 |
| `StreamFinish` | `(reqID, streamID, finalContent string) error` | 完成流式回复（10分钟内必须完成） |
| `RespondTemplateCard` | `(reqID string, card *TemplateCard) error` | 回复模板卡片 |
| `RespondFile` | `(reqID, mediaID string) error` | 回复文件 |
| `RespondImage` | `(reqID, mediaID string) error` | 回复图片 |
| `RespondVoice` | `(reqID, mediaID string) error` | 回复语音 |
| `RespondVideo` | `(reqID, mediaID string) error` | 回复视频 |

**注意**: 所有回复方法必须在收到消息回调后的 **24 小时内** 调用。

### 回复欢迎语（响应 enter_chat 事件）

| 方法 | 签名 | 说明 |
|------|------|------|
| `RespondWelcome` | `(reqID, content string) error` | 发送文本欢迎语 |
| `RespondWelcomeMarkdown` | `(reqID, content string) error` | 发送 Markdown 欢迎语 |

**注意**: 必须在收到进入会话事件的 **5 秒内** 调用。

### 更新模板卡片（响应 template_card_event 事件）

| 方法 | 签名 | 说明 |
|------|------|------|
| `UpdateTemplateCard` | `(reqID string, card *TemplateCard) error` | 更新卡片内容 |

**注意**: 必须在收到点击事件的 **5 秒内** 调用。

### 主动推送消息（aibot_send_msg）

| 方法 | 签名 | 说明 |
|------|------|------|
| `SendText` | `(chatID string, chatType uint32, content string) error` | 推送文本消息 |
| `SendMarkdown` | `(chatID string, chatType uint32, content string) error` | 推送 Markdown 消息 |
| `SendTemplateCard` | `(chatID string, chatType uint32, card *TemplateCard) error` | 推送模板卡片 |

**前提条件**: 目标会话中必须有用户先给机器人发过消息。
**频率限制**: 单个会话 30 条/分钟，1000 条/小时。

ChatType 常量:
```go
wecom.ChatTypeSingle = 1  // 单聊（填 userid）
wecom.ChatTypeGroup  = 2  // 群聊（填 chatid）
wecom.ChatTypeAuto   = 0  // 自动识别
```

### 临时素材上传

| 方法 | 签名 | 说明 |
|------|------|------|
| `UploadFromFile` | `(mediaType, filePath string) (*UploadResult, error)` | 上传本地文件 |
| `UploadFromReader` | `(mediaType, filename string, size int64, r io.Reader) (*UploadResult, error)` | 从 Reader 上传 |

Media 类型:
```go
wecom.MediaTypeFile   // 普通文件, max 20MB
wecom.MediaTypeImage  // 图片 PNG/JPG/GIF, max 10MB
wecom.MediaTypeVoice  // 语音 AMR, max 2MB
wecom.MediaTypeVideo  // 视频 MP4, max 10MB
```

返回的 MediaID 有效期 **3 天**，可用于回复消息中的 media_id 参数。

## 流式回复完整示例

```go
client.OnMessage(func(reqID string, msg *wecom.MsgCallbackBody) {
    // 1. 开始流式回复（显示"思考中..."）
    streamID, err := client.StreamStartWithContent(reqID, "正在查询数据，请稍候...")
    if err != nil { return }

    // 2. 分步更新内容（可选多次）
    client.StreamUpdate(reqID, streamID, "已获取部分结果...")

    // 3. 完成流式回复（最终内容）
    reply := "查询完成：今日订单共 128 笔，总金额 ¥52,340.00"
    client.StreamFinish(reqID, streamID, reply)
})
```

## 依赖

| 包 | 用途 |
|----|------|
| `github.com/gorilla/websocket` | WebSocket 客户端 |

Go 版本要求: **>= 1.23**

## 协议说明

本 SDK 对应企业微信文档路径:
**消息接收与发送 → 智能机器人 → 智能机器人长连接**

核心命令对照表:

| 功能 | 命令 (cmd) |
|------|-------------|
| 订阅认证 | `aibot_subscribe` |
| 收到消息 | `aibot_msg_callback` |
| 收到事件 | `aibot_event_callback` |
| 回复欢迎语 | `aibot_respond_welcome_msg` |
| 回复消息 | `aibot_respond_msg` |
| 更新卡片 | `aibot_respond_update_msg` |
| 主动推送 | `aibot_send_msg` |
| 心跳 | `ping` |
| 上传初始化 | `aibot_upload_media_init` |
| 上传分片 | `aibot_upload_media_chunk` |
| 上传完成 | `aibot_upload_media_finish` |
