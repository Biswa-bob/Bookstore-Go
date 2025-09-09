package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/Biswa-bob/bookstore/internal/middleware"
	"github.com/Biswa-bob/bookstore/internal/store"
	"github.com/Biswa-bob/bookstore/internal/utils"
)

type registrationUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userstore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userstore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registrationUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if len(req.Name) > 50 {
		return errors.New("name cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// password validations
	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registrationUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// How do we deal with their passwords
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: registering user %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (h *UserHandler) HandleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		// h.logger.Printf("ERROR: decodingCreateWorkout: %v")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in to update"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": currentUser})
}
