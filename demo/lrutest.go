package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wh-data/lru_cache"
)

func main() {
	go func() {
		//start pprof
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()
	lruCache := initTest()
	count := 0
	for count < 10 { //36367 {
		count++
		key := strconv.Itoa(count)
		lruCache.Set(key, "test", 20)
		//fmt.Println("count: ", count)
		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println("waiting for gc")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "loadbig":
			loadBig(lruCache)
		case "loadsmall":
			loadSmall(lruCache)
			fmt.Println("finish load small")
		case "pointer":
			pointer(lruCache)
		case "getbig":
			getbig(lruCache)
		case "getsmall":
			getsmall(lruCache)
		case "gc":
			runtime.GC()
		case "view":
			fmt.Println(lruCache.ViewMap())
			fmt.Println(lruCache.ViewLinkedList())
		default:
			//lruCache.Get(text)
			fmt.Println(lruCache.Get(text))
		}
		fmt.Println(lruCache.GetSize())

	}

}

func initTest() *lru_cache.LRUCache {
	//test
	lruCache := lru_cache.NewLRUCache(36367)
	fmt.Println(lruCache.GetCapacity())
	fmt.Println("===========")
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Print("size and err: ")
	fmt.Println(lruCache.GetSize())
	fmt.Print("get key : 1, res and err: ")
	fmt.Println(lruCache.Get("1"))
	fmt.Print("get key : 2, res and err: ")
	fmt.Println(lruCache.Get("2"))
	fmt.Println("===========")
	lruCache.Set("1", "a", -1)
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Print("size and err: ")
	fmt.Println(lruCache.GetSize())
	fmt.Print("get key : 1, res and err: ")
	fmt.Println(lruCache.Get("1"))
	fmt.Print("get key : 2, res and err: ")
	fmt.Println(lruCache.Get("2"))
	fmt.Println("===========")
	lruCache.Set("2", "b", 3)
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Print("size and err: ")
	fmt.Println(lruCache.GetSize())
	fmt.Print("get key : 1, res and err: ")
	fmt.Println(lruCache.Get("1"))
	fmt.Print("get key : 2, res and err: ")
	fmt.Println(lruCache.Get("2"))
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Print("get key : 1, res and err: ")
	fmt.Println(lruCache.Get("1"))
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Println("===========")
	fmt.Print("key 1 expire: ")
	fmt.Println(lruCache.ViewExpire("1"))
	fmt.Print("key 2 expire: ")
	fmt.Println(lruCache.ViewExpire("2"))
	fmt.Println("===========")
	//time.Sleep(40 * time.Second)
	fmt.Print("key 1 expire: ")
	fmt.Println(lruCache.ViewExpire("1"))
	fmt.Print("key 2 expire: ")
	fmt.Println(lruCache.ViewExpire("2"))
	fmt.Println("===========")
	fmt.Print("list and err: ")
	fmt.Println(lruCache.ViewLinkedList())
	fmt.Print("map and err: ")
	fmt.Println(lruCache.ViewMap())
	fmt.Print("size and err: ")
	fmt.Println(lruCache.GetSize())
	fmt.Print("get key : 1, res and err: ")
	fmt.Println(lruCache.Get("1"))
	fmt.Print("get key : 2, res and err: ")
	fmt.Println(lruCache.Get("2"))
	return lruCache
}

func loadBig(cache *lru_cache.LRUCache) {
	content, err := ioutil.ReadFile("bigfile.txt")
	for i := 0; i < 1000; i++ {
		suf := strconv.Itoa(i)
		if err != nil {
			fmt.Println("load err", err)
		}
		if len(content) < 1 {
			fmt.Println("load err: empty content")
			return
		}
		var newcontent []byte
		newcontent = append(newcontent, content...)
		newcontent = append(newcontent, []byte("big"+suf)...)
		cache.Set("big"+suf, newcontent, 10)
	}
}

func loadSmall(cache *lru_cache.LRUCache) {
	c := 0
	for c < 10000 {
		suf := strconv.Itoa(c)
		cache.Set("small"+suf, "small content"+suf, 10)
		c++
		time.Sleep(1 * time.Millisecond)
	}
}

func pointer(cache *lru_cache.LRUCache) {
	cache.Set("pointer", "pointer content", 10)
}

func getbig(cache *lru_cache.LRUCache) {
	c := 0
	for c < 1000 {
		suf := strconv.Itoa(c)
		cache.Get("big" + suf)
		c++
		time.Sleep(1 * time.Millisecond)
	}
}

func getsmall(cache *lru_cache.LRUCache) {
	c := 0
	for c < 10000 {
		suf := strconv.Itoa(c)
		cache.Get("small" + suf)
		c++
		time.Sleep(1 * time.Millisecond)
	}
}
