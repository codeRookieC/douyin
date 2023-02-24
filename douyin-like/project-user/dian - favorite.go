package main

import (
	"context"
	"test.com/project-user/pkg/model"
)

type DianfavoriteService struct {
	UnimplementedDianfavoriteServiceServer
}

var DianfavoriteService = &dianfavoriteService{}

func (u *dianfavoriteService) FavoriteAction(context context.Context, request *DouyinFavoriteActionRequest) (*DouyinFavoriteActionResponse, error) {
	dbManager, _ := cli.GetDB() // 获得db
	_, err := token.CheckAuth(context)

	statusMsg := "Verification failed!"
	if err != nil {
		return &DouyinUserResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
		}, nil
	}
	video_id := request.Video_id
	action_type := request.Action_type
	user_id := request.User_id

	if action_type == 1 {
		favorite := model.Favorite{Video_id: video_id, Uer_id: user_id}
		result := dbManager.Create(&favorite)
		statusMsg := ""
		if result.Error != nil {
			statusMsg = "failed! dianzan failed."
			return &DouyinFavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  &statusMsg,
			}, result.Error
		} else {
			statusMsg = "success"
			return &DouyinFavoriteActionResponse{
				StatusCode: 0,
				StatusMsg:  &statusMsg,
			}, nil
		}
	}
	if action_type == 2 {
		favorite := model.Favorite{}
		result := dbManager.where("video_id=? and user_id=?", video_id, user_id).Delete(&favorite)
		statusMsg := ""
		if result.Error != nil {
			statusMsg = "failed! quxiaodianzan failed."
			return &DouyinFavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  &statusMsg,
			}, result.Error
		} else {
			statusMsg = "success"
			return &DouyinFavoriteActionResponse{
				StatusCode: 0,
				StatusMsg:  &statusMsg,
			}, nil
		}
	}
	if action_type != 1 || action_type != 2 {
		statusMsg := "chuan ru chu cuo"
		return &DouyinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
		}, nil
	}

}

func (u *dianfavoriteService) FavoriteList(context context.Context, request *Douyin_favorite_list_request) (*Douyin_favorite_list_response, error) {

	dbManager, _ := cli.GetDB() // 获得db

	_, err := token.CheckAuth(context)

	statusMsg := "Verification failed!"
	if err != nil {
		return &DouyinUserResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
			Video:      nil,
		}, nil
	}
	user_id := request.User_id
	var video_id []int64
	result := dbManager.model(Favorite{}).where("user_id=?", user_id).Pluck("video_id", &video_id)

	if result.Error != nil {
		statusMsg = "failed! query userID failed."
		return &DouyinUserResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
			Video:      nil,
		}, nil
	}
	video := model.Video{}
	result := dbManager.where("video_id in (?)", video_id).Take(&vide)
	if result.Error != nil {
		statusMsg = "failed! chaxun failed."
		return &DouyinUserResponse{
			StatusCode: -1,
			StatusMsg:  &statusMsg,
			Video:      nil,
		}, nil
	} else {
		statusMsg = "success!"
		return &DouyinUserResponse{
			StatusCode: 0,
			StatusMsg:  &statusMsg,
			Video: &Video{
				id:          video.video_id,
				play_url:    video.play_url,
				cover_url:   video.cover_url,
				is_favorite: ture,
				title:       video.title,
			},
		}, nil
	}
}
