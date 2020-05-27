package controllers

import (
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

// TODO: Add validator to the configs. Maybe this interface needs to be moved.
type Validator interface {
	Struct(s interface{}) error
}

type baseController struct {
	logger    log.Logger
	validator Validator
}
