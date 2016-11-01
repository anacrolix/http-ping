package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/anacrolix/missinggo"
)

type httpPinger struct {
	strAddr string
	port    int
}

func schemePort(scheme string) int {
	switch scheme {
	case "http":
		return 80
	case "https":
		return 443
	}
	return 0
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
	hmp := missinggo.SplitHostMaybePort(u.Host)
	if hmp.NoPort {
		hmp.Port = schemePort(u.Scheme)
		hmp.NoPort = false
	}
	strAddr := hmp.String()
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
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "torrent.express",
		NextProtos:         []string{"h2"},
	})
	t = time.Now()
	err = tlsConn.Handshake()
	if err != nil {
		log.Fatalln(time.Since(t), "error during TLS handshake:", err)
	}
	log.Println(time.Since(t), "TLS handshake")
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	tlsConnBufR := bufio.NewReader(tlsConn)
	t = time.Now()
	err = req.Write(tlsConn)
	resp, err := http.ReadResponse(tlsConnBufR, req)
	if err != nil {
		log.Fatalln(time.Since(t), err)
	}
	log.Println(time.Since(t), "HTTP round trip")
	resp.Write(os.Stderr)
}
