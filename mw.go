package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ValidateToken is the middleware for validating session token before calling the handlers
// Session tokens must be included in the 'x-api-jwt' header.
// A master key bypass is implemented for testing purposes.
func ValidateToken(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		// If x-master-key header has master key, bypass validation
		if masterKey := c.Request.Header.Get("x-master-key"); masterKey == LocalConfig.masterKey {
			c.Next()
			return
		}

		// /device and /open don't require a session token
		if c.Request.URL.Path != "/device" && c.Request.URL.Path != "/open" {
			token := c.Request.Header.Get("x-api-jwt")
			res, err := s.Redis.Do("get", token)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "Error in validating token"})
				c.Abort()
				return
			}

			if res == nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "Token timedout or is not valid"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// EnsureJSONBody ensures that all POST requests have content type of JSON
func EnsureJSONBody(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.ContentType() != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "All POST request must have JSON Content-Type"})
			c.Abort()
			return
		}

		c.Next()
	}
}
