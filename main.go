package main

import (
	"fmt"
	"log"

	"github.com/0tsuki/amazon/jp"
)

func AmazonScrape() {
	p, err := jp.NewProduct("B019GNUT0C")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Price: %d\n", p.Price)
}

func main() {
	AmazonScrape()
}
