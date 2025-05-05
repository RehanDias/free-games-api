package services

import (
	"encoding/json"
	"net/http"
	"time"

	"free-games-epic/internal/models"
	"free-games-epic/internal/utils"
)

const (
	EPIC_API = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=ID&allowCountries=ID"
)

type EpicService struct{}

func NewEpicService() *EpicService {
	return &EpicService{}
}

func (s *EpicService) GetFreeGames() (*models.GamesData, error) {
	resp, err := http.Get(EPIC_API)
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
		Current:  s.formatGames(currentGames, "FREE NOW"),
		Upcoming: s.formatGames(upcomingGames, "COMING SOON"),
	}, nil
}

func (s *EpicService) filterFreeGames(data models.EpicResponse) ([]models.Game, []models.Game) {
	var currentGames, upcomingGames []models.Game

	for _, game := range data.Data.Catalog.SearchStore.Elements {
		if len(game.Promotions.PromotionalOffers) == 0 && len(game.Promotions.UpcomingPromotionalOffers) == 0 {
			continue
		}

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

func (s *EpicService) formatGames(games []models.Game, status string) []models.FormattedGame {
	var formattedGames []models.FormattedGame

	for _, game := range games {
		var endDate string

		if status == "FREE NOW" && len(game.Promotions.PromotionalOffers) > 0 && len(game.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.PromotionalOffers[0].PromotionalOffers[0].EndDate
		} else if len(game.Promotions.UpcomingPromotionalOffers) > 0 && len(game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 {
			endDate = game.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].EndDate
		}

		endDateTime, _ := time.Parse(time.RFC3339, endDate)
		effectiveTime, _ := time.Parse(time.RFC3339, game.EffectiveDate)

		formatted := s.createFormattedGame(game, status, effectiveTime, endDateTime)
		formattedGames = append(formattedGames, formatted)
	}

	return formattedGames
}

func (s *EpicService) createFormattedGame(game models.Game, status string, effectiveTime, endDateTime time.Time) models.FormattedGame {
	var wideURL, thumbnailURL string
	for _, img := range game.KeyImages {
		if img.Type == "OfferImageWide" {
			wideURL = img.URL
		}
		if img.Type == "Thumbnail" {
			thumbnailURL = img.URL
		}
	}

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
			Discount:               "100%",
			Current:                s.getCurrentPrice(status),
		},
	}

	formatted.Images.Wide = wideURL
	formatted.Images.Thumbnail = thumbnailURL
	formatted.URLs.Product = "https://store.epicgames.com/en-US/p/" + pageSlug
	formatted.Availability.EndDate = utils.FormatDate(endDateTime)

	return formatted
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
	return "Unknown"
}

func (s *EpicService) getCurrentPrice(status string) string {
	if status == "FREE NOW" {
		return "FREE"
	}
	return "COMING SOON"
}
