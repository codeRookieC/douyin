package service

import (
	"context"
	"github.com/codeRookieC/douyin/douyin-user/cli"
	"github.com/codeRookieC/douyin/douyin-user/model"
	"github.com/codeRookieC/douyin/douyin-user/service/token"
)

type userService struct {
	UnimplementedUserServiceServer
}

var UserService = &userService{}

func (u *userService) Register(context context.Context, request *DouyinUserRegisterRequest) (*DouyinUserRegisterResponse, error) {
	dbManager, _ := cli.GetDB() // 获得db
	username := request.Username
	password := request.Password
	user := model.User{Username: username, Password: password}
	result := dbManager.Create(&user)

	statusMsg := ""
	if result.Error != nil {
		statusMsg = "failed! insert failed."
		return &DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
			Token:      "",
			UserId:     -1,
		}, result.Error
	}

	tokenString := token.CreateToken(username)
	statusMsg = "success! Registered."
	return &DouyinUserRegisterResponse{
		UserId:     user.UserID,
		StatusMsg:  &statusMsg,
		Token:      tokenString,
		StatusCode: 0,
	}, nil

}
func (u *userService) Login(context context.Context, request *DouyinUserLoginRequest) (*DouyinUserLoginResponse, error) {
	dbManager, _ := cli.GetDB() // 获得db

	username := request.Username
	password := request.Password
	user := model.User{Username: username, Password: password}
	result := dbManager.Where("username=? and password=?", username, password).First(&user)

	statusMsg := ""
	if result.Error != nil {
		statusMsg = "failed! query username failed."
		return &DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
			Token:      "",
			UserId:     -1,
		}, result.Error
	}

	tokenString := token.CreateToken(user.Username)
	statusMsg = "success! Login."
	return &DouyinUserLoginResponse{
		UserId:     user.UserID,
		StatusMsg:  &statusMsg,
		Token:      tokenString,
		StatusCode: 0,
	}, nil

}

// UserInfo 这个接口应该是我去看别人的信息，而不是查询自己的信息
func (u *userService) UserInfo(context context.Context, request *DouyinUserRequest) (*DouyinUserResponse, error) {

	dbManager, _ := cli.GetDB() // 获得db

	_, err := token.CheckAuth(context)

	statusMsg := "Verification failed!"
	if err != nil {
		return &DouyinUserResponse{
			StatusCode: -1,
			User:       nil,
			StatusMsg:  &statusMsg,
		}, nil
	}

	user := model.User{}
	result := dbManager.Where("user_id=", request.UserId).First(&user)

	if result.Error != nil {
		statusMsg = "failed! query userID failed."
		return &DouyinUserResponse{
			StatusCode: -1,
			User:       nil,
			StatusMsg:  &statusMsg,
		}, nil
	}
	var followCount int64 = 200
	var followerCount int64 = 300
	statusMsg = "success!"
	return &DouyinUserResponse{
		StatusCode: -1,
		User: &User{
			Id:            user.UserID,
			Name:          user.Username,
			FollowCount:   &followCount,
			FollowerCount: &followerCount,
			IsFollow:      false,
		},
		StatusMsg: &statusMsg,
	}, nil
}
