package main

import (
	"encoding/json"
	"testing"
	"time"
)

var testData = make([]string, 0, 10)
var now int64 = time.Now().Unix()

func init() {
	for i := 0; i < 10; i++ {
		data := LevelGaugeRedisData{
			Time:  now + int64(i*10),
			Event: uint8(i % 2),
			Level: uint8(i),
		}

		dataJSONStr, _ := json.Marshal(data)
		testData = append(testData, string(dataJSONStr))
	}
}

func TestLevelGaugeDataFilterNoQuery(t *testing.T) {
	timeFilter := []int64{0, -1}
	filteredData, _ := LevelGaugeDataFilter(testData, "test_deviceid", timeFilter, -1)

	if length := len(filteredData); length != 10 {
		t.Errorf("Filtered data should contain all data, instead got length: %v", length)
	}
}

func TestLevelGaugeDataFilterDate(t *testing.T) {
	timeFilter := []int64{now, now + 41}
	filteredData, _ := LevelGaugeDataFilter(testData, "test_deviceid", timeFilter, -1)

	if length := len(filteredData); length != 5 {
		t.Errorf("Data filtered by date should have 5 data, instead got length: %v", length)
	}

	for _, d := range filteredData {
		if d.Time < timeFilter[0] || d.Time > timeFilter[1] {
			t.Errorf("Filtered data should not have time that is not in range, got: %v", d.Time)
		}
	}
}

func TestLevelGaugeDataFilterEvent(t *testing.T) {
	timeFilter := []int64{0, -1}
	filteredData, _ := LevelGaugeDataFilter(testData, "test_deviceid", timeFilter, 1)

	if length := len(filteredData); length != 5 {
		t.Errorf("Data filtered by event should have 5 data, instead got length: %v", length)
	}

	for _, d := range filteredData {
		if d.Event == 0 {
			t.Errorf("Filtered data should not have event == 0, got: %v", d.Time)
		}
	}
}

func TestLevelGaugeDataFilterFullQuery(t *testing.T) {
	timeFilter := []int64{now, now + 41}
	filteredData, _ := LevelGaugeDataFilter(testData, "test_deviceid", timeFilter, 1)

	if length := len(filteredData); length != 2 {
		t.Errorf("Data filtered by event should have 5 data, instead got length: %v", length)
	}

	for _, d := range filteredData {
		if d.Event == 0 || d.Time < timeFilter[0] || d.Time > timeFilter[1] {
			t.Errorf("Filtered data should not have event == 0 or time not in interval, got time: %v and event: %v",
				d.Time, d.Event)
		}
	}
}
