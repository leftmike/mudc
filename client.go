package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leftmike/mudc/telnet"
)

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
