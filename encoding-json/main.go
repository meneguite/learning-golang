package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		fmt.Fprintf(w, "%s %s is %d years old!", user.FirstName, user.LastName, user.Age)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		user := User{
			FirstName: "Ronaldo",
			LastName:  "Meneguite",
			Age:       33,
		}

		json.NewEncoder(w).Encode(user)
	})

	http.ListenAndServe(":8080", nil)
}

// Examples

// Decode
// curl -s -XPOST -d'{"firstname":"Ronaldo","lastname":"Meneguite","age":33}' http://localhost:8080/decode

// Encode
// curl -s http://localhost:8080/encode
