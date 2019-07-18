package main

import (
	"fmt"
	"github.com/jacoduplessis/golf/ss/server"
	"log"
	"net/http"
	"time"
)

func main() {

	c := http.Client{
		Timeout: time.Second * 10,
	}

	srv := server.GetServer(c)
	fmt.Printf("Listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
