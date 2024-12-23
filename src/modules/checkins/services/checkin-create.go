package services

import (
	"api-gym-on-go/src/config/errors"
	"api-gym-on-go/src/config/utils"
	"api-gym-on-go/src/models"
	"api-gym-on-go/src/modules/checkins/interfaces"
	"api-gym-on-go/src/modules/checkins/schemas"
)

type CheckinCreate struct {
	checkinsRepository interfaces.CheckinsRepository
}

func NewCheckinCreateService(checkinsRepository interfaces.CheckinsRepository) *CheckinCreate {
	return &CheckinCreate{checkinsRepository: checkinsRepository}
}

func (cc *CheckinCreate) CreateCheckin(body *schemas.CheckinCreateBody) error {
	checkinAlrightExistsToday, err := cc.checkinsRepository.FindCheckinByIdOnDate(body.IDUser)
	if err != nil {
		return err
	} else if checkinAlrightExistsToday != nil && checkinAlrightExistsToday.IDUser == body.IDUser {
		return &errors.MaxNumberOfCheckinsError{}
	}

	gym, err := cc.checkinsRepository.FindGymByID(body.IDGym)
	if err != nil {
		return err
	} else if gym == nil {
		return &errors.ResourceNotFoundError{}
	}

	from := utils.Coordinate{Latitude: body.UserLatitude, Longitude: body.UserLongitude}
	to := utils.Coordinate{Latitude: gym.Latitude, Longitude: gym.Longitude}

	distance := utils.GetDistanceBetweenCoordinates(from, to)

	if distance > 1 {
		return &errors.InvalidCoordinatesError{}
	}

	checkin := models.Checkin{
		IDUser: body.IDUser,
		IDGym:  body.IDGym,
	}

	return cc.checkinsRepository.CreateCheckin(&checkin)
}
