package main

import (
	"github.com/gin-gonic/gin"
  "net/http"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"
)

func StoreData(s *Server) func (*gin.Context) {
	return func (c *gin.Context) {
		var data LevelGaugeData
		if err := c.BindJSON(&data); err == nil {
      if err := data.Validate(); err == nil {
        jsonData, _ := json.Marshal(LevelGaugeRedisData{
          Time: data.Time,
          Event: data.Event,
          Level: data.Level,
        })
        s.Redis.Do("rpush", data.DeviceId, string(jsonData))
      } else {
        c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body has invalid data"})
      }
		} else {
      c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
    }
	}
}

func RetrieveData(s *Server) func (*gin.Context) {
  return func (c *gin.Context) {

  }
}

func GenerateToken(s *Server) func (*gin.Context) {
  return func (c *gin.Context) {

  }
}

func OpenSession(s * Server) func (*gin.Context) {
  return func (c *gin.Context) {

  }
}
