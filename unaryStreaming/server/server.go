package unarystreaming

import (
	"context"

	userpb "github.com/centodiechi/unary_streams/protos/user/v1"
	storage "github.com/centodiechi/unary_streams/unaryStreaming/storage_provider"
	"github.com/centodiechi/unary_streams/unaryStreaming/utils"
	"go.uber.org/zap"
)

var Logger *zap.Logger

type RegisterService struct {
	userpb.UnimplementedRegisterServiceServer
}

type LoginService struct {
	userpb.UnimplementedAuthServiceServer
}

func (rs *RegisterService) SignUp(ctx context.Context, req *userpb.SignUpRequest) (*userpb.SignUpResponse, error) {
	uid := utils.Hash(req.User.Email + req.User.Password)
	req.User.Uid = uid
	req.User.Password = utils.Hash(req.User.Password)
	user := storage.User{
		UID:      req.User.Uid,
		Email:    req.User.Email,
		Password: req.User.Password,
		Profile:  storage.Profile{FirstName: req.User.Profile.Firstname, LastName: req.User.Profile.Lastname},
		Role:     storage.Role(req.User.Role),
	}
	Logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	Logger.With(zap.Any("User", user)).Info("User Created")
	err = storage.DataBase.CreateUser(user)
	if err != nil {
		return &userpb.SignUpResponse{Message: "Error Occured"}, err
	}
	return &userpb.SignUpResponse{Message: "Signed Up Succesfully"}, nil
}

func (ls *LoginService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	err := storage.DataBase.AuthenticateUser(utils.Hash(req.Email + req.Password))
	if err != nil {
		return &userpb.LoginResponse{Message: "Error Occured"}, err
	}
	Logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	Logger.With(zap.Any("User", req.Email)).Info("User Logged In")
	return &userpb.LoginResponse{Message: "logged in"}, nil
}
