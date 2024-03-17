package quest

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

type AddQuest struct{}

func (addQuest *AddQuest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in add quest: ", err)
		return
	}

	var quest Quest
	err = json.Unmarshal(raw, &quest)
	if err != nil || !strings.Contains(string(raw), "name") || !strings.Contains(string(raw), "cost") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		if err != nil {
			log.Print("Error in unmarshalling quest: ", err)
		} else {
			log.Print("Error in unmarshalling quest: ", "incorrect json")
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in load environments: ", err)
		return
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

	query := `insert into quest (name, cost) 
		values ($1, $2)`
	result, err := db.Exec(query,
		quest.Name, quest.Cost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in insert data quest: ", err)
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
