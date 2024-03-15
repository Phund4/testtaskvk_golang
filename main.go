package main

import (
	"log"
	"net/http"
	client "github.com/Phund4/testtaskvk_golang/client"
	quest "github.com/Phund4/testtaskvk_golang/quest"
	"github.com/gorilla/mux"
)

func main() {
	log.Print("Hello world sample started.")
	r := mux.NewRouter()

	r.Handle("/addclient", &client.AddClient{}).Methods("POST").Headers("Content-Type", "application/json");
	r.Handle("/addquest", &quest.AddQuest{}).Methods("POST").Headers("Content-Type", "application/json");
	r.Handle("/completequest", &quest.CompleteQuest{}).Methods("POST").Headers("Content-Type", "application/json");
	r.Handle("/getclientinfo", &client.GetClientInfo{}).Methods("GET").Headers("Content-Type", "application/json");

	http.ListenAndServe(":8080", r)
}