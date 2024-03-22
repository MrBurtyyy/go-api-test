package main

import (
	"github.com/MrBurtyyy/go-api-test/internal/server"
)

func main() {
	server.Init()
	server.ListenAndServe(":3000")
}
