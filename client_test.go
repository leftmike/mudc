package main

import (
	"reflect"
	"testing"
)

func TestOutputParser(t *testing.T) {
	cases := []struct {
		s      string
		output []string
	}{
		{
			s: "\x1b[49;49m\x1b[0;10mYou enter the mass-transit tube.\x1b[49;49m\x1b[0;10m\n" +
				"\x1b[49;49m\x1b[0;10mNorth Mine - Southern Transit Station\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
				"Despite the constant cleaning efforts, the maintenance crews haven't been able to\x1b[49;49m\x1b[0;10m\n" +
				"\x1b[49;49m\x1b[0;10mExits:\x1b[49;49m\x1b[0;10m south, northwest, north\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
				"\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10mTwo Mine Guards\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
				"Mine Guard Captain\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n",
			output: []string{
				"You enter the mass-transit tube.",
				"North Mine - Southern Transit Station",
				"Despite the constant cleaning efforts, the maintenance crews haven't been able to",
				"Exits: south, northwest, north",
				"Two Mine Guards",
				"Mine Guard Captain",
			},
		},
	}

	for _, c := range cases {
		var op outputParser
		for _, b := range []byte(c.s) {
			op.writeByte(b)
		}
		if !reflect.DeepEqual(op.output, c.output) {
			t.Errorf("outputParser.writeByte(%s) got %q, want %q", c.s, op.output, c.output)
		}
	}
}
