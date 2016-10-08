package main

//go build -x github.com/x1m3/relContent

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"crypto/md5"
	"encoding/hex"
)

type Response struct {
	Status string
	Output string
	Message string
}

type Like struct {
	userId string
	contentId string
}

func newLike(userId string, contentId string, nonce string, time string) (*Like,bool) {
const secret string="ELPERRODESANROQUENOTIENERABO"

	hasher := md5.New()
	hasher.Write([]byte(secret+time+contentId))
	newNonce := hex.EncodeToString(hasher.Sum(nil))

	if newNonce!=nonce {
		return nil, false
	}

	like := new(Like)
	like.contentId = contentId
	like.userId = userId
	return like, true;
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


var userLikes UserLikeList

func main() {
	userLikes.init()

	router := mux.NewRouter()
	router.HandleFunc("/like", registerLikeAction).Methods("POST");
	router.HandleFunc("/like", viewLikeAction).Methods("GET");

	userLikes.installCronThatConsolidates()

	log.Fatal(http.ListenAndServe(":8000", router))
}
