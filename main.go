package main

import (
	"link/server"
	"link/storage"
)

func main() {
	db := storage.InitDatabase("localhost")
	defer db.Session.Close()

	server.CreateLinkServer(db, 8888)
}
