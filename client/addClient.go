package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	env "github.com/joho/godotenv"
	rabbit "github.com/Phund4/testtaskvk_golang/rabbit/RabbitTest"
	_ "github.com/lib/pq"
)

type AddClient struct{}

func (addClient *AddClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in add client: %s", err.Error()))
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil || !strings.Contains(string(raw), "name") || !strings.Contains(string(raw), "balance") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		if err != nil {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling client: %s", err.Error()))
		} else {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling client: %s", "incorrect json"))
		}
		return
	}

	err = env.Load(".env")
	if err != nil {
		err = env.Load("../.env") // чтобы работали тесты
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected error\n"))
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in load environments: %s", err.Error()))
		}
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in connection to database: %s", err.Error()))
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
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in insert data client: %s", err.Error()))
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in response: %s", err.Error()))
		return
	}

	rabbit.SendRabbitMessage(fmt.Sprintf("Rows affected: %v", num))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("complete\n"))
}
