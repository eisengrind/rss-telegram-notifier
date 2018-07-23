package main

import (
	"bytes"
	"log"

	"github.com/playnet-public/flagenv"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	tgAPIKey   = flagenv.String("tg-api-key", "", "This is your telegram bot-api-token")
	chatID     = flagenv.Int("chat-id", 0, "The chat to use.")
	rssFeedURL = flagenv.String("rss-feed-url", "", "The RSS-feed to use for updates")
)

type message struct {
	Title       string
	Description string
	Date        string
	Link        string
	Author      string
}

func main() {
	flagenv.Parse()

	if *chatID == 0 {
		log.Fatal("you have to enter an chat id")
	}

	if *rssFeedURL == "" {
		log.Fatal("you have to enter a RSS feed url")
	}

	bot, err := tgbotapi.NewBotAPI(*tgAPIKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	updator := NewRSSFeed(*rssFeedURL, "guids.json")
	updates := updator.GetUpdateChan()

	for update := range updates {
		var b bytes.Buffer

		if err := parsedTemplateMessage.Execute(&b, &message{
			Title:       update.Title,
			Description: update.Description,
			Date:        update.PublishedParsed.Format("02.01.2006 @ 15:04"),
			Link:        update.Link,
			Author:      update.Author.Name,
		}); err != nil {
			log.Fatal(err.Error())
		}

		msg := tgbotapi.NewMessage(
			int64(*chatID),
			b.String(),
		)
		msg.ParseMode = "html"

		if _, err = bot.Send(msg); err != nil {
			log.Fatal(err.Error())
		}
	}
}
