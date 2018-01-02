package rater

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Rater struct {
	redisc *redis.Client
}

func New(client *redis.Client) *Rater {
	return &Rater{
		redisc: client,
	}
}

func (r *Rater) AllowMinute(key string, max int64) (int64, int64, bool) {
	return r.AllowN(key, max, 60, 1)
}

func (r *Rater) AllowHour(key string, max int64) (int64, int64, bool) {
	return r.AllowN(key, max, 60*60, 1)
}

func (r *Rater) AllowDay(key string, max int64) (int64, int64, bool) {
	return r.AllowN(key, max, 60*60*24, 1)
}

// 不使用GET + INCR + EXPIRE : 并发且key已存在时,实际执行的次数可能多于max
func (r *Rater) AllowN(key string, max, d, step int64) (remaining, reset int64, isAllow bool) {
	t := time.Now().Unix()
	slot := t / d

	reset = (slot+1)*d - t // next change for new windwo(s)

	count, err := r.incr(newKeyBySlot(key, slot), d, step)
	if err == nil {
		isAllow = count <= max
		if isAllow {
			remaining = max - count
		}
	}

	return
}

func (r *Rater) incr(key string, d, step int64) (int64, error) {
	var incr *redis.IntCmd
	_, err := r.redisc.Pipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.IncrBy(key, step)
		pipe.Expire(key, time.Duration(d)*time.Second)

		return nil
	})

	count, err := incr.Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func newKeyBySlot(key string, slot int64) string {
	return fmt.Sprintf("%s-%d", key, slot)
}

func (r *Rater) Incr(key string, d, step int64) (int64, error) {
	return r.incr(key, d, step)
}

func (r *Rater) AllowMax(key string, max int64) (int64, bool) {
	isAllow := false

	num, err := r.redisc.Get(key).Int64()
	if err == nil {
		isAllow = num <= max
	} else if err == redis.Nil {
		isAllow = true
	}

	return num, isAllow
}
