// +build linux,darwin
// The main package handles connecting to the door lock and receiving various
// command line flags, including the JSON users database file, serial port location
// and the baud rate.
package main

import (
	"flag"
	"github.com/chimera/go-inside/door"
)

// Available command line flags with sane-ish defaults.
var users_file = flag.String("db", "users.json", "The users JSON file to use.")
var port_path = flag.String("port", "/dev/tty.usbmodem621", "The serial port that the Arduino is running on.")
var baud = flag.Int("baud", 19200, "The baudrate to connect to the serial port with.")

func main() {
	// Parse any command line flags.
	flag.Parse()

	// Create a new connection to the door lock
	door := &door.DoorLock{
		Baud:           *baud,
		SerialPortPath: *port_path,
		UsersFile:      *users_file,
	}

	// Handle inputting of user RFID codes
	door.Listen()
}
