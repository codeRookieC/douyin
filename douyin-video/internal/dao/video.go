package dao

import (
	"context"
	"fmt"
	"time"
	"video_server/internal/database/gorms"
	"video_server/internal/model"

	"gorm.io/gorm"
)

//	type VideoRepo interface {
//		GetVideoListByLatestTime(context.Context, int64) (*[]data.Video, error)
//		GetVideoListByUserId(ctx context.Context, user_id int64) (*[]data.Video, error)
//		GetUserInfoById(ctx context.Context, userid int64) (data.User, error)
//		GetFollowCountById(ctx context.Context, userid int64) (int64, error)
//		GetFollowerCountById(ctx context.Context, userid int64) (int64, error)
//		GetFavoriteCountByVideoId(ctx context.Context, VideoId int64) (int64, error)
//		GetCommentCountByVideoId(ctx context.Context, VideoId int64) (int64, error)
//		IsFavorite(ctx context.Context, userid int64, videoid int64) (bool, error)
//		CreateVideo(ctx context.Context, id int64, author_id int64, play_url string, cover_url string, title string) (bool, error)
//	}
type VideoDao struct {
	conn *gorms.GormConn
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: gorms.New(),
	}
}
func (v *VideoDao) GetVideoListByLatestTime(ctx context.Context, LatestTime int64) ([]model.Video, error) {
	var videos []model.Video
	Time := time.Unix(LatestTime, 0)
	err := v.conn.Session(ctx).Model(&model.Video{}).Where("created_at<?", Time).Order("created_at DESC").Limit(30).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return nil, err
	}
	return videos, nil
}

// ok
func (vs *VideoDao) GetVideoListByUserId(ctx context.Context, user_id int64) ([]model.Video, error) {
	var videos []model.Video
	err := vs.conn.Session(ctx).Model(&model.Video{}).Where("author_id =?", user_id).Order("created_at").Limit(30).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return nil, err
	}
	return videos, nil
}

// ok
func (vs *VideoDao) GetUserInfoByUserId(ctx context.Context, userid int64) (*model.User, error) {
	var user model.User
	err := vs.conn.Session(ctx).Model(&model.User{}).Where("Id=?", userid).First(&user).Error
	if err == gorm.ErrRecordNotFound {

		return &user, err
	}
	return &user, nil
}

// ok
func (vs *VideoDao) GetFollowCountByUserId(ctx context.Context, userid int64) (int64, error) {
	var count int64

	err := vs.conn.Session(ctx).Model(&model.Follow{}).Where("follower_id=?", userid).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return 0, err
	}
	return count, nil
}

// ok
func (vs *VideoDao) GetFollowerCountByUserId(ctx context.Context, userid int64) (int64, error) {
	var count int64

	err := vs.conn.Session(ctx).Model(&model.Follow{}).Where("follow_id=?", userid).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return 0, err
	}
	return count, nil
}

// ok
func (vs *VideoDao) GetFavoriteCountByVideoId(ctx context.Context, VideoId int64) (int64, error) {
	var count int64

	err := vs.conn.Session(ctx).Model(&model.Favorite{}).Where("video_id=?", VideoId).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return count, err
	}
	return count, nil
}

// ok
func (vs *VideoDao) GetCommentCountByVideoId(ctx context.Context, VideoId int64) (int64, error) {
	var count int64

	err := vs.conn.Session(ctx).Model(&model.Comment{}).Where("video_id=?", VideoId).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		//没有找到对应的数据
		return count, err
	}
	return count, nil
}

// ok
func (vs *VideoDao) IsFavorite(ctx context.Context, userid int64, videoid int64) (bool, error) {
	is := false
	var count int64
	err := vs.conn.Session(ctx).Where("user_id=? AND video_id=?", userid, videoid).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return is, err
	}
	if count == 1 {
		is = true
	}
	return is, nil
}

// ok
func (vs *VideoDao) IsFollow(ctx context.Context, follower_id int64, follow_id int64) (bool, error) {
	is := false
	var count int64
	err := vs.conn.Session(ctx).Where("follower_id=? AND follow_id=?", follower_id, follow_id).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return is, err
	}
	if count == 1 {
		is = true
	}
	return is, nil
}

// ok
func (vs *VideoDao) CreateVideo(ctx context.Context, id int64, author_id int64, play_url string, cover_url string, title string) (bool, error) {
	video := model.Video{
		Id:        id,
		Author_id: author_id,
		Play_url:  play_url,
		Cover_url: cover_url,
		Title:     title,
	}
	fmt.Print(video)
	err := vs.conn.Session(ctx).Model(&model.Video{}).Create(&video).Error
	if err == gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}
