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

type Users struct {
	usersService      services.UsersService
	encryptionService services.EncryptionService
	tokensService     services.TokensService
	emailsService     services.EmailsService
	oauth2Service     services.OAuth2Service
	baseController
}

func NewUsers(
	usersService services.UsersService,
	encryptionService services.EncryptionService,
	tokensService services.TokensService,
	emailsService services.EmailsService,
	oauth2Service services.OAuth2Service,
	logger log.Logger,
	validator Validator,
) *Users {
	return &Users{
		usersService:      usersService,
		encryptionService: encryptionService,
		tokensService:     tokensService,
		emailsService:     emailsService,
		oauth2Service:     oauth2Service,
		baseController: baseController{
			logger:    logger,
			validator: validator,
		},
	}
}

func (uc *Users) Register(res http.ResponseWriter, req *http.Request) {
	userRequest := transfermodels.CreateUserRequest{}
	if err := json.NewDecoder(req.Body).Decode(&userRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	if err := uc.validator.Struct(userRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	_, err := uc.usersService.GetByEmail(userRequest.Email)
	if err != services.ErrUserNotFound {
		if err == nil {
			http.Error(res, "User already exists", http.StatusConflict)
			return
		}

		uc.logger.WithError(err).Warnln("could not get user by email")
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		return
	}

	saltedHash, err := uc.encryptionService.GenerateSaltedHash(&userRequest.Password)
	if err != nil {
		uc.logger.WithError(err).Warnln("Could not generate salted hash")
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	confirmationToken := uc.emailsService.GenerateRandomEmailToken()

	user := &models.User{
		Email:                  userRequest.Email,
		EmailConfirmed:         false,
		EmailConfirmationToken: &confirmationToken,
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
		if err = uc.emailsService.SendConfirmationEmail(user.Email, *user.EmailConfirmationToken); err != nil {
			uc.logger.WithError(err).Warnln("could not send confirmation email")
		}
	}()

	uc.returnJsonResponse(res, transfermodels.CreateUserResponse{
		Ok: true,
	})
}

func (uc *Users) Login(res http.ResponseWriter, req *http.Request) {
	loginRequest := transfermodels.LoginRequest{}
	if err := json.NewDecoder(req.Body).Decode(&loginRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	user, err := uc.usersService.GetByEmail(loginRequest.Email)
	if err != nil {
		if err == services.ErrUserNotFound {
			http.Error(res, InvalidCredentials, http.StatusNotFound)
			return
		}

		uc.logger.WithError(err).Warnln("Could not get username by email")
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !uc.encryptionService.PasswordsMatch(&loginRequest.Password, &user.HashedPassword) {
		http.Error(res, InvalidCredentials, http.StatusNotFound)
		return
	}

	// TODO: Uncomment this line when done with testing
	// if !user.EmailConfirmed {
	// 	http.Error(res, EmailNotConfirmed, http.StatusBadRequest)
	// 	return
	// }

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

	uc.returnJsonResponse(res, resp)
}

func (uc *Users) RedirectToFacebookAuth(res http.ResponseWriter, req *http.Request) {
	redirectionURL := uc.oauth2Service.GenerateAuthURL()
	http.Redirect(res, req, redirectionURL, http.StatusTemporaryRedirect)
}

func (uc *Users) FacebookLogin(res http.ResponseWriter, req *http.Request) {
	fbLoginRequest := struct {
		State string `json:"state"`
		Code  string `json:"code"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&fbLoginRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	token, err := uc.oauth2Service.GetToken(fbLoginRequest.State, fbLoginRequest.Code)
	if err != nil {
		uc.logger.WithError(err).Warnln("Failed to get oauth2 token")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	userInfo, err := uc.oauth2Service.GetUserInfo(token)
	if err != nil {
		uc.logger.WithError(err).Warnln("Failed to get oauth2 user info")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	// Some internal server error occurred
	dbUser, err := uc.usersService.GetByEmail(userInfo.Email)
	if err != nil && err != services.ErrUserNotFound {
		uc.logger.WithError(err).Warnln("Failed to get user by email")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	// The user is not found
	if err == services.ErrUserNotFound {
		user := &models.User{
			Email:                  userInfo.Email,
			EmailConfirmed:         true,
			EmailConfirmationToken: nil,
			HashedPassword:         "",
			Role:                   models.Regular,
		}

		err = uc.usersService.CreateUser(user)
		if err != nil {
			uc.logger.WithError(err).Warnln("Failed to create new facebook user")
			http.Error(res, InternalServerError, http.StatusInternalServerError)
			return
		}

		dbUser = user
	} else { // The user is found
		// The user is found but they haven't registered through facebook
		if dbUser.HashedPassword != "" {
			http.Error(res, "You have already registered with this email from the basic registration form!", http.StatusConflict)
			return
		}
	}

	jwt, claims, err := uc.tokensService.GenerateSignedToken(&services.UserClaims{
		Id:   dbUser.Id,
		Role: dbUser.Role.String(),
	})
	if err != nil {
		uc.logger.WithError(err).Warnln("Failed to generate signed token")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	resp := transfermodels.LoginResponse{
		Token:   jwt,
		Expires: time.Unix(claims.ExpiresAt.Unix(), 0),
		Email:   dbUser.Email,
		Role:    claims.Role,
	}

	uc.returnJsonResponse(res, resp)
}

func (uc *Users) ConfirmEmail(res http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")
	token := req.URL.Query().Get("token")

	if email == "" || token == "" {
		http.NotFound(res, req)
		return
	}

	user, err := uc.usersService.GetByEmail(email)
	if err != nil {
		if err == services.ErrUserNotFound {
			http.NotFound(res, req)
			return
		}

		uc.logger.WithError(err).Warnln("Could not get username by email")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	if user.EmailConfirmed {
		http.Error(res, "Email already confirmed", http.StatusConflict)
		return
	}

	if *user.EmailConfirmationToken != token {
		http.NotFound(res, req)
		return
	}

	if err = uc.usersService.ConfirmEmail(user.Id); err != nil {
		uc.logger.WithError(err).Warnln("could not confirm email")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "https://google.com?confirmation_successful=true", http.StatusSeeOther)
}
