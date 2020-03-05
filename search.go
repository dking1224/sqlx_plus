package sqlx_plus

import (
	"fmt"
	"strings"
)

type Search struct {
	TableName   string   `desc:"表名"`
	SelectField string   `desc:"需要查询的字段"`
	Filter      []Filter `desc:"条件,Field Condition Value 拼接,多个使用and连接"`
	Where       string   `desc:"如果where存在则Filter失效"`
	Order       string   `desc:"排序,如 order by id"`
	PageIndex   int64    `desc:"页数"`
	PageSize    int64    `desc:"分页单页大小,如果此数据大于0则进行分页处理"`
	PK          string   `desc:"主键"`
	Limit       string
	DriverName  string
}

type Filter struct {
	Field     string
	Condition string
	Value     string
}

func NewSearch(dest interface{}, driverName string) *Search {
	return &Search{
		TableName:  GetTableName(dest),
		PK:         GetPK(dest),
		DriverName: driverName,
	}
}

func NewTableSearch(tableName string, PK string, driverName string) *Search {
	return &Search{
		TableName:  tableName,
		PK:         PK,
		DriverName: driverName,
	}
}

func NewInitSearch(driverName string) *Search {
	return &Search{DriverName: driverName}
}

func (search *Search) Model(dest interface{}) *Search {
	search.TableName = GetTableName(dest)
	search.PK = GetPK(dest)
	return search
}

func (search *Search) Field(field string) *Search {
	search.SelectField = field
	return search
}

func (search *Search) FilterString(where string) *Search {
	search.Where = where
	return search
}

func (search *Search) Filters(filters []Filter) *Search {
	search.Filter = filters
	return search
}

func (search *Search) OrderBy(order string) *Search {
	search.Order = order
	return search
}

func (search *Search) Page(pageIndex int64, pageSize int64) *Search {
	search.PageIndex = pageIndex
	search.PageSize = pageSize
	return search
}

func (search *Search) Select(SelectField string, where string, order string) *Search {
	search.SelectField = SelectField
	search.Where = where
	search.Order = order
	return search
}

func (search *Search) SelectFilter(SelectField string, filter []Filter, order string) *Search {
	search.SelectField = SelectField
	search.Filter = filter
	search.Order = order
	return search
}

func (search *Search) SelectById() *Search {
	search.Where = search.PK + "=? "
	return search
}

func (search *Search) SelectBatchIds() *Search {
	search.Where = search.PK + " in (?)"
	return search
}

func (search *Search) GetCount(key string) *Search {
	search.SelectField = fmt.Sprintf(" count(%s) ", key)
	return search
}

//支持mysql sqlite3 pg oracle
func (search *Search) SelectOne() *Search {
	if search.DriverName == "mysql" || search.DriverName == "sqlite3" || search.DriverName == "postgres" {
		search.Limit = "limit 1"
	} else {
		search.Where = fmt.Sprintf("1=1 and %s and %s", search.Where, "rownum=1")
	}
	return search
}

func (search *Search) GetSql() string {
	builder := strings.Builder{}
	builder.WriteString("select ")
	if search.SelectField != "" {
		builder.WriteString(search.SelectField)
	} else {
		builder.WriteString("*")
	}
	builder.WriteString(" from ")
	builder.WriteString(search.TableName)
	if search.Where != "" {
		builder.WriteString(" where ")
		builder.WriteString(search.Where)
	} else {
		if len(search.Filter) != 0 {
			if checkFilter(search.Filter) {
				builder.WriteString(" where ")
				for i := 0; i < len(search.Filter); i++ {
					if i == 0 {
						builder.WriteString("(")
						builder.WriteString(filterString(search.Filter[i]))
						builder.WriteString(")")
					} else {
						builder.WriteString(" and (")
						builder.WriteString(filterString(search.Filter[i]))
						builder.WriteString(")")
					}
				}
			}
		}
	}
	if search.Order != "" {
		builder.WriteString(" ")
		builder.WriteString(search.Order)
		builder.WriteString(" ")
	}
	if search.Limit != "" {
		builder.WriteString(" ")
		builder.WriteString(search.Limit)
		builder.WriteString(" ")
	}
	if search.PageSize > 0 {
		if search.DriverName == "mysql" || search.DriverName == "sqlite3" || search.DriverName == "postgres" {
			builder.WriteString(fmt.Sprintf(" limit %d,%d", search.PageIndex, search.PageSize))
		} else {
			sql := search.pageOracleSql(builder.String())
			Log.Debug("select sql:", sql)
			return sql
		}
	}
	Log.Debug("select sql:", builder.String())
	return builder.String()
}

func filterString(filter Filter) string {
	return filter.Field + filter.Condition + filter.Value
}

func checkFilter(filter []Filter) bool {
	for _, item := range filter {
		if item.Field == "" || item.Condition == "" || item.Value == "" {
			return false
		}
	}
	return true
}

func (search *Search) pageOracleSql(sql string) string {
	startSql := strings.ReplaceAll(sql, "select",
		"SELECT * FROM ( SELECT A.*, ROWNUM RNFROM (SELECT *")
	endSql := fmt.Sprintf(
		") A WHERE ROWNUM <= (%d) ) WHERE RN > (%d)",
		search.PageIndex*search.PageSize, (search.PageIndex-1)*search.PageSize)
	return startSql + endSql
}
