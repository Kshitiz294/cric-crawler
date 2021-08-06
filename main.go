package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Batsman struct {
	Name string
	Runs string
	Balls string
	Fours string
	Sixes string
	SR string
}

type Bowler struct {
	Name string
	Overs string
	Maidens string
	Runs string
	Wickets string
	Economy string
}

type Summary struct {
	Score string
	RR string
}

func main() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan bool)

	// Start initial crawl
	crawl(quit)

	// After that every 10s
	for {
		select {
			case <- ticker.C:
				crawl(quit)
			case <- quit:
				ticker.Stop()
				return
		}
	}
}

func crawl(quit chan bool) {
	crawler := colly.NewCollector(
		colly.AllowedDomains("cricbuzz.com","www.cricbuzz.com"),
		colly.MaxDepth(1),
		colly.Async(true),
		colly.Debugger(&debug.LogDebugger{}),
	)

	summary := Summary{}
	batsmen := []Batsman{}
	bowlers := []Bowler{}

	crawler.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL:", r.URL)
	})

	crawler.OnHTML("div.cb-min-bat-rw", func(h *colly.HTMLElement) {
		summary.Score = h.ChildText(".cb-font-20.text-bold")
		summary.RR = h.ChildText(".cb-font-12.cb-text-gray")
	})

	crawler.OnHTML("div.cb-min-inf", func(table *colly.HTMLElement) {
		header := table.ChildText("div.cb-bg-gray")[:6]
		table.ForEach("div.cb-min-itm-rw", func(_ int, row *colly.HTMLElement) {
			entries := []string{}
			row.ForEach("div", func(_ int, e *colly.HTMLElement) {
				entries = append(entries, e.Text)
			})
			if header == "Batter" {
				batsman := Batsman{
					Name: entries[0],
					Runs: entries[1],
					Balls: entries[2],
					Fours: entries[3],
					Sixes: entries[4],
					SR: entries[5],
				}
				batsmen = append(batsmen, batsman)
			} else {
				bowler := Bowler{
					Name: entries[0],
					Overs: entries[1],
					Maidens: entries[2],
					Runs: entries[3],
					Wickets: entries[4],
					Economy: entries[5],
				}
				bowlers = append(bowlers, bowler)
			}
		})
	})

	crawler.OnError(func(_ *colly.Response, _ error) {
		quit <- true
	})

	crawler.Visit("https://www.cricbuzz.com/live-cricket-scores/32047/eng-vs-ind-1st-test")
	crawler.Wait()
	fmt.Println(batsmen)
	fmt.Println(bowlers)
	fmt.Println(summary)
}