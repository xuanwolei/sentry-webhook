package main

import (
	"fmt"
	"github.com/xuanwolei/sentry_webhook/internal"
	"net/http"
)

func main() {
	addr := "0.0.0.0:8960"
	fmt.Printf("listening addr:%s", addr)
	http.HandleFunc("/", internal.ServeHandle())
	http.ListenAndServe(addr, nil)
	fmt.Println("down")
}
