package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var Enforcer *casbin.Enforcer

// Enforce Casbin 权限控制 RBAC
func Enforce(adapter *gormadapter.Adapter, role string, c *fiber.Ctx) (bool, error) {
	var err error
	Enforcer, err = casbin.NewEnforcer(viper.GetString("casbin.config-path"), adapter)
	if err != nil {
		return false, err
	}
	// 启用casbin功能
	if !viper.GetBool("casbin.enable") {
		Enforcer.EnableEnforce(false)
	}
	// 开启casbin权限检查日志
	if viper.GetBool("casbin.log") {
		Enforcer.EnableLog(true)
	}
	err = Enforcer.LoadPolicy()
	if err != nil {
		return false, err
	}
	sub := role
	obj := c.Path()
	act := c.Method()
	return Enforcer.Enforce(sub, obj, act)
}
