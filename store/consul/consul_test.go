package consul

import (
	"github.com/stretchr/testify/assert"
	"github.com/wondenge/stor"
	"github.com/wondenge/stor/store"
	"testing"
	"time"
)

var (
	client = "localhost:8500"
)

func makeConsulClient(t *testing.T) store.Store {

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

	kv, err := stor.NewStore(store.CONSUL, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*Consul); !ok {
		t.Fatal("Error registering and initializing consul")
	}
}

func TestConsulStore(t *testing.T) {
	kv := makeConsulClient(t)
	lockKV := makeConsulClient(t)
	ttlKV := makeConsulClient(t)

	stor.RunTestCommon(t, kv)
	stor.RunTestAtomic(t, kv)
	stor.RunTestWatch(t, kv)
	stor.RunTestLock(t, kv)
	stor.RunTestLockTTL(t, kv, lockKV)
	stor.RunTestTTL(t, kv, ttlKV)
	stor.RunCleanup(t, kv)
}

func TestGetActiveSession(t *testing.T) {
	kv := makeConsulClient(t)

	consul := kv.(*Consul)

	key := "foo"
	value := []byte("bar")

	// Put the first key with the Ephemeral flag
	err := kv.Put(key, value, &store.WriteOptions{TTL: 2 * time.Second})
	assert.NoError(t, err)

	// Session should not be empty
	session, err := consul.getActiveSession(key)
	assert.NoError(t, err)
	assert.NotEqual(t, session, "")

	// Delete the key
	err = kv.Delete(key)
	assert.NoError(t, err)

	// Check the session again, it should return nothing
	session, err = consul.getActiveSession(key)
	assert.NoError(t, err)
	assert.Equal(t, session, "")
}