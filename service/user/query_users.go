package user

import (
	"vid-msa/model"
	"vid-msa/pkg/orm"
	"vid-msa/serializer"
)

// QueryUsersService 用户分页查询的服务
type QueryUsersService struct {
	Page  int `form:"page" json:"page" query:"page,required"`
	Limit int `form:"limit" json:"limit" query:"limit,required"`
}

// QueryUsers 分页查询用户
func (q *QueryUsersService) QueryUsers() *serializer.Response {
	users := make([]*model.User, 0)
	var total int64 = 0
	model.DB.Model(&model.User{}).Count(&total)
	model.DB.Scopes(orm.Paginate(q.Page, q.Limit)).Model(&model.User{}).Find(&users)
	res := serializer.BuildListResponse(int(total), q.Page, q.Limit, serializer.BuildUsersResponse(users))
	return res
}
