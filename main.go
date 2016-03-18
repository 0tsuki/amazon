package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/0tsuki/amazon/jp"
	"github.com/k0kubun/pp"
)

func AmazonScrape(asin string) {
	p, err := jp.NewProduct(asin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", p.Title)
	fmt.Printf("Price: %d\n", p.Price)
	pp.Print(p.Review)
}

func main() {
	var asin = flag.String("asin", "", "ID of an amazon product to find")
	flag.Parse()
	AmazonScrape(*asin)
}
