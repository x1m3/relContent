package main

import (
	"time"
	"sync"
	"bytes"
	"log"
)

type UserLikeList struct {
	sync.Mutex
	list map[string]UserLike
}

func (userLikeList *UserLikeList) init() {
	userLikeList.list = make(map[string]UserLike)
}

func (userLikeList *UserLikeList) toJson() string {
	var output bytes.Buffer

	userLikeList.Lock()
	defer userLikeList.Unlock()

	i:=1
	output.WriteString("{\"userIds\":[")
	for _,userLike := range userLikeList.list {
		output.WriteString(userLike.toJson())
		if i<len(userLikeList.list) {
			output.WriteString(",")
		}
		i++
	}
	output.WriteString("]}")
	return output.String()
}



func (userLikeList *UserLikeList) annotateLike( likeRequest Like) {

	userLikeList.Lock()
	defer userLikeList.Unlock()

	entry,userExists := userLikeList.list[likeRequest.userId]
	if !userExists {
		entry = UserLike{userId: likeRequest.userId, contentIds:make([]string, 0) , lastViewed: time.Now(), consolidated:false}
	}
	// Evitamos entradas duplicadas
	duplicatedContentId:=false
	for _, contentId := range entry.contentIds {
		if contentId == likeRequest.contentId {
			duplicatedContentId=true
		}
	}
	if (!duplicatedContentId) {
		newContentIds := make([]string, len(entry.contentIds) + 1)
		copy(newContentIds, entry.contentIds)
		newContentIds[len(entry.contentIds)] = likeRequest.contentId
		entry.contentIds = newContentIds
		userLikeList.list[likeRequest.userId] = entry
	}
}

func (userLikeList *UserLikeList) consolidate(consChannel chan string, config *Config) int {
var count int

	count=0
	now :=time.Now()
	log.Println("Consolidating obsolete sessions")

	userLikeList.Lock()
	defer userLikeList.Unlock()

	for userId, userLike := range userLikeList.list {
		if !userLike.consolidated &&
		len(userLike.contentIds)>=config.Runtime.MinSessionSize &&
		int(now.Sub(userLike.lastViewed).Seconds())>config.Runtime.SessionClosedAfterSeconds {

			graph.annotateAllRelationsBidirectional(&userLike);

			tmp := userLikeList.list[userId]
			tmp.consolidated =true
			userLikeList.list[userId]=tmp

			go userLikeList.list[userId].save( consChannel )
			count++
		}
	}
	return count
}

func (userLikeList *UserLikeList) removeConsolidated( toRemoveChannel chan string, count int) {

	for i:=0; i<count; i++ {
		userId := <-toRemoveChannel
		userLikeList.Lock()
		delete(userLikeList.list, userId)
		userLikeList.Unlock()
	}
}

func (userLikeList *UserLikeList) consolidateCron(config *Config) {
	for {
		time.Sleep(time.Second *  time.Duration(config.Runtime.RunConsolidatorEverySeconds))
		consolidatedChannel := make(chan string, 20)
		countRemoved:=userLikeList.consolidate(consolidatedChannel, config)
		userLikeList.removeConsolidated(consolidatedChannel, countRemoved)
		close(consolidatedChannel)
	}
}

func (userLikeList *UserLikeList) installCronThatConsolidates(config *Config) {
	go userLikeList.consolidateCron( config)
}
