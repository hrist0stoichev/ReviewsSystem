package models

// TODO: Do i need the RatingsTotal, RatingsCount? They are used for internal purposes.
type Restaurant struct {
	Id            string
	OwnerId       string
	Owner         *User
	Name          string
	City          string
	Address       string
	RatingsTotal  int
	RatingsCount  int
	AverageRating float32
}
