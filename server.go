package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
)

// Server struct contains different instances and fields necessary for all handlers
// It registers the routes for handlers and has the Redis client instance
type Server struct {
	Router *gin.Engine
	Redis  redis.Conn
}

// InitServer creates a server instance with routes registered and Redis connected
// If 'useMiddleware' is false, it will not validate session tokens and headers
func InitServer(useMiddleware bool) *Server {
	server := Server{}

	// Make Redis connection
	redisConn, err := redis.Dial("tcp", LocalConfig.redisPort)

	if err != nil {
		log.Fatalln(err)
	}

	server.Redis = redisConn

	// Register handlers
	gin.SetMode(gin.ReleaseMode)
	server.Router = gin.New()

	if useMiddleware {
		server.Router.Use(ValidateToken(&server), EnsureJSONBody(&server))
	}

	server.Router.POST("/device", GenerateToken(&server))
	server.Router.POST("/close", CloseSession(&server))
	server.Router.POST("/open", OpenSession(&server))

	server.Router.POST("/store", StoreData(&server))
	server.Router.POST("/retrieve", RetrieveData(&server))

	return &server
}

// Start starts the server to listen on the configured port
func (s *Server) Start() {
	s.Router.Run(LocalConfig.port)
}

// Teardown is called when main function ends and performs necessary teardown operations
func (s *Server) Teardown() {
	if s.Redis != nil {
		s.Redis.Close()
	}
}
