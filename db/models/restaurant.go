package models

type Restaurant struct {
	Id            string
	OwnerId       string
	Owner         *User
	MinReviewId   *string
	MinReview     *Review
	MaxReviewId   *string
	MaxReview     *Review
	Name          string
	City          string
	Address       string
	Img           string
	Description   string
	RatingsTotal  int
	RatingsCount  int
	AverageRating float32
}
