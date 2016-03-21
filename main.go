package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0tsuki/amazon/jp"
	"github.com/0tsuki/amazonpa"
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

func ItemLookup() {
	var asin = flag.String("asin", "", "ID of an amazon product to find")
	var group = flag.String("resg", "ItemAttributes,Offers,SalesRank,BrowseNodes", "Response group, comma separeted")
	flag.Parse()

	c := amazonpa.Client{
		AssociateTag: os.Getenv("AMAZONPA_ASSOCIATE_TAG"),
		AccessKeyId:  os.Getenv("AMAZONPA_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("AMAZONPA_SECRET_KEY"),
		Host:         os.Getenv("AMAZONPA_HOST"),
	}

	g := strings.Split(*group, ",")
	resp, err := c.ItemLookup(*asin, "ASIN", g)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(resp)
}

func main() {
	asins := jp.GetBestsellers("electronics")
	fmt.Printf("%v\n", asins)
}
