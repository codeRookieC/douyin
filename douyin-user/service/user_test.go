package service

import (
	"fmt"
	"testing"
	"user/cli"
)

func TestRegister(t *testing.T) {
	// 连接数据库
	err := cli.InitDB()
	if err != nil {
		panic(err)
	}
	douyinUserRegisterRequset := &DouyinUserRegisterRequest{
		Username: "白歌",
		Password: "123456",
	}
	douyinUserRegisterResponse, err := UserService.Register(nil, douyinUserRegisterRequset)

	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", *douyinUserRegisterResponse)

}

func TestLogin(t *testing.T) {
	// 连接数据库
	err := cli.InitDB()
	if err != nil {
		panic(err)
	}
	douyinUserLoginRequset := &DouyinUserLoginRequest{
		Username: "白歌",
		Password: "123456",
	}
	douyinUserLoginResponse, err := UserService.Login(nil, douyinUserLoginRequset)

	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", *douyinUserLoginResponse)

}

func TestUserInfo(t *testing.T) {
	fmt.Println("")
}
