package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	time.Sleep(time.Second * 10)
	for _ = range time.Tick(time.Second) {
		resp, err := http.Get("http://echo:80")
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("✅ Request complete - %s", string(bodyBytes))
	}
}
