package sender

import (
	"image"
	"image/png"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"time"
	"os/signal"
	"fmt"
	"crypto/rand"
	"encoding/hex"
)

var randomFilename string

func getFilename() string {
	if randomFilename != "" {
		return randomFilename
	}

	// Clean up previous runs
	bytes, err := exec.Command("sh", "-c", "rm -f /tmp/screensend-*.png").CombinedOutput()
	if err != nil {
		os.Stderr.Write(bytes)
		os.Stderr.WriteString("\n")
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error removing previous .png files"))
	}
	var b [8]byte
	rand.Read(b[:])
	randomFilename = fmt.Sprintf("/tmp/screensend-%s.png", hex.EncodeToString(b[:]))
	return randomFilename
}

func deleteOnShutdown() {
	ch := make(chan os.Signal)

	signal.Notify(ch, os.Interrupt)
	<-ch
	if randomFilename != "" {
		os.Remove(randomFilename)
	}
}

func SendImage(connectStr string) (err error, suggestSleep time.Duration) {
	conn, err := net.Dial("tcp", connectStr)
	if err != nil {
		return errors.Wrap(err, "dial remote"), 2500*time.Millisecond
	}
	defer conn.Close()

	c := exec.Command("screencapture", "-x", "-C", getFilename())
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		return errors.Wrap(err, "pipe()"), 0
	}
	c.ExtraFiles = []*os.File{pipeW}

	err = c.Start()
	if err != nil {
		return errors.Wrap(err, "screencapture exec"), 0
	}
	// Need to close our side of pipe before trying to read from it
	pipeW.Close()

	var img image.Image
	var goErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		rawImg, err := png.Decode(pipeR)
		if err == nil {
			img = resize.Resize(1280, 720, rawImg, resize.Bilinear)
		} else {
			goErr = err
		}
		wg.Done()
	}()
	err = c.Wait() // screencapture
	if err != nil {
		return errors.Wrap(err, "screencapture exit"), 0
	}
	wg.Wait() // resize
	if goErr != nil {
		return errors.Wrap(err, "decoding png"), 0
	}
	var enc png.Encoder
	err = enc.Encode(conn, img)
	if err != nil {
		return errors.Wrap(err, "writing png"), 0
	}

	return nil, 0
}
