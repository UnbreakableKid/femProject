package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/unbreakablekid/femProject/internal/middleware"
	"github.com/unbreakablekid/femProject/internal/store"
	"github.com/unbreakablekid/femProject/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid workout id",
		})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)

	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "doesn't exist",
		})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decoding create workout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request sent",
		})
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"error": "must be logged in",
		})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "failed to create workout",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid workout update id",
		})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)

	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request",
		})
		return

	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"error": "must be logged in",
		})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
				"error": "workout does not exist",
			})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{
			"error": "you are not authorized to update this",
		})
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		wh.logger.Printf("ERROR: updatingWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internalservererror",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: deleteParsingId: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid workout id",
		})
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"error": "must be logged in",
		})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
				"error": "workout does not exist",
			})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{
			"error": "you are not authorized to update this",
		})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)

	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: deleteWorkoutNotFound: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "workout not found",
		})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server errro",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{})

}
