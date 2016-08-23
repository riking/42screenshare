package receiver

import (
	"bytes"
	"fmt"
	"image/png"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

func ReceiveImages(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go getImage(conn)
	}
}

func getImage(c net.Conn) {
	fmt.Println("Connected to", c.RemoteAddr())
	var buf bytes.Buffer

	n, err := buf.ReadFrom(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "read from socket"))
		c.Close()
		return
	}
	c.Close()

	if n == 0 {
		fmt.Println("Got 0-byte connection, ignoring")
		return
	}

	b := buf.Bytes()
	_, err := png.DecodeConfig(bytes.NewReader(b))
	if err != nil {
		fmt.Println("Invalid PNG:", err)
		return
	}

	imageMutex.Lock()
	defer imageMutex.Unlock()
	imageLastUpdated = time.Now()
	imageBytes = b
	fmt.Println("Received new image of size", len(b), "at", imageLastUpdated)
}
