package wsclient

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

const host = "kiloex-pythnet-ae42.mainnet.pythnet.rpcpool.com"
const origin = "https://app.kiloex.io"

type kiloexPriceFeed struct {
	ctx context.Context
	c   *websocket.Conn

	msgCh chan Msg
	errCh chan error
}

func (k *kiloexPriceFeed) Connect() error {
	var wssUrl = url.URL{Scheme: "wss", Host: host, Path: "/ws"}
	var headers = k.buildHeaders()

	c, _, err := websocket.DefaultDialer.DialContext(k.ctx, wssUrl.String(), headers)
	if err != nil {
		return errors.New("dial:" + err.Error())
	}

	k.c = c

	go func() {
		defer c.Close()
		<-k.ctx.Done()
	}()

	go func() {
		for {
			mt, message, err := k.c.ReadMessage()
			if err != nil {
				k.errCh <- err
				return
			}

			k.msgCh <- Msg{
				Type: mt,
				Body: message,
			}
		}
	}()

	return nil
}

func (k *kiloexPriceFeed) Send(msg Msg) error {
	if err := k.c.WriteMessage(msg.Type, msg.Body); err != nil {
		return err
	}

	return nil
}

func (k *kiloexPriceFeed) Feed() (<-chan Msg, <-chan error) {
	return k.msgCh, k.errCh
}

func (k *kiloexPriceFeed) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Host", host)
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	headers.Set("Origin", origin)
	headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	headers.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")

	return headers
}

func NewPriceFeed(ctx context.Context) WsClient {
	return &kiloexPriceFeed{
		ctx:   ctx,
		msgCh: make(chan Msg, 1),
		errCh: make(chan error, 1),
	}
}
