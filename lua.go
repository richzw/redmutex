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

var GetScript = redis.NewScript(`
	if (redis.call('exists', KEYS[1]) == 0) then
	  redis.call('hincrby', KEYS[1], ARGV[2], 1);
	  redis.call('pexpire', KEYS[1], ARGV[1]);
	  return nil;
	end;
	
	if (redis.call('hexists', KEYS[1], ARGV[2]) == 1) then
	  redis.call('hincrby', KEYS[1], ARGV[2], 1);
	  redis.call('pexpire', KEYS[1], ARGV[1]);
	  return nil;
	end;
	
	return redis.call('pttl', KEYS[1]);
`)

var DelScript = redis.NewScript(`
	if (redis.call('hexists', KEYS[1], ARGV[3]) == 0) then
	 return nil;
	end;
	
	local counter = redis.call('hincrby', KEYS[1], ARGV[3], -1); 
	
	if (counter > 0) then
	 redis.call('pexpire', KEYS[1], ARGV[2]);
	return 0;
	else
	 redis.call('del', KEYS[1]);
	 redis.call('publish', KEYS[2], ARGV[1]);
	return 1;
	end;
	return nil;
`)

var RefreshScript = redis.NewScript(`
	if (redis.call('hexists', KEYS[1], ARGV[2]) == 1) then 
	 redis.call('pexpire', KEYS[1], ARGV[1]); 
	  return 1; 
	end;  
	return 0;
`)
