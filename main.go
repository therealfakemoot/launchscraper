package main

import (
	"fmt"
	"strings"
	// "log"
	// "time"

	"github.com/gocolly/colly/v2"
)

type Launch struct {
	Time        string
	Date        string
	Site        string
	Mission     string
	Vehicle     string
	Description string
}

type Extractor struct {
	Launches []Launch
}

func lex() {

}

func (ex *Extractor) Extract(e *colly.HTMLElement) {
	launches := make([]Launch, 0)

	e.ForEach("div.datename", func(n int, e *colly.HTMLElement) {
		launches = append(launches, Launch{})
		data := strings.Split(e.Text, "•")
		if len(data) != 2 {
			return
		}
		launches[len(launches)-1].Vehicle = strings.ReplaceAll(data[0], "\n\t", "")
		launches[len(launches)-1].Mission = strings.ReplaceAll(data[1], "\t\n", "")
	})

	i := 0
	e.ForEach("div.missiondata", func(n int, e *colly.HTMLElement) {
		data := strings.Split(e.Text, "\n\n\n")
		if len(data) != 2 {
			return
		}
		launches[i].Time = strings.ReplaceAll(data[0], "\n\t", "")
		launches[i].Site = strings.ReplaceAll(data[1], "\n\t", "")
		i++
	})

	i = 0
	e.ForEach("div.missdescrip", func(n int, e *colly.HTMLElement) {
		launches[i].Description = strings.ReplaceAll(e.Text, "\n\t", "")
		i++
	})

	ex.Launches = launches

}

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	// datename
	// missiondata
	// missdescrip
	ex := Extractor{}
	c.OnHTML("div.mh-main.clearfix", func(e *colly.HTMLElement) {
		ex.Extract(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://spaceflightnow.com/launch-schedule/")

	for _, l := range ex.Launches {
		if strings.Contains(l.Site, "Wallop") {
			fmt.Printf("%#+v\n", l)
		}
	}
}
