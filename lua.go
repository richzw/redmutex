package redmutex

import (
	"github.com/go-redis/redis/v8"
)


var DelScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)


var ExpScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

var PTTLScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then 
		return redis.call("PTTL", KEYS[1]) 
	else 
		return -1
	end
`)