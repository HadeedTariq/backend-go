package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
func nameHandler(w http.ResponseWriter, r *http.Request) {
	age := r.URL.Query().Get("age")
	vars := mux.Vars(r)
	name := vars["name"]
	fmt.Fprintf(w, "User name: %s and age is %s", name, age)
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)

}

type Product struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbUser := os.Getenv("USER")
	dbHost := os.Getenv("HOST")
	dbPort := os.Getenv("DB_PORT")
	database := os.Getenv("DATABASE")

	conn, err := pgx.Connect(context.Background(), "postgres://user:hadeed#12896@localhost:5432/dbname")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
		return
	}

	fmt.Println(greeting)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(product)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/user/{name}", nameHandler)

	r.HandleFunc("/create-user", createUser).Methods("POST")
	r.HandleFunc("/create-product", createProduct).Methods("POST")

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
