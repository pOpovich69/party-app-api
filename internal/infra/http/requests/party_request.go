package requests

import (
	"go-rest-api/internal/domain"
	"time"
)

type CreatePartyRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Price       int32     `json:"price" validate:"required"`
	StartDate   time.Time `json:"startDate" validate:"required"`
}

type UpdatePartyRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Price       int32     `json:"price" validate:"required"`
	StartDate   time.Time `json:"startDate" validate:"required"`
}

func (upr UpdatePartyRequest) ToDomainModel() (interface{}, error) {
	return domain.Party{
		Title:       upr.Title,
		Description: upr.Description,
		Image:       upr.Image,
		Price:       upr.Price,
		StartDate:   upr.StartDate,
	}, nil
}

func (cpr CreatePartyRequest) ToDomainModel() (interface{}, error) {
	return domain.Party{
		Title:       cpr.Title,
		Description: cpr.Description,
		Image:       cpr.Image,
		Price:       cpr.Price,
		StartDate:   cpr.StartDate,
	}, nil
}
