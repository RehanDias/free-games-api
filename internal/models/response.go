package models

type BaseResponse struct {
	Success   bool   `json:"success"`
	Timestamp string `json:"timestamp"`
}

type ApiResponse struct {
	BaseResponse
	Data GamesData `json:"data"`
}

type ErrorResponse struct {
	BaseResponse
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type GamesData struct {
	Current  []FormattedGame `json:"current"`
	Upcoming []FormattedGame `json:"upcoming"`
}

type FormattedGame struct {
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	OfferType     string    `json:"offerType"`
	EffectiveDate string    `json:"effectiveDate"`
	Seller        string    `json:"seller"`
	Price         GamePrice `json:"price"`
	Images        struct {
		Wide      string `json:"wide"`
		Thumbnail string `json:"thumbnail"`
	} `json:"images"`
	URLs struct {
		Product string `json:"product"`
	} `json:"urls"`
	Availability struct {
		EndDate string `json:"endDate"`
	} `json:"availability"`
}

type GamePrice struct {
	OriginalPrice          float64 `json:"originalPrice"`
	FormattedOriginalPrice string  `json:"formattedOriginalPrice"`
	Discount               string  `json:"discount"`
	Current                string  `json:"current"`
}
