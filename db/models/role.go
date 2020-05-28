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
	valueByte, ok := value.([]byte)
	if !ok {
		return errors.New("role is not a byte array")
	}

	valueString := string(valueByte)
	for i, role := range roles {
		if role == valueString {
			*r = Role(uint8(i))
			return nil
		}
	}

	return errors.New("invalid role")
}
