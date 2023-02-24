package cli

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"log"
)

var gDB *gorm.DB

const dsn = "root:123456@(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"

func InitDB() error {
	//dsn := viper.GetString("app.DbDsn")
	log.Println("dsn:", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         64,    // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置

	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `userServer` with this option enabled
		},
	})

	if err != nil {
		return err
	}

	gDB = db
	return nil
}

func GetDB() (*gorm.DB, error) {
	if gDB == nil {
		return nil, fmt.Errorf("GetDB err db is nil")
	}
	return gDB, nil
}
