package service

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const salt = "qwerty"

type ValidateService struct {
	validate *validator.Validate
}

func NewValidateService() *ValidateService {
	return &ValidateService{validate: validator.New()}
}

func (s *ValidateService) validateStruct(data interface{}) error {
	return s.validate.Struct(data)
}

func (s *ValidateService) validateVariable(data interface{}, tag string) error {
	return s.validate.Var(data, tag)
}

func generatePassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes)
}

func comparePasswordHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
