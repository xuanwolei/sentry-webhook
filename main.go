package main

import (
	"fmt"
	"github.com/xuanwolei/sentry_webhook/internal"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0:80"
	}

	fmt.Println("listening addr:", addr)
	http.HandleFunc("/", internal.ServeHandle())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("listen error:", err)
	}

	fmt.Println("down")
}
