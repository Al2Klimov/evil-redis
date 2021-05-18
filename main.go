package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var listen = flag.String("listen", "[::]:6379", "HOST:PORT")
var errmsg = flag.String("error", "I'm sorry, Dave. I'm afraid I can't do that.", "MESSAGE")
var sleep = flag.Duration("sleep", time.Second, "DURATION")

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	srv, err := net.Listen("tcp", *listen)
	if err != nil {
		return err
	}
	defer srv.Close()

	for {
		clt, err := srv.Accept()
		if err != nil {
			return err
		}

		go handle(clt)
	}
}

func handle(conn net.Conn) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go read(conn, wg)
	go write(conn, wg)
	go clos(conn, wg)
}

func read(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	io.Copy(ioutil.Discard, conn)
}

func write(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if _, err := fmt.Fprintf(conn, "-%s\r\n", *errmsg); err != nil {
			break
		}

		time.Sleep(*sleep)
	}
}

func clos(conn net.Conn, wg *sync.WaitGroup) {
	wg.Wait()
	conn.Close()
	log.Print(conn.RemoteAddr().String())
}
