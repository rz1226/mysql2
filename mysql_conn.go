package mysql2

import (
	"database/sql"
	"fmt"
	"time"
)

type DBConf struct {
	user   string
	pass   string
	host   string
	port   string
	dbName string
	maxCon int
}

func NewDBConf(user, pass, host, port, dbName string, maxCon int) DBConf {
	conf := DBConf{}
	conf.user = user
	conf.pass = pass
	conf.host = host
	conf.port = port
	conf.dbName = dbName
	conf.maxCon = maxCon
	return conf
}

// 获取配置串
func (c *DBConf) Str() string {
	str := fmt.Sprint(c.user, ":", c.pass, "@tcp(", c.host, ":", c.port, ")/", c.dbName, "?charset=utf8")
	return str
}

func Connect(c DBConf) *DB {
	db, err := newDB(c.Str(), c.maxCon)
	db.err = err
	return db
}

// 代表事务
type DBTx struct {
	realtx *sql.Tx
	err    error
}

// 提交事务
func (tx *DBTx) Commit() error {
	return tx.realtx.Commit()
}

// 事务回滚
func (tx *DBTx) Rollback() error {
	return tx.realtx.Rollback()
}

// 代表数据库操作，自带池子
type DB struct {
	realPool *sql.DB
	err      error
}

const MAXLIFETIME = 1200
const MAXOPENCONNS = 1000
const DEFAULTOPENCONNS = 20
const MAXIDLECONNSRATIO = 5

// "gechengzhen:123456@tcp(172.16.1.61:3306)/userdata?charset=utf8"
func newDB(conStr string, maxOpenConns int) (*DB, error) {
	if maxOpenConns <= 0 || maxOpenConns >= MAXOPENCONNS {
		maxOpenConns = DEFAULTOPENCONNS
	}
	pool, err := sql.Open("mysql", conStr)
	if err == nil {
		pool.SetMaxOpenConns(maxOpenConns)
		pool.SetMaxIdleConns(maxOpenConns / MAXIDLECONNSRATIO)
		pool.SetConnMaxLifetime(time.Second * MAXLIFETIME)
		p := &DB{}
		p.realPool = pool
		return p, nil
	}
	return nil, err
}

// 获取*sql.DB
func (p *DB) DB() *sql.DB {
	return p.realPool
}

func (p *DB) Begin() *DBTx {
	realtx, err := p.realPool.Begin()

	t := &DBTx{}
	t.realtx = realtx
	t.err = err
	return t

}
