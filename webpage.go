package gobase

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

func GetWebPageTitle(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr, Timeout: 3 * time.Second}

	resp, err := client.Get(url)

	if err != nil {
		return "", genWebPageError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", genWebPageError(fmt.Errorf("status code is not 200: %d", resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", genWebPageError(err)
	}

	title := doc.Find("title").Text()
	if title != "" {
		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			if t, e := s.Attr("property"); e {
				if t == "og:title" || t == "twitter:title" {
					if t, _ = s.Attr("content"); t != "" {
						title = t
					}
				}
			}
		})
	}
	return title, nil
}

func genWebPageError(err error) error {
	return fmt.Errorf("process webpage failed. %s", err)
}
