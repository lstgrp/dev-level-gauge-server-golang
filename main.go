package main

func main() {
	server := InitServer()
	defer server.Teardown()
	server.Start()
}
