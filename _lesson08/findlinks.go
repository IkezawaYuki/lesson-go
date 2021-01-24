package main

import (
	"flag"
	"fmt"
	"gopl.io/ch5/links"
	"log"
	"math"
)

var tokens = make(chan struct{}, 20)

type leveledList struct {
	depth int
	lists []string
}

var depthFlag = flag.Int("depth", math.MaxInt32, "depth of links")

func crawl(depth int, url string) *leveledList {
	if depth > *depthFlag {
		return &leveledList{depth: depth + 1, lists: nil}
	}
	fmt.Printf("%3d: %s\n", depth, url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return &leveledList{depth: depth + 1, lists: list}
}

func main() {
	flag.Parse()
	worklist := make(chan *leveledList)
	var n int
	n++
	go func() {
		worklist <- &leveledList{0, flag.Args()}
	}()
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		llist := <-worklist
		for _, link := range llist.lists {
			if !seen[link] {
				seen[link] = true
				n++
				go func(depth int, link string) {
					worklist <- crawl(depth, link)
				}(llist.depth, link)
			}
		}
	}
}
