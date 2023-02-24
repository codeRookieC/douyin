package gorms

import (
	"context"
	"fmt"
	"video_server/config"
	"video_server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB

func init() {
	//配置MySQL连接参数
	username := config.C.MysqlConfig.Username //账号
	password := config.C.MysqlConfig.Password //密码
	host := config.C.MysqlConfig.Host         //数据库地址，可以是Ip或者域名
	port := config.C.MysqlConfig.Port         //数据库端口
	Dbname := config.C.MysqlConfig.Db         //数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	fmt.Println(dsn)
	var err error
	//连接数据库
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	MigratorDB()
	//adddata()
}
func adddata() {
	GetDB().Create(&model.User{Id: 1, Username: "user3", Password: "123456"})
	// GetDB().Create(&model.Video{Id: 1, Author_id: 1, Play_url: "url123", Cover_url: "./cccc", Title: "video1"})
	// GetDB().Create(&model.Follow{Id: 1, Follower_id: 2, Follow_id: 1})
	// GetDB().Create(&model.Favorite{Id: 1, Video_id: 1, User_id: 2})
	// GetDB().Create(&model.Comment{Id: 1, Video_id: 1, User_id: 2, Content: "good goood"})
}
func MigratorDB() {
	if (!_db.Migrator().HasTable(&model.User{})) {
		_db.AutoMigrate(&model.User{})
	}
	if (!_db.Migrator().HasTable(&model.Video{})) {
		_db.AutoMigrate(&model.Video{})
	}
	if (!_db.Migrator().HasTable(&model.Follow{})) {
		_db.AutoMigrate(&model.Follow{})
	}
	if (!_db.Migrator().HasTable(&model.Favorite{})) {
		_db.AutoMigrate(&model.Favorite{})
	}
	if (!_db.Migrator().HasTable(&model.Comment{})) {
		_db.AutoMigrate(&model.Comment{})
	}
}

func GetDB() *gorm.DB {
	return _db
}

type GormConn struct {
	db *gorm.DB
	tx *gorm.DB
}

func (g *GormConn) Begin() {
	g.tx = GetDB().Begin()
}

func New() *GormConn {
	return &GormConn{db: GetDB()}
}

func NewTran() *GormConn {
	return &GormConn{db: GetDB(), tx: GetDB()}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) RollBack() {
	g.tx.Rollback()
}

func (g *GormConn) Commit() {
	g.tx.Commit()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}
