package main

import "fmt"

// App - the struct which contains things like pointers to db connections
type App struct{}

// Run - sets up our application
func (app *App) Run() error{
	fmt.Println("setting up our app")
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
