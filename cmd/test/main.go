package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

const subMsg = `{"ids":["ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace"],"type":"subscribe","binary":false}`
const host = "kiloex-pythnet-ae42.mainnet.pythnet.rpcpool.com"
const origin = "https://app.kiloex.io"

func main() {
	var wssUrl = url.URL{Scheme: "wss", Host: host, Path: "/ws"}
	var headers = http.Header{}
	//headers.Set("Host", host)
	//headers.Set("Pragma", "no-cache")
	//headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	headers.Set("Origin", origin)
	//headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	//headers.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")

	c, _, err := websocket.DefaultDialer.Dial(wssUrl.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s, type: %d", message, mt)
		}
	}()

	if err := c.WriteMessage(websocket.TextMessage, []byte(subMsg)); err != nil {
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

	<-done

}
