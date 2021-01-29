package redmutex

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrNotLock = errors.New("redmutex: not lock")
	ErrNotHeld = errors.New("redmutex: not held")
)

type Mutex struct {
	key  string
	val  string
	ttl  time.Duration

	backOff Retry
	valFunc ValFunc

	client redis.Client
}

func (m *Mutex) Lock(ctx context.Context) error {
	val, err := m.valFunc()
	if err != nil {
		return err
	}

	deadlineCtx, cancel := context.WithDeadline(ctx, time.Now().Add(m.ttl))
	defer cancel()

	var t *time.Timer
	for {
		ok, err := m.client.SetNX(deadlineCtx, m.key, val, m.ttl).Result()
		if err != nil {
			return err
		} else if ok {
			return nil
		}

		backOff := m.backOff.NextBackOff()
		if backOff < 1 {
			return ErrNotLock
		}

		if t == nil {
			t = time.NewTimer(backOff)
		} else {
			t.Reset(backOff)
		}

		select {
		case <- deadlineCtx.Done():
			return ErrNotLock
		case <- t.C:
		}
	}

}

func (m *Mutex) UnLock(ctx context.Context) error {
	res, err := DelScript.Run(ctx, m.client, []string{m.key}, m.val).Result()
	if err == redis.Nil {
		return ErrNotHeld
	} else if err != nil {
		return err
	}

	if r, ok := res.(int64); !ok || r != 1 {
		return ErrNotHeld
	}
	return nil
}

// Refresh the lock through extent expire time
func (m *Mutex) Refresh(ctx context.Context, ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	res, err := ExpScript.Run(ctx, m.client, []string{m.key}, m.val, ttlVal).Result()
	if err == redis.Nil {
		return ErrNotHeld
	} else if err != nil {
		return err
	}

	if r, ok := res.(int64); !ok || r != 1 {
		return ErrNotHeld
	}
	return nil
}

// Get the remaining ttl
func (m *Mutex) TTL(ctx context.Context) (time.Duration, error) {
	res, err := PTTLScript.Run(ctx, m.client, []string{m.key}, m.val).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if r, ok := res.(int64); ok || r > 0 {
		return time.Duration(r) * time.Millisecond, nil
	}
	return 0, nil
}

