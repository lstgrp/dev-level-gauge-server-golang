package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

func GetDeviceId(key string) string {
	id := uuid.NewV5(LocalConfig.UUIDNamespace, key)
	return id.String()
}

func GetTokenString(deviceId string) string {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: LocalConfig.tokenLifetime,
		Subject:   deviceId,
	})
	tokenStr, _ := tokenObj.SignedString([]byte(LocalConfig.tokenKey))

	return tokenStr
}
