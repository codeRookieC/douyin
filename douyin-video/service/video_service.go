package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"video_server/api/proto/pb"
	"video_server/internal/dao"
	"video_server/internal/oss"
	"video_server/internal/repo"
)

type VideoService struct {
	pb.UnimplementedVideoServiceServer
	dataDao repo.VideoRepo
	ossDao  repo.OssRepo
}

func New() *VideoService {
	return &VideoService{
		dataDao: dao.NewVideoDao(),
		ossDao:  oss.NewOssDao(),
	}
}
func (vs *VideoService) DouyinFeed(c context.Context, req *pb.DouyinFeedRequest) (*pb.DouyinFeedResponse, error) {

	var _StatusCode int32
	var _StatusMsg string
	var latestTime int64
	_StatusCode = 0
	_StatusMsg = "succeed"
	resp := pb.DouyinFeedResponse{}
	latestTime = *req.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	videoInfo_list, err := vs.dataDao.GetVideoListByLatestTime(c, latestTime)
	latestTime = videoInfo_list[len(videoInfo_list)-1].CreatedAt.Unix()
	fmt.Println(videoInfo_list)
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}

	for _, value := range videoInfo_list {

		_Video := pb.Video{}
		_User := pb.User{}
		Author_id := value.Author_id
		userInfo, err := vs.dataDao.GetUserInfoByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		followCount, err := vs.dataDao.GetFollowCountByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		followerCount, err := vs.dataDao.GetFollowerCountByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		Id := value.Id
		favoriteCount, err := vs.dataDao.GetFavoriteCountByVideoId(c, Id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()

		}
		commentCount, err := vs.dataDao.GetCommentCountByVideoId(c, Id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}

		isFollow := false
		isFavorite := false
		user_Id := userInfo.Id
		_User.Id = &user_Id
		_User.Name = &userInfo.Username
		_User.FollowerCount = &followerCount
		_User.FollowCount = &followCount
		_User.IsFollow = &isFollow

		value_Id := value.Id
		_Video.Id = &value_Id
		_Video.Author = &_User
		_Video.PlayUrl = &value.Play_url
		_Video.CoverUrl = &value.Cover_url
		_Video.Title = &value.Title
		_Video.CommentCount = &commentCount
		_Video.FavoriteCount = &favoriteCount
		_Video.IsFavorite = &isFavorite

		resp.VideoList = append(resp.VideoList, &_Video)

	}
	resp.StatusCode = &_StatusCode
	resp.StatusMsg = &_StatusMsg
	resp.NextTime = &latestTime
	return &resp, err

}
func (vs *VideoService) DouyinPublishAction(c context.Context, req *pb.DouyinPublishActionRequest) (*pb.DouyinPublishActionResponse, error) {
	// message douyin_publish_action_request {
	// 	required string token = 1; // 用户鉴权token
	// 	required bytes data = 2; // 视频数据
	// 	required string title = 3; // 视频标题
	//   }
	var ID int64
	var author_id int64
	var play_url string
	var cover_url string
	var title string
	var _StatusCode int32
	var _StatusMsg string
	_StatusCode = 0
	_StatusMsg = "succeed"

	author_id, _ = strconv.ParseInt(*req.Token, 10, 64) // Token to UserId

	title = *req.Title
	ID = time.Now().Unix()
	// play_url
	f, err := os.Create("./resource/" + title + ".mp4")
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}
	defer f.Close()
	_, err = f.Write(req.Data)
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}

	err = vs.ossDao.PutObject(title+strconv.FormatInt(ID, 10)+".mp4", "./resource/"+title+".mp4")
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}
	play_url, err = vs.ossDao.GetObjectUrl(title+strconv.FormatInt(ID, 10)+".mp4", 300)
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}
	os.Remove("./resource/" + title + ".mp4")
	// cover_url
	// common.GetSnapshot("./resource/"+title+".mp4", "./resource/"+title+".png", 1)
	// err = vs.ossDao.PutObject(title+strconv.FormatInt(ID, 10)+".png", "./resource/"+title+".png")
	// if err != nil {
	// 	log.Println(err)
	// 	_StatusCode = -1
	// 	_StatusMsg = err.Error()
	// }
	// cover_url, err = vs.ossDao.GetObjectUrl(title+strconv.FormatInt(ID, 10)+".png", 300)
	// if err != nil {
	// 	log.Println(err)
	// 	_StatusCode = -1
	// 	_StatusMsg = err.Error()
	// }

	cover_url, err = vs.ossDao.GetObjectUrl("go.PNG", 300)
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}
	// ....
	_, err = vs.dataDao.CreateVideo(c, ID, author_id, play_url, cover_url, title)

	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}

	resp := pb.DouyinPublishActionResponse{
		StatusCode: &_StatusCode,
		StatusMsg:  &_StatusMsg,
	}

	return &resp, err

}

func (vs *VideoService) DouyinPublishList(c context.Context, req *pb.DouyinPublishListRequest) (*pb.DouyinPublishListResponse, error) {

	resp := pb.DouyinPublishListResponse{}
	var _StatusCode int32
	var _StatusMsg string
	_StatusCode = 0
	_StatusMsg = "succeed"

	videoInfo_list, err := vs.dataDao.GetVideoListByUserId(c, *req.UserId)
	if err != nil {
		log.Println(err)
		_StatusCode = -1
		_StatusMsg = err.Error()
	}

	for _, value := range videoInfo_list {

		_Video := pb.Video{}
		_User := pb.User{}
		Author_id := value.Author_id
		userInfo, err := vs.dataDao.GetUserInfoByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		followCount, err := vs.dataDao.GetFollowCountByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		followerCount, err := vs.dataDao.GetFollowerCountByUserId(c, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		Video_Id := value.Id
		favoriteCount, err := vs.dataDao.GetFavoriteCountByVideoId(c, Video_Id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		commentCount, err := vs.dataDao.GetCommentCountByVideoId(c, Video_Id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}

		isFollow, err := vs.dataDao.IsFollow(c, *req.UserId, Author_id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		isFavorite, err := vs.dataDao.IsFavorite(c, *req.UserId, Video_Id)
		if err != nil {
			log.Println(err)
			_StatusCode = -1
			_StatusMsg = err.Error()
		}
		user_Id := userInfo.Id
		_User.Id = &user_Id
		_User.Name = &userInfo.Username
		_User.FollowerCount = &followerCount
		_User.FollowCount = &followCount
		_User.IsFollow = &isFollow

		_Video.Id = &Video_Id
		_Video.Author = &_User
		_Video.PlayUrl = &value.Play_url
		_Video.CoverUrl = &value.Cover_url
		_Video.Title = &value.Title
		_Video.CommentCount = &commentCount
		_Video.FavoriteCount = &favoriteCount
		_Video.IsFavorite = &isFavorite

		resp.VideoList = append(resp.VideoList, &_Video)

	}
	resp.StatusCode = &_StatusCode
	resp.StatusMsg = &_StatusMsg

	return &resp, err
}
