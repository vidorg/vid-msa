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

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	// validate username and password
	if req.GetUsername() == "" {
		return nil, errors.New("user.LoginRequest.username is empty")
	}
	if req.GetPassword() == "" {
		return nil, errors.New("user.LoginRequest.password is empty")
	}
	var userModel model.User
	// find in db
	err := l.svcCtx.DB.Model(&model.User{}).
		Where("username = ?", req.Username).
		Where("password = ?", req.Password).
		First(&userModel).Error
	if err != nil {
		return nil, err
	}
	// generate jwt
	token, err := jwt.GenerateToken(userModel.ID, []byte(l.svcCtx.Secret))
	if err != nil {
		return nil, err
	}
	// store to redis
	pattern := jwt.RedisTokenConcat(strconv.FormatInt(userModel.ID, 10), token, "vid")
	if err = l.svcCtx.Rdb.Set(pattern, "1"); err != nil {
		return nil, err
	}
	return &user.LoginResponse{
		Token: pattern,
	}, nil
}
