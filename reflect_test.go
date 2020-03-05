package sqlx_plus

import (
	"fmt"
	"testing"
)

type User struct {
	ID   uint64 `db:"id" orm:"pk" table:"user"`
	Name uint64 `db:"name"`
}

func TestGetPK(t *testing.T) {
	fmt.Println(GetPK(&User{ID: 1, Name: 12}))
	fmt.Println(GetTableName(&User{}))
}

func TestGetPK2(t *testing.T) {
	users := make([]User, 0)
	fmt.Println(GetPK(&users))
	fmt.Println(GetTableName(&users))
}

func TestGetUpdateCol(t *testing.T) {
	fmt.Println(GetUpdateCol(&User{ID: 1, Name: 12}))
}
