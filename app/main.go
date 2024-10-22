package main

import (
	"flarrocca/go-distribuited-lock-poc/app/handler"
	"flarrocca/go-distribuited-lock-poc/app/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	redisSync := tools.NewRedisSync()
	handler := handler.NewHanler(redisSync)

	r := gin.Default()
	r.GET("/lock", handler.DoRequiredLockOperation)

	r.Run(":8080")
}
