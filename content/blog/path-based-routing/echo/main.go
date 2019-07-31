package main

import (
	"fmt"
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
	fmt.Fprintf(w, "%s - %s!", echo, r.URL.Path[1:])
}
