package redmutex

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
	"github.com/go-redis/redis/v8"
)

type RedMutex struct {
	client redis.Client
}

func NewRedMutex(client redis.Client) *RedMutex {
	return &RedMutex{client: client}
}

func (rm *RedMutex) NewMutex(key string, options...Options) *Mutex {
	m := &Mutex{
		key:     key,
		ttl:     0,
		backOff: NoRetry(),
		valFunc: DefaultValFunc,

		client: rm.client,
	}

	for _, o := range options {
		o.Apply(m)
	}

	return m
}

type Options interface {
	Apply(*Mutex)
}

type OptionFunc func(*Mutex)

func (f OptionFunc) Apply(mutex *Mutex) {
	f(mutex)
}

func WithTTL(ttl time.Duration) Options {
	return OptionFunc(func (m *Mutex) {
		m.ttl = ttl
	})
}

func WithBackOff(bf Retry) Options {
	return OptionFunc(func (m *Mutex) {
		m.backOff = bf
	})
}

func WithValFunc(vf ValFunc) Options {
	return OptionFunc(func (m *Mutex) {
		m.valFunc = vf
	})
}


type ValFunc func() (string, error)

func DefaultValFunc() (string, error) {
	b := make([]byte, 16)

	// TODO: compare with rand.Read(b)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
