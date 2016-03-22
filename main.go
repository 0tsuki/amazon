package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"container/list"
	"time"
	"regexp"

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

func extractId(url string) (string, error) {
	re := regexp.MustCompile(`/(electronics/[0-9]+)/ref=`)
	res := re.FindStringSubmatch(url)
	if res == nil {
		return "", fmt.Errorf("cannot extract %s", url)
	}
	return res[1], nil
}

type Task struct {
	Todo *list.List
	Done map[string]string
	Check map[string]string
}

func NewTask() *Task {
	return &Task{
		Todo: list.New(),
		Done: make(map[string]string),
		Check: make(map[string]string),
	}
}

func (task *Task) PushBackIfNon(url string) (ok bool, err error) {
	id, err := extractId(url)
	if err != nil {
		return false, err
	}
	_, exists := task.Check[id]
	if exists {
		log.Printf("already exists %s", id)
		return false, nil
	}
	task.Check[id] = url
	task.Todo.PushBack(url)
	return true, nil
}

func (todo *Task) MergeUrls(urls []string) error {
	for _, u := range urls {
		todo.PushBackIfNon(u)
	}
	return nil
}

func main() {
	urlsToSearch := jp.GetBestsellerUrl("http://www.amazon.co.jp/gp/bestsellers/electronics/ref=zg_bs_nav_0/375-5078344-1017226")
	task := NewTask()
	for _, url := range urlsToSearch {
		task.PushBackIfNon(url)
	}

	for e := task.Todo.Front(); e != nil; e = e.Next() {
		time.Sleep(5 * time.Second)
		url, ok := e.Value.(string)
		if !ok {
			continue
		}

		log.Printf("searching ... %s", url)
		urls := jp.GetBestsellerUrl(url)
		id, err := extractId(url)
		if err != nil {
			log.Fatal(err)
		}
		task.Done[id] = url
		task.MergeUrls(urls)

		log.Printf("todo: %d, done: %d", task.Todo.Len(), len(task.Done))
	}

	fmt.Println("--------- todo -------")
	pp.Println(task.Todo)

	fmt.Println("--------- done -----------")
	pp.Println(task.Done)
}
