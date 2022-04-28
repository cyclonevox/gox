package database

import `xorm.io/xorm`

type txFunc func(session *xorm.Session) (int64, error)

type Transaction interface {
	Do(txFunc) (int64, error)
	Rollback() error
	Commit() error
}

type tx struct {
	session *xorm.Session
}

// NewTx 创建一个数据库事务接口
func NewTx(engine *xorm.Engine) (Transaction, error) {
	session := engine.NewSession()
	if err := session.Begin(); err != nil {
		return nil, err
	}

	return &tx{session: session}, nil
}

// Do 执行事务操作，出现错误会进行Rollback
func (t *tx) Do(f txFunc) (int64, error) {
	var (
		affected int64
		err      error
	)

	if affected, err = f(t.session); err != nil {
		_ = t.Rollback()

		return 0, err
	}

	return affected, nil
}

// Rollback 回滚事务
func (t *tx) Rollback() error {
	return t.session.Rollback()
}

// Commit 提交事务
func (t *tx) Commit() error {
	return t.session.Commit()
}
