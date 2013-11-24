package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var users_file = flag.String("users-file", "users.json", "The users JSON file to use.")
var port = flag.Int("port", 3000, "The port to run on.")

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

	// Check if the file does not exist yet.
	if _, err := os.Stat(*users_file); os.IsNotExist(err) {

		// Let user know we created a user file for them.
		log.Print("No users JSON file, creating one now: ", *users_file)

		// Create file.
		file, err := os.Create(*users_file)
		if err != nil {
			log.Fatal("Error opening/creating users file: ", err.Error())
		}
		defer file.Close()

		// Make the file only readable/writable by the current user.
		err = os.Chmod(*users_file, 0600)
		if err != nil {
			log.Fatal("Could not update permissions to user file: ", err.Error())
		}

		// Add an empty JSON list to the file so it can be parsed by the JSON marshaller.
		_, err = file.WriteString("[]")
		if err != nil {
			log.Fatal("Could not add empty JSON hash to file: ", err.Error())
		}

		// Close the file now that we're done with.
		file.Close()
	}

	// Read file contents into a list of bytes
	outputAsBytes, err := ioutil.ReadFile(*users_file)
	if err != nil {
		log.Fatal("Error reading file contents: ", err.Error())
	}

	// Parse the JSON in the file and return a slice of User structs.
	var output []User
	err = json.Unmarshal(outputAsBytes, &output)
	if err != nil {
		log.Fatal("Error unmarshalling JSON file: ", err.Error())
	}

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
		// TODO: Return a struct converted to JSON.
		fmt.Fprintf(w, `{ "message": "Bad code there Hass!" }`)
		log.Print("Invalid code input: ", code)
		return
	}

	// No code passed in.
	w.WriteHeader(http.StatusNotAcceptable)
	// TODO: Return a struct converted to JSON.
	fmt.Fprintf(w, `{ "message": "Gotta give me a code, bro!" }`)
	log.Print("Missing code input: ", code)
	return
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	flag.Parse()

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
	p := ":" + strconv.Itoa(*port)
	log.Print("Starting server on port", p)
	http.ListenAndServe(p, nil)
}
