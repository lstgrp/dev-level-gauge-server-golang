// +build !test

package main

func main() {
	server := InitServer(true)
	defer server.Teardown()
	server.Start()
}
