package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alsosee/finder/structs"

	"github.com/PuerkitoBio/goquery"
)

type WikipediaResponse struct {
	Parse WikipediaResponseParse `json:"parse"`
}

type WikipediaResponseParse struct {
	Title string `json:"title"`
	Text  struct {
		Content string `json:"*"`
	} `json:"text"`
}

func NewWikipedia() Scraper {
	return &Wikipedia{}
}

type Wikipedia struct{}

func (wikipedia *Wikipedia) Scrape(page string) (map[string]structs.Content, error) {
	resp, err := wikipedia.GetPageResponse(page)
	if err != nil {
		return nil, fmt.Errorf("getting page response: %w", err)
	}

	text := resp.Parse.Text.Content

	infobox, err := ParseInfobox(text)
	if err != nil {
		return nil, fmt.Errorf("parsing infobox: %w", err)
	}

	content := structs.Content{
		Name:      page,
		Wikipedia: "https://en.wikipedia.org/wiki/" + page,
	}
	for k, v := range infobox {
		switch k {
		case "Directed by":
			content.Directors = []string{v}
		case "Produced by":
			content.Producers = []string{v}
		case "Written by":
			content.Writers = []string{v}
		case "Cinematography":
			content.Cinematography = []string{v}
		case "Edited by":
			content.Editors = []string{v}
		case "Music by":
			content.Music = []string{v}
		case "Production companies":
			content.Production = []string{v}
		case "Distributed by":
			content.Distributors = []string{v}
		case "Release dates":
			content.Released = v
		case "Running time":
			// parse time.Duration in format like "100 minutes"
			v = strings.Replace(v, " minutes", "m", 1)

			duration, err := time.ParseDuration(v)
			if err != nil {
				log.Printf("parsing duration: %v", err)
				continue
			}
			content.Length = duration
		}
	}

	// fmt.Println(text)
	result := map[string]structs.Content{}
	result[page] = content
	return result, nil
}

func (wikipedia *Wikipedia) GetPageResponse(page string) (*WikipediaResponse, error) {
	baseURL, err := url.Parse("https://en.wikipedia.org/w/api.php")
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	params := url.Values{}
	params.Add("action", "parse")
	params.Add("page", page)
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	response, err := http.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("getting page: %w", err)
	}
	defer response.Body.Close()

	var wikipediaResponse WikipediaResponse
	err = json.NewDecoder(response.Body).Decode(&wikipediaResponse)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &wikipediaResponse, nil
}

func ParseInfobox(text string) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	infobox := map[string]string{}

	doc.Find(".infobox tr").Each(func(i int, s *goquery.Selection) {
		key := s.Find(".infobox-label").Text()
		value := s.Find(".infobox-data").Text()
		infobox[key] = value
	})

	return infobox, nil
}
