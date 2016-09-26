package main

type Config struct {
  port string
  redisPort string
}

var LocalConfig = Config{
  port: ":5656",
  redisPort: ":6379",
}

