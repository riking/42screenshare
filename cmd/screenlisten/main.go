package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/riking/42screenshare/receiver"
)

func main() {
	listenPort := flag.Int("port", 4242, "port to listen on")
	httpPort := flag.Int("http", 3000, "port to serve HTTP on")

	flag.Parse()

	imgListen := fmt.Sprintf(":%d", *listenPort)
	httpListen := fmt.Sprintf("localhost:%d", *httpPort)

	receiver.SetupHTTP()
	l, err := net.Listen("tcp", imgListen)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Could not bind to localhost:%d", *listenPort))
		return
	}

	ch := make(chan struct{})
	go func() {
		fmt.Println("Listening for images on", imgListen)
		err := receiver.ReceiveImages(l)
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "accept()"))
		ch <- struct{}{}
	}()
	go func() {
		fmt.Println("Serving HTTP at http://" + httpListen)
		err := http.ListenAndServe(httpListen, nil)
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "http.Serve()"))
		ch <- struct{}{}
	}()

	<-ch
	return // exit()
}
