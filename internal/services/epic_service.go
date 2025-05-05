package services

import (
	"encoding/json"
	"net/http"
	"time"

	"free-games-epic/internal/models"
	"free-games-epic/internal/utils"
)

const (
	EpicAPIURL         = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions"
	EpicStoreBaseURL   = "https://store.epicgames.com/en-US/p/"
	StatusFreeNow      = "FREE NOW"
	StatusComingSoon   = "COMING SOON"
	DefaultDiscount    = "100%"
	DefaultSeller      = "Unknown"
	ImageTypeWide      = "OfferImageWide"
	ImageTypeThumbnail = "Thumbnail"
)

type EpicService struct {
	client *http.Client
}

func NewEpicService() *EpicService {
	return &EpicService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *EpicService) GetFreeGames() (*models.GamesData, error) {
	req, err := http.NewRequest(http.MethodGet, EpicAPIURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", "en-US")
	q.Add("country", "ID")
	q.Add("allowCountries", "ID")
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var epicResp models.EpicResponse
	if err := json.NewDecoder(resp.Body).Decode(&epicResp); err != nil {
		return nil, err
	}

	currentGames, upcomingGames := s.filterFreeGames(epicResp)

	return &models.GamesData{
		Current:  s.formatGames(currentGames, StatusFreeNow),
		Upcoming: s.formatGames(upcomingGames, StatusComingSoon),
	}, nil
}

func (s *EpicService) filterFreeGames(data models.EpicResponse) (current, upcoming []models.Game) {
	for _, game := range data.Data.Catalog.SearchStore.Elements {
		if !s.hasValidPromotions(game) {
			continue
		}

		if s.isCurrentlyFree(game) {
			current = append(current, game)
		} else if s.isUpcomingFree(game) {
			upcoming = append(upcoming, game)
		}
	}
	return current, upcoming
}

func (s *EpicService) hasValidPromotions(game models.Game) bool {
	return len(game.Promotions.PromotionalOffers) > 0 || len(game.Promotions.UpcomingPromotionalOffers) > 0
}

func (s *EpicService) isCurrentlyFree(game models.Game) bool {
	return len(game.Promotions.PromotionalOffers) > 0 &&
		len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 &&
		game.Price.TotalPrice.DiscountPrice == 0
}

func (s *EpicService) isUpcomingFree(game models.Game) bool {
	return len(game.Promotions.UpcomingPromotionalOffers) > 0 &&
		len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 &&
		game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].DiscountSetting.DiscountPercentage == 0
}

func (s *EpicService) formatGames(games []models.Game, status string) []models.FormattedGame {
	formattedGames := make([]models.FormattedGame, 0, len(games))

	for _, game := range games {
		endDate := s.getEndDate(game, status)
		endDateTime, _ := time.Parse(time.RFC3339, endDate)
		effectiveTime, _ := time.Parse(time.RFC3339, game.EffectiveDate)

		formatted := s.createFormattedGame(game, status, effectiveTime, endDateTime)
		formattedGames = append(formattedGames, formatted)
	}

	return formattedGames
}

func (s *EpicService) getEndDate(game models.Game, status string) string {
	if status == StatusFreeNow && len(game.Promotions.PromotionalOffers) > 0 &&
		len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 {
		return game.Promotions.PromotionalOffers[0].PromotionalOffers[0].EndDate
	}
	if len(game.Promotions.UpcomingPromotionalOffers) > 0 &&
		len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 {
		return game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].EndDate
	}
	return ""
}

func (s *EpicService) createFormattedGame(game models.Game, status string, effectiveTime, endDateTime time.Time) models.FormattedGame {
	images := s.getGameImages(game)
	pageSlug := s.getPageSlug(game)

	formatted := models.FormattedGame{
		Title:         game.Title,
		Description:   game.Description,
		Status:        status,
		OfferType:     game.OfferType,
		EffectiveDate: utils.FormatDate(effectiveTime),
		Seller:        s.getSellerName(game),
		Price: models.GamePrice{
			OriginalPrice:          game.Price.TotalPrice.OriginalPrice,
			FormattedOriginalPrice: game.Price.TotalPrice.FmtPrice.OriginalPrice,
			Discount:               DefaultDiscount,
			Current:                s.getCurrentPrice(status),
		},
	}

	formatted.Images.Wide = images.wide
	formatted.Images.Thumbnail = images.thumbnail
	formatted.URLs.Product = EpicStoreBaseURL + pageSlug
	formatted.Availability.EndDate = utils.FormatDate(endDateTime)

	return formatted
}

type gameImages struct {
	wide      string
	thumbnail string
}

func (s *EpicService) getGameImages(game models.Game) gameImages {
	var images gameImages
	for _, img := range game.KeyImages {
		switch img.Type {
		case ImageTypeWide:
			images.wide = img.URL
		case ImageTypeThumbnail:
			images.thumbnail = img.URL
		}
	}
	return images
}

func (s *EpicService) getPageSlug(game models.Game) string {
	if len(game.CatalogNs.Mappings) > 0 {
		return game.CatalogNs.Mappings[0].PageSlug
	}
	if game.ProductSlug != "" {
		return game.ProductSlug
	}
	return game.UrlSlug
}

func (s *EpicService) getSellerName(game models.Game) string {
	if game.Seller.Name != "" {
		return game.Seller.Name
	}
	return DefaultSeller
}

func (s *EpicService) getCurrentPrice(status string) string {
	if status == StatusFreeNow {
		return "FREE"
	}
	return StatusComingSoon
}
