package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
)

func CrawlContest(id, start, end int) {
	contestants := make([]*Contestant, 0, (end-start)*100)
	for i := start; i < end; i++ {
		url := fmt.Sprintf("https://codeforces.com/contest/%d/ratings/page/%d", id, i)
		contestants = append(contestants, CrawlPage(url)...)
	}

	bytes, err := json.Marshal(contestants)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fmt.Sprintf("%d.json", id))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(string(bytes))
}

func CrawlPage(url string) []*Contestant {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("url(%s): %v\n", url, doc.Find("title").Text())

	contestants := make([]*Contestant, 0, 100)
	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		rank := cast.ToInt(strings.TrimSpace(row.Find("td").Eq(0).Text()))
		before := cast.ToInt(strings.TrimSpace(row.Find("td").Eq(4).Find("span").Eq(0).Text()))
		after := cast.ToInt(strings.TrimSpace(row.Find("td").Eq(4).Find("span").Eq(1).Text()))

		if rank != 0 {
			contestants = append(contestants, &Contestant{
				Rank:       rank,
				Rating:     before,
				NeedRating: after,
			})
		}
	})
	return contestants
}
