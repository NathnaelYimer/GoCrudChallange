package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Person struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

var persons []Person

func getAllPersons(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json.NewEncoder(w).Encode(persons)
}

func getPersonByID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	for _, p := range persons {
		if p.ID == params["id"] {
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = uuid.New().String() 
	persons = append(persons, person)
	json.NewEncoder(w).Encode(persons)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	for index, p := range persons {
		if p.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			var updatedPerson Person
			_ = json.NewDecoder(r.Body).Decode(&updatedPerson)
			updatedPerson.ID = params["id"]
			persons = append(persons, updatedPerson)
			json.NewEncoder(w).Encode(updatedPerson)
			return
		}
	}
	json.NewEncoder(w).Encode(persons)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	for index, p := range persons {
		if p.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(persons)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/person", getAllPersons).Methods("GET")
	router.HandleFunc("/person/{id}", getPersonByID).Methods("GET")
	router.HandleFunc("/person", createPerson).Methods("POST")
	router.HandleFunc("/person/{id}", updatePerson).Methods("PUT")
	router.HandleFunc("/person/{id}", deletePerson).Methods("DELETE")

	// Start server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
