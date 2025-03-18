package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/hugo0706/femProject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	
	r.Group(func (r chi.Router){
		r.Use(app.Middleware.Authenticate)
		r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
		
		r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)	
		r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	})
	
	r.Get("/health", app.HealthCheck)

	r.Post("/users", app.UserHandler.HandleRegisterUser)
	
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	return r
}