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

var (
	forceTLS = flag.Bool("tls", false, "force TLS connection")
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

	return conn
}

func client(args []string) {
	conn := telnet.NewConn(connect(args[0] + ":" + args[1]))

	go func(conn net.Conn) {
		b := make([]byte, 1)
		for {
			_, err := conn.Read(b)
			if err == io.EOF {
				os.Exit(0)
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Print(string(b))
		}
	}(conn)

	b := make([]byte, 1)
	for {
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

func copyBytes(dst io.Writer, src io.Reader, s string) {
	b := make([]byte, 1)
	for {
		_, err := src.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		dst.Write(b)
		//fmt.Printf("%s%d ", s, b[0])
	}
}

func proxy(args []string) {
	clnt := telnet.NewConn(connect(args[0] + ":" + args[1]))

	var port string
	if len(args) == 3 {
		port = args[2]
	} else {
		port = args[1]
	}
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		svr, err := l.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		svr = telnet.NewConn(svr)
		go copyBytes(svr, clnt, "c>s:")
		go copyBytes(clnt, svr, "s>c:")
	}
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
