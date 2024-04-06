package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	helpers "github.com/Phund4/testtaskvk_golang/helpers"
)

type AddClient struct{}

func (addClient *AddClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.SendMessages(w, http.StatusInternalServerError,
			"unexpected error\n", fmt.Sprintf("Error in add client: %s", err.Error()))
		return
	}

	var client Client
	err = json.Unmarshal(raw, &client)
	if err != nil {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling client: %s", err.Error()))
	} else if !strings.Contains(string(raw), "name") || !strings.Contains(string(raw), "balance") {
		helpers.SendMessages(w, http.StatusBadRequest,
			"incorrect data\n", fmt.Sprintf("Error in unmarshalling client: %s", "incorrect json"))
	}

	query := `insert into client (name, balance) 
		values ($1, $2)`
	result, httpStatus, msg, err := helpers.DBExec(query, client.Name, client.Balance)
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

	helpers.SendMessages(w, http.StatusOK, "complete\n", fmt.Sprintf("Rows affected: %v", num))
}