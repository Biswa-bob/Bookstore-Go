package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Biswa-bob/bookstore/internal/store"
	"github.com/Biswa-bob/bookstore/internal/tokens"
	"github.com/Biswa-bob/bookstore/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Printf("ERROR: getUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if user == nil {
		h.logger.Printf("ERROR: user not found for email: %s", req.Email)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Match(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: PasswordHash.Matches %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	// Create legacy random token for compatibility
	_, err = h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: Creating Token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Generate access and refresh JWTs
	jwtPair, err := tokens.GenerateJWTPair(user.ID, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		h.logger.Printf("ERROR: Generating JWT pair %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Persist the refresh token hash for revocation/rotation
	refreshToken := &tokens.Token{
		Hash:   tokens.HashToken(jwtPair.RefreshToken),
		UserID: user.ID,
		Expiry: jwtPair.RefreshExpiry,
		Scope:  tokens.ScopeRefresh,
	}
	if err := h.tokenStore.Insert(refreshToken); err != nil {
		h.logger.Printf("ERROR: Inserting refresh token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"access_token":       jwtPair.AccessToken,
		"refresh_token":      jwtPair.RefreshToken,
		"access_expires_at":  jwtPair.AccessExpiry,
		"refresh_expires_at": jwtPair.RefreshExpiry,
		"token_type":         "Bearer",
	})
}
