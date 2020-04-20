package zookeeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/wondenge/stor"
	"github.com/wondenge/stor/store"
	"testing"
	"time"
)

var (
	client = "localhost:2181"
)

func makeZkClient(t *testing.T) store.Store {
	kv, err := New(
		[]string{client},
		&store.Config{
			ConnectionTimeout: 3 * time.Second,
		},
	)

	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	return kv
}

func TestRegister(t *testing.T) {
	Register()

	kv, err := stor.NewStore(store.ZK, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*Zookeeper); !ok {
		t.Fatal("Error registering and initializing zookeeper")
	}
}

func TestZkStore(t *testing.T) {
	kv := makeZkClient(t)
	ttlKV := makeZkClient(t)

	stor.RunTestCommon(t, kv)
	stor.RunTestAtomic(t, kv)
	stor.RunTestWatch(t, kv)
	stor.RunTestLock(t, kv)
	stor.RunTestTTL(t, kv, ttlKV)
	stor.RunCleanup(t, kv)
}
