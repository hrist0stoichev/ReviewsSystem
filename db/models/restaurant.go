package models

// TODO: Do I need the RatingsTotal, RatingsCount? They are used for internal purposes.
// TODO: Do I need the OwnerPtr
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
