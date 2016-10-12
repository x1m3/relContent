package main

import (
	"time"
	"encoding/json"
	"fmt"
	"os"
	"log"
	"bytes"
)

type UserLike struct {
	userId       string
	contentIds   []string
	lastViewed   time.Time
	consolidated bool
}

func (userLike UserLike)  toJson() string {
var consolidated string

	consolidated="false"
	if  userLike.consolidated ==true {
		consolidated="true"
	}
	jsonContentIds,_ := json.Marshal(userLike.contentIds)
	return fmt.Sprintf("{\"user_id\":\"%s\",\"content_ids\":%s,\"lastViewed\":\"%s\", \"consolidated\":%s}",
		userLike.userId,
		jsonContentIds,
		userLike.lastViewed.String(),
		consolidated,
	)
}

/*
 * TODO: Inject a repo instead of writing directly to a file
 */
func (userLike UserLike) save( consChannel chan string ) {
var err error
var file *os.File
var output bytes.Buffer

	file, err = os.OpenFile("lala.txt", os.O_APPEND | os.O_WRONLY| os.O_EXCL, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.OpenFile("lala.txt", os.O_CREATE| os.O_WRONLY| os.O_EXCL, 0666)
		}
		if err != nil {
			log.Fatal(err)
		}

	}
	output.WriteString(userLike.toJson())
	output.WriteString("\r\n")

	_, err = file.Write(output.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Printf("key[%s] value[%s]\n", userLike.userId, userLike.toJson())
	consChannel <- userLike.userId
}
