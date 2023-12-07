package log

import (
	"github.com/stretchr/testify/assert"
	config "golinkcut/internal/config"
	"testing"
)

func TestNew(t *testing.T) {
	l := New()
	if l == nil {
		t.Error("logger is nil, zap Logger expected")
	}
}

func TestNewForTest(t *testing.T) {
	logger, entries := NewForTest()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	entries.TakeAll()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
}

func TestNewWithConfig(t *testing.T) {
	config := config.Config{"debug": true}
	l := NewWithConfig(config)
	if l == nil {
		t.Error("logger is nil, zap Logger expected")
	}
}
