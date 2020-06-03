package models

type Restaurant struct {
	Id            string
	OwnerId       string
	Owner         *User
	Name          string
	City          string
	Address       string
	Img           string
	Description   string
	RatingsTotal  int
	RatingsCount  int
	AverageRating float32
}
