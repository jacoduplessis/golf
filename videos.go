package main

import (
	"encoding/json"
	"fmt"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"log"
	"strings"
)

type SSVideo struct {
	ID          string `json:"id"`
	Tournament  string `json:"tournament"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Images      []struct {
		URL string `json:"url"`
	}
}

type Video struct {
	Title        string
	Description  string
	URL          string
	SRC          string
	ThumbnailSRC string
}

var videoListingSelector = cascadia.MustCompile("videos-listing")
var videoSourceSelector = cascadia.MustCompile("videos-single")

func cleanVideoSource(s string) string {
	s = html.UnescapeString(s)
	s = strings.Replace(s, `"`, "", -1)
	s = strings.Replace(s, `\/`, `/`, -1)
	s = strings.Replace(s, ".m3u8", "", 1)
	return s
}

func getVideoListing() []*Video {

	res, err := client.Get("https://supersport.com/golf/video")

	if err != nil {
		return nil
	}

	defer res.Body.Close()

	root, err := html.Parse(res.Body)

	if err != nil {
		return nil
	}

	var listingData string

	listingNode := videoListingSelector.MatchFirst(root)

	if listingNode == nil {
		return nil
	}

	for _, attr := range listingNode.Attr {
		if attr.Key == "listing" {
			listingData = attr.Val
			break
		}
	}

	listingData = html.UnescapeString(listingData)

	var d []*SSVideo

	err = json.NewDecoder(strings.NewReader(listingData)).Decode(&d)
	if err != nil {
		fmt.Printf("Error extracting JSON from HTML: %v", err)
		return nil
	}

	var videos []*Video

	for _, vidData := range d {
		vid := &Video{}

		vid.URL = "https://supersport.com" + vidData.URL
		vid.Title = vidData.Name
		vid.Description = vidData.Description
		vid.ThumbnailSRC = vidData.Images[0].URL

		videos = append(videos, vid)
	}

	// reverse
	// for i, j := 0, len(videos)-1; i < j; i, j = i+1, j-1 {
	// 	videos[i], videos[j] = videos[j], videos[i]
	// }

	// print(len(videos))

	ch := make(chan string, len(videos))

	for _, v := range videos {

		if v.URL == "" {
			continue
		}

		go func(video *Video) {

			res, err := client.Get(video.URL)
			if err != nil {
				log.Printf("Error retrieving video page: %v", err)
				ch <- ""
				return
			}

			defer res.Body.Close()

			root, err := html.Parse(res.Body)

			vidElement := videoSourceSelector.MatchFirst(root)

			if vidElement == nil {
				log.Println("vidElement is nil")
				ch <- ""
				return
			}

			for _, attr := range vidElement.Attr {
				if attr.Key == "video" {
					ch <- cleanVideoSource(attr.Val)
					return
				}
			}

		}(v)

	}

	for _, v := range videos {
		v.SRC = <-ch
	}

	return videos
}
