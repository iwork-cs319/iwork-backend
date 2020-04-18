package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-api/utils"
	"log"
	"strings"
	"time"
)

type RedisCache struct {
	pool *redis.Pool
}

func NewRedis(redisUrl string) (Cache, error) {
	redisCache := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
	cache := &RedisCache{pool: redisCache}
	err := cache.Ping()
	return cache, err
}

func (c *RedisCache) Ping() error {
	conn := c.pool.Get()
	// Send PING command to Redis
	pong, err := conn.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	_, err = redis.String(pong, err)
	if err != nil {
		return err
	}

	return nil
}

const WorkspaceLockKey = "wLockKey"
const LockExpiry = 20 * time.Minute

func (c *RedisCache) CreateWorkspaceLock(workspaceId string, start, end time.Time) error {
	conn := c.pool.Get()
	defer conn.Close()

	hashKey := fmt.Sprintf("%s:%s", WorkspaceLockKey, workspaceId)
	key := fmt.Sprintf("%d", time.Now().Unix())
	val := fmt.Sprintf("%d-%d", start.Unix(), end.Unix())

	if _, err := conn.Do("HSET", hashKey, key, val); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) CheckWorkspaceLock(workspaceId string, start, end time.Time) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	hashKey := fmt.Sprintf("%s:%s", WorkspaceLockKey, workspaceId)
	values, err := redis.Values(conn.Do("HGETALL", hashKey))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	expireThreshold := time.Now().UTC()
	for i := 0; i < len(values); i += 2 {
		createdVal, _ := redis.String(values[i], nil)
		lockRangeVal, _ := redis.String(values[i+1], nil)
		lockRange := strings.Split(lockRangeVal, "-")
		createdTime, err := utils.TimeStampToTime(createdVal)
		if err != nil {
			log.Printf("CheckWorkpsaceLock: failed to parse timeStamp k=%s", createdVal)
			continue
		}
		if createdTime.Add(LockExpiry).Before(expireThreshold) {
			// Expired Key
			// Delete
			_, err = conn.Do("HDEL", hashKey, createdVal)
			if err != nil {
				log.Printf("CheckWorkspaceLock: failed to delete expired key, hashKey=%s, key=%s", hashKey, createdVal)
			}
			continue
		}
		start1, err := utils.TimeStampToTime(lockRange[0])
		if err != nil {
			log.Printf("CheckWorkpsaceLock: failed to parse timeStamp v=%s", lockRangeVal)
			continue
		}
		end1, err := utils.TimeStampToTime(lockRange[1])
		if err != nil {
			log.Printf("CheckWorkpsaceLock: failed to parse timeStamp v=%s", lockRangeVal)
			continue
		}
		if timeIntersects(start1, end1, start, end) { // Lock is set for this time period
			return true, nil
		}
	}
	return false, nil
}

func timeIntersects(start1 time.Time, end1 time.Time, start2 time.Time, end2 time.Time) bool {
	return (start1.Equal(start2) && end1.Equal(end2)) ||
		(start2.Before(start1) && end2.After(end1)) ||
		(start1.Before(start2) && end1.After(end2)) ||
		(start2.Before(start1) && end2.After(start1)) ||
		(start2.Before(end1) && end2.After(end1))
}
