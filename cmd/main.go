package main

import (
	"github.com/polyanskdev/todo-app/pkg/app"
)

// @title Todo App API
// @version 1.0
// @description Api server for Todo Application
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
