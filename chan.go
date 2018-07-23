package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSFeed interface {
	GetUpdateChan() chan gofeed.Item
}

type rssFeed struct {
	url        string
	guids      []string
	backupFile string
}

func NewRSSFeed(url, backupFile string) RSSFeed {
	return &rssFeed{
		url:        url,
		backupFile: backupFile,
		guids:      []string{},
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

	backupFile, err := os.OpenFile(r.backupFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := json.NewDecoder(backupFile).Decode(&r.guids); err != nil {
		log.Println(err.Error())
	}
	backupFile.Close()

	for {
		feed, err := feedParser.ParseURL(r.url)
		if err != nil {
			log.Println(err.Error())
		} else {
			for _, item := range feed.Items {
				if !r.guidSent(item.GUID) {
					log.Printf(`New entry: "%s"`, item.Title)
					ch <- *item
					r.guids = append(r.guids, item.GUID)
					backupFile, err := os.OpenFile(r.backupFile, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0777)
					if err != nil {
						log.Println(err.Error())
					} else {
						if err := json.NewEncoder(backupFile).Encode(&r.guids); err != nil {
							log.Println(err.Error())
						}
					}
					backupFile.Close()
				}
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func (r *rssFeed) guidSent(guid string) bool {
	for _, g := range r.guids {
		if guid == g {
			return true
		}
	}

	return false
}
