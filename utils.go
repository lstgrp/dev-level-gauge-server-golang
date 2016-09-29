package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

// GetDeviceId creates a unique device id for the given key
// The key here is the device serial, and a v5 uuid is created
func GetDeviceId(key string) string {
	id := uuid.NewV5(LocalConfig.UUIDNamespace, key)
	return id.String()
}

// GetTokenString generates a session token for the given device id
// It uses the HS256 signing method, and although we give it a expiration time
// the source of truth lies in the Redis DB with ttl, so jwt expiration check is not performed
func GetTokenString(deviceId string) string {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: LocalConfig.tokenLifetime,
		Subject:   deviceId,
	})
	tokenStr, _ := tokenObj.SignedString([]byte(LocalConfig.tokenKey))

	return tokenStr
}
