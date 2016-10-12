package main

import (
)
import (
	"crypto/md5"
	"encoding/hex"
)

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


