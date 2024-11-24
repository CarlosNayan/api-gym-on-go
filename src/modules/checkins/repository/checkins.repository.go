package repository

import (
	"api-gym-on-go/models"
	"api-gym-on-go/src/config/utils"
	"log"

	"gorm.io/gorm"
)

type CheckinRepository struct {
	DB *gorm.DB
}

func NewCheckinRepository(db *gorm.DB) *CheckinRepository {
	return &CheckinRepository{DB: db}
}

func (cr *CheckinRepository) CreateCheckin(checkin *models.Checkin) error {
	return cr.DB.Create(checkin).Error
}

func (cr *CheckinRepository) FindCheckinByIdOnDate(id_user string) (*models.Checkin, error) {
	var checkin models.Checkin

	now, err := utils.NewMoment()
	if err != nil {
		log.Fatalf("Erro ao criar o data: %v", err)
	}

	startOfDay := now.StartOf("day").Format()
	endOfDay := now.EndOf("day").Format()

	result := cr.DB.
		Where("id_user = ? AND created_at BETWEEN ? AND ?", id_user, startOfDay, endOfDay).
		First(&checkin)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &checkin, nil
}

func (cr *CheckinRepository) FindCheckinById(id_checkin string) (*models.Checkin, error) {
	var checkin models.Checkin

	result := cr.DB.
		Where("id_checkin = ?", id_checkin).
		First(&checkin)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &checkin, nil
}

func (cr *CheckinRepository) UpdateCheckin(id_checkin string) (*models.Checkin, error) {
	var updatedCheckin models.Checkin

	err := cr.DB.Model(&models.Checkin{}).
		Where("id_checkin = ?", id_checkin).
		Updates(map[string]interface{}{
			"validated_at": "now()",
		}).
		First(&updatedCheckin, "id_checkin = ?", id_checkin).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &updatedCheckin, nil
}

func (cr *CheckinRepository) CountByUserId(id_user string) (int64, error) {
	var count int64
	err := cr.DB.Model(&models.Checkin{}).
		Where("id_user = ?", id_user).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return count, err
}

func (cr *CheckinRepository) ListAllCheckinsHistoryOfUser(id_user string, page int) ([]models.Checkin, error) {
	var checkins []models.Checkin
	err := cr.DB.Where("id_user = ?", id_user).Limit(10).
		Offset((page - 1) * 10).
		Find(&checkins).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return checkins, err
}

func (cr *CheckinRepository) FindGymByID(id_gym string) (*models.Gym, error) {
	var gym models.Gym
	err := cr.DB.Where("id_gym = ?", id_gym).First(&gym).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &gym, nil
}
