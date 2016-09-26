package main

import (
  "github.com/satori/go.uuid"
)

func GetDeviceId (key string) string {
  id := uuid.NewV5(LocalConfig.UUIDNamespace, key)
  return id.String()
}
