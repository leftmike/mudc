package main

import (
	"reflect"
	"testing"
)

func TestStripAnsi(t *testing.T) {
	s := "\x1b[49;49m\x1b[0;10mYou enter the mass-transit tube.\x1b[49;49m\x1b[0;10m\n" +
		"\x1b[49;49m\x1b[0;10mNorth Mine - Southern Transit Station\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
		"Despite the constant cleaning efforts, the maintenance crews haven't been able to\x1b[49;49m\x1b[0;10m\n" +
		"\x1b[49;49m\x1b[0;10mExits:\x1b[49;49m\x1b[0;10m south, northwest, north\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
		"\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10mTwo Mine Guards\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n" +
		"Mine Guard Captain\x1b[49;49m\x1b[0;10m\x1b[49;49m\x1b[0;10m\n"
	output := []string{
		"You enter the mass-transit tube.",
		"North Mine - Southern Transit Station",
		"Despite the constant cleaning efforts, the maintenance crews haven't been able to",
		"Exits: south, northwest, north",
		"Two Mine Guards",
		"Mine Guard Captain",
	}

	var op outputParser
	for _, b := range []byte(s) {
		op.writeByte(b)
	}
	if !reflect.DeepEqual(op.output, output) {
		t.Errorf("outputParser.writeByte(%s) got %q, want %q", s, op.output, output)
	}
}

func TestOutputParser(t *testing.T) {
	cases := []struct {
		s      string
		output []string
	}{
		{
			s: `abc def ghi
jkl mno pqr
> 
`,
			output: []string{
				`abc def ghi
jkl mno pqr`,
			},
		},
		{
			s: `North Mine - Southern Transit Station
Exits: south, northwest, north
Two Mine Guards
Mine Guard Captain
> North Mine - Central Transit Station
Exits: up, southwest, southeast, south, northwest, northeast, north, down
Miner service station, Waste Receptacle, and mining shift schedule are here.
> This mine shaft is a testimony to the courage and determination of The Company's workers.
Exits: west, up, south, north, east, down
> The 'walls' of the shaft have been roughly hewn from solid granite and bedrock.
Exits: west, up, south, north, east, down
> 
`,
			output: []string{
				"North Mine - Southern Transit Station\nExits: south, northwest, north\nTwo Mine Guards\nMine Guard Captain",
				"North Mine - Central Transit Station\nExits: up, southwest, southeast, south, northwest, northeast, north, down\nMiner service station, Waste Receptacle, and mining shift schedule are here.",
				"This mine shaft is a testimony to the courage and determination of The Company's workers.\nExits: west, up, south, north, east, down",
				"The 'walls' of the shaft have been roughly hewn from solid granite and bedrock.\nExits: west, up, south, north, east, down",
			},
		},
		{
			s: `Room number one
Exits: north, east, south, west
> [CoreMUD] PlayerOne enters the game.
Room number two
Exits: north, east, south, west
> 
`,
			output: []string{
				"Room number one\nExits: north, east, south, west",
				"[CoreMUD] PlayerOne enters the game.",
				"Room number two\nExits: north, east, south, west",
			},
		},
		{
			s: `Room number one
Exits: north, east, south, west
> [CoreMUD] PlayerOne enters the game.
[CoreMUD] PlayerTwo leaves the game.
Room number two
Exits: north, east, south, west
> Comm rl: [Root] Hello, World!
Comm rl: [Root] Is there anyone out there?
Room number three
> 
`,
			output: []string{
				"Room number one\nExits: north, east, south, west",
				"[CoreMUD] PlayerOne enters the game.",
				"[CoreMUD] PlayerTwo leaves the game.",
				"Room number two\nExits: north, east, south, west",
				"Comm rl: [Root] Hello, World!",
				"Comm rl: [Root] Is there anyone out there?",
				"Room number three",
			},
		},
	}

	for _, c := range cases {
		var output []string
		op := outputParser{
			outputFn: func(s string) {
				output = append(output, s)
			},
		}

		for _, b := range []byte(c.s) {
			op.writeByte(b)
		}

		if !reflect.DeepEqual(output, c.output) {
			t.Errorf("outputParser.writeByte(%s) got %q, want %q", c.s, output, c.output)
		}
	}
}
