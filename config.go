package main

import "github.com/satori/go.uuid"

type Config struct {
  port string
  redisPort string
  UUIDNamespace uuid.UUID
}

var uuidName, _ = uuid.FromString("27d03927-7c8f-469e-8ba1-68a376d43cc9")

var LocalConfig = Config{
  port: ":5656",
  redisPort: ":6379",
  UUIDNamespace: uuidName,
}

