package main

import (
	"github.com/codeRookieC/douyin/douyin-api/config"
	"github.com/codeRookieC/douyin/douyin-api/router"
	common "github.com/codeRookieC/douyin/douyin-common"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	common.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}
