package quest

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

type CompleteQuest struct{}

func (completeQuest *CompleteQuest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		log.Print("Error in complete quest: ", err)
		return
	}

	type completeInfo struct {
		ClientID int `json:"client_id"`
		QuestID int `json:"quest_id"`
	}

	var complete completeInfo;

	err = json.Unmarshal(raw, &complete)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		log.Print("Error in unmarshalling complete info: ", err)
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

	query := `select *
		from complete_quests
		where client_id = $1 and quest_id = $2`;
	selectCompletes, err := db.Exec(query, 
		complete.ClientID, complete.QuestID);
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in validation data complete_quests: ", err)
		return
	}
	validationResult, _ := selectCompletes.RowsAffected();
	if validationResult != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Client already completed this quest\n"))
		log.Print("Error in complete_quests (client already completed): ", err)
		return
	}

	query = `insert into complete_quests (client_id, quest_id) 
		values ($1, $2)`;
	result, err := db.Exec(query,
		complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in insert data complete_quests: ", err)
		return
	}

	query = `update client 
		set balance = (select cost from quest where id = $2) + balance 
		where id = $1`;
	_, err = db.Exec("update client set balance = (select cost from quest where id = $2) + balance where id = $1",
		complete.ClientID, complete.QuestID);
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		log.Print("Error in replenish balance complete_quests: ", err)
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