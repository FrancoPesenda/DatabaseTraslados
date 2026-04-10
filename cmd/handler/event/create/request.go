package create

import domain "github.com/FrancoPesenda/eventra/internal/domain"

type request struct {
	Name      string `json:"name"`
	LocationID int    `json:"location_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Image     string `json:"image,omitempty"`
}

func (r request) eventDomain() domain.Event {
	return domain.Event{
		Name:       r.Name,
		LocationID: r.LocationID,
		StartDate:  r.StartDate,
		EndDate:    r.EndDate,
		Image:      r.Image,
	}
}
