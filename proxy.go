package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leftmike/mudc/telnet"
)

func copyBytes(dst io.Writer, src io.Reader) {
	b := make([]byte, 1)
	for {
		_, err := src.Read(b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		_, err = dst.Write(b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
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

		fmt.Println("Connection from", svr.RemoteAddr())

		//svr = telnet.NewConn(svr)
		go copyBytes(svr, clnt) // XXX: io.Copy?
		go copyBytes(clnt, svr)
	}
}
