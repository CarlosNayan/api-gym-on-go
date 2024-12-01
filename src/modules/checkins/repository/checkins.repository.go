package repository

import (
	"api-gym-on-go/models"
	"api-gym-on-go/src/config/utils"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CheckinRepository struct {
	DB *sql.DB
}

func NewCheckinRepository(db *sql.DB) *CheckinRepository {
	return &CheckinRepository{DB: db}
}

func (cr *CheckinRepository) CreateCheckin(checkin *models.Checkin) error {
	id := uuid.New()

	query := `
		INSERT INTO checkins
		(id_user, id_gym)
		VALUES
		($1, $2, $3)
	`

	_, err := cr.DB.Exec(query, id, checkin.IDUser, checkin.IDGym)
	if err != nil {
		return fmt.Errorf("error inserting checkin: %w", err)
	}

	return nil
}

func (cr *CheckinRepository) FindCheckinByIdOnDate(id_user string) (*models.Checkin, error) {
	var checkin models.Checkin

	now, err := utils.NewMoment()
	if err != nil {
		log.Fatalf("Erro ao criar o data: %v", err)
	}

	startOfDay := now.StartOf("day").Format()
	endOfDay := now.EndOf("day").Format()

	query := `
		SELECT id, id_user, id_gym, created_at 
		FROM checkins 
		WHERE id_user = $1 
		AND created_at BETWEEN $2 AND $3
	`

	row := cr.DB.QueryRow(query, id_user, startOfDay, endOfDay)
	err = row.Scan(&checkin.ID, &checkin.IDUser, &checkin.IDGym, &checkin.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching checkin: %w", err)
	}

	return &checkin, nil
}

func (cr *CheckinRepository) FindCheckinById(id_checkin string) (*models.Checkin, error) {
	var checkin models.Checkin

	query := `
		SELECT id, id_user, id_gym, created_at 
		FROM checkins 
		WHERE id_checkin = $1
	`

	row := cr.DB.QueryRow(query, id_checkin)
	err := row.Scan(&checkin.ID, &checkin.IDUser, &checkin.IDGym, &checkin.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching checkin: %w", err)
	}

	return &checkin, nil
}

func (cr *CheckinRepository) UpdateCheckin(id_checkin string) (*models.Checkin, error) {
	var updatedCheckin models.Checkin

	query := `
		UPDATE checkins
		SET validated_at = now()
		WHERE id_checkin = $1
	`

	row := cr.DB.QueryRow(query, id_checkin)
	if row.Err() != nil {
		return nil, fmt.Errorf("error updating checkin: %w", row.Err())
	}

	err := row.Scan(&updatedCheckin.ID, &updatedCheckin.IDUser, &updatedCheckin.IDGym, &updatedCheckin.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning checkin row: %w", err)
	}

	return &updatedCheckin, nil
}

func (cr *CheckinRepository) CountByUserId(id_user string) (int64, error) {
	var count int64

	query := `
		SELECT COUNT(*) FROM checkins 
		WHERE id_user = $1
	`

	row := cr.DB.QueryRow(query, id_user)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error scanning checkin row: %w", err)
	}

	return count, err
}

func (cr *CheckinRepository) ListAllCheckinsHistoryOfUser(id_user string, page int) ([]models.Checkin, error) {
	var checkins []models.Checkin

	query := `
		SELECT id, id_user, id_gym, created_at 
		FROM checkins 
		WHERE id_user = $1
		LIMIT 10
		OFFSET $2
	`

	rows, err := cr.DB.Query(query, id_user, (page-1)*10)
	if err != nil {
		return nil, fmt.Errorf("error fetching checkins: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var checkin models.Checkin
		err = rows.Scan(&checkin.ID, &checkin.IDUser, &checkin.IDGym, &checkin.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning checkin row: %w", err)
		}
		checkins = append(checkins, checkin)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return checkins, err
}

func (cr *CheckinRepository) FindGymByID(id_gym string) (*models.Gym, error) {
	var gym models.Gym

	query := `
		SELECT id, gym_name, description, latitude, longitude FROM gyms
		WHERE id = $1
	`

	row := cr.DB.QueryRow(query, id_gym)
	err := row.Scan(&gym.ID, &gym.GymName, &gym.Description, &gym.Latitude, &gym.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching gym: %w", err)
	}

	return &gym, nil
}
