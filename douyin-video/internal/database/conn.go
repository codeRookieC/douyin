package database

type DbConn interface { //必要的操作 开启事务 回滚和提交
	Begin()
	RollBack()
	Commit()
}
