package goretryhandler

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Backoff interface {
	Next(retry int) time.Duration
}

type constantBackOff struct {
	backOffInterval       int64
	maximumJitterInterval int64
}

// Next returns the next time for retrying operation with constant strategy
func (c *constantBackOff) Next(retry int) time.Duration {
	return (time.Duration(c.backOffInterval) * time.Millisecond) +
		(time.Duration(rand.Int63n(c.maximumJitterInterval+1)) * time.Millisecond)
}

// NewConstantBackOff returns an instance of constantBackOff
func NewConstantBackOff(backOffInterval, maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}
	return &constantBackOff{
		backOffInterval:       int64(backOffInterval / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

type exponentialBackOff struct {
	exponentFactor        float64
	initialTimeout        float64
	maxTimeout            float64
	maximumJitterInterval int64
}

func (eb *exponentialBackOff) Next(retry int) time.Duration {
	if retry < 0 {
		retry = 0
	}
	return time.Duration(math.Min(eb.initialTimeout*math.Pow(eb.exponentFactor, float64(retry)),
		eb.maxTimeout)+
		float64(rand.Int63n(eb.maximumJitterInterval+1))) * time.Millisecond
}

func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}
	return &exponentialBackOff{
		exponentFactor:        exponentFactor,
		initialTimeout:        float64(initialTimeout / time.Millisecond),
		maxTimeout:            float64(maxTimeout / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}
