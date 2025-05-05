package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	EPIC_API = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=ID&allowCountries=ID"
)

// Structs for Epic Games API response
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

// Response structs
type ApiResponse struct {
	Success   bool      `json:"success"`
	Timestamp string    `json:"timestamp"`
	Data      GamesData `json:"data"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Timestamp string `json:"timestamp"`
	Error     struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
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

// Handler function for Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request for CORS preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// For any other non-GET request
	if r.Method != "GET" {
		handleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get data from Epic Games API
	resp, err := http.Get(EPIC_API)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var epicResp EpicResponse
	if err := json.NewDecoder(resp.Body).Decode(&epicResp); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter free games
	currentGames, upcomingGames := filterFreeGames(epicResp)

	// Format response
	response := ApiResponse{
		Success:   true,
		Timestamp: formatDate(time.Now()),
		Data: GamesData{
			Current:  formatGames(currentGames, "FREE NOW"),
			Upcoming: formatGames(upcomingGames, "COMING SOON"),
		},
	}

	// Send response
	json.NewEncoder(w).Encode(response)
}

func handleError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	errResp := ErrorResponse{
		Success:   false,
		Timestamp: formatDate(time.Now()),
		Error: struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: message,
			Code:    code,
		},
	}
	json.NewEncoder(w).Encode(errResp)
}

func filterFreeGames(data EpicResponse) ([]Game, []Game) {
	var currentGames, upcomingGames []Game

	for _, game := range data.Data.Catalog.SearchStore.Elements {
		// Skip if no promotions
		if len(game.Promotions.PromotionalOffers) == 0 && len(game.Promotions.UpcomingPromotionalOffers) == 0 {
			continue
		}

		// Check for current free games
		if len(game.Promotions.PromotionalOffers) > 0 &&
			len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 &&
			game.Price.TotalPrice.DiscountPrice == 0 {
			currentGames = append(currentGames, game)
		} else if len(game.Promotions.UpcomingPromotionalOffers) > 0 &&
			len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 &&
			game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].DiscountSetting.DiscountPercentage == 0 {
			upcomingGames = append(upcomingGames, game)
		}
	}

	return currentGames, upcomingGames
}

func formatGames(games []Game, status string) []FormattedGame {
	var formattedGames []FormattedGame

	for _, game := range games {
		var endDate string

		if status == "FREE NOW" && len(game.Promotions.PromotionalOffers) > 0 && len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.PromotionalOffers[0].PromotionalOffers[0].EndDate
		} else if len(game.Promotions.UpcomingPromotionalOffers) > 0 && len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].EndDate
		}

		// Parse date as time.Time for formatting
		endDateTime, _ := time.Parse(time.RFC3339, endDate)

		// Find images
		var wideURL, thumbnailURL string
		for _, img := range game.KeyImages {
			if img.Type == "OfferImageWide" {
				wideURL = img.URL
			}
			if img.Type == "Thumbnail" {
				thumbnailURL = img.URL
			}
		}

		// Find page slug
		var pageSlug string
		if len(game.CatalogNs.Mappings) > 0 {
			pageSlug = game.CatalogNs.Mappings[0].PageSlug
		}
		if pageSlug == "" {
			pageSlug = game.ProductSlug
		}
		if pageSlug == "" {
			pageSlug = game.UrlSlug
		}

		// Parse effective date
		effectiveTime, _ := time.Parse(time.RFC3339, game.EffectiveDate)

		// Format game data
		formatted := FormattedGame{
			Title:         game.Title,
			Description:   game.Description,
			Status:        status,
			OfferType:     game.OfferType,
			EffectiveDate: formatDate(effectiveTime),
			Seller:        getSellerName(game),
			Price: GamePrice{
				OriginalPrice:          game.Price.TotalPrice.OriginalPrice,
				FormattedOriginalPrice: game.Price.TotalPrice.FmtPrice.OriginalPrice,
				Discount:               "100%",
				Current: func() string {
					if status == "FREE NOW" {
						return "FREE"
					} else {
						return "COMING SOON"
					}
				}(),
			},
		}

		formatted.Images.Wide = wideURL
		formatted.Images.Thumbnail = thumbnailURL
		formatted.URLs.Product = "https://store.epicgames.com/en-US/p/" + pageSlug
		formatted.Availability.EndDate = formatDate(endDateTime)

		formattedGames = append(formattedGames, formatted)
	}

	return formattedGames
}

func getSellerName(game Game) string {
	if game.Seller.Name != "" {
		return game.Seller.Name
	}
	return "Unknown"
}

func formatDate(date time.Time) string {
	// Convert to Jakarta time (UTC+7)
	jakartaTime := date.UTC().Add(7 * time.Hour)

	// Format date
	year, month, day := jakartaTime.Date()
	hour, min, sec := jakartaTime.Clock()

	// Format month as short name
	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthStr := monthNames[month-1]

	return fmt.Sprintf("%s %d, %d %02d:%02d:%02d WIB",
		monthStr, day, year, hour, min, sec)
}
