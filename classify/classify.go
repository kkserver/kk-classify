package classify

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"time"
)

type Classify struct {
	Id    int64  `json:"id"`
	Alias string `json:"alias"` //别名
	Pid   int64  `json:"pid"`   //上级分类
	Name  string `json:"name"`  //名称
	Path  string `json:"path"`  // pid/pid
	Tags  string `json:"tags"`  //搜索标签
	Oid   int64  `json:"oid"`
}

type IClassifyApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetClassifyTable() *kk.DBTable
}

type ClassifyApp struct {
	app.App

	DB *app.DBConfig

	Remote *remote.Service

	Classify      *ClassifyService
	ClassifyTable kk.DBTable
}

const twepoch = int64(1424016000000)

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func NewOid() int64 {
	return milliseconds() - twepoch
}

func (C *ClassifyApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *ClassifyApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *ClassifyApp) GetClassifyTable() *kk.DBTable {
	return &C.ClassifyTable
}
