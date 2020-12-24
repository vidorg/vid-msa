package logic

import (
	"context"
	"github.com/pkg/errors"
	"user/internal/model"
	"user/internal/pkg/jwt"
	"user/internal/svc"
	"user/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(req *user.GetUserRequset) (*user.GetUserResponse, error) {
	// validate token empty
	if req.GetToken() == "" {
		return nil, errors.New("user.GetUserRequset.token is empty")
	}
	// redis
	_, err := l.svcCtx.Rdb.Get(req.Token)
	if err != nil {
		return nil, errors.Wrap(err, "rdb.getToken err")
	}
	// split token
	token := jwt.GetTokenFromRedisPattern(req.Token)
	// get user
	claims, err := jwt.ParseToken(token, l.svcCtx.Secret)
	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseToken err")
	}
	userModel := model.User{}
	if err = l.svcCtx.DB.Model(&model.User{}).
		First(&userModel, claims.UID).Error; err != nil {
		return nil, err
	}
	return &user.GetUserResponse{
		Id:         userModel.ID,
		Username:   userModel.UserName,
		Avatar:     userModel.Avatar,
		Gender:     userModel.Gender,
		CreateTime: userModel.CreatedAt,
	}, nil
}
