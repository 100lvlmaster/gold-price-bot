package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
	u := launcher.New().
		Headless(true).
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	page := browser.Context(ctx).MustPage()

	log.Println("Navigating to target site...")
	err := rod.Try(func() {
		page.MustNavigate("https://www.ibja.co/")
		page.MustWaitLoad()
	})
	if err != nil {
		return 0, fmt.Errorf("navigation failed: %v", err)
	}

	log.Println("Page loaded. Waiting for stability...")
	time.Sleep(2 * time.Second)

	log.Println("Locating element...")
	var priceStr string
	err = rod.Try(func() {
		el := page.MustElement("#lblFineGold999")
		el.MustWaitVisible()
		priceStr = el.MustText()
	})
	if err != nil {
		return 0, fmt.Errorf("element rendering timed out: %v", err)
	}

	numVal, err := extractNum(priceStr)
	if err != nil {
		return 0, err
	}

	return float64(numVal) * gstMultiplier, nil
}
