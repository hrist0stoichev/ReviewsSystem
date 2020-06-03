package transfermodels

import (
	"time"
)

type CreateReviewRequest struct {
	RestaurantId string `json:"restaurant_id" validate:"required,uuid"`
	Rating       uint8  `json:"rating" validate:"required,min=1,max=5"`
	Comment      string `json:"comment" validate:"required,min=30,max=300"`
}

type ReviewSimpleResponse struct {
	Id        string    `json:"id"`
	Reviewer  string    `json:"reviewer"`
	Rating    uint8     `json:"rating"`
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`
	Answer    *string   `json:"answer"`
}
