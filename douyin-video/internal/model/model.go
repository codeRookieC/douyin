package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       int64
	Username string
	Password string
}

type Follow struct {
	gorm.Model
	Id          int64
	Follower_id int64
	Follow_id   int64
}

type Video struct {
	gorm.Model
	Id        int64
	Author_id int64
	Play_url  string
	Cover_url string
	Title     string
}

type Favorite struct {
	gorm.Model
	Id       int64
	Video_id int64
	User_id  int64
}

type Comment struct {
	gorm.Model
	Id       int64
	Video_id int64
	User_id  int64
	Content  string
}
