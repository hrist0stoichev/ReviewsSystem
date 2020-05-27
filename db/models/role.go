package models

import (
	"database/sql/driver"

	"github.com/pkg/errors"
)

type Role uint8

const (
	Regular Role = iota
	Owner
	Admin
)

var roles = [...]string{
	"regular",
	"owner",
	"admin",
}

func (r Role) String() string {
	return roles[r]
}

func (r Role) Value() (driver.Value, error) {
	return r.String(), nil
}

func (r *Role) Scan(value interface{}) error {
	valueString, ok := value.(string)
	if !ok {
		return errors.New("role is not a string")
	}

	for i, role := range roles {
		if role == valueString {
			*r = Role(uint8(i))
			return nil
		}
	}

	return errors.New("invalid role")
}
