package main

import (
	"flag"
	"log"
	"net"
	"net/url"
	"runtime"
	"time"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"

	"github.com/gorilla/websocket"
)

var (
	flagURL  = flag.String("to", "wss://hostname.tld?token=test", "host to connect")
	flagUser = flag.String("nats-user", "", "nats user")
	flagPass = flag.String("nats-pass", "", "nats password")
)

type customDialer struct{}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {

	u, _ := url.Parse(*flagURL)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.UnderlyingConn(), nil
}

func main() {
	flag.Parse()

	opts := []nats.Option{
		nats.SetCustomDialer(&customDialer{}),
		nats.ReconnectWait(1 * time.Second),
	}

	if *flagUser != "" {
		if *flagPass == "" {
			opts = append(opts, nats.Token(*flagUser))
		} else {
			opts = append(opts, nats.UserInfo(*flagUser, *flagPass))
		}
	}

	nc, err := nats.Connect(*flagURL, opts...)
	if err != nil {
		panic(err)
	}

	sc, err := stan.Connect("test-cluster", "test", stan.NatsConn(nc))
	if err != nil {
		panic(err)
	}

	sc.Subscribe("test", func(msg *stan.Msg) {
		log.Printf("%s", msg.Data)
	})

	go func() {
		for {
			sc.Publish("test", []byte("test message"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	runtime.Goexit()

}
