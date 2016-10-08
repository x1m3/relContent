package main

import (
	"time"
	"sync"
	"bytes"
	"fmt"
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

func (userLikeList *UserLikeList) consolidate(consChannel chan string) int {
var count int

	count=0
	now :=time.Now()
	fmt.Println("Consolidando")

	userLikeList.Lock()
	defer userLikeList.Unlock()

	for userId, userLike := range userLikeList.list {
		if !userLike.consolidated && now.Sub(userLike.lastViewed).Seconds()>10 {
			tmp := userLikeList.list[userId]
			tmp.consolidated=true
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

func (userLikeList *UserLikeList) consolidateCron(d time.Duration) {
	for {
		time.Sleep(d)
		consolidatedChannel := make(chan string, 20)
		countRemoved:=userLikeList.consolidate(consolidatedChannel)
		userLikeList.removeConsolidated(consolidatedChannel, countRemoved)
		close(consolidatedChannel)
	}
}

func (userLikeList *UserLikeList) installCronThatConsolidates() {
	go userLikeList.consolidateCron( 3 * time.Second )
}
