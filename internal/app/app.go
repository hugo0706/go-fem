package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hugo0706/femProject/internal/api"
	"github.com/hugo0706/femProject/internal/middleware"
	"github.com/hugo0706/femProject/internal/migrations"
	"github.com/hugo0706/femProject/internal/store"
)

type Application struct {
	Logger		   	*log.Logger
	WorkoutHandler 	*api.WorkoutHandler
	UserHandler 	*api.UserHandler
	TokenHandler	*api.TokenHandler
	Middleware		middleware.UserMiddleware
	DB				*sql.DB	   
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userStore := store.NewPostgresUserStore(pgDB)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger: 		logger,
		WorkoutHandler: workoutHandler,
		UserHandler: 	userHandler,
		TokenHandler: 	tokenHandler,
		Middleware: 	middlewareHandler,
		DB:				pgDB,
	}
	
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}