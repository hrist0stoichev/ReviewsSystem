package services

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type EncryptionService interface {
	GenerateSaltedHash(rawPassword *string) (string, error)
	PasswordsMatch(rawPassword, hashedPassword *string) bool
}

var DefaultEncryptionCost = bcrypt.DefaultCost

type encryptionService struct {
	cost int
}

func NewEncryptionService(cost int) EncryptionService {
	return &encryptionService{
		cost: cost,
	}
}

func (us *encryptionService) GenerateSaltedHash(rawPassword *string) (string, error) {
	saltedHash, err := bcrypt.GenerateFromPassword([]byte(*rawPassword), us.cost)

	// Remove the raw password from memory as fast as possible
	*rawPassword = ""

	if err != nil {
		return "", errors.Wrap(err, "could not generate salted hash from password")
	}

	return string(saltedHash), nil
}

func (us *encryptionService) PasswordsMatch(rawPassword, hashedPassword *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(*rawPassword))
	*rawPassword = ""
	return err == nil
}
