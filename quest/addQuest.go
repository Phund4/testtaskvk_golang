package quest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
)

type AddQuest struct{}

func (addQuest *AddQuest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.SendMessages(w, http.StatusInternalServerError,
			"unexpected error\n", fmt.Sprintf("Error in add quest: %s", err.Error()))
		return
	}

	var quest Quest
	err = json.Unmarshal(raw, &quest)
	if err != nil {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling quest: %s", err.Error()))
	} else if !strings.Contains(string(raw), "name") || !strings.Contains(string(raw), "cost") {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling quest: %s", "incorrect json"))
	}

	query := `insert into quest (name, cost) 
		values ($1, $2)`
	result, httpStatus, msg, err := helpers.DBExec(query, quest.Name, quest.Cost)
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
