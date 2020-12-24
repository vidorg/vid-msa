package logic

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"user/internal/model"
	"user/internal/pkg/jwt"

	"user/internal/svc"
	"user/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// validate username and password
	if req.GetUsername() == "" {
		return nil, errors.New("user.LoginRequest.username is empty")
	}
	if req.GetPassword() == "" {
		return nil, errors.New("user.LoginRequest.password is empty")
	}
	// store to db
	userModel := model.User{
		UserName: req.Username,
		Password: req.Password,
	}
	if err := l.svcCtx.DB.Where("username = ?", userModel.UserName).Error; err != nil {
		return nil, errors.New("user is existed")
	}
	if err := l.svcCtx.DB.Create(&userModel).Error; err != nil {
		return nil, err
	}
	// generate jwt
	token, err := jwt.GenerateToken(userModel.ID, l.svcCtx.Secret)
	if err != nil {
		return nil, err
	}
	// store to redis
	pattern := jwt.RedisTokenConcat(strconv.FormatInt(userModel.ID, 10), token, "vid")
	if err = l.svcCtx.Rdb.Set(pattern, "1"); err != nil {
		return nil, err
	}
	return &user.RegisterResponse{
		Username: req.Username,
		Token:    pattern,
	}, nil
}
