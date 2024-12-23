package services

import (
	"api-gym-on-go/src/config/errors"
	"api-gym-on-go/src/models"
	"api-gym-on-go/src/modules/users/interfaces"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UsersCreateService struct {
	UserRepository interfaces.UserRepository
}

func NewUsersCreateService(userRepo interfaces.UserRepository) *UsersCreateService {
	return &UsersCreateService{UserRepository: userRepo}
}

func (ucs *UsersCreateService) CreateUser(user *models.User) (createdUser *models.User, err error) {
	emailExist, err := ucs.UserRepository.UserEmailVerify(user.Email)
	fmt.Println(emailExist != nil && *emailExist == user.Email)
	if emailExist != nil && *emailExist == user.Email {
		return nil, &errors.UserAlreadyExistsError{}
	} else if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 6)
	if err != nil {
		return nil, &errors.CustomError{Message: "Failed to hash password", Code: 500}
	}

	user.Password = string(hash)

	createdUser, err = ucs.UserRepository.CreateUser(user)

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
