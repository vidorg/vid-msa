package controller

import (
	"vid-msa/service/file"

	"github.com/gofiber/fiber/v2"
)

// GetUploadToken get oss file token
func GetUploadToken(c *fiber.Ctx) error {
	service := &file.GetQiniuFileTokenService{}
	response := service.GenerateSimpleToken()
	return c.JSON(response)
}

// UploadFile upload file to server
func UploadFile(c *fiber.Ctx) error {
	service := &file.UploadToLocalService{}
	response := service.UploadFileToLocal(c)
	return c.JSON(response)
}
