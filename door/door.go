package door

import (
	"fmt"
	"github.com/chimera/go-inside/rs232"
	"github.com/chimera/go-inside/users"
	"log"
	// "github.com/distributed/sers"
)

type DoorLock struct {
	Serial         rs232.Port
	Baud           int
	SerialPortPath string
	UsersFile      string
}

func (d *DoorLock) Connect() {
	log.Println("Connecting to door lock via serial...")

	// Configure the serial connection.
	options := rs232.Options{
		BitRate:  uint32(d.Baud),
		DataBits: 8,
		StopBits: 1,
		Parity:   rs232.PARITY_NONE,
		Timeout:  0,
	}

	p, err := rs232.Open(d.SerialPortPath, options)
	if err != nil {
		log.Printf("rs232.Open(): %s\n", err)
		e := err.(*rs232.Error)
		errType := ""
		switch e.Code {
		case rs232.ERR_DEVICE:
			errType = "ERR_DEVICE"
		case rs232.ERR_ACCESS:
			errType = "ERR_ACCESS"
		case rs232.ERR_PARAMS:
			errType = "ERR_PARAMS"
		}
		log.Fatalf("Failed to connect to serial port with error code: %d (%s)\n", e.Code, errType)
	}

	log.Printf("Opened serial port %s\n", p.String())

	// Attach a reference of the serial port to the DoorLock struct.
	d.Serial = *p
}

func (d *DoorLock) Unlock() {
	log.Println("Unlocking door...")
	_, err := d.Serial.Write([]byte("1"))
	if err != nil {
		log.Fatalf("Could not unlock door, with error: %s", err)
	}
	log.Println("Door unlocked!")
}

func (d *DoorLock) Disconnect() {
	log.Println("Disconnecting from door lock...")
	d.Serial.Close()
}

func (d *DoorLock) Listen() {

	// Make sure to connect to the door lock.
	d.Connect()

	// Listen for incoming RFID codes.
	for {
		log.Print("Please input your RFID code for access: ")

		// Check for incoming RFID codes.
		var code string
		_, err := fmt.Scan(&code)
		if err != nil {
			log.Fatal(err)
		}

		// If a code is received, send it to get authenticated.
		err = users.AuthenticateCode(code, d.UsersFile)
		if err != nil {
			log.Println(err.Error())
		} else {
			// Log them in if their code is valid.
			log.Printf("Congrats, your code '%s' is valid, come on in!\n", code)
			d.Unlock()
		}
	}

	// Make sure to disconnect from the door when we're done.
	defer d.Disconnect()
}
