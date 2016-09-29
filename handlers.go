package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StoreData is the handler for /store
// It receives the data to save in JSON body, validates it and saves it to redis
// as JSON string in a list
func StoreData(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data LevelGaugeData

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		if err := data.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body has invalid data"})
			return
		}

		jsonData, _ := json.Marshal(LevelGaugeRedisData{
			Time:  data.Time,
			Event: data.Event,
			Level: data.Level,
		})
		_, err := s.Redis.Do("rpush", data.DeviceId, string(jsonData))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}

// RetrieveData is the handler for /retrieve
// It takes the query object as JSON body and returns a list of data according to the query
// Querying data for all device id is forbidden, and query filters data only within the device id given
func RetrieveData(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data = struct {
			DeviceId string `json:"deviceid"`
		}{
			DeviceId: "",
		}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		dataSlice, err := s.Redis.Do("lrange", data.DeviceId, 0, -1)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		// If no data is stored for the given id, return an empty array
		if dataSlice == nil {
			c.JSON(http.StatusOK, gin.H{"result": "[]"})
			return
		}

		// Convert slice of JSON strings from redis to slice of LevelGaugeData
		assertedDataSlice := dataSlice.([]interface{})
		finalData := make([]LevelGaugeData, 0)
		for _, d := range assertedDataSlice {
			var redisData LevelGaugeRedisData
			assertedD := d.([]byte)
			json.Unmarshal(assertedD, &redisData)
			finalData = append(finalData, LevelGaugeData{
				DeviceId: data.DeviceId,
				Time:     redisData.Time,
				Event:    redisData.Event,
				Level:    redisData.Level,
			})
		}

		finalDataJSONBytes, err := json.Marshal(finalData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": string(finalDataJSONBytes)})
	}
}

// GenerateToken is the handler for /device
// It takes the device serial to generate a device id, and it returns
// the device id, session token and token lifetime
func GenerateToken(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data TokenParameter

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		if err := data.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body has invalid data"})
			return
		}

		deviceId := GetDeviceId(data.Device.Serial)
		tokenStr := GetTokenString(deviceId)

		if _, err := s.Redis.Do("SETEX", tokenStr, LocalConfig.tokenLifetime, tokenStr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"deviceid": deviceId,
			"token":    tokenStr,
			"ttl":      LocalConfig.tokenLifetime,
		})
	}
}

// OpenSession is the handler for /open
// It takes the device id given and generates a new session token
func OpenSession(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data = struct {
			DeviceId string `json:"deviceid" binding:"required"`
		}{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		deviceId := data.DeviceId
		tokenStr := GetTokenString(deviceId)

		if _, err := s.Redis.Do("SETEX", tokenStr, LocalConfig.tokenLifetime, tokenStr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenStr,
			"ttl":   LocalConfig.tokenLifetime,
		})
	}
}

// CloseSession is the handler for /close
// It takes the session token and erases it from redis, closing the session
func CloseSession(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data = struct {
			Token string `json:"token" binding:"required"`
		}{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		if _, redisErr := s.Redis.Do("DEL", data.Token); redisErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}
