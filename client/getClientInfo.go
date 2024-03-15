package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	env "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type GetClientInfo struct{}

func (*GetClientInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in get info client: ", err)
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		log.Print("Error in unmarshalling client: ", err)
		return
	}

	err = env.Load("./.env");
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

	query := `select name, balance from client
		where id = $1`;
	row, err := db.Query(query, client.ID);
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in get info client: ", err)
		return
	}
	row.Next()
	err = row.Scan(&client.Name, &client.Balance);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in response: ", err)
		return
	}
	str := fmt.Sprintf("Client ID: %d, Client name: %s, Client Balance: %f\n",
			client.ID, client.Name, client.Balance)
	w.Write([]byte(str))


	query = `select q.name, q.cost 
		from quest as q inner join complete_quests as cq
		on q.id = cq.quest_id
		where cq.client_id = $1`
	rows, err := db.Query(query, client.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in get info client: ", err)
		return
	}
	defer rows.Close()

	type getInfoStruct struct {
		QuestName  string
		Cost       float32
	}

	result := []getInfoStruct{}

	for rows.Next() {
		r := getInfoStruct{}
		err := rows.Scan(&r.QuestName, &r.Cost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected error\n"))
			log.Print("Error in response: ", err)
			return
		}

		result = append(result, r)
	}

	for _, el := range result {
		str := fmt.Sprintf("Quest name: %s, Quest cost: %f\n",
			el.QuestName, el.Cost);
		w.Write([]byte(str));
	}

	log.Print("Rows affected: ", len(result));
	w.Write([]byte("complete\n"));
}