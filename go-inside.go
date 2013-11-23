package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tarm/goserial"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
)

// Open the users.json file
// + Add/Remove/List users (extra!)

type User struct {
	Name, Code string
}

// type apiResponse struct {
// 	Users []User
// }

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func GetUsers() []User {

	// Read file contents into a list of bytes
	outputAsBytes, err := ioutil.ReadFile("users.json")

	// Parse the JSON in the file and return a slice of User structs.
	var output []User
	err = json.Unmarshal(outputAsBytes, &output)
	if err != nil {
		log.Fatal(err)
	}

	// Available codes to query against.
	return output
}

func CheckinHandler(w http.ResponseWriter, r *http.Request) {

	// Create a new serial configuration.
	c := &serial.Config{Name: "/dev/tty.usbmodem", Baud: 9600}

	// Open a connection based on the serial config.
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// Write to the serial conneciton.
	n, err := s.Write([]byte("1"))
	if err != nil {
		log.Fatal(err)
	}

	users := GetUsers()

	// Get the code from the request.
	code := r.FormValue("code")

	// See if the code matches our list of authenticated users.
	if code != "" {
		log.Print("Received RFID code: ", code)
		for _, user := range users {
			if code == user.Code {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Nice one %s!", user.Name)
				return
			}
		}

		// Code didn't match an authenticated user.
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Bad code there Hass!")
		log.Print("Invalid code input: ", code)
		return
	}

	// No code passed in.
	w.WriteHeader(http.StatusNotAcceptable)
	fmt.Fprintf(w, "Gotta give me a code bro!")
	log.Print("No code input: ", code)
	return
}

func main() {

	// Create the webserver
	r := mux.NewRouter()

	// Handle the checkin form.
	r.HandleFunc("/", HomeHandler)

	// Listen for submissions of the checkin form
	r.HandleFunc("/checkin", CheckinHandler).Methods("POST")

	// Handle the routes.
	http.Handle("/", r)

	// Run the server.
	port := ":3000"
	log.Print("Starting server on port", port)
	http.ListenAndServe(port, nil)
}
