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

// TODO: Add/Remove/List users (extra!)

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

	// Make sure to return JSON.
	// w.Header().Set("Content-Type", "applicaiton/json")

	// Fetch the list of authenticated users.
	users := GetUsers()

	// Get the code from the request.
	code := r.FormValue("code")

	// See if the code matches our list of authenticated users.
	if code != "" {
		log.Print("Received RFID code: ", code)
		for _, user := range users {
			if code == user.Code {
				w.WriteHeader(http.StatusOK)
				// TODO: Return a struct converted to JSON.
				fmt.Fprintf(w, `{ "message": "Nice one %s!" }`, user.Name)
				return
			}
		}

		// Code didn't match an authenticated user.
		w.WriteHeader(http.StatusUnauthorized)

		fmt.Fprintf(w, `{ "message": "Bad code there Hass!" }`)
		log.Print("Invalid code input: ", code)
		return
	}

	// No code passed in.
	w.WriteHeader(http.StatusNotAcceptable)
	fmt.Fprintf(w, `{ "message": "Gotta give me a code, bro!" }`)
	log.Print("Missing code input: ", code)
	return
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	// Create the webserver
	r := mux.NewRouter()

	// Handle the checkin form.
	r.HandleFunc("/", HomeHandler)

	// Redirect people if they get to checking via GET.
	r.HandleFunc("/checkin", RedirectHandler).Methods("GET")

	// Listen for submissions of the checkin form
	r.HandleFunc("/checkin", CheckinHandler).Methods("POST")

	// Handle the routes.
	http.Handle("/", r)

	// Run the server.
	p := ":3000"
	log.Print("Starting server on port", p)
	http.ListenAndServe(p, nil)
}
