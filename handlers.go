package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StoreData(s *Server) func(*gin.Context) {
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

func RetrieveData(s *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		var data LevelGaugeDataQuery

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
			return
		}

		dataSlice, err := s.Redis.Do("lrange", data.DeviceId, 0, -1)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		if dataSlice == nil {
			c.JSON(http.StatusOK, gin.H{"result": "[]"})
			return
		}

		dataByteSlice := dataSlice.([]interface{})
		finalData := make([]string, 0)
		for _, d := range dataByteSlice {
			dByteSlice := d.([]byte)
			finalData = append(finalData, string(dByteSlice))
		}

		filteredData, err := LevelGaugeDataFilter(finalData, data.DeviceId, data.Date, data.Event)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		filteredDataJSONString, err := json.Marshal(filteredData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": string(filteredDataJSONString)})
	}
}

func GenerateToken(s *Server) func(*gin.Context) {
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

func OpenSession(s *Server) func(*gin.Context) {
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

func CloseSession(s *Server) func(*gin.Context) {
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
