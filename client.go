package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/peterh/liner"

	"github.com/leftmike/mudc/telnet"
)

var (
	wantTrace   = flag.String("trace", "", "trace client input and/or output to mudc.log: input, output, both")
	inputTrace  io.Writer
	outputTrace io.Writer
)

/*
var (
	promptRegex  = regexp.MustCompile(`^> `)
	coremudRegex = regexp.MustCompile(`^\[CoreMUD\] `)
	commRegex    = regexp.MustCompile(`^Comm [a-z]+: \[[a-zA-Z]+\]`)
)
*/

func clientOutput(conn net.Conn) {
	//var buf bytes.Buffer
	//var output []string

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

		if outputTrace != nil {
			outputTrace.Write(b)
		}
		/*
			if b[0] == 10 {
				if coremudRegex.Match(buf.Bytes()) {
					fmt.Println("coremudRegex", buf.String())
				} else if commRegex.Match(buf.Bytes()) {
					fmt.Println("commRegex", buf.String())
				} else {
					output = append(output, buf.String())
					//fmt.Println("$$", buf.String())
				}

				buf.Reset()
			} else if promptRegex.Match(buf.Bytes()) {
				//fmt.Println("promptRegex", buf.String())
				if len(output) > 0 {
					for _, s := range output {
						fmt.Println("#", s)
					}
				}

				output = nil
			} else {
				buf.Write(b)
			}
		*/
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
