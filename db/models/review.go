package models

import (
	"time"
)

type Review struct {
	Id           string
	RestaurantId string
	Restaurant   *Restaurant
	ReviewerId   string
	Reviewer     *User
	Rating       uint8
	Timestamp    time.Time
	Comment      string
	Answer       *string
}
