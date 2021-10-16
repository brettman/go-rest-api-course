package main

import (

	"fmt"
	transportHTTP "github.com/brettman/go-rest-api-course/internal/transport/http"
	"net/http"
)

// App - the struct which contains things like pointers to db connections
type App struct{}

// Run - sets up our application
func (app *App) Run() error{
	fmt.Println("setting up our app")
	handler := transportHTTP.NewHandler()
	handler.SetupRoutes()
	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to setup server")
		return err
	}
	return nil
}

func main() {
	fmt.Println("Go REST API course")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting REST API")
		fmt.Println(err)
	}
}
