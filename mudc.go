package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leftmike/mudc/telnet"
)

func main() {
	forceTLS := flag.Bool("tls", false, "force TLS connection")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: mudc [options] <host> <port>")
		os.Exit(1)
	}

	addr := args[0] + ":" + args[1]

	var err error
	var conn net.Conn
	conn, err = tls.Dial("tcp", addr, &tls.Config{})
	if err != nil {
		if *forceTLS {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	conn = telnet.NewConn(conn)

	go func(conn net.Conn) {
		for {
			b := make([]byte, 1)
			_, err := conn.Read(b)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Print(string(b))
		}
	}(conn)

	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if b[0] == 13 {
			continue
		}

		_, err = conn.Write(b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
