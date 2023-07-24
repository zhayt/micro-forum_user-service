package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zhayt/user-service/model"
	"github.com/zhayt/user-service/model/dto"
	pb "github.com/zhayt/user-service/proto"
	"github.com/zhayt/user-service/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const _defaultContextTimeout = 5 * time.Second

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  *storage.Storage
	validate *ValidateService
	l        *zap.Logger
}

func NewUserService(storage *storage.Storage, validate *ValidateService, l *zap.Logger) *UserService {
	return &UserService{storage: storage, validate: validate, l: l}
}

func (s *UserService) CreateUser(ctx context.Context, userPB *pb.User) (*pb.UserProfileDTO, error) {
	// convert proto struct to golang struct
	user := model.NewUser(userPB)

	// validate struct data
	if err := s.validate.validateStruct(user); err != nil {
		s.l.Error("validateStruct error", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to validate struct: %w", err)
	}

	user.Password = generatePassword(user.Password)
	// try to create user
	userID, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		s.l.Error("CreateUser error", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create user: %w", err)
	}

	s.l.Info("User created", zap.Uint64("id", userID))
	return &pb.UserProfileDTO{
		Id:    userID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *pb.GetUserByIDReq) (*pb.User, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		s.l.Error("GetUserByID error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("%s", err))
		}

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	s.l.Info("User found", zap.Uint64("id", user.ID))
	userPB := &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return userPB, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailReq) (*pb.User, error) {
	user, err := s.storage.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.l.Error("GetUserByEmail error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("%s", err))
		}

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	userPB := &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return userPB, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, passDTO *pb.ChangeUserPasswordDTO) (*pb.UserUpdateResponse, error) {
	// convert proto struct to my struct
	userPassDTO := dto.NewChangeUserPasswordDTO(passDTO)

	// data validate
	if err := s.validate.validateStruct(userPassDTO); err != nil {
		s.l.Error("validateStruct error", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%s", err))
	}

	// get user for check password
	user, err := s.storage.GetUserByEmail(ctx, userPassDTO.Email)
	if err != nil {
		s.l.Error("GetUserByEmail error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%s", err))
		}

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	// compare password
	if err = comparePasswordHash(user.Password, userPassDTO.OldPassword+salt); err != nil {
		s.l.Error("comparePasswordHash error", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%s", err))
	}

	// update user password
	if err = s.storage.UpdateUserPassword(ctx, userPassDTO); err != nil {
		s.l.Error("UpdateUserPassword error", zap.Error(err))
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	// return response
	s.l.Info("User password updated", zap.Uint64("id", user.ID))
	return &pb.UserUpdateResponse{
		Success: true,
		Message: "Password updated",
	}, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, nameDTO *pb.ChangeUserNameDTO) (*pb.UserUpdateResponse, error) {
	// convert proto struct to my struct
	userNameUpdate := dto.NewChangeUserNameDTO(nameDTO)

	user, err := s.storage.GetUserByEmail(ctx, nameDTO.Email)
	if err != nil {
		s.l.Error("GetUserByEmail", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%s", err))
		}

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	if err = s.storage.UpdateUserName(ctx, userNameUpdate); err != nil {
		s.l.Error("UpdateUserName error", zap.Error(err))
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s", err))
	}

	s.l.Info("User name updated", zap.Uint64("id", user.ID))
	return &pb.UserUpdateResponse{
		Success: true,
		Message: "User name updated",
	}, nil
}
