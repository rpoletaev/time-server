package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = 37
const epochDelta = 2208988800
const readDuration = 2 * time.Second

func main() {
	port := determinePort()
	errLogger := log.New(os.Stderr, "error:", log.LstdFlags)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		errLogger.Fatal(err)
	}

	log.Println("listen on:", listener.Addr().String())

	for {
		// waiting user
		conn, err := listener.Accept()
		if err != nil {
			errLogger.Println(err)
			continue
		}

		go func(c net.Conn) {
			// close connection:
			// 1.on writing error
			// 2.after user close connection
			// 3.if the user does not close the connection for too long
			defer c.Close()

			t := time.Now()
			log.Printf("Connection accepted:%s unix:%d\n", conn.RemoteAddr().String(), t.Unix())

			// send time info
			if err := sendTime(conn, t); err != nil {
				// log.Logger.Output allow using it on goroutines https://golang.org/src/log/log.go?s=6348:6390#L153
				errLogger.Println(err)
				return
			}

			// waiting for the user to close connection
			conn.SetReadDeadline(time.Now().Add(readDuration))
			conn.Read(nil)
		}(conn)
	}
}

func sendTime(w io.Writer, t time.Time) error {
	uintT := uintTime(t)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uintT) // network byte order

	if _, err := w.Write(buf); err != nil {
		return err
	}

	return nil
}

func determinePort() int {

	var port int
	flag.IntVar(&port, "p", defaultPort, "listening port")
	flag.Parse()

	return port
}

func uintTime(t time.Time) uint32 {
	// add delta between 1900 and 1970
	return uint32(t.Unix() + epochDelta)
}
