package main

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

//RSSFeed represents just an interface which provides a channel for feed updates
type RSSFeed interface {
	GetUpdateChan() chan gofeed.Item
}

type rssFeed struct {
	url     string
	timeout time.Duration
}

//NewRSSFeed creates a new atom RSS feed instance implementing the RSSFeed interface
func NewRSSFeed(url string, timeout time.Duration) RSSFeed {
	return &rssFeed{
		url:     url,
		timeout: timeout,
	}
}

//GetUpdatesChan returns a channel for updates of an rss feed
func (r *rssFeed) GetUpdateChan() chan gofeed.Item {
	ch := make(chan gofeed.Item, 50)

	go r.updator(ch)

	return ch
}

func (r *rssFeed) updator(ch chan gofeed.Item) {
	feedParser := gofeed.NewParser()
	updatedTime := time.Now()

	for {
		feed, err := feedParser.ParseURL(r.url)

		if err != nil {
			log.Println(err.Error())
		} else if feed.UpdatedParsed.Sub(updatedTime) > 0 {
			for _, item := range feed.Items {
				if item.PublishedParsed.Sub(updatedTime) >= 0 {
					log.Printf(`New entry: "%s"`, item.Title)
					ch <- *item
				}
			}
			updatedTime = *feed.UpdatedParsed
		}

		time.Sleep(r.timeout)
	}
}
