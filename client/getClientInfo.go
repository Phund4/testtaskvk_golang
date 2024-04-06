package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
)

type GetClientInfo struct{}

func (*GetClientInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.SendMessages(w, http.StatusInternalServerError,
			"unexpected error\n", fmt.Sprintf("Error in get info client: %s", err.Error()))
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling client: %s", err.Error()))
	} else if !strings.Contains(string(raw), "client_id") {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling client: %s", "incorrect json"))
	}

	query := `select name, balance from client
		where id = $1`
	row, httpStatus, msg, err := helpers.DBQuery(query, client.ID)
	if err != nil {
		helpers.SendMessages(w, httpStatus, msg, err.Error())
		return
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&client.Name, &client.Balance)
	if err != nil {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect client id\n", fmt.Sprintf("Error in response: %s", err.Error()))
		return
	}

	str := fmt.Sprintf("Client ID: %d, Client name: %s, Client Balance: %f\n",
		client.ID, client.Name, client.Balance)
	w.Write([]byte(str))

	query = `select q.name, q.cost 
		from quest as q inner join complete_quests as cq
		on q.id = cq.quest_id
		where cq.client_id = $1`
	rows, httpStatus, msg, err := helpers.DBQuery(query, client.ID)
	if err != nil {
		helpers.SendMessages(w, httpStatus, msg, err.Error())
		return
	}
	defer rows.Close()

	type getInfoStruct struct {
		QuestName string
		Cost      float32
	}
	result := []getInfoStruct{}

	for rows.Next() {
		r := getInfoStruct{}
		err := rows.Scan(&r.QuestName, &r.Cost)
		if err != nil {
			helpers.SendMessages(w, http.StatusInternalServerError,
				"unexpected error\n", fmt.Sprintf("Error in response: %s", err.Error()))
			return
		}

		result = append(result, r)
	}

	for _, el := range result {
		str := fmt.Sprintf("Quest name: %s, Quest cost: %f\n",
			el.QuestName, el.Cost)
		w.Write([]byte(str))
	}

	helpers.SendMessages(w, http.StatusOK, "complete\n", fmt.Sprintf("Rows affected: %d", len(result)))
}
