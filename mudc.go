package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/reiver/go-telnet"
)

func main() {
	forceTLS := flag.Bool("tls", false, "force TLS connection")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: mudc [options] <host> <port>")
		os.Exit(1)
	}

	conn, err := telnet.DialToTLS(flag.Arg(0)+":"+flag.Arg(1), &tls.Config{})
	if err != nil {
		if *forceTLS {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		conn, err = telnet.DialTo(flag.Arg(0) + ":" + flag.Arg(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	go func(conn *telnet.Conn) {
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
		fmt.Printf("%v ", b[0])
		_, err = conn.Write(b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	/*
		var caller telnet.Caller = telnet.StandardCaller

		//@TOOD: replace "example.net:5555" with address you want to connect to.
		//telnet.DialToAndCall("coremud.org:4000", caller)
		telnet.DialToAndCallTLS("coremud.org:4022", caller, &tls.Config{})
	*/
	/*
		// connect to a socket
		conn, err := net.Dial("tcp", "coremud.org:4000")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// send a string
		fmt.Fprintf(conn, "Hello, world!\n")
		// read the response
		for {
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println(status)
		}
	*/
}
