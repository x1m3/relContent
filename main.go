package main

//go build -x github.com/x1m3/relContent

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"flag"
)



var userLikes UserLikeList
var graph *Graph
var config Config

func main() {

	configFile := flag.String("config", "config.conf","Path to the config file")
	flag.Parse()

	log.Print("Loading config file")
	ok:=config.load(*configFile)
	if (!ok) {
		log.Fatal("Error reading config file");
	}
	log.Println("[DONE]")

	userLikes.init()
	graph = newGraph()

	log.Print("Loading Database")
	graph.loadData(&config)
	log.Println("[DONE]")

	router := mux.NewRouter()
	router.HandleFunc("/like", registerLikeAction).Methods("POST");
	router.HandleFunc("/like", viewLikeAction).Methods("GET");
	router.HandleFunc("/related/{id}", viewRelatedAction).Methods("GET");

	userLikes.installCronThatConsolidates(&config)

	log.Println("Listening for connections.")
	log.Fatal(http.ListenAndServe(":8000", router))
}
