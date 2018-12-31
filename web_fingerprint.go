package main

import (
	"log"
	"net/http"
	"strings"
	"golang.org/x/net/html"
	//"os"
	"time"
)

func web_fingerprint(ip string) string {
	client := http.Client{
		Timeout: time.Duration(3*time.Second),
	}
	res, err := client.Get("http://" + ip)
	if err != nil {
		log.Printf("web fingerprint: %s", err)
		return("")
	} else {
		defer res.Body.Close()
	}

	tokenizer := html.NewTokenizer(res.Body)
	found_title := false
	for {
		tok := tokenizer.Next()
		if tok == html.ErrorToken {
			log.Printf("HTML Parsing Error: %s", tokenizer.Err())
			break
		}
		//log.Printf("%s", tokenizer.Token())
		if found_title {
			return(tokenizer.Token().String())
		}
		if strings.ToLower(tokenizer.Token().String()) == "<title>" {
			found_title=true
		}
	}
	return("")
}

//func main() {
//	fp := web_fingerprint(os.Args[1])
//	log.Printf(fp)
//}