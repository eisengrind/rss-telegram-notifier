package main

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

//GetUpdatesChan returns a channel for updates of an rss feed
func GetUpdatesChan(url string) chan gofeed.Item {
	ch := make(chan gofeed.Item, 50)

	go func() {
		fp := gofeed.NewParser()
		lastTime := time.Now().UTC()
		for {
			feed, err := fp.ParseURL(url)
			if err != nil {
				log.Println(err.Error())
			} else {
				for _, item := range feed.Items {
					if item.PublishedParsed.After(lastTime) {
						log.Println("New entry: ", item.Title)
						ch <- *item
					}
				}
			}

			lastTime = time.Now().UTC()
			time.Sleep(time.Second * 30)
		}
	}()

	return ch
}
