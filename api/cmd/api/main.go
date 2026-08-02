package main

import "github.com/RoyerHernandez/pos-system/api/cmd/api/server"

func main() {
	srv := server.New()
	srv.Run()
}
