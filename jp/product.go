package jp

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BaseUrl = "http://www.amazon.co.jp/dp"

type Product struct {
	Doc   *goquery.Document
	Asin  string
	Price int
}

func NewProduct(asin string) (*Product, error) {
	p := Product{Asin: asin}
	p.getDocument()
	if p.Doc == nil {
		return &p, fmt.Errorf("product not found. ASIN:%s", asin)
	}

	r := strings.NewReplacer(" ", "", "￥", "", ",", "")
	var err error
	p.Doc.Find("#priceblock_saleprice,#priceblock_ourprice").Each(func(i int, s *goquery.Selection) {
		price, err := strconv.ParseInt(r.Replace(s.Text()), 10, 0)
		if err == nil {
			p.Price = int(price)
		}
	})
	if err != nil {
		return &p, err
	}
	return &p, nil
}

func (p *Product) getDocument() {
	var err error
	p.Doc, err = goquery.NewDocument(fmt.Sprintf("%s/%s", BaseUrl, p.Asin))
	if err != nil {
		log.Fatal(err)
	}
}
