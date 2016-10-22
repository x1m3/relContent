package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
)

type Response struct {
	Status string
	Output string
	Message string
}


func registerLikeAction( w http.ResponseWriter, request *http.Request) {
	var response Response
	likeRequest, ok := newLike(
		request.FormValue("user_id"),
		request.FormValue("content_id"),
		request.FormValue("nonce"),
		request.FormValue("time"))
	if ok  {
		userLikes.annotateLike( *likeRequest )
		response = Response{Status:"OK", Output:"", Message:"Enqueued" }
	}else {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{Status:"ERROR", Output:"", Message:"Invalid nonce" }
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func viewLikeAction( w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json := userLikes.toJson();
	w.Write([]byte(json))
}

func viewRelatedAction( w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	nodes := graph.relatedNodes(id)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))

	if len(nodes)>0 {
		i := 0
		for _, node := range nodes {
			w.Write([]byte(node.toJson()))
			if (i!=len(nodes)-1) {
				w.Write([]byte(","))
			}
			i++
		}
	}
	w.Write([]byte("]"))
	log.Printf("%s %s", request.Method, request.RequestURI)
}
