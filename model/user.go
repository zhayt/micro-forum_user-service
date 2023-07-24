package model

import pb "github.com/zhayt/user-service/proto"

type User struct {
	ID       uint64
	Name     string `validate:"required,alpha,min=3,max=50"`
	Email    string `validate:"required,lowercase"`
	Password string `validate:"required"`
}

func NewUser(user *pb.User) *User {
	return &User{
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}
