package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hugo0706/femProject/internal/middleware"
	"github.com/hugo0706/femProject/internal/store"
	"github.com/hugo0706/femProject/internal/utils"
)

type WorkoutHandler struct{
	workoutStore store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger: logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid_workout_id"})
		return
	}
	
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "bad_request"})
		return
	}
	
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser{
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}
	
	workout.UserID = currentUser.ID
	
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: CreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	
	exisitingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	if exisitingWorkout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "not_found"})
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
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "missing or invalid payload"})
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
	
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser{
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}
	
	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "not authorized"})
		return
	}
	
	err = wh.workoutStore.UpdateWorkout(exisitingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: UpdateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": exisitingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		wh.logger.Printf("ERROR: parseworkoutID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser{
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}
	
	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "not authorized"})
		return
	}
	
	err = wh.workoutStore.DeleteWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: DeleteWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	w.Write(nil)
}