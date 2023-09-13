package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/peterh/liner"

	"github.com/leftmike/mudc/telnet"
)

var (
	wantTrace   = flag.String("trace", "", "trace client input and/or output to mudc.log: input, output, both")
	inputTrace  io.Writer
	outputTrace io.Writer
)

var (
	promptRegex  = regexp.MustCompile(`^> `)
	coremudRegex = regexp.MustCompile(`^\[CoreMUD\] `)
	commRegex    = regexp.MustCompile(`^Comm [a-z]+: \[[a-zA-Z]+\]`)
)

type outputParser struct {
	buf      bytes.Buffer
	output   []string
	ansi     bool
	outputFn func(s string)
}

func (op *outputParser) writeByte(b byte) {
	if op.ansi {
		if b == 'm' {
			op.ansi = false
		}
		return
	}

	if b == 27 { // escape
		op.ansi = true
		return
	}

	if b != '\n' {
		op.buf.WriteByte(b)
		return
	}

	buf := op.buf.Bytes()
	if promptRegex.Match(buf) {
		op.outputFn(strings.Join(op.output, "\n"))
		op.output = op.output[:0]
		buf = buf[2:] // XXX: take length of prompt into account
	}

	if coremudRegex.Match(buf) || commRegex.Match(buf) {
		op.outputFn(string(buf))
	} else {
		op.output = append(op.output, string(buf))
	}

	op.buf.Reset()
}

func clientOutput(conn net.Conn) {
	var op outputParser
	buf := make([]byte, 1)

	for {
		_, err := conn.Read(buf)
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Print(string(buf))

		if outputTrace != nil {
			outputTrace.Write(buf)
		}

		for _, b := range buf {
			op.writeByte(b)
		}
	}
}

func client(args []string) {
	conn := telnet.NewConn(connect(args[0] + ":" + args[1]))

	if *wantTrace == "input" || *wantTrace == "output" || *wantTrace == "both" {
		w, err := os.OpenFile("mudc.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if *wantTrace == "input" || *wantTrace == "both" {
			inputTrace = w
		}
		if *wantTrace == "output" || *wantTrace == "both" {
			outputTrace = w
		}

		fmt.Fprintf(w, "\nmudc %d %s\n\n", os.Getpid(), time.Now().Format(time.UnixDate))
	}

	go clientOutput(conn)

	line := liner.NewLiner()
	defer line.Close()

	for {
		s, err := line.Prompt("")
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if inputTrace != nil {
			inputTrace.Write([]byte(s))
			inputTrace.Write([]byte{'\n'})
		}

		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte{'\n'})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
