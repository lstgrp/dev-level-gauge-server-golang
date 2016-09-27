package main

import (
	"encoding/json"
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

// LevelGaugeDataFilter takes a slice of JSON strings given from the Redis server
// and filters them according to the given query parameters.
func LevelGaugeDataFilter(jsonStrData []string, deviceid string, date []int64, event int) ([]LevelGaugeData, error) {
	dataSlice := make([]LevelGaugeData, 0)

	for _, data := range jsonStrData {
		parsedData := LevelGaugeData{DeviceId: deviceid}
		if err := json.Unmarshal([]byte(data), &parsedData); err != nil {
			return nil, err
		}

		if date[1] == -1 && event == -1 {
			dataSlice = append(dataSlice, parsedData)
		} else if date[1] == -1 && event != -1 {
			if int(parsedData.Event) == event {
				dataSlice = append(dataSlice, parsedData)
			}
		} else if date[1] != -1 && event == -1 {
			if date[0] <= parsedData.Time && parsedData.Time <= date[1] {
				dataSlice = append(dataSlice, parsedData)
			}
		} else {
			if date[0] <= parsedData.Time && parsedData.Time <= date[1] && int(parsedData.Event) == event {
				dataSlice = append(dataSlice, parsedData)
			}
		}
	}

	return dataSlice, nil
}
