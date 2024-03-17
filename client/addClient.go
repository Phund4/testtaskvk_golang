package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type AddClient struct{}

func (addClient *AddClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in add client: ", err)
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil || !strings.Contains(string(raw), "name") || !strings.Contains(string(raw), "balance") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		if err != nil {
			log.Print("Error in unmarshalling client: ", err)
		} else {
			log.Print("Error in unmarshalling client: ", "incorrect json")
		}
		return
	}

	err = env.Load(".env")
	if err != nil {
		err = env.Load("../.env") // чтобы работали тесты
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected error\n"))
			log.Print("Error in load environments: ", err)
		}
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in connection to database: ", err)
		return
	}
	defer db.Close()

	query := `insert into client (name, balance) 
		values ($1, $2)`
	result, err := db.Exec(query,
		client.Name, client.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in insert data client: ", err)
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in response: ", err)
		return
	}

	log.Print("Rows affected: ", num)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("complete\n"))
}
