package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var echo string

func main() {
	echo = os.Getenv("ECHO")
	if echo == "" {
		echo = "UNDEFINED"
	}
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Println("received GET request /")
	fmt.Fprintf(w, "%s - %s!", echo, r.URL.Path[1:])
}
