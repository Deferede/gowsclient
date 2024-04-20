package main

import (
	"context"
	"cryptowsclient/internal/wsclient"
	"github.com/gorilla/websocket"
	"log"
)

const subMsg = `{"ids":["ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace"],"type":"subscribe","binary":false}`

func main() {
	ctx := context.Background()

	feed := wsclient.NewPriceFeed(ctx)

	if err := feed.Connect(); err != nil {
		panic(err)
	}

	if err := feed.Send(wsclient.Msg{
		Type: websocket.TextMessage,
		Body: []byte(subMsg),
	}); err != nil {
		panic(err)
	}

	msgChan, errChan := feed.Feed()

	for {
		select {
		case msg := <-msgChan:
			log.Println("receive ", string(msg.Body))
		case err := <-errChan:
			panic(err)
		}
	}
}
