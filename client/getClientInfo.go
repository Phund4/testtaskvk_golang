package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
	rabbit "github.com/Phund4/testtaskvk_golang/rabbit/RabbitTest"
)

type GetClientInfo struct{}

func (*GetClientInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unexpected error\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in get info client: %s", err.Error()))
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil || !strings.Contains(string(raw), "client_id") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect data\n"))
		if err != nil {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling client: %s", err.Error()))
		} else {
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in unmarshalling client: %s", "incorrect json"))
		}
		return
	}

	query := `select name, balance from client
		where id = $1`
	row, httpStatus, msg, err := helpers.DBQuery(query, client.ID)
	if err != nil {
		w.WriteHeader(httpStatus)
		w.Write([]byte(msg))
		rabbit.SendRabbitMessage(err.Error())
		return
	}
	defer row.Close()

	row.Next();
	err = row.Scan(&client.Name, &client.Balance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect client id\n"))
		rabbit.SendRabbitMessage(fmt.Sprintf("Error in response: %s", err.Error()))
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
		w.WriteHeader(httpStatus)
		w.Write([]byte(msg))
		rabbit.SendRabbitMessage(err.Error())
		return
	}
	defer rows.Close();

	type getInfoStruct struct {
		QuestName string
		Cost      float32
	}

	result := []getInfoStruct{}

	for rows.Next() {
		r := getInfoStruct{}
		err := rows.Scan(&r.QuestName, &r.Cost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected error\n"))
			rabbit.SendRabbitMessage(fmt.Sprintf("Error in response: %s", err.Error()))
			return
		}

		result = append(result, r)
	}

	for _, el := range result {
		str := fmt.Sprintf("Quest name: %s, Quest cost: %f\n",
			el.QuestName, el.Cost)
		w.Write([]byte(str))
	}

	rabbit.SendRabbitMessage(fmt.Sprintf("Rows affected: %d", len(result)))
	w.WriteHeader(http.StatusOK);
	w.Write([]byte("complete\n"))
}