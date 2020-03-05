package sqlx_plus

import (
	"testing"
)

func TestGetRandomInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Log.Info(GetRandomInt(5))
	}
}
