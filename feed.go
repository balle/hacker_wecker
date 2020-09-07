// Handle RSS feeds
package hackerWecker

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Url   string
	Title string
	Items []string
}

func FetchFeeds(outputChan chan<- Feed) {
	// Fetch the contents of all feeds, parse and filter them
	inputChan := make(chan fetchResult)

	for url, _ := range config.Feeds {
		go fetchUrl(url, inputChan)
	}

	for i := 0; i < len(config.Feeds); i++ {
		input := <-inputChan
		feed, err := parseFeed(input.Url, input.Content)
		var result Feed

		if err == nil {
			result.Title = feed.Title
			result.Url = input.Url

			for _, item := range feed.Items {
				if filterFeed(input.Url, item.Title) {
					result.Items = append(result.Items, item.Title)
				}
			}
		}

		outputChan <- result
	}
}

func parseFeed(url string, content string) (*gofeed.Feed, error) {
	// Parse an RSS or Atom feed
	parser := gofeed.NewParser()

	feed, err := parser.ParseString(content)

	if err != nil {
		Speak(fmt.Sprintf("Cannot parse feed from %s: %v\n", url, err))
		feed = nil
	}

	return feed, err
}

func filterFeed(url, text string) bool {
	// Check if the feed should be read regarding to the include and exclude filters
	readFeed := true

	if config.Feeds[url]["exclude"] != nil {
		for i := 0; i < len(config.Feeds[url]["exclude"]); i++ {
			if strings.Contains(strings.ToLower(text), strings.ToLower(config.Feeds[url]["exclude"][i])) {
				readFeed = false
				break
			}
		}
	}

	if readFeed && config.Feeds[url]["include"] != nil {
		readFeed = false

		for i := 0; i < len(config.Feeds[url]["include"]); i++ {

			if strings.Contains(strings.ToLower(text), strings.ToLower(config.Feeds[url]["include"][i])) {
				readFeed = true
				break
			}
		}
	}

	return readFeed
}

func ReadFeed(feed Feed) {
	// Read all feed items
	fmt.Printf("////[ %s\n\n", feed.Title)
	Speak(feed.Title)

	for _, item := range feed.Items {
		fmt.Printf("Speak %s\n", item)
		Speak(item)
	}

	fmt.Println()
}
