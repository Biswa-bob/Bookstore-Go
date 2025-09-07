package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadIDParam(r *http.Request) (primitive.ObjectID, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return primitive.NilObjectID, errors.New("invalid id parameter")

	}

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid id parameter type")
	}
	return objectID, nil
}
