package telnet

import (
	"bufio"
	"fmt"
	"net"
)

const (
	iacByte  = 255
	sbByte   = 250
	seByte   = 240
	willByte = 251
	wontByte = 252
	doByte   = 253
	dontByte = 254
)

type telnetConn struct {
	net.Conn
	r *bufio.Reader
}

func NewConn(conn net.Conn) net.Conn {
	return &telnetConn{
		Conn: conn,
		r:    bufio.NewReader(conn),
	}
}

func (tc *telnetConn) Read(p []byte) (int, error) {
	n := 0

	for n < len(p) {
		b, err := tc.r.ReadByte()
		if err != nil {
			return 0, err
		}

		if b == iacByte {
			b, err = tc.r.ReadByte()
			if err != nil {
				return 0, err
			}
			switch b {
			case iacByte:
				p[n] = b
				n += 1
			case sbByte:
				fmt.Println("sbByte")
			case seByte:
				fmt.Println("seByte")
			case willByte, wontByte, doByte, dontByte:
				_, err = tc.r.ReadByte()
				if err != nil {
					return 0, err
				}
			default:
				fmt.Printf("unknown byte following IAC: %d\n", b)

			}
		} else {
			p[n] = b
			n += 1
		}
	}

	return n, nil
}

func (tc *telnetConn) Write(p []byte) (int, error) {
	return tc.Conn.Write(p)
}
