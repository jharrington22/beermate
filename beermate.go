package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: beermate slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("beermate ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			fmt.Println("Message received")
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) >= 2 {

				// Splice ID off message
				askedBeerName := strings.Join(parts[1:len(parts)], " ")

				fmt.Println(askedBeerName)

				result := getBeerAdvocateData(askedBeerName)

				beerAdvURL := "https://www.beeradvocate.com"

				if result[1] != "" {
					m.Text = "Here you go, " + result[0] + " " + beerAdvURL + result[1]
				} else {
					m.Text = result[0]
				}

				postMessage(ws, m)
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
}

func getBeerAdvocateData(beerQuery string) [2]string {
	fmt.Println("Running func")

	v := url.Values{}

	v.Set("q", beerQuery)

	var result [2]string

	searchURL := "http://beeradvocate.com/search/?" + v.Encode() + "&qt=beer"

	fmt.Println("Search URL: " + searchURL)

	doc, err := goquery.NewDocument(searchURL)

	if err != nil {
		fmt.Println("goquery load failed!")
		result[0] = "Cannot retrieve URL"
		result[1] = searchURL
	} else {
		fmt.Println("goquery load ok")
	}

	fmt.Println("Error from goquery?")

	// Check for error
	noResultsCheck := doc.Find("#ba-content").Find("li").Text()

	r, err := regexp.Compile("(?i)" + "No results. Try being more specific")

	if err != nil {
		fmt.Println("Regex compile error")
	}

	if r.MatchString(noResultsCheck) {
		fmt.Println("Error no results found")
		result[0] = "I'm sorry, BA has no results for that beer! (" + searchURL + ")"
	} else {
		fmt.Println("Results found")
	}

	sel := doc.Find("#ba-content").Find("ul").Find("li")

	for i := range sel.Nodes {
		fmt.Println("Running loop")

		beerName := sel.Eq(i).Find("a").Find("b").Text()
		brewery := sel.Eq(i).Find("a").Last().Text()
		beerLink, ok := sel.Eq(i).Find("a").Attr("href")

		if !ok {
			fmt.Println("Error with beerLink")
		}

		fmt.Println("Brewery: " + brewery)
		fmt.Println("Beer name: " + beerName)
		fmt.Println("Beer search URL: " + searchURL)

		if ok {

			r, err := regexp.Compile("(?i)" + beerQuery)

			if err != nil {
				fmt.Println("Regex compile error")
			}

			// Need to filter out the brewery, eg. Firestone union jack IPA would return No match as beerName will be Union Jack IPA
			if r.MatchString(beerName) == true {
				fmt.Printf("Match! " + result[0] + " " + result[1] + "\n")
				fmt.Printf("Replaced: " + strings.Replace(beerName, brewery, "", -1))
				result[0] = beerName
				result[1] = beerLink
				return result
			} else {
				result[0] = "Couldn't get href for " + beerName
				result[1] = ""
				fmt.Printf("No match: " + result[0] + " " + result[1] + "\n")
			}
		} else {
			fmt.Println("Error fetching beerlink: " + beerLink)
		}

		fmt.Println("Something bad happened")
	}

	return result
}
