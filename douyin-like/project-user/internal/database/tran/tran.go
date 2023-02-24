package tran

import "test.com/project-user/internal/database"

// Transaction  事务的操作一定与数据库有关 那么我们这里要注入数据库的连接 gorm.db 这里采取接口的方式注入 解耦 之后可能会使用别的连接方式
type Transaction interface {
	Action(func(conn database.DbConn) error) error
}