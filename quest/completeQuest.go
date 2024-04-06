package quest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
)

type CompleteQuest struct{}

func (completeQuest *CompleteQuest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.SendMessages(w, http.StatusInternalServerError,
			"unexpected error\n", fmt.Sprintf("Error in complete quest: %s", err.Error()))
		return
	}

	type completeInfo struct {
		ClientID int `json:"client_id"`
		QuestID  int `json:"quest_id"`
	}
	var complete completeInfo

	err = json.Unmarshal(raw, &complete)
	if err != nil {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling complete info: %s", err.Error()))
	} else if !strings.Contains(string(raw), "client_id") || !strings.Contains(string(raw), "quest_id") {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling complete info: %s", "incorrect json"))
	}

	query := `select *
		from complete_quests
		where client_id = $1 and quest_id = $2`
	selectCompletes, httpStatus, msg, err := helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		helpers.SendMessages(w, httpStatus, msg, err.Error())
		return
	}

	validationResult, _ := selectCompletes.RowsAffected()
	if validationResult != 0 {
		helpers.SendMessages(w, http.StatusBadRequest,
			"Client already completed this quest\n", "Error in complete quests (client already completed this quest)")
		return
	}

	query = `update client 
		set balance = (select cost from quest where id = $2) + balance 
		where id = $1`
	_, httpStatus, msg, err = helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		helpers.SendMessages(w, httpStatus, msg, err.Error())
		return
	}

	query = `insert into complete_quests
		values ($1, $2)`
	result, httpStatus, msg, err := helpers.DBExec(query, complete.ClientID, complete.QuestID)
	if err != nil {
		helpers.SendMessages(w, httpStatus, msg, err.Error())
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		helpers.SendMessages(w, http.StatusInternalServerError,
			"unexpected error\n", fmt.Sprintf("Error in response: %s", err.Error()))
		return
	}

	helpers.SendMessages(w, http.StatusOK, "complete\n", fmt.Sprintf("Rows affected: %d", num))
}
