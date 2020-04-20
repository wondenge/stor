package redis

import (
	"github.com/stretchr/testify/assert"
	"github.com/wondenge/stor"
	"github.com/wondenge/stor/store"
	"testing"
)

var (
	client = "localhost:6379"
)

func makeRedisClient(t *testing.T) store.Store {
	kv, err := newRedis([]string{client}, "", nil)
	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	// NOTE: please turn on redis's notification
	// before you using watch/watchtree/lock related features
	kv.client.ConfigSet("notify-keyspace-events", "KA")

	return kv
}

func TestRegister(t *testing.T) {
	Register()

	kv, err := stor.NewStore(store.REDIS, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*Redis); !ok {
		t.Fatal("Error registering and initializing redis")
	}
}

func TestRedisStore(t *testing.T) {
	kv := makeRedisClient(t)
	lockTTL := makeRedisClient(t)
	kvTTL := makeRedisClient(t)

	stor.RunTestCommon(t, kv)
	stor.RunTestAtomic(t, kv)
	stor.RunTestWatch(t, kv)
	stor.RunTestLock(t, kv)
	stor.RunTestLockTTL(t, kv, lockTTL)
	stor.RunTestTTL(t, kv, kvTTL)
	stor.RunCleanup(t, kv)
}
