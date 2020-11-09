package file

import (
	"vid-msa/serializer"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
)

// GetQiniuFileTokenService qiniu token
type GetQiniuFileTokenService struct {
}

// GenerateSimpleToken
func (s *GetQiniuFileTokenService) GenerateSimpleToken() *serializer.Response {
	putPolicy := storage.PutPolicy{
		Scope: viper.GetString("qiniu.bucket"),
	}
	mac := qbox.NewMac(viper.GetString("qiniu.accessKey"), viper.GetString("qiniu.secretKey"))
	uploadToken := putPolicy.UploadToken(mac)
	return &serializer.Response{
		Code: 200,
		Msg:  "获取七牛token成功",
		Data: uploadToken,
	}
}
