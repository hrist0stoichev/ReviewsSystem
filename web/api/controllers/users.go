package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/transfermodels"
)

type usersController struct {
	usersService services.UsersService
	baseController
}

func NewUsers(usersService services.UsersService, logger log.Logger, validator Validator) *usersController {
	return &usersController{
		usersService: usersService,
		baseController: baseController{
			logger:    logger,
			validator: validator,
		},
	}
}

func (uc *usersController) Register(res http.ResponseWriter, req *http.Request) {
	userRequest := transfermodels.CreateUserRequest{}
	if err := json.NewDecoder(req.Body).Decode(&userRequest); err != nil {
		uc.logger.WithError(err).Warnln("Could not decode request body")
		http.Error(res, "a problem occurred while reading user request", http.StatusBadRequest)
		return
	}

	err := uc.validator.Struct(userRequest)
	userRequest.ConfirmPassword = ""

	// TODO: should we cast err.(validator.ValidationErrors)? https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:          userRequest.Email,
		EmailConfirmed: false,
		Role:           models.Regular,
	}

	if userRequest.IsOwner {
		user.Role = models.Owner
	}

	if err = uc.usersService.CreateUser(user, &userRequest.Password); err != nil {
		uc.logger.WithError(err).Warnln("Could not create user")
		http.Error(res, "a problem occurred while creating a user", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
