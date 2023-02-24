package dao

import (
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DbConn
}

func (t TransactionImpl) Action(f func(conn database.DbConn) error) error {
	//开启事务
	t.conn.Begin()
	err := f(t.conn)
	if err != nil { //操作错误 回滚
		t.conn.RollBack()
		return err
	}
	//操作成功 提交
	t.conn.Commit()
	return nil
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTran(),
	}
}