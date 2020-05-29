package transfermodels

type CreateRestaurantRequest struct {
	Name    string `json:"name" validate:"required,min=5,max=60"`
	City    string `json:"city" validate:"required,min=5,max=30"`
	Address string `json:"address" validate:"required,min=5,max=100"`
}

type RestaurantSimpleResponse struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	City          string  `json:"city"`
	Address       string  `json:"address"`
	AverageRating float32 `json:"average_rating"`
}
