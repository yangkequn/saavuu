package saavuu

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"
)

func (scvCtx *ServiceContext) RedisGet(key string, param interface{}) (err error) {
	cmd := Config.rds.Get(scvCtx.ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, &param)
}
func (scvCtx *ServiceContext) RedisSet(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := Config.rds.Set(scvCtx.ctx, key, bytes, expiration)
	return status.Err()
}
func (scvCtx *ServiceContext) RedisHGet(key string, field string, param interface{}) (err error) {
	cmd := Config.rds.HGet(scvCtx.ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, &param)
}

func (scvCtx *ServiceContext) RedisHSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := Config.rds.HSet(scvCtx.ctx, key, field, bytes)
	return status.Err()
}
func RedisDo(c context.Context, rds *redis.Client, ServiceKey string, paramIn map[string]interface{}) (result []byte, err error) {
	var (
		b       []byte
		BackTo  string = strconv.FormatInt(rand.Int63(), 36)
		results []string
	)
	paramIn["BackTo"] = BackTo

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	ppl := rds.Pipeline()
	ppl.RPush(c, ServiceKey, b)
	//长期不执行的任务，抛弃
	ppl.Expire(c, ServiceKey, time.Second*60)
	if _, err := ppl.Exec(c); err != nil {
		return nil, err
	}
	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = rds.BLPop(c, time.Second*20, BackTo).Result(); err != nil {
		return nil, err
	}
	return []byte(results[1]), nil
}
