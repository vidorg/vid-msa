package util

import (
	"strings"

	"github.com/seefs001/seefslib-go/xstring"

	"github.com/gofiber/fiber/v2"
)

// IdsStrToIdsIntGroup 获取URL中批量id并解析
func IdsStrToIdsIntGroup(key string, c *fiber.Ctx) []int {
	return idsStrToIdsIntGroup(c.Params(key))
}

func idsStrToIdsIntGroup(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := xstring.StringToInt(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
