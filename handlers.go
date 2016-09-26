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
        _, err := s.Redis.Do("rpush", data.DeviceId, string(jsonData))

        if err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
        }

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
    var data TokenParameter
    if err := c.BindJSON(&data); err == nil {
      deviceId := GetDeviceId(data.Device.Serial)
      token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
        ExpiresAt: LocalConfig.tokenLifetime,
        Subject: deviceId,
      })
      tokenStr, _ := token.SignedString(LocalConfig.tokenKey)
      _, err := s.Redis.Do("set", tokenStr, LocalConfig.tokenLifetime, "")

      if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal error"})
      }

      c.JSON(http.StatusOK, gin.H{
        "deviceId": deviceId,
        "token": tokenStr,
        "ttl": LocalConfig.tokenLifetime,
      })
    } else {
      c.JSON(http.StatusBadRequest, gin.H{"status": "JSON Body is missing fields"})
    }
  }
}

func OpenSession(s * Server) func (*gin.Context) {
  return func (c *gin.Context) {

  }
}
