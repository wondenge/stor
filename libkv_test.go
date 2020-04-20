package stor

import (
	"github.com/stretchr/testify/assert"
	"github.com/wondenge/stor/store"
	"testing"
	"time"
)

func TestNewStoreUnsupported(t *testing.T) {
	client := "localhost:9999"

	kv, err := NewStore(
		"unsupported",
		[]string{client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	assert.Error(t, err)
	assert.Nil(t, kv)
	assert.Equal(t, "Backend storage not supported yet, please choose one of ", err.Error())
}
