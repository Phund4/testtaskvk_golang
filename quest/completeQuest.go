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

	rabbit "github.com/Phund4/testtaskvk_golang/rabbit/RabbitTest"
	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type CompleteQuest struct{}

func (completeQuest *CompleteQuest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in complete quest: %s", err.Error()))
		return
	}

	type completeInfo struct {
		ClientID int `json:"client_id"`
		QuestID  int `json:"quest_id"`
	}

	var complete completeInfo

	err = json.Unmarshal(raw, &complete)
	if err != nil || !strings.Contains(string(raw), "client_id") || !strings.Contains(string(raw), "quest_id") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		if err != nil {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling complete info: %s", err.Error()))
		} else {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling complete info: %s", "incorrect json"))
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

	query := `select *
		from complete_quests
		where client_id = $1 and quest_id = $2`
	selectCompletes, err := db.Exec(query,
		complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in validation data complete_quests: %s", err.Error()))
		return
	}

	validationResult, _ := selectCompletes.RowsAffected()
	if validationResult != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Client already completed this quest\n"))
		rabbit.SendRabbitMessage("Error in complete quests (client already completed this quest)")
		return
	}

	query = `insert into complete_quests (client_id, quest_id) 
		values ($1, $2)`
	result, err := db.Exec(query,
		complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in insert data complete_quests: %s", err.Error()))
		return
	}

	query = `update client 
		set balance = (select cost from quest where id = $2) + balance 
		where id = $1`
	_, err = db.Exec("update client set balance = (select cost from quest where id = $2) + balance where id = $1",
		complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect values\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in replenish balance complete_quests: %s", err.Error()))
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in response: %s", err.Error()))
		return
	}

	log.Print("Rows affected: ", num)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("complete\n"))
}