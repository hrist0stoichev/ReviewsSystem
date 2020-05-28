package services

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UsersService interface {
	CreateUser(user *models.User, rawPassword *string) error
	GetByEmail(email string) (*models.User, error)
	PasswordsMatch(rawPassword, hashedPassword *string) bool
	GenerateToken(user *models.User) (*JWT, error)
}

type JWT struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Role      string `json:"role"`
}

type usersService struct {
	db db.Manager
}

func NewUserService(dbManager db.Manager) UsersService {
	return &usersService{
		db: dbManager,
	}
}

func (us *usersService) CreateUser(user *models.User, rawPassword *string) error {
	saltedHash, err := bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)

	// Remove the raw password from memory as fast as possible
	*rawPassword = ""

	if err != nil {
		return errors.Wrap(err, "could not generate salted hash from password")
	}

	user.HashedPassword = string(saltedHash)

	// TODO: Return custom error if email already exists
	if err = us.db.Users().Insert(user); err != nil {
		return errors.Wrap(err, "could not insert user in database")
	}

	return nil
}

func (us *usersService) GetByEmail(email string) (*models.User, error) {
	user, err := us.db.Users().GetByEmail(email)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrUserNotFound
		}

		return nil, errors.Wrap(err, "could not get user by email")
	}

	return user, nil
}

func (us *usersService) PasswordsMatch(rawPassword, hashedPassword *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(*rawPassword))
	*rawPassword = ""
	return err == nil
}

func (us *usersService) GenerateToken(user *models.User) (*JWT, error) {
	expiresAt := time.Now().Add(8 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Subject":   user.Id,
		"Role":      user.Role.String(),
		"ExpiresAt": expiresAt,
	})

	strToken, err := token.SignedString([]byte("pass"))
	if err != nil {
		return nil, errors.Wrap(err, "could not sign token")
	}

	return &JWT{
		Token:     strToken,
		ExpiresAt: expiresAt,
		Role:      user.Role.String(),
	}, nil
}
