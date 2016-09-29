package main

import "errors"

// LevelGaugeData is the data model for level gauge data from client
// JSON body in POST request must follow this model and it will be validated
type LevelGaugeData struct {
	Time     int64  `json:"time" binding:"required"`
	Event    uint8  `json:"event"`
	Level    uint8  `json:"level"`
	DeviceId string `json:"deviceid" binding:"required"`
}

// LevelGaugeRedisData is the data model that will be saved in Redis as JSON string
type LevelGaugeRedisData struct {
	Time  int64 `json:"time" binding:"required"`
	Event uint8 `json:"event"`
	Level uint8 `json:"level"`
}

// Validate function validates the JSON body given in /store POST requests
func (data LevelGaugeData) Validate() error {
	if data.DeviceId == "" {
		return errors.New("Invalid field 'deviceid'")
	}

	if data.Time == 0 {
		return errors.New("Invalid field 'time'")
	}

	if data.Level == 0 {
		return errors.New("Invalid field 'level'")
	}

	return nil
}

// TokenParameter is the data model for generating a session token
// Currently only the device field is necessary
type TokenParameter struct {
	Device struct {
		Name   string `json:"name" binding:"required"`
		Serial string `json:"serial" binding:"required"`
	} `json:"device"`
	App  struct{} `json:"app"`
	User struct{} `json:"user"`
}

// Validate validates the given token parameters for valid data
func (data TokenParameter) Validate() error {
	if data.Device.Name == "" {
		return errors.New("Invalid field `Device.Name`")
	}

	if data.Device.Serial == "" {
		return errors.New("Invalid field `Device.Serial`")
	}

	return nil
}
