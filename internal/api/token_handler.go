package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hugo0706/femProject/internal/store"
	"github.com/hugo0706/femProject/internal/tokens"
	"github.com/hugo0706/femProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore store.UserStore
	logger *log.Logger
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore: userStore,
		logger: logger,
	}
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (t *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request params"})
		return
	}

	user, err := t.userStore.GetUserByUsername(req.Username)
	if err != nil {
		t.logger.Printf("ERROR: GetUserByUsername %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	if user == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}
	
	passwordsMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		t.logger.Printf("ERROR: PasswordHash.Matches %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	if !passwordsMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}
	
	token, err := t.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		t.logger.Printf("ERROR: CreateNewToken %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal_server_error"})
		return
	}
	
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}