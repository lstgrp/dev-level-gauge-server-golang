package main

import "errors"

type LevelGaugeData struct {
	Time     int64 `json:"time"`
	Event uint8 `json:"event"`
	Level    uint8  `json:"level"`
	DeviceId string `json:"deviceId"`
}

type LevelGaugeRedisData struct {
	Time     int64 `json:"time"`
	Event uint8 `json:"event"`
	Level    uint8  `json:"level"`
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
