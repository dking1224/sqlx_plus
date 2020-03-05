package sqlx_plus

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var Log = logrus.New()

func init() {
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Log.Level = logrus.DebugLevel
}

type Config struct {
	DataSourceName  string
	ConnMaxLifetime time.Duration
	MaxIdle         int
	MaxOpen         int
}

type DBConfig struct {
	Master []Config
	Slave  []Config
}

type DBBuilder struct {
	DriverName string
	Config     DBConfig
	master     []*sqlx.DB
	slave      []*sqlx.DB
	TplDir     string
	Method     map[string]interface{}
}

func New(config DBConfig, driverName string, tplDir string, method map[string]interface{}) (*DBBuilder, error) {
	if len(config.Master) == 0 {
		return nil, errors.New("配置信息为空")
	}
	builder := new(DBBuilder)
	builder.DriverName = driverName
	builder.createMaster(config)
	builder.createSlave(config)
	if tplDir != "" {
		InitTpl(tplDir, method)
	}
	return builder, nil
}

func (builder *DBBuilder) SetLogger(logger *logrus.Logger) {
	Log = logger
}

func (builder *DBBuilder) createMaster(config DBConfig) error {
	listDB, err := builder.createDBList(config.Master)
	if err != nil {
		return err
	}
	builder.master = listDB
	return nil
}

func (builder *DBBuilder) createSlave(config DBConfig) error {
	listDB, err := builder.createDBList(config.Slave)
	if err != nil {
		return err
	}
	builder.slave = listDB
	return nil
}

func (builder *DBBuilder) createDBList(config []Config) ([]*sqlx.DB, error) {
	size := len(config)
	if size == 0 {
		return nil, nil
	}
	listDB := make([]*sqlx.DB, 0, size)
	for _, v := range config {
		db, err := builder.creatDB(v)
		if err != nil {
			return nil, err
		}
		listDB = append(listDB, db)
	}
	return listDB, nil
}

func (builder *DBBuilder) creatDB(config Config) (db *sqlx.DB, err error) {
	db, err = sqlx.Open(builder.DriverName, config.DataSourceName)
	if err = db.Ping(); err != nil {
		return
	}
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetMaxOpenConns(config.MaxOpen)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	return
}

func (builder *DBBuilder) getQueryDB() *sqlx.DB {
	if len(builder.slave) == 0 {
		return builder.getExecuteDB()
	}
	var randInt = GetRandomInt(len(builder.master))
	return builder.slave[randInt]
	return nil
}

func (builder *DBBuilder) getExecuteDB() *sqlx.DB {
	if len(builder.master) == 0 {
		return nil
	}
	var randInt = GetRandomInt(len(builder.master))
	return builder.master[randInt]
}

func (builder *DBBuilder) NewOrm() *Session {
	return NewSession(builder)
}

func (builder *DBBuilder) NewTxOrm() *Session {
	session := NewSession(builder)
	session.Tx = session.WriteDB.MustBegin()
	return session
}

func GetRandomInt(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}
