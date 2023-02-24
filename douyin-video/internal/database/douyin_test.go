package database

import (
	"testing"
	"video_server/internal/database/gorms"
	"video_server/internal/model"
)

func addusers() {
	// var users = []model.User{{Username: "user1", Password: "123456"},
	// 	{Username: "user2", Password: "123456"},
	// 	{Username: "user3", Password: "123456"}}
	// gorms.GetDB().Create(&users)
	gorms.GetDB().Create(&model.User{Username: "user2", Password: "123456"})
}

func TestAddUsers(t *testing.T) {
	addusers()
}
