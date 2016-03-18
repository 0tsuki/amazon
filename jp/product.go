package jp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BaseUrl = "http://www.amazon.co.jp/dp"

type Product struct {
	Asin   string
	Title  string
	Price  int
	Review *Review
	Doc    *goquery.Document
}

type Review struct {
	Value     float64
	Customers int
}

func NewProduct(asin string) (*Product, error) {
	p := Product{Asin: asin}
	p.getDocument()
	if p.Doc == nil {
		return &p, fmt.Errorf("product not found. ASIN:%s", asin)
	}

	// Title
	p.Doc.Find("#productTitle").Each(func(i int, s *goquery.Selection) {
		p.Title = s.Text()
	})

	// Price
	var err error
	r := strings.NewReplacer(" ", "", "￥", "", ",", "")
	p.Doc.Find("#priceblock_saleprice,#priceblock_ourprice").Each(func(i int, s *goquery.Selection) {
		price, err := strconv.ParseInt(r.Replace(s.Text()), 10, 0)
		if err != nil {
			log.Fatal("can not parse price. ASIN: " + asin)
		} else {
			p.Price = int(price)
		}
	})
	if err != nil {
		return &p, err
	}

	// Review
	var review Review
	vreg := regexp.MustCompile(`[0-9\.]+\z`)
	p.Doc.Find("#acrPopover").Each(func(i int, s *goquery.Selection) {
		trimed := strings.TrimSpace(s.Text())
		if vreg.MatchString(trimed) {
			review.Value, err = strconv.ParseFloat(vreg.FindString(trimed), 64)
			if err != nil {
				log.Fatalf("can not parse review. string: %s, error: %s", trimed, err)
			}
		} else {
			log.Fatalf("can not parse review. string: %s", trimed)
		}
	})

	creg := regexp.MustCompile(`^[0-9\.]+`)
	p.Doc.Find("#acrCustomerReviewText").Each(func(i int, s *goquery.Selection) {
		trimed := strings.TrimSpace(s.Text())
		if creg.MatchString(trimed) {
			num, err := strconv.ParseInt(creg.FindString(trimed), 10, 0)
			if err != nil {
				log.Fatalf("can not parse review. string: %s, error: %s", trimed, err)
			} else {
				review.Customers = int(num)
			}
		} else {
			log.Fatalf("can not parse review. string: %s", trimed)
		}
	})
	p.Review = &review

	return &p, nil
}

func (p *Product) getDocument() {
	var err error
	p.Doc, err = goquery.NewDocument(fmt.Sprintf("%s/%s", BaseUrl, p.Asin))
	if err != nil {
		log.Fatal(err)
	}
}
