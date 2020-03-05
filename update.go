package sqlx_plus

import (
	"strings"
)

type Delete struct {
	TableName string
	PK        string
	Where     string
}

func NewDelete() *Delete {
	return &Delete{}
}

func NewTableDelete(tableName string, pk string) *Delete {
	return &Delete{
		TableName: tableName,
		PK:        pk,
	}
}

func (delete *Delete) Model(dest interface{}) *Delete {
	delete.TableName = GetTableName(dest)
	delete.PK = GetPK(dest)
	return delete
}

func (delete *Delete) FilterString(where string) *Delete {
	delete.Where = where
	return delete
}

func (delete *Delete) DeleteById() *Delete {
	delete.Where = delete.PK + "=?"
	return delete
}

func (delete *Delete) DeleteBatchIds() *Delete {
	delete.Where = delete.PK + " in(?) "
	return delete
}

func (delete *Delete) GetSql() string {
	builder := strings.Builder{}
	builder.WriteString("delete from ")
	builder.WriteString(delete.TableName)
	if delete.Where != "" {
		builder.WriteString(" where ")
		builder.WriteString(delete.Where)
	}
	Log.Debug("delete sql:", builder.String())
	return builder.String()
}

type Update struct {
	TableName string
	PK        string
	Where     string
	Dest      interface{}
}

func NewUpdate() *Update {
	return &Update{}
}

func (update *Update) Model(dest interface{}) *Update {
	update.PK = GetPK(dest)
	update.TableName = GetTableName(dest)
	update.Dest = dest
	return update
}

func (update *Update) UpdateById() *Update {
	update.Where = update.PK + "=?"
	return update
}

//只支持单表updateById
func (update *Update) GetSql() (string, []interface{}) {
	builder := strings.Builder{}
	builder.WriteString("update ")
	builder.WriteString(update.TableName)
	builder.WriteString(" set ")
	setSql, data := GetUpdateCol(update.Dest)
	builder.WriteString(setSql)
	if update.Where != "" {
		builder.WriteString(" where ")
		builder.WriteString(update.Where)
	}
	Log.Debug("update sql:", builder.String())
	return builder.String(), data
}

type Insert struct {
	TableName string
	Dest      interface{}
}

func NewInsert() *Insert {
	return &Insert{}
}

func (insert *Insert) Model(dest interface{}) *Insert {
	insert.Dest = dest
	insert.TableName = GetTableName(dest)
	return insert
}

func (insert *Insert) GetSql() (string, []interface{}) {
	builder := strings.Builder{}
	builder.WriteString("insert into ")
	builder.WriteString(insert.TableName)
	setSql, data := GetInsertCol(insert.Dest)
	builder.WriteString(setSql)
	Log.Debug("insert sql:", builder.String())
	return builder.String(), data
}
