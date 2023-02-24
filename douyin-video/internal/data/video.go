package data

type Video struct {
	Id        int64  // 视频唯一标识
	Author_id int64  // 视频作者id
	PlayUrl   string // 视频播放地址
	CoverUrl  string // 视频封面地址
	Title     string // 视频标题
}

type User struct {
	Id       int64  // 用户id
	Name     string // 用户名称
	IsFollow bool   // true-已关注，false-未关注
}
