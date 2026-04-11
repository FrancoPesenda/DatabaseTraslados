package create

import domain "github.com/FrancoPesenda/eventra/internal/domain"

type createResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LocationID int    `json:"location_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Image      string `json:"image,omitempty"`
}

func newCreateResponse(e domain.Event) createResponse {
	return createResponse{
		ID:         e.ID,
		Name:       e.Name,
		LocationID: e.LocationID,
		StartDate:  e.StartDate,
		EndDate:    e.EndDate,
		Image:      e.Image,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
