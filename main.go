package main

//go build -x github.com/x1m3/relContent

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"flag"
	"os"
	"fmt"
	"strconv"
)



var userLikes UserLikeList
var graph *Graph
var config Config

func main() {

	configFile := flag.String("config", "config.conf","Path to the config file")
	flag.Parse()
	ok:=config.load(*configFile)
	if (!ok) {
		log.Fatal("Error reading config file");
	}

	userLikes.init()
	graph = newGraph()
	graph.loadData(&config)

	for i:=0; i<=3804; i++ {
		nodes := graph.relatedNodes(strconv.Itoa(i))
		if len(nodes)>0 {
			fmt.Printf("%d=>{", i)
			for _, node := range nodes {
				fmt.Printf(" %s(%d) ", node.edgeId, node.weight)
			}
			fmt.Printf("}\r\n")
		}
	}

	os.Exit(0)

	router := mux.NewRouter()
	router.HandleFunc("/like", registerLikeAction).Methods("POST");
	router.HandleFunc("/like", viewLikeAction).Methods("GET");

	userLikes.installCronThatConsolidates(&config)

	log.Fatal(http.ListenAndServe(":8000", router))
}
