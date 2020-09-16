package initialization

import "github.com/go-redis/redis"

var Client *redis.Client

func init() {
	if Configuration.RedisSetting.Exists {
		// 自动监听redis
		redisIp := Configuration.RedisSetting.Ip
		redisPort := Configuration.RedisSetting.Port
		if redisIp == "" {
			redisIp = "localhost"
		}
		if redisPort == "" {
			redisPort = "6379"
		}
		Client = redis.NewClient(&redis.Options{
			Addr:     redisIp + ":" + redisPort,
			Password: Configuration.RedisSetting.Password,
			DB:       Configuration.RedisSetting.Db,
		})
	}
}