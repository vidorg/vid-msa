package session

//var Session *session.Session
//
//// InitSession 初始化session
//func InitSession() *session.Session {
//
//	// redis存储session
//	provider,_ := redis.New(redis.Config{
//		KeyPrefix:   "session",
//		Addr:        viper.GetString("redis.address"),
//		PoolSize:    8,
//		IdleTimeout: 30 * time.Second,
//		Password:    viper.GetString("redis.password"),
//		DB:          viper.GetInt("redis.db"),
//	})
//
//	return session.New(session.Config{
//		//Session expiration time, possible values : 0 means no expiry (24 years), -1 means when the browser closes, >0 is the time.Duration which the session cookies should expire.
//		Expiration: -1,
//		Provider:   provider,
//		// chrome新特性,调试时关闭
//		SameSite: "None",
//		Secure:   true,
//	})
//}
//
//// Set 设置session key-value
//func Set(c *fiber.Ctx, key string, value interface{}) {
//	store := Session.Get(c)
//	store.Set(key, value)
//	err := store.Save()
//	if err != nil {
//		logger.Error("store set err", err)
//	}
//}
//
//// Get 获取session key-value
//func Get(c *fiber.Ctx, key string) interface{} {
//	store := Session.Get(c)
//	data := store.Get(key)
//	return data
//}
//
//// Delete 删除k-v
//func Delete(c *fiber.Ctx, key string) {
//	store := Session.Get(c)
//	store.Delete(key)
//	err := store.Save()
//	if err != nil {
//		logger.Error("delete session k-v err", err)
//	}
//}
