package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const schiitDealsURL = "https://www.schiit.com/b-stocks"
const webhookURL = "https://maker.ifttt.com/trigger/schiit-sale/with/key/kx0E2XfhHj0FJtUKDnGcYxZQn0mOiT3PFpUAP3sQ1M4"

func callWebhook(productName, price string) {
	requestBody, err := json.Marshal(map[string]string{
		"value1": productName,
		"value2": price,
	})

	if err != nil {
		log.Fatal(err)
	}

	request, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("GET", schiitDealsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "BStockCheckerBot v1.0 - currently on the lookout for Modi 3 and Loki B-stocks :)")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	document.Find(".product").Each(func(i int, s *goquery.Selection) {
		productName := strings.TrimSpace(s.Find("div.title").Text())
		price := strings.TrimSpace(s.Find("div.price").Text())

		if productName == "Modi 3" {
			callWebhook(productName, price)
		} else if productName == "Loki" {
			callWebhook(productName, price)
		} else if productName == "Magni 3" {
			callWebhook(productName, price)
		}
	})
}
