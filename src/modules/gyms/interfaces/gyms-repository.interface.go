package interfaces

import "api-gym-on-go/models"

type GymsRepository interface {
	CreateGym(gym *models.Gym) error
	GymsNearby(latitude, longitude float64) ([]models.Gym, error)
	SearchGyms(query string) ([]models.Gym, error)
}
