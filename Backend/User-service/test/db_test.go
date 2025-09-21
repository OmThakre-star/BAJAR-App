package test

import (
	"testing"
	"user-service/db"
)

func TestInitDB(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("InitDB panicked: %v", r)
		}
	}()
	db.InitDB("invalid-dsn") // Should not panic
}
