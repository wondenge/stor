package etcd

import (
	"github.com/stretchr/testify/assert"
	"github.com/wondenge/stor"
	"github.com/wondenge/stor/store"
	"testing"
	"time"
)

var (
	client = "localhost:4001"
)

func makeEtcdV3Client(t *testing.T) store.Store {
	kv, err := New(
		[]string{client},
		&store.Config{
			ConnectionTimeout: 3 * time.Second,
			Username:          "test",
			Password:          "very-secure",
		},
	)

	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	return kv
}

func TestRegister(t *testing.T) {
	Register()

	kv, err := stor.NewStore(store.ETCDV3, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*EtcdV3); !ok {
		t.Fatal("Error registering and initializing etcd with v3 client")
	}
}

func TestEtcdV3Store(t *testing.T) {
	kv := makeEtcdV3Client(t)
	lockKV := makeEtcdV3Client(t)
	ttlKV := makeEtcdV3Client(t)

	stor.RunTestCommon(t, kv)
	stor.RunTestAtomic(t, kv)
	stor.RunTestWatch(t, kv)
	stor.RunTestLock(t, kv)
	stor.RunTestLockTTL(t, kv, lockKV)
	stor.RunTestListLock(t, kv)
	stor.RunTestTTL(t, kv, ttlKV)
	stor.RunCleanup(t, kv)
}