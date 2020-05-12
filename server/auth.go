package server

import (
	"crypto/md5"
	"encoding/hex"
)

type authFunc interface {
	Validate(token string, url string) bool
}

type accessTokenValidator struct {
	AccessToken string
}

func (v accessTokenValidator) Validate(accessToken string, url string) bool {
	hash := md5.New()
	hash.Write([]byte(url + v.AccessToken))
	return accessToken == hex.EncodeToString(hash.Sum(nil))
}
