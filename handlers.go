package main

import (
	"github.com/gin-gonic/gin"
)

func StoreData(s *Server) func (*gin.Context) {
	return func (c *gin.Context) {
		var data LevelGaugeData
		if c.BindJSON(&data) == nil {
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
