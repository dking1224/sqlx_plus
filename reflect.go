package sqlx_plus

import (
	"reflect"
	"strings"
)

func GetTableName(dest interface{}) string {
	tag := getFieldNameByTag(dest, "table")
	if tag != "" && tag.Get("table") != "" {
		return tag.Get("table")
	}
	return ""
}

func GetPK(dest interface{}) string {
	tag := getFieldNameByTag(dest, "orm")
	if tag != "" && tag.Get("orm") == "pk" {
		return tag.Get("db")
	}
	return ""
}

func GetUpdateCol(dest interface{}) (string, []interface{}) {
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil
	}
	builder := strings.Builder{}
	data := make([]interface{}, 0, t.NumField())
	var pkData interface{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("db") != "" {
			builder.WriteString(tag.Get("db"))
			builder.WriteString("=?,")
			data = append(data, reflect.ValueOf(dest).Elem().FieldByName(t.Field(i).Name).Interface())
		}
		if tag.Get("orm") != "" && tag.Get("orm") == "pk" {
			pkData = reflect.ValueOf(dest).Elem().FieldByName(t.Field(i).Name).Interface()
		}
	}
	data = append(data, pkData)
	return strings.TrimRight(builder.String(), ","), data
}

func GetInsertCol(dest interface{}) (string, []interface{}) {
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil
	}
	builder := strings.Builder{}
	builder.WriteString("(")
	builder1 := strings.Builder{}
	builder1.WriteString("(")
	data := make([]interface{}, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("db") != "" {
			builder.WriteString(tag.Get("db"))
			builder.WriteString(",")
			builder1.WriteString("?,")
			data = append(data, reflect.ValueOf(dest).Elem().FieldByName(t.Field(i).Name).Interface())
		}
	}
	str1 := strings.TrimRight(builder.String(), ",") + ")"
	str2 := strings.TrimRight(builder1.String(), ",") + ")"
	return str1 + " values " + str2, data
}

func getFieldNameByTag(dest interface{}, key string) reflect.StructTag {
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get(key) != "" {
			return tag
		}
	}
	return ""
}
