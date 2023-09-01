package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	wantTLS = flag.String("tls", "yes", "TLS connection: yes, no, force")
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "    mudc [options] <host> <port>")
	fmt.Fprintln(os.Stderr, "    mudc proxy [options] <host> <port> [<local-port>]")
	os.Exit(1)
}

func connect(addr string) net.Conn {
	var err error
	var conn net.Conn

	if *wantTLS == "yes" || *wantTLS == "force" {
		conn, err = tls.Dial("tcp", addr, &tls.Config{})
		if err != nil {
			if *wantTLS == "force" {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			conn = nil
		}
	}

	if conn == nil {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return conn
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 3 && len(args) <= 4 && args[0] == "proxy" {
		proxy(args[1:])
	} else if len(args) == 2 {
		client(args)
	} else {
		usage()
	}
}
