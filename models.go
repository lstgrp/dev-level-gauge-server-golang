package main

import "errors"

type LevelGaugeData struct {
	DeviceId string `json:"deviceId" binding:"required"`
	Time     int64 `json:"time" binding:"required"`
	Level    uint8  `json:"level" binding:"required"`
}

func (data LevelGaugeData) Validate() error {
	if data.DeviceId == "" {
		return errors.New("Invalid field 'deviceId'")
	}

	if data.Time == 0 {
		return errors.New("Invalid field 'time'")
	}

	if data.Level == 0 {
		return errors.New("Invalid field 'level'")
	}

	return nil
}
