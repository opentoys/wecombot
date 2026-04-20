package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opentoys/wecombot/webot"
)

func main() {
	key := os.Getenv("WECOM_WEBHOOK_KEY")
	if key == "" {
		log.Fatal("Set WECOM_WEBHOOK_KEY env var to your webhook key")
	}

	bot := webot.New(key)

	// 1. Send a simple text message
	resp, err := bot.SendText("Hello from webot SDK!")
	if err != nil {
		log.Fatalf("SendText: %v", err)
	}
	fmt.Printf("SendText: errcode=%d errmsg=%s\n", resp.ErrCode, resp.ErrMsg)

	// 2. Send text with @mention
	resp, err = bot.SendText("Attention please!", webot.WithMentionList("wangqing"), webot.MentionAll())
	if err != nil {
		log.Fatalf("SendText @mention: %v", err)
	}
	fmt.Printf("SendText @mention: errcode=%d\n", resp.ErrCode)

	// 3. Send markdown
	resp, err = bot.SendMarkdown(`## Deploy Notice
> Deployment completed successfully!

| Service | Status |
| ------- | ------ |
| API     | <font color="info">Running</font> |
| Worker  | <font color="info">Running</font> |`)
	if err != nil {
		log.Fatalf("SendMarkdown: %v", err)
	}
	fmt.Printf("SendMarkdown: errcode=%d\n", resp.ErrCode)

	// 4. Send news article
	resp, err = bot.SendSingleArticle(
		"Weekly Report",
		"This week's summary report is ready for review.",
		"https://example.com/report",
		"https://example.com/report-cover.png",
	)
	if err != nil {
		log.Fatalf("SendNews: %v", err)
	}
	fmt.Printf("SendNews: errcode=%d\n", resp.ErrCode)

	// 5. Send template card (text_notice)
	card := &webot.TextNoticeCard{
		CardType: webot.CardTypeTextNotice,
		MainTitle: &webot.TemplateCardMainTitle{
			Title: "Deploy Success",
			Desc:  "Production environment v2.3.0 deployed",
		},
		Source: &webot.TemplateCardSource{
			IconURL:   "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0",
			Desc:      "CI/CD",
			DescColor: 0,
		},
		EmphasisContent: &webot.TemplateCardEmphasisContent{
			Title: "3",
			Desc:  "services updated",
		},
		HorizontalContentList: []*webot.TemplateCardHorizontalContent{
			{Keyname: "Version", Value: "v2.3.0"},
			{Keyname: "Environment", Value: "production"},
			{Keyname: "Details", Value: "View logs", Type: 1, URL: "https://example.com/logs"},
		},
		SubTitleText: "All health checks passed.",
		CardAction: &webot.TemplateCardAction{
			Type: 1,
			URL:  "https://example.com/deploy-details",
		},
	}

	resp, err = bot.SendTemplateCard(card)
	if err != nil {
		log.Fatalf("SendTemplateCard: %v", err)
	}
	fmt.Printf("SendTemplateCard: errcode=%d\n", resp.ErrCode)

	// 6. Upload and send file
	uploadResp, err := bot.UploadFile("/path/to/file.pdf")
	if err != nil {
		log.Printf("UploadFile: %v (skipping file send)", err)
	} else {
		resp, err = bot.SendFile(uploadResp.MediaID)
		if err != nil {
			log.Fatalf("SendFile: %v", err)
		}
		fmt.Printf("SendFile: errcode=%d media_id=%s\n", resp.ErrCode, uploadResp.MediaID)
	}
}
