package dto

import pb "github.com/zhayt/user-service/proto"

type ChangeUserPasswordDTO struct {
	Email              string `validator:"required"`
	OldPassword        string `validator:"required"`
	NewPassword        string `validator:"required,eqfield=ConfirmNewPassword""`
	ConfirmNewPassword string `validate:"required"`
}

func NewChangeUserPasswordDTO(dto *pb.ChangeUserPasswordDTO) *ChangeUserPasswordDTO {
	return &ChangeUserPasswordDTO{
		Email:              dto.Email,
		OldPassword:        dto.OldPassword,
		NewPassword:        dto.NewPassword,
		ConfirmNewPassword: dto.ConfirmNewPassword,
	}
}

type ChangeUserNameDTO struct {
	Email string `validator:"required"`
	Name  string `validator:"required,min=3,max=50"`
}

func NewChangeUserNameDTO(dto *pb.ChangeUserNameDTO) *ChangeUserNameDTO {
	return &ChangeUserNameDTO{
		Email: dto.Email,
		Name:  dto.NewName,
	}
}
