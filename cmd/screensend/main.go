package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/riking/42screenshare/sender"
)

var help = `
Usage: screensend [-port 4242] [-sleep 0.1s] <hostname>

    hostname:  The hostname of the machine to send the screen to.
               E.g. e1z1r2p20
    port:      The port to send screenshots on. E.g. 4256
    sleep:     How long to wait between sending images
`

func main() {
	sendPort := flag.Int("port", 4242, "port to listen on")
	sleepTime := flag.Duration("sleep", 100*time.Millisecond, "time to sleep between images")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println(help)
		return
	}
	hostname := flag.Arg(0)
	connectStr := fmt.Sprintf("%s:%d", hostname, *sendPort)

	for {
		fmt.Println("Sending...")
		err, sleepAdd := sender.SendImage(connectStr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Sent")
		}
		if sleepAdd != 0 {
			time.Sleep(sleepAdd)
		} else {
			time.Sleep(*sleepTime)
		}
	}
}
