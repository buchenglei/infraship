package connector

import (
	"net/http"

	"github.com/buchenglei/infraship/connector/cache"
	"github.com/buchenglei/infraship/connector/database"
	"github.com/buchenglei/infraship/connector/web"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	_ Connector[*gorm.DB]      = &database.MysqlConnector{}
	_ Connector[*redis.Client] = &cache.RedisConnector{}
	_ Connector[*http.Client]  = &web.StdHttpConnector{}
)
