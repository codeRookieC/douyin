package user

import (
	"context"
	common "github.com/codeRookieC/douyin/douyin-common"
	"github.com/codeRookieC/douyin/douyin-common/errs"
	"github.com/codeRookieC/douyin/douyin-grpc/user/userServer"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}

func (*HandlerUser) Register(ctx *gin.Context) {
	result := &common.Result{}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := UserServerClient.Register(c, &userServer.DouyinUserRegisterRequest{Username: username, Password: password})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.StatusCode))
}

func (*HandlerUser) Login(ctx *gin.Context) {
	result := &common.Result{}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//token := common.CreateToken(username)
	_, err := UserServerClient.Login(c, &userServer.DouyinUserLoginRequest{Username: username, Password: password})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(""))
}
