package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hugo0706/femProject/internal/store"
)

type WorkoutHandler struct{
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}
	
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch workout", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}
	
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not create workout", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}
	
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	exisitingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "failed to retrieve workout", http.StatusInternalServerError)
		return
	}
	
	if exisitingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	
	var updateWorkoutRequest struct {
		Title 			*string				 	`json:"title"`
		Description 	*string 				`json:"description"`
		DurationMinutes *int		 			`json:"duration_minutes"`
		CaloriesBurned 	*int 					`json:"calories_burned"`
		Entries 		[]store.WorkoutEntry	`json:"entries"`
	}
	
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if updateWorkoutRequest.Title != nil {
		exisitingWorkout.Title = *updateWorkoutRequest.Title
	}
	
	if updateWorkoutRequest.Description != nil {
		exisitingWorkout.Description = *updateWorkoutRequest.Description
	}
	
	if updateWorkoutRequest.DurationMinutes != nil {
		exisitingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	
	if updateWorkoutRequest.CaloriesBurned != nil {
		exisitingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	
	if updateWorkoutRequest.Entries != nil {
		exisitingWorkout.Entries = updateWorkoutRequest.Entries
	}
	
	err = wh.workoutStore.UpdateWorkout(exisitingWorkout)
	if err != nil {
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exisitingWorkout)
}