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

/**
 * [频率限制](https://github.com/SFLAQiu/web-develop/blob/master/Redis%E5%AE%9E%E6%88%98%E4%B9%8B%E9%99%90%E5%88%B6%E6%93%8D%E4%BD%9C%E9%A2%91%E7%8E%87.md)
 * d 时间范围X秒内
 * step 自增步长
 * max 限制操作数Y次
 * limitTTL 超出封印时间Z
 *
 * 场景:
 * 1. 留言功能限制，30秒内只能评论10次，超出次数不让能再评论 LimitWithTTL(key,30,1,10,0)
 * 1. 点赞功能限制，10秒内只能点赞10次，超出次数后不能再点赞，并封印1个小时 LimitWithTTL(key,30,1,10,60*60)
 * 1. 上传记录功能，需要限制一天只能上传 100次，超出次数不让能再上传 LimitWithTTL(key,24*60*60,1,100,nextDay-now)
 */
func (r *Rater) LimitWithTTL(key string, d, step, max, limitTTL int64) bool {
	current, err := r.redisc.Get(key).Int64()
	if err != nil && err != redis.Nil {
		return false
	}

	if current >= max {
		return false
	}

	newCurrent := r.redisc.IncrBy(key, step).Val()
	if newCurrent <= current { // incr failed
		return false
	}

	if newCurrent == 1 {
		r.redisc.Expire(key, time.Duration(d)*time.Second)
	}
	// 超出后根据需要重新设置过期失效时间
	if newCurrent == max && limitTTL > 0 {
		r.redisc.Expire(key, time.Duration(limitTTL)*time.Second)
	}

	return newCurrent <= max
}
