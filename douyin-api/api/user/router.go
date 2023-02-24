package user

import (
	"github.com/codeRookieC/douyin/douyin-api/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RouterUser struct {
}

func init() {
	log.Println("init userServer router")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Route(r *gin.Engine) {
	InitRpcUserClient()
	h := New()
	r.POST("/douyin/user/register", h.Register)
}
