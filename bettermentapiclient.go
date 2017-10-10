package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
)

// BettermentAPIClient ...
type BettermentAPIClient struct {
	Email     string
	Password  string
	csrfToken string
}

// Summary - ize the account earning
func (bmc *BettermentAPIClient) Summary() {
	dom, login := bmc.login()
	if !login {
		panic("Failed to log in!")
	}

	// TODO: need to think about how I could do this more reliably
	bmBalance := dom.Find("h5:contains('Total Betterment Balance') + p").First().Text()
	totalNetWorth := dom.Find("h5:contains('Total Net Worth') + p").First().Text()
	totalEarnings := dom.Find("h5:contains('Total Earnings') + p").First().Text()
	taxLossesHarvested := dom.Find("h5:contains('Tax Losses Harvested') + p").First().Text()

	fmt.Printf("Betterment balance: %s\n", bmBalance)
	fmt.Printf("Total Earnings: %s\n", totalEarnings)
	fmt.Printf("Tax Losses Harvested: %s\n", taxLossesHarvested)
	fmt.Printf("Total Net Worth: %s\n", totalNetWorth)
}

// TODO: This should not return bow()
// This will need to change to support getting info from pages other than /app/summary
func (bmc *BettermentAPIClient) login() (*goquery.Selection, bool) {
	// Request login page
	bow := surf.NewBrowser()
	err := bow.Open("https://wwws.betterment.com/app/login")
	if err != nil {
		panic(err)
	}

	// Find CSRF token we will need
	token, found := bow.Dom().Find("meta[name=csrf-token]").First().Attr("content")
	if !found {
		panic("Failed to find csrf-token!")
	}

	// Add it as header to next request
	bmc.csrfToken = token
	bow.AddRequestHeader("X-CSRF-Token", bmc.csrfToken)

	// Fill out login form
	fm, err := bow.Form("form.new_web_authentication")
	if err != nil {
		panic(err)
	}

	// Authenticate
	fm.Input("web_authentication[email]", bmc.Email)
	fm.Input("web_authentication[password]", bmc.Password)
	if fm.Submit() != nil {
		panic(err)
	}
	if bow.StatusCode() != 200 {
		panic("Bad HTTP status code from login page!")
	}

	// Look for errors authenticating
	bow.Find("div.sc-Flash-message").Each(func(_ int, s *goquery.Selection) {
		panic("LOGIN FAILURE: " + strings.Trim(s.Text(), " \r\n\t"))
	})

	return bow.Dom(), true
}
