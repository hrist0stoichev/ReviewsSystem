package transfermodels

type CreateRestaurantRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=60"`
	City        string `json:"city" validate:"required,min=5,max=30"`
	Address     string `json:"address" validate:"required,min=5,max=100"`
	Img         string `json:"img" validate:"required,url"`
	Description string `json:"description" validate:"required,min=30,max=500"`
}

type RestaurantSimpleResponse struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	City          string  `json:"city"`
	Address       string  `json:"address"`
	Img           string  `json:"img"`
	Description   string  `json:"description"`
	AverageRating float32 `json:"average_rating"`
}

type RestaurantDetailedResponse struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	City          string                `json:"city"`
	Address       string                `json:"address"`
	Img           string                `json:"img"`
	Description   string                `json:"description"`
	AverageRating float32               `json:"average_rating"`
	MinReview     *ReviewSimpleResponse `json:"min_review"`
	MaxReview     *ReviewSimpleResponse `json:"max_review"`
}
