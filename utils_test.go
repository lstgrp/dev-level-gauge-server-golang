package main

import (
	"sync"
	"testing"
)

func TestAddCustomHandler(t *testing.T) {
	AddCustomHandler("test_handler", func(data LevelGaugeData) error {
		return nil
	})

	if length := len(CustomHandlers.Handlers); length != 1 {
		t.Errorf("Expected handler map length to be %v, instead got %v", 1, length)
	}
}

func TestExecuteAllHandlers(t *testing.T) {
	testFlag := false

	var wg sync.WaitGroup

	AddCustomHandler("test_handler", func(data LevelGaugeData) error {
		testFlag = true
		return nil
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		ExecuteAllHandlers(LevelGaugeData{})
	}()

	wg.Wait()

	if !testFlag {
		t.Error("Expected handler to be called and change test flag")
	}
}
