package goretryhandler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff(t *testing.T) {
	backoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 0*time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, backoff.Next(0))
	assert.Equal(t, 200*time.Millisecond, backoff.Next(1))
	assert.Equal(t, 400*time.Millisecond, backoff.Next(2))
	assert.Equal(t, 800*time.Millisecond, backoff.Next(3))
	assert.Equal(t, 1*time.Second, backoff.Next(4))
	assert.Equal(t, 1*time.Second, backoff.Next(5))
}
