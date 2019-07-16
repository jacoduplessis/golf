package main

import (
	"fmt"
	"github.com/jacoduplessis/golf/ss/server"
	"log"
	"net/http"
	"os"
)

func getListenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8000"
}

func main() {

	c := http.Client{}
	mux := server.GetMux(c)
	listenAddr := getListenAddr()
	fmt.Printf("Listening on %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
