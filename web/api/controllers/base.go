package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

const (
	ModelDecodeError    = "Could not decode request body"
	InternalServerError = "Something went wrong"
)

// TODO: Add validator to the configs. Maybe this interface needs to be moved.
type Validator interface {
	Struct(s interface{}) error
}

type baseController struct {
	logger    log.Logger
	validator Validator
}

func (bc *baseController) returnJsonResponse(w http.ResponseWriter, res interface{}) {
	if err := json.NewEncoder(w).Encode(res); err != nil {
		bc.logger.WithError(err).Warnln("Could not encode response")
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}
}
