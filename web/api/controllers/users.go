package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/transfermodels"
)

const (
	InvalidCredentials = "Invalid username or password"
	EmailNotConfirmed  = "Email is not confirmed"
)

type usersController struct {
	usersService      services.UsersService
	encryptionService services.EncryptionService
	tokensService     services.TokensService
	emailsService     services.EmailsService
	baseController
}

func NewUsers(
	usersService services.UsersService,
	encryptionService services.EncryptionService,
	tokensService services.TokensService,
	emailsService services.EmailsService,
	logger log.Logger,
	validator Validator,
) *usersController {
	return &usersController{
		usersService:      usersService,
		encryptionService: encryptionService,
		tokensService:     tokensService,
		emailsService:     emailsService,
		baseController: baseController{
			logger:    logger,
			validator: validator,
		},
	}
}

func (uc *usersController) Register(res http.ResponseWriter, req *http.Request) {
	userRequest := transfermodels.CreateUserRequest{}
	if err := json.NewDecoder(req.Body).Decode(&userRequest); err != nil {
		uc.logger.WithError(err).Warnln(ModelDecodeError)
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	err := uc.validator.Struct(userRequest)
	userRequest.ConfirmPassword = ""

	// TODO: should we cast err.(validator.ValidationErrors)? https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	saltedHash, err := uc.encryptionService.GenerateSaltedHash(&userRequest.Password)
	if err != nil {
		uc.logger.WithError(err).Warnln("Could not generate salted hash")
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:                  userRequest.Email,
		EmailConfirmed:         false,
		EmailConfirmationToken: uc.emailsService.GenerateRandomEmailToken(),
		HashedPassword:         saltedHash,
		Role:                   models.Regular,
	}

	if userRequest.IsOwner {
		user.Role = models.Owner
	}

	if err = uc.usersService.CreateUser(user); err != nil {
		uc.logger.WithError(err).Warnln("Could not create user")
		http.Error(res, "a problem occurred while creating a user", http.StatusBadRequest)
		return
	}

	// Send the confirmation email async. A functionality for resending emails should be implemented.
	go func() {
		if err = uc.emailsService.SendConfirmationEmail(user.Email, user.EmailConfirmationToken); err != nil {
			uc.logger.WithError(err).Warnln("could not send confirmation email")
		}
	}()

	res.WriteHeader(http.StatusOK)
}

func (uc *usersController) Login(res http.ResponseWriter, req *http.Request) {
	loginRequest := transfermodels.LoginRequest{}
	if err := json.NewDecoder(req.Body).Decode(&loginRequest); err != nil {
		uc.logger.WithError(err).Warnln(ModelDecodeError)
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	if err := uc.validator.Struct(loginRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.usersService.GetByEmail(loginRequest.Email)
	if err != nil {
		if err == services.ErrUserNotFound {
			http.Error(res, InvalidCredentials, http.StatusUnauthorized)
			return
		}

		uc.logger.WithError(err).Warnln("Could not get username by email")
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !uc.encryptionService.PasswordsMatch(&loginRequest.Password, &user.HashedPassword) {
		http.Error(res, InvalidCredentials, http.StatusUnauthorized)
		return
	}

	if !user.EmailConfirmed {
		http.Error(res, EmailNotConfirmed, http.StatusUnauthorized)
		return
	}

	jwt, claims, err := uc.tokensService.GenerateSignedToken(&services.UserClaims{
		Id:   user.Id,
		Role: user.Role.String(),
	})
	if err != nil {
		uc.logger.WithError(err).Warnln("Could not generate jwt")
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		return
	}

	resp := transfermodels.LoginResponse{
		Token:   jwt,
		Expires: time.Unix(claims.ExpiresAt.Unix(), 0),
		Email:   user.Email,
		Role:    claims.Role,
	}

	if err = json.NewEncoder(res).Encode(resp); err != nil {
		uc.logger.WithError(err).Warnln("Could not encode response")
	}
}
