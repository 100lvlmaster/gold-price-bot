package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func extractNum(input string) (int, error) {
	var buf bytes.Buffer
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			buf.WriteByte(input[i])
		}
	}
	if buf.Len() == 0 {
		return 0, errors.New("could not parse digits")
	}
	result, err := strconv.Atoi(buf.String())
	if err != nil {
		return 0, fmt.Errorf("could not convert to int: %v", err)
	}
	return result, nil
}

func ScrapeGoldPrice(gstMultiplier float64) (float64, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"),
	)

	var priceStr string
	var err error

	c.OnHTML("#lblFineGold999", func(e *colly.HTMLElement) {
		priceStr = strings.TrimSpace(e.Text)
	})

	c.OnError(func(r *colly.Response, e error) {
		err = fmt.Errorf("request failed with status %d: %v", r.StatusCode, e)
	})

	log.Println("Navigating to target site with Colly...")
	visitErr := c.Visit("https://www.ibja.co/")
	if visitErr != nil {
		return 0, visitErr
	}

	if err != nil {
		return 0, err
	}

	if priceStr == "" {
		return 0, errors.New("could not find price element on page")
	}

	log.Printf("Found price string: %s", priceStr)

	numVal, err := extractNum(priceStr)
	if err != nil {
		return 0, err
	}

	return float64(numVal) * gstMultiplier, nil
}
