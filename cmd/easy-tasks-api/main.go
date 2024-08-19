package main

import (
	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/server"
	_ "github.com/JosueMolinaMorales/EasyTasksAPI/pkg/env"
)

func main() {
	server.RunServer()
}
