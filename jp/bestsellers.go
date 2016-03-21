package jp

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const urlFormat = "http://www.amazon.co.jp/gp/bestsellers/%s/ref=zg_bs_nav_0"

// GetBestsellers returns the ASINs of Amazon Bestsellers
func GetBestsellers(tag string) (asins []string) {
	for _, i := range []int{1, 2, 3, 4, 5} {
		var url string
		if i == 1 {
			url = fmt.Sprintf(urlFormat, tag)
		} else {
			url = fmt.Sprintf(urlFormat, tag) + fmt.Sprintf("&pg=%d", i)
		}
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}

		// Find ASINs
		re := regexp.MustCompile(`/dp/([0-9A-Z]+)/`)
		doc.Find(".zg_title a").Each(func(i int, s *goquery.Selection) {
			str, exists := s.Attr("href")
			if exists {
				res := re.FindStringSubmatch(str)
				asins = append(asins, res[1])
			} else {
				log.Fatal("not found href")
			}
		})
	}
	return asins
}

// URL http://www.amazon.co.jp/gp/bestsellers/electronics/ref=zg_bs_nav_0
