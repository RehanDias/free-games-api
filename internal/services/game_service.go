package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"epic-games-free/internal/models"
)

type GameService struct {
	client *http.Client
}

func NewGameService() *GameService {
	return &GameService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *GameService) GetFreeGames() (*models.GameResponse, error) {
	url := "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=ID&allowCountries=ID"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return &models.GameResponse{
			Success:   false,
			Timestamp: formatDate(time.Now()),
			Error: &models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		}, nil
	}
	defer resp.Body.Close()

	var epicResp models.EpicResponse
	if err := json.NewDecoder(resp.Body).Decode(&epicResp); err != nil {
		return &models.GameResponse{
			Success:   false,
			Timestamp: formatDate(time.Now()),
			Error: &models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		}, nil
	}

	games := filterFreeGames(epicResp.Data.Catalog.SearchStore.Elements)

	return &models.GameResponse{
		Success:   true,
		Timestamp: formatDate(time.Now()),
		Data: models.GameData{
			Current:  formatGames(games.current, "FREE NOW"),
			Upcoming: formatGames(games.upcoming, "COMING SOON"),
		},
	}, nil
}

type gameCategories struct {
	current  []models.EpicGame
	upcoming []models.EpicGame
}

func filterFreeGames(elements []models.EpicGame) gameCategories {
	result := gameCategories{
		current:  make([]models.EpicGame, 0),
		upcoming: make([]models.EpicGame, 0),
	}

	for _, game := range elements {
		if game.Promotions.PromotionalOffers != nil && len(game.Promotions.PromotionalOffers) > 0 &&
			len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 &&
			game.Price.TotalPrice.DiscountPrice == 0 {
			result.current = append(result.current, game)
		} else if game.Promotions.UpcomingPromotionalOffers != nil &&
			len(game.Promotions.UpcomingPromotionalOffers) > 0 &&
			len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 &&
			game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].DiscountSetting.DiscountPercentage == 0 {
			result.upcoming = append(result.upcoming, game)
		}
	}

	return result
}

func formatGames(games []models.EpicGame, status string) []models.Game {
	formatted := make([]models.Game, len(games))

	for i, game := range games {
		var wide, thumbnail string
		for _, img := range game.KeyImages {
			switch img.Type {
			case "OfferImageWide":
				wide = img.URL
			case "Thumbnail":
				thumbnail = img.URL
			}
		}

		pageSlug := game.URLSlug
		if len(game.CatalogNs.Mappings) > 0 {
			pageSlug = game.CatalogNs.Mappings[0].PageSlug
		} else if game.ProductSlug != "" {
			pageSlug = game.ProductSlug
		}

		var endDate string
		if status == "FREE NOW" && len(game.Promotions.PromotionalOffers) > 0 &&
			len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.PromotionalOffers[0].PromotionalOffers[0].EndDate
		} else if status == "COMING SOON" && len(game.Promotions.UpcomingPromotionalOffers) > 0 &&
			len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].EndDate
		}

		formatted[i] = models.Game{
			Title:         game.Title,
			Description:   game.Description,
			Status:        status,
			OfferType:     game.OfferType,
			EffectiveDate: formatDate(time.Now()),
			SellerName:    getSellerName(game.Seller.Name),
			Price: models.Price{
				OriginalPrice:          float64(game.Price.TotalPrice.OriginalPrice),
				FormattedOriginalPrice: game.Price.TotalPrice.FmtPrice.OriginalPrice,
				Discount:               "100%",
				Current:                getCurrent(status),
			},
			Images: models.Images{
				Wide:      wide,
				Thumbnail: thumbnail,
			},
			URLs: models.URLs{
				Product: fmt.Sprintf("https://store.epicgames.com/en-US/p/%s", pageSlug),
			},
			Availability: models.Availability{
				EndDate: formatEndDate(endDate),
			},
		}
	}

	return formatted
}

func getSellerName(name string) string {
	if name == "" {
		return "Unknown"
	}
	return name
}

func getCurrent(status string) string {
	if status == "FREE NOW" {
		return "FREE"
	}
	return "COMING SOON"
}

func formatDate(date time.Time) string {
	jakartaLoc, _ := time.LoadLocation("Asia/Jakarta")
	jakartaTime := date.In(jakartaLoc)
	return fmt.Sprintf("%s WIB",
		jakartaTime.Format("Jan 2, 2006 15:04:05"))
}

func formatEndDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	endDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return ""
	}
	return formatDate(endDate)
}
