package main

import "fmt"

var CustomHandlers = struct {
	Handlers map[string]func(data LevelGaugeData) error
}{
	make(map[string]func(data LevelGaugeData) error),
}

func AddCustomHandler(name string, handler func(data LevelGaugeData) error) {
	CustomHandlers.Handlers[name] = handler
}

func RemoveCustomHandler(name string) {
	delete(CustomHandlers.Handlers, name)
}

func ExecuteAllHandlers(data LevelGaugeData) {
	for name, handler := range CustomHandlers.Handlers {
		if err := handler(data); err != nil {
			fmt.Printf("Error while executing handler %v\nError: %v\n", name, err)
		}
	}
}
