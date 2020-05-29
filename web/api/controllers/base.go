package controllers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

const (
	ModelDecodeError    = "Could not decode request body"
	InternalServerError = "Something went wrong"
	DefaultTop          = 20
	MinTop              = 1
	MaxTop              = 50
	DefaultSkip         = 0
	MinSkip             = 0
	MaxSkip             = math.MaxInt32
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

// parseIntParam parses an int parameter from the URI putting it inside the boundary [min, max].
// If the param cannot be parsed to int, a default value is used.
func (bc *baseController) parseIntParam(req *http.Request, param string, def, min, max int) int {
	intParam, err := strconv.Atoi(req.URL.Query().Get(param))
	if err != nil {
		intParam = def
	}

	if intParam > max {
		intParam = max
	}

	if intParam < min {
		intParam = min
	}

	return intParam
}
