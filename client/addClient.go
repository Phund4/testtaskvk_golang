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

	query := `insert into client (name, balance) 
		values ($1, $2)`
	result, httpStatus, msg, err := helpers.DBExec(query, client.Name, client.Balance);
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

	rabbit.SendRabbitMessage(fmt.Sprintf("Rows affected: %v", num))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("complete\n"))
}