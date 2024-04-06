package quest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
	rabbit "github.com/Phund4/testtaskvk_golang/rabbit/RabbitTest"
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

	query := `select *
		from complete_quests
		where client_id = $1 and quest_id = $2`
	selectCompletes, httpStatus, msg, err := helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(httpStatus)
		w.Write([]byte(msg))
		rabbit.SendRabbitMessage(err.Error())
		return
	}

	validationResult, _ := selectCompletes.RowsAffected()
	if validationResult != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Client already completed this quest\n"))
		rabbit.SendRabbitMessage("Error in complete quests (client already completed this quest)")
		return
	}

	query = `update client 
		set balance = (select cost from quest where id = $2) + balance 
		where id = $1`
	_, httpStatus, msg, err = helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(httpStatus)
		w.Write([]byte(msg))
		rabbit.SendRabbitMessage(err.Error())
		return
	}

	query = `insert into complete_quests
		values ($1, $2)`
	result, httpStatus, msg, err := helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		w.WriteHeader(httpStatus)
		w.Write([]byte(msg))
		rabbit.SendRabbitMessage(err.Error())
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in response: %s", err.Error()))
		return
	}

	rabbit.SendRabbitMessage(fmt.Sprintf("Rows affected: %d", num));
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("complete\n"))
}
