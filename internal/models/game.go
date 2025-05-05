package models

type GameResponse struct {
	Success   bool           `json:"success"`
	Timestamp string         `json:"timestamp"`
	Data      GameData       `json:"data"`
	Error     *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type GameData struct {
	Current  []Game `json:"current"`
	Upcoming []Game `json:"upcoming"`
}

type Game struct {
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Status        string       `json:"status"`
	OfferType     string       `json:"offerType"`
	EffectiveDate string       `json:"effectiveDate"`
	SellerName    string       `json:"seller"`
	Price         Price        `json:"price"`
	Images        Images       `json:"images"`
	URLs          URLs         `json:"urls"`
	Availability  Availability `json:"availability"`
}

type Price struct {
	OriginalPrice          float64 `json:"originalPrice"`
	FormattedOriginalPrice string  `json:"formattedOriginalPrice"`
	Discount               string  `json:"discount"`
	Current                string  `json:"current"`
}

type Images struct {
	Wide      string `json:"wide"`
	Thumbnail string `json:"thumbnail"`
}

type URLs struct {
	Product string `json:"product"`
}

type Availability struct {
	EndDate string `json:"endDate"`
}

// Epic Games API response structures
type EpicResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []EpicGame `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

type EpicGame struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	EffectiveDate string `json:"effectiveDate"`
	OfferType     string `json:"offerType"`
	Seller        struct {
		Name string `json:"name"`
	} `json:"seller"`
	Price struct {
		TotalPrice struct {
			OriginalPrice int64 `json:"originalPrice"`
			DiscountPrice int64 `json:"discountPrice"`
			FmtPrice      struct {
				OriginalPrice string `json:"originalPrice"`
				DiscountPrice string `json:"discountPrice"`
			} `json:"fmtPrice"`
		} `json:"totalPrice"`
	} `json:"price"`
	KeyImages []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"keyImages"`
	CatalogNs struct {
		Mappings []struct {
			PageSlug string `json:"pageSlug"`
		} `json:"mappings"`
	} `json:"catalogNs"`
	ProductSlug string `json:"productSlug"`
	URLSlug     string `json:"urlSlug"`
	Promotions  struct {
		PromotionalOffers []struct {
			PromotionalOffers []struct {
				StartDate       string `json:"startDate"`
				EndDate         string `json:"endDate"`
				DiscountSetting struct {
					DiscountType       string  `json:"discountType"`
					DiscountPercentage float64 `json:"discountPercentage"`
				} `json:"discountSetting"`
			} `json:"promotionalOffers"`
		} `json:"promotionalOffers"`
		UpcomingPromotionalOffers []struct {
			PromotionalOffers []struct {
				StartDate       string `json:"startDate"`
				EndDate         string `json:"endDate"`
				DiscountSetting struct {
					DiscountType       string  `json:"discountType"`
					DiscountPercentage float64 `json:"discountPercentage"`
				} `json:"discountSetting"`
			} `json:"promotionalOffers"`
		} `json:"upcomingPromotionalOffers"`
	} `json:"promotions"`
}
