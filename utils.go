package main

import (
	"encoding/json"
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

func LevelGaugeDataFilter(jsonStrData []string, deviceid string, date []int, event int) ([]LevelGaugeRedisData, error) {
	dataSlice := make([]LevelGaugeData, 0)

	for _, data := range jsonStrData {
		parsedData := LevelGaugeData{DeviceId: deviceid}
		if err := json.Unmarshal([]byte(data), &parsedData); err != nil {
			return nil, err
		}

		if date[1] == -1 && event == -1 {
			dataSlice = append(dataSlice, parsedData)
		} else if date[1] == -1 && event != -1 {
			if parsedData.Event == event {
				dataSlice = append(dataSlice, parsedData)
			}
		} else if date[1] != -1 && event == -1 {
			if date[0] < parsedData.Time < date[1] {
				dataSlice = append(dataSlice, parsedData)
			}
		} else {
			if date[0] < parsedData.Time < date[1] && parsedData.Event == event {
				dataSlice = append(dataSlice, parsedData)
			}
		}
	}

	return dataSlice, nil
}
