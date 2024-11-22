package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"koriebruh/arc/config"
	"testing"
	"time"
)

var ctx = context.Background()

func TestDataBase(t *testing.T) {

	db := config.GetDataBase()

	sqlDB, err := db.DB()
	assert.Nil(t, err)

	err = sqlDB.Ping()
	assert.Nil(t, err)
}

// DONT FORGE BEFORE START THIS DO "redis-cli monitor" OR after this test "redis-cli get art"
func TestRedis(t *testing.T) {

	rdb := config.GetRedis()

	err := rdb.Set(ctx, "art", "explosion", 15*time.Minute).Err()
	assert.Nil(t, err)

	result, err := rdb.Get(ctx, "art").Result()
	assert.Nil(t, err)
	assert.Equal(t, "explosion", result)
}
