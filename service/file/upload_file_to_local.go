package file

import (
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
	"github.com/seefs001/seefslib-go/xfile"
	"github.com/seefs001/seefslib-go/xstring"
)

// UploadToLocalService
type UploadToLocalService struct {
}

// UploadFileToLocal 上传文件到本地static目录
func (s *UploadToLocalService) UploadFileToLocal(c *fiber.Ctx) *serializer.Response {
	file, err := c.FormFile("file")
	if err != nil {
		return serializer.UploadFileErr("上传文件出错", err)
	}
	ext := xfile.GetFileExt(file.Filename)
	filename := xstring.RandStringRunes(10) + ext
	fileLocation := "./static/" + filename
	err = c.SaveFile(file, fileLocation)
	if err != nil {
		return serializer.UploadFileErr("文件保存出错", err)
	}
	return &serializer.Response{
		Code: 200,
		Msg:  "上传文件成功",
		Data: filename,
	}
}
