package sqlx_plus

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var (
	t        *template.Template
	allFile  []string
	PrintLog = true
)

func getAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			getAllFile(pathname + fi.Name() + fileSeparator())
		} else {
			if path.Ext(fi.Name()) == ".tpl" {
				filePath := filepath.Join(pathname, fi.Name())
				allFile = append(allFile, filePath)
			}
		}
	}
	return err
}

func InitTpl(path string, methods map[string]interface{}) error {
	if path == "" {
		return errors.New("dir path is nil")
	}
	t = template.New("template sql")
	funcMap := template.FuncMap{}
	if methods != nil {
		for k, v := range methods {
			funcMap[k] = v
		}
	}
	t.Funcs(funcMap)
	if err := getAllFile(path); err != nil {
		return err
	}
	if len(allFile) > 0 {
		t.ParseFiles(allFile...)
	}
	return nil
}

func fileSeparator() string {
	osType := runtime.GOOS
	if osType == "windows" {
		return "\\"
	}
	return "/"
}

type SqlResult struct {
	Sql   string
	Error error
}

//sql拼接
func SqlStr(name string, data interface{}) SqlResult {
	buf := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		return SqlResult{Sql: "", Error: err}
	}
	sql := strings.TrimSpace(strings.Replace(buf.String(), "\r\n", " ", -1))
	if PrintLog {
		log.Print(name, "[", sql, "]")
		log.Print(name, data)
	}
	return SqlResult{Sql: sql, Error: nil}
}
