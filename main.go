package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my-backend/controller"
	"my-backend/middlewares"
	"my-backend/utils"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("USER")
	dbHost := os.Getenv("HOST")
	dbPort := os.Getenv("DB_PORT")
	database := os.Getenv("DATABASE")
	database_password := os.Getenv("DB_PASSWORD")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser,
		database_password,
		dbHost,
		dbPort,
		database,
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Printf("Unable to connect to the database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	posts, err := utils.GetChapters(context.Background(), conn)

	if err != nil {
		fmt.Printf("Something went wrong fetching posts: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(posts)

}

func main() {
	r := mux.NewRouter()
	os.MkdirAll("./uploads", os.ModePerm)

	loggedMux := middlewares.LoggingMiddleware(r)
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/upload", utils.UploadFile).Methods("POST")
	r.HandleFunc("/user/{name}", nameHandler)

	r.HandleFunc("/create-user", createUser).Methods("POST")
	r.HandleFunc("/create-product", createProduct).Methods("POST")
	r.HandleFunc("/register-user", controller.RegisterUser).Methods("POST")

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
