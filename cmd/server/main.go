package main

import (
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/brettman/go-rest-api-course/internal/database"
	transportHTTP "github.com/brettman/go-rest-api-course/internal/transport/http"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// App - the struct which contains info about the app
type App struct{
	Name string
	Version string
}

// Run - sets up our application
func (app *App) Run() error{
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName": app.Name,
			"AppVersion":app.Version,
		}).Info("Setting up application")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()
	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		log.Error("Failed to setup server")
		return err
	}
	return nil
}

func main() {
	app := App{
		Name: "Commenting Service",
		Version: "1.0.0",
	}

	if err := app.Run(); err != nil {
		log.Error("Error starting REST API")
		log.Fatal(err)
	}
}
