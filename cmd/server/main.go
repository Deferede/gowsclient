package main

import (
	"cryptowsclient/internal/wsclient"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var host = "0.0.0.0:8090"

var upgrader = websocket.Upgrader{} // use default options

func kiloexPriceFeed(w http.ResponseWriter, r *http.Request) {
	var priceFeed = wsclient.NewPriceFeed(r.Context())
	var msgCh <-chan wsclient.Msg
	var errCh <-chan error

	if err := priceFeed.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			panic(err)
		}
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	msgCh, errCh = priceFeed.Feed()

	go func() {
		for {
			select {
			case msg := <-msgCh:
				err = c.WriteMessage(msg.Type, msg.Body)
				if err != nil {
					log.Println("write error:", err)
					break
				}
			case err := <-errCh:
				err = c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				if err != nil {
					log.Println("write error:", err)
					break
				}
			}
		}
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s, type: %d", message, mt)

		msg := wsclient.Msg{
			Type: mt,
			Body: message,
		}

		if err := priceFeed.Send(msg); err != nil {
			if err = c.WriteMessage(mt, message); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/kiloexPriceFeed", kiloexPriceFeed)
	log.Fatal(http.ListenAndServe(host, nil))
}
