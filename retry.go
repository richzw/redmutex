package redmutex

import (
	"math/rand"
	"time"
)

type Retry interface {
	NextBackOff() time.Duration
}

type Linear struct {
	d time.Duration
}

func (l *Linear) NextBackOff() time.Duration {
	return l.d
}

func NoRetry() Retry  {
	return &Linear{d: 0}
}

func LinearRetry(d time.Duration) Retry {
	return &Linear{d: d}
}

type Exponential struct {
	baseDelay time.Duration
	maxDelay  time.Duration

	curRetries int64
}

func ExponentialRetry(baseDelay, maxDelay time.Duration) Retry {
	return &Exponential{baseDelay: baseDelay, maxDelay: maxDelay}
}

func (bc *Exponential) NextBackOff() time.Duration {
	if  bc.curRetries == 0 {
		return bc.baseDelay
	}

	expBase := int64(2)
	backoff := int64(bc.baseDelay)
	backoff += rand.Int63n(expBase*bc.curRetries) + expBase*bc.curRetries
	if backoff > int64(bc.maxDelay) {
		return bc.maxDelay
	}
	bc.curRetries++

	return time.Duration(backoff)
}
