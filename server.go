package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	Router *gin.Engine
	Redis  redis.Conn
}

func InitServer() *Server {
	server := Server{}

	// Make Redis connection
	redisConn, err := redis.Dial("tcp", LocalConfig.redisPort)

	if err != nil {
		log.Fatalln(err)
	}

	server.Redis = redisConn

	// Register handlers
	server.Router = gin.New()
	server.Router.POST("/store", StoreData(&server))
	server.Router.POST("/device", GenerateToken(&server))

	return &server
}

func (s *Server) Start() {
	s.Router.Run(LocalConfig.port)
}

func (s *Server) Teardown() {
	if s.Redis != nil {
		s.Redis.Close()
	}
}
