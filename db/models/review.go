package models

import (
	"time"
)

type Review struct {
	Id           string
	RestaurantId string
	ReviewerId   string
	Rating       uint8
	Timestamp    time.Time
	Comment      string
	Answer       *string
}
