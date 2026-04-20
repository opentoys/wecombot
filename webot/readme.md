# webot — 企业微信群机器人 Webhook Go SDK

基于 [消息推送 API](https://developer.work.weixin.qq.com/document/path/99110) 的 Go 实现。

## 快速开始

```go
import "github.com/opentoys/wecombot/webot"

bot := webot.New("YOUR_WEBHOOK_KEY")
resp, err := bot.SendText("hello world")
```

## 消息类型

### 文本 (text)

```go
// 纯文本
bot.SendText("广州今日天气：29度")

// @指定成员 + @所有人
bot.SendText("紧急通知！",
    webot.WithMentionList("zhangsan", "lisi"),
    webot.MentionAll(),
)

// @手机号
bot.SendText("请查收",
    webot.WithMentionMobile("13800001111"),
)
```

### Markdown / Markdown V2

```go
// 标准 Markdown（支持 @成员、字体颜色）
err := bot.SendMarkdown(`## 实时告警
> 新增用户反馈<font color="warning">132例</font>
> 普通用户反馈:<font color="comment">117例</font>`)

// Markdown V2（支持表格、代码块、引用等增强语法）
err := bot.SendMarkdownV2(`# 标题
| 姓名 | 尺寸 | 地址 |
| :--: | :--: | :---: |
| 张三  | S    | 广州  |`)
```

### 图片 (image)

```go
// base64 方式
b64, md5 := encodeBase64MD5(imageBytes) // 自行编码
bot.SendImage(b64, md5)

// 直接发送本地文件
bot.SendImageFile("/path/to/image.png")
```

### 图文 (news)

```go
// 单条图文
bot.SendSingleArticle(
    "中秋节礼品领取",
    "今年中秋节公司有豪礼相送",
    "https://example.com/detail",
    "https://example.com/pic.jpg",
)

// 多条图文 (1-8 条)
bot.SendNews([]webot.Article{
    {Title: "文章1", Description: "...", URL: "...", PicURL: "..."},
    {Title: "文章2", Description: "...", URL: "...", PicURL: "..."},
})
```

### 文件 / 语音 (需要先上传获取 media_id)

```go
// 上传文件 -> 获取 media_id -> 发送
uploadResp, _ := bot.UploadFile("/path/to/report.pdf")
bot.SendFile(uploadResp.MediaID)

// 一键上传并发送
bot.SendFileFromFile("/path/to/report.pdf")

// 语音 (AMR 格式, <=2MB, <=60s)
uploadResp, _ := bot.UploadVoice("/path/to/voice.amr")
bot.SendVoice(uploadResp.MediaID)
```

### 模板卡片 (template_card)

**文本通知模板卡片：**

```go
card := &webot.TextNoticeCard{
    CardType: webot.CardTypeTextNotice,
    Source: &webot.TemplateCardSource{
        IconURL:   "https://...",
        Desc:      "企业微信",
        DescColor: 0,
    },
    MainTitle: &webot.TemplateCardMainTitle{
        Title: "欢迎使用企业微信",
        Desc:  "您的好友正在邀请您加入",
    },
    EmphasisContent: &webot.TemplateCardEmphasisContent{
        Title: "100",
        Desc:  "数据含义",
    },
    HorizontalContentList: []*webot.TemplateCardHorizontalContent{
        {Keyname: "邀请人", Value: "张三"},
        {Keyname: "官网", Value: "点击访问", Type: 1, URL: "https://..."},
    },
    CardAction: &webot.TemplateCardAction{Type: 1, URL: "https://..."},
}
bot.SendTemplateCard(card)
```

**图文展示模板卡片：**

```go
card := &webot.NewsNoticeCard{
    CardType: webot.CardTypeNewsNotice,
    MainTitle: &webot.TemplateCardMainTitle{
        Title: "欢迎使用企业微信",
        Desc:  "您的好友正在邀请您加入企业微信",
    },
    CardImage: &webot.NewsNoticeCardImage{
        URL:         "https://...",
        AspectRatio: 1.78,
    },
    VerticalContentList: []*webot.NewsNoticeVerticalContent{
        {Title: "惊喜红包等你来拿", Desc: "下载还能抢红包！"},
    },
    CardAction: &webot.TemplateCardAction{Type: 1, URL: "https://..."},
}
bot.SendTemplateCard(card)
```

## 频率限制

每个 webhook 限 **20 条/分钟**。SDK 提供内置速率限制：

```go
// 启用自动限流（超限时阻塞等待）
bot := webot.New(key).WithRateLimit()
// 后续 Send* 调用会自动控制频率
```

## 自定义配置

```go
// 自定义 Webhook URL
bot := webot.NewWithURL("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx")

// 自定义 HTTP Client（如设置代理）
customClient := &http.Client{
    Timeout: 15 * time.Second,
    Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
}
bot := webot.NewWithClient(key, customClient)
```

## API 参考

### Bot 构造

| 方法 | 说明 |
|------|------|
| `New(key)` | 用 webhook key 创建 Bot |
| `NewWithURL(url)` | 用完整 URL 创建 Bot |
| `NewWithClient(key, hc)` | 带自定义 HTTP Client |
| `WithRateLimit()` | 启用 20条/min 自动限流 |

### 发送消息

| 方法 | 返回值 | 说明 |
|------|--------|------|
| `SendText(content, opts...)` | `(*Response, error)` | 发送文本（支持 @提及） |
| `SendMarkdown(content)` | `(*Response, error)` | 发送 Markdown |
| `SendMarkdownV2(content)` | `(*Response, error)` | 发送 Markdown V2 |
| `SendImage(base64, md5)` | `(*Response, error)` | 发送图片 (base64) |
| `SendImageFile(path)` | `(*Response, error)` | 发送本地图片文件 |
| `SendNews(articles[])` | `(*Response, error)` | 发送图文 (1-8篇) |
| `SendSingleArticle(...)` | `(*Response, error)` | 单条图文快捷方法 |
| `SendFile(mediaID)` | `(*Response, error)` | 发送文件 |
| `SendVoice(mediaID)` | `(*Response, error)` | 发送语音 |
| `SendFileFromFile(path)` | `(*Response, error)` | 上传+发送文件一步完成 |
| `SendTemplateCard(card)` | `(*Response, error)` | 发送模板卡片 |

### 媒体上传

| 方法 | 返回值 | 说明 |
|------|--------|------|
| `UploadMedia(reader, name, type)` | `(*UploadResp, error)` | 通用上传 |
| `UploadFile(path)` | `(*UploadResp, error)` | 上传文件 (<=20MB) |
| `UploadVoice(path)` | `(*UploadResp, error)` | 上传语音 (<=2MB, AMR) |

### TextOption 选项

| 选项 | 说明 |
|------|------|
| `WithMentionList(userIDs...)` | @userid 列表 |
| `WithMentionMobile(mobiles...)` | @手机号列表 |
| `MentionAll()` | @所有人 |

### 结构体

| 类型 | 说明 |
|------|------|
| `WebhookRequest` | 发送请求外层结构 |
| `WebhookResponse` | 发送响应 `{errcode, errmsg}` |
| `TextPayload` | 文本消息体 |
| `MarkdownPayload` | Markdown 消息体 |
| `ImagePayload` | 图片消息体 (base64+md5) |
| `Article` | 图文条目 |
| `TextNoticeCard` | 文本通知模板卡片 |
| `NewsNoticeCard` | 图文展示模板卡片 |
| `UploadMediaResponse` | 上传响应 `{media_id, type, created_at}` |

## 依赖

- `golang.org/x/net` (context)
- Go >= 1.23
