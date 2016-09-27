package main

import "errors"

type LevelGaugeData struct {
	Time     int64  `json:"time" binding:"required"`
	Event    uint8  `json:"event"`
	Level    uint8  `json:"level"`
	DeviceId string `json:"deviceid" binding:"required"`
}

type LevelGaugeRedisData struct {
	Time  int64 `json:"time" binding:"required"`
	Event uint8 `json:"event"`
	Level uint8 `json:"level"`
}

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

type LevelGaugeDataQuery struct {
	DeviceId string  `json:"deviceid" binding:"required"`
	Date     []int64 `json:"date" binding:"required"`
	Event    int     `json:"event"`
}

type TokenParameter struct {
	Device struct {
		Name   string `json:"name" binding:"required"`
		Serial string `json:"serial" binding:"required"`
	} `json:"device"`
	App  struct{} `json:"app"`
	User struct{} `json:"user"`
}

func (data TokenParameter) Validate() error {
	if data.Device.Name == "" {
		return errors.New("Invalid field `Device.Name`")
	}

	if data.Device.Serial == "" {
		return errors.New("Invalid field `Device.Serial`")
	}

	return nil
}
