package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/hugo0706/femProject/internal/store"
	"github.com/hugo0706/femProject/internal/utils"
)

type registerUserRequest struct {
	Username 	string `json:"username"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
	Bio 		string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger: logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	
	if req.Email == "" {
		return errors.New("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	
	if req.Password == "" {
		return errors.New("password is required")
	}
	
	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %w", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
	}
	
	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w,http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	
	user := &store.User{
		Username: req.Username,
		Email: req.Email,
	}
	
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	err = h.userStore.CreateUser(user)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}