package main

import "github.com/satori/go.uuid"

// Config contains all the configurations for the sample server
type Config struct {
  // port is the port number where the server will run
	port          string

  // redisPort is the port number where redis is running
	redisPort     string

  // UUIDNamespace is the UUID used to generate v5 uuid.
  // We use v5 uuid to ensure we get the same device id for given device serial code.
  // For more information about how namespaces are used to create v5 uuid,
  // check this link (http://stackoverflow.com/questions/10867405/generating-v5-uuid-what-is-name-and-namespace)
	UUIDNamespace uuid.UUID

  // tokenKey is a string used to generate signed jwt tokens
	tokenKey      string

  // tokenLifetime is the lifetime during which the token is valid
	tokenLifetime int64

  // masterKey is the master token used to bypass session token validation
  // In production you should keep this key secret or disable the feature totally
	masterKey     string
}

// UUIDNamespace used to create v5 uuid
var uuidName, _ = uuid.FromString("27d03927-7c8f-469e-8ba1-68a376d43cc9")

// LocalConfig used for running the server in a local environment
var LocalConfig = Config{
	port:          ":5656",
	redisPort:     ":6379",
	UUIDNamespace: uuidName,
	tokenKey:      "lshjdfsdjhf",
	tokenLifetime: 3600,
	masterKey:     "master-key",
}
