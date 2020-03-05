package sqlx_plus

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Session struct {
	ReadDB  *sqlx.DB
	WriteDB *sqlx.DB
	Tx      *sqlx.Tx
	builder *DBBuilder
}

func NewSession(builder *DBBuilder) *Session {
	return &Session{
		ReadDB:  builder.getQueryDB(),
		WriteDB: builder.getExecuteDB(),
		builder: builder,
	}
}

func (session *Session) Commit() error {
	if session.Tx != nil {
		return session.Tx.Commit()
	}
	return nil
}

func (session *Session) Rollback() error {
	if session.Tx != nil {
		return session.Tx.Rollback()
	}
	return nil
}

func (session *Session) Search(dest interface{}) *Search {
	return NewSearch(dest, session.builder.DriverName)
}

func (session *Session) Get(dest interface{}, sql string, args ...interface{}) error {
	if args[0] == nil {
		if session.Tx != nil {
			return session.Tx.Get(dest, sql)
		} else {
			return session.ReadDB.Get(dest, sql)
		}
	} else {
		if session.Tx != nil {
			return session.Tx.Get(dest, sql, args...)
		} else {
			return session.ReadDB.Get(dest, sql, args...)
		}
	}
	return nil
}

func (session *Session) Select(dest interface{}, sql string, args ...interface{}) error {
	if args[0] == nil {
		if session.Tx != nil {
			return session.Tx.Select(dest, sql)
		} else {
			return session.ReadDB.Select(dest, sql)
		}
	} else {
		if session.Tx != nil {
			return session.Tx.Select(dest, sql, args...)
		} else {
			return session.ReadDB.Select(dest, sql, args...)
		}
	}
	return nil
}

//通过主键查询
func (session *Session) SelectById(dest interface{}, id ...interface{}) error {
	search := session.Search(dest)
	return session.Get(dest, search.SelectById().GetSql(), id...)
}

//通过主键列表查询
func (session *Session) SelectBatchIds(dest interface{}, id ...interface{}) error {
	search := session.Search(dest)
	query, args, err := sqlx.In(search.SelectBatchIds().GetSql(), id)
	if err != nil {
		return err
	}
	query = session.ReadDB.Rebind(query)
	return session.Select(dest, query, args...)
}

//根据条件查询单条数据
func (session *Session) SelectOne(dest interface{}, where string, args interface{}) error {
	return session.SelectOneField(dest, "", where, args)
}

//根据条件查询单条数据(包含自定义查询列名)
func (session *Session) SelectOneField(dest interface{}, selectField string, where string, args interface{}) error {
	search := session.Search(dest)
	search.Where = where
	search.SelectField = selectField
	return session.Get(dest, search.SelectOne().GetSql(), args)
}

//根据条件查询列表
func (session *Session) SelectList(dest interface{}, where string, args interface{}) error {
	return session.SelectListField(dest, "", where, args)
}

//根据条件查询列表(包含自定义查询列名)
func (session *Session) SelectListField(dest interface{}, selectField string, where string, args interface{}) error {
	search := session.Search(dest)
	search.Where = where
	search.SelectField = selectField
	return session.Select(dest, search.GetSql(), args)
}

//查询条数
func (session *Session) SelectCount(dest interface{}, where string, args interface{}) (int64, error) {
	return session.SelectCountKey(dest, "*", where, args)
}

//根据字段查询条数 count(key)
func (session *Session) SelectCountKey(dest interface{}, key string, where string, args interface{}) (int64, error) {
	search := session.Search(dest)
	search.Where = where
	var count int64
	err := session.Get(&count, search.GetCount(key).GetSql(), args)
	return count, err
}

func (session *Session) SelectPage(dest interface{}, where string, pageIndex int64, pageSize int64, args interface{}) error {
	return session.SelectPageField(dest, "", where, pageIndex, pageSize, args)
}

func (session *Session) SelectPageField(dest interface{}, selectField string, where string, pageIndex int64, pageSize int64, args interface{}) error {
	search := session.Search(dest)
	search.Where = where
	search.SelectField = selectField
	search.PageIndex = pageIndex
	search.PageSize = pageSize
	return session.Select(dest, search.GetSql(), args)
}

func (session *Session) DeleteById(dest interface{}, id interface{}) (sql.Result, error) {
	if session.Tx != nil {
		return session.Tx.Exec(NewDelete().Model(dest).DeleteById().GetSql(), id)
	} else {
		return session.WriteDB.Exec(NewDelete().Model(dest).DeleteById().GetSql(), id)
	}
	return nil, nil
}

func (session *Session) DeleteBatchIds(dest interface{}, id ...interface{}) (sql.Result, error) {
	query, args, err := sqlx.In(NewDelete().Model(dest).DeleteBatchIds().GetSql(), id)
	if err != nil {
		return nil, err
	}
	query = session.WriteDB.Rebind(query)
	if session.Tx != nil {
		return session.Tx.Exec(query, args...)
	} else {
		return session.WriteDB.Exec(query, args...)
	}
	return nil, nil
}

func (session *Session) DeleteAll(dest interface{}) (sql.Result, error) {
	if session.Tx != nil {
		return session.Tx.Exec(NewDelete().Model(dest).GetSql())
	} else {
		return session.WriteDB.Exec(NewDelete().Model(dest).GetSql())
	}
	return nil, nil
}

func (session *Session) UpdateById(dest interface{}) (sql.Result, error) {
	sql, data := NewUpdate().Model(dest).UpdateById().GetSql()
	if session.Tx != nil {
		return session.Tx.Exec(sql, data...)
	} else {
		return session.WriteDB.Exec(sql, data...)
	}
	return nil, nil
}

func (session *Session) Insert(dest interface{}) (sql.Result, error) {
	sql, data := NewInsert().Model(dest).GetSql()
	if session.Tx != nil {
		return session.Tx.Exec(sql, data...)
	} else {
		return session.WriteDB.Exec(sql, data...)
	}
	return nil, nil
}
