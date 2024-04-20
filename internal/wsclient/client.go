package wsclient

type WsClient interface {
	Connect() error
	Send(msg Msg) error
	Feed() (<-chan Msg, <-chan error)
}
type Msg struct {
	Type int
	Body []byte
}
