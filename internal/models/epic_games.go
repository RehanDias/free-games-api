package models

type EpicResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []Game `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

type Game struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OfferType   string `json:"offerType"`
	Seller      struct {
		Name string `json:"name"`
	} `json:"seller"`
	EffectiveDate string `json:"effectiveDate"`
	Price         struct {
		TotalPrice struct {
			OriginalPrice float64 `json:"originalPrice"`
			DiscountPrice float64 `json:"discountPrice"`
			FmtPrice      struct {
				OriginalPrice string `json:"originalPrice"`
			} `json:"fmtPrice"`
		} `json:"totalPrice"`
	} `json:"price"`
	KeyImages   []KeyImage `json:"keyImages"`
	CatalogNs   CatalogNs  `json:"catalogNs"`
	ProductSlug string     `json:"productSlug"`
	UrlSlug     string     `json:"urlSlug"`
	Promotions  struct {
		PromotionalOffers         []PromotionOffer `json:"promotionalOffers"`
		UpcomingPromotionalOffers []PromotionOffer `json:"upcomingPromotionalOffers"`
	} `json:"promotions"`
}

type KeyImage struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type CatalogNs struct {
	Mappings []struct {
		PageSlug string `json:"pageSlug"`
	} `json:"mappings"`
}

type PromotionOffer struct {
	PromotionalOffers []struct {
		StartDate       string `json:"startDate"`
		EndDate         string `json:"endDate"`
		DiscountSetting struct {
			DiscountPercentage int `json:"discountPercentage"`
		} `json:"discountSetting"`
	} `json:"promotionalOffers"`
}
