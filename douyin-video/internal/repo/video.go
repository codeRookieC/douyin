package repo

import (
	"context"
	"video_server/internal/model"
)

type VideoRepo interface {
	GetVideoListByLatestTime(context.Context, int64) ([]model.Video, error)
	GetVideoListByUserId(context.Context, int64) ([]model.Video, error)
	GetUserInfoByUserId(context.Context, int64) (*model.User, error)
	GetFollowCountByUserId(context.Context, int64) (int64, error)
	GetFollowerCountByUserId(context.Context, int64) (int64, error)
	GetFavoriteCountByVideoId(context.Context, int64) (int64, error)
	GetCommentCountByVideoId(context.Context, int64) (int64, error)
	IsFavorite(ctx context.Context, userid int64, videoid int64) (bool, error)
	IsFollow(ctx context.Context, follower_id int64, follow_id int64) (bool, error)
	CreateVideo(ctx context.Context, id int64, author_id int64, play_url string, cover_url string, title string) (bool, error)
}
