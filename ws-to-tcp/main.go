package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	flagBind    = flag.String("bind", ":8081", "bind to")
	flagBackend = flag.String("be", "127.0.0.1:4222", "nats server")
	flagToken   = flag.String("token", "test", "secret token for http connection")
	flagPath    = flag.String("http-path", "mq", "http path to websockets")
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func copyWorker(dst io.Writer, src io.Reader, doneCh chan bool) {
	io.Copy(dst, src)
	doneCh <- true
}

func main() {
	flag.Parse()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/"+*flagPath, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != *flagToken {
			http.NotFound(w, r)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		ncon, err := net.Dial("tcp", *flagBackend)
		if err != nil {
			log.Printf("%v", err)
			conn.Close()
			return
		}
		doneCh := make(chan bool)

		conn.UnderlyingConn().(*net.TCPConn).SetKeepAlivePeriod(1 * time.Second)

		go copyWorker(ncon, conn.UnderlyingConn(), doneCh)
		go copyWorker(conn.UnderlyingConn(), ncon, doneCh)

		<-doneCh

		ncon.Close()
		conn.Close()

		<-doneCh
	})

	http.ListenAndServe(*flagBind, nil)
}
