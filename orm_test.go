package sqlx_plus

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

type TUser struct {
	Id         int64        `db:"id" json:"id" orm:"pk" table:"t_user"`
	Name       string       `db:"name" json:"name"`
	Age        int          `db:"age" json:"age"`
	CreateTime sql.NullTime `db:"time" json:"create_time"`
}

var builder *DBBuilder

func init() {
	builder, _ = New(DBConfig{
		Master: []Config{
			{
				DataSourceName: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true",
				MaxIdle:        30,
				MaxOpen:        10},
		},
	}, "mysql", "", nil)
	//设置logger 默认使用logrus logger
	/*logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.Level = logrus.DebugLevel
	builder.SetLogger(logger)*/
}

func TestSession_SelectById(t *testing.T) {
	session := builder.NewOrm()
	var user TUser
	err := session.SelectById(&user, 1)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(user)
}

func TestSession_SelectBatchIds(t *testing.T) {
	session := builder.NewOrm()
	user := make([]TUser, 0)
	err := session.SelectBatchIds(&user, 1, 2)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(user)
}

func TestSession_SelectCount(t *testing.T) {
	session := builder.NewOrm()
	count, err := session.SelectCount(&TUser{}, "", nil)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(count)
	count, err = session.SelectCountKey(&TUser{}, "id", "age=?", 12)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(count)
}

func TestSession_SelectList(t *testing.T) {
	session := builder.NewTxOrm()
	defer session.Rollback()
	user := make([]TUser, 0)
	err := session.SelectList(&user, "age=?", 12)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(user)
	session.Commit()
}

func TestSession_SelectPage(t *testing.T) {
	session := builder.NewOrm()
	user := make([]TUser, 0)
	err := session.SelectPage(&user, "age=?", 3, 3, 12)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(user)
}

func TestSession_Get(t *testing.T) {
	session := builder.NewOrm()
	var user TUser
	search := NewTableSearch("t_user", "id", "mysql")
	err := session.Get(&user, search.Select("", "id=?", "").GetSql(), 2)
	if err != nil {
		Log.Error(err)
	}
	Log.Info(user)
}

func TestSession_UpdateById(t *testing.T) {
	user := &TUser{
		Id:   1,
		Name: "test12",
		Age:  29,
		CreateTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	session := builder.NewOrm()
	_, err := session.UpdateById(user)
	if err != nil {
		Log.Error(err)
	}
}

func TestSession_Insert(t *testing.T) {
	user := &TUser{
		Id:   9,
		Name: "test19",
		Age:  29,
		CreateTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	session := builder.NewOrm()
	_, err := session.Insert(user)
	if err != nil {
		Log.Error(err)
	}
}

func TestSession_DeleteById(t *testing.T) {
	session := builder.NewOrm()
	_, err := session.DeleteById(&TUser{}, 6)
	if err != nil {
		Log.Error(err)
	}
}

func TestSession_DeleteBatchIds(t *testing.T) {
	session := builder.NewOrm()
	_, err := session.DeleteBatchIds(&TUser{}, 7, 8, 9)
	if err != nil {
		Log.Error(err)
	}
}
