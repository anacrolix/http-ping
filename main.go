package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/anacrolix/missinggo"
)

type httpPinger struct {
	strAddr string
	port    int
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	urlStr := flag.Arg(0)
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	strAddr := u.Host
	missinggo.SplitHostMaybePort(u.Host)
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	if err != nil {
		log.Fatalln(time.Since(t), err)
	}
	log.Println(time.Since(t), "resolved", strAddr, "to", tcpAddr)
	t = time.Now()
	conn, err := net.Dial("tcp", tcpAddr.String())
	if err != nil {
		log.Fatalln(time.Since(t), err)
	}
	defer conn.Close()
	log.Println(time.Since(t), "dialed", tcpAddr)
	t = time.Now()
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "torrent.express",
		// MinVersion:         tls.VersionTLS12,
	})
	log.Println(time.Since(t), "TLS client")
	err = tlsConn.Handshake()
	if err != nil {
		log.Fatalln(time.Since(t), "error during TLS handshake:", err)
	}
	log.Println(time.Since(t), "TLS handshake")
	http.NewRequest("GET", urlStr, nil)
}
