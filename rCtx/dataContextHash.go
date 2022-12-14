package rCtx

import (
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

func (dc *DataCtx) HGet(key string, field string, param interface{}) (err error) {
	//use reflect to check if param is a pointer
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		logger.Lshortfile.Fatal("param must be a pointer")
	}

	cmd := dc.Rds.HGet(dc.Ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) HSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.HSet(dc.Ctx, key, field, bytes)
	return status.Err()
}

func (dc *DataCtx) HExists(key string, field string) (ok bool) {
	cmd := dc.Rds.HExists(dc.Ctx, key, field)
	return cmd.Val()
}

func (dc *DataCtx) HGetAll(key string, mapOut interface{}) (err error) {
	mapElem := reflect.TypeOf(mapOut)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Fatal("mapOut must be a map[string] struct/interface{}")
	}
	cmd := dc.Rds.HGetAll(dc.Ctx, key)
	data, err := cmd.Result()
	if err != nil {
		return err
	}
	//append all data to mapOut
	var result error
	structSupposed := mapElem.Elem()
	for k, v := range data {
		//make a copy of stru , to obj
		obj := reflect.New(structSupposed).Interface()
		if err = msgpack.Unmarshal([]byte(v), &obj); err != nil {
			result = err
		} else {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(obj).Elem())
		}
	}
	return result
}
func (dc *DataCtx) HSetAll(key string, _map interface{}) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Fatal("mapOut must be a map[string] struct/interface{}")
	}
	//HSet each element of _map to redis
	var result error
	pipe := dc.Rds.Pipeline()
	for _, k := range reflect.ValueOf(_map).MapKeys() {
		v := reflect.ValueOf(_map).MapIndex(k)
		if bytes, err := msgpack.Marshal(v.Interface()); err != nil {
			result = err
		} else {
			pipe.HSet(dc.Ctx, key, k.String(), bytes)
		}
	}
	pipe.Exec(dc.Ctx)
	return result
}

func (dc *DataCtx) HGetAllDefault(key string) (param map[string]interface{}, err error) {
	cmd := dc.Rds.HGetAll(dc.Ctx, key)
	data, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	param = make(map[string]interface{})
	//make a copoy of valueStruct
	// unmarshal value of data to the copy
	// store unmarshaled result to param
	for k, v := range data {
		var obj interface{}
		if err = msgpack.Unmarshal([]byte(v), &obj); err == nil {
			param[k] = obj
		}
	}
	return param, nil
}
func (dc *DataCtx) HLen(key string) (length int64) {
	cmd := dc.Rds.HLen(dc.Ctx, key)
	return cmd.Val()
}
func (dc *DataCtx) HDel(key string, field string) (err error) {
	status := dc.Rds.HDel(dc.Ctx, key, field)
	return status.Err()
}
func (dc *DataCtx) HKeys(key string) (fields []string) {
	cmd := dc.Rds.HKeys(dc.Ctx, key)
	return cmd.Val()
}
func (dc *DataCtx) HVals(key string) (values []string) {
	cmd := dc.Rds.HVals(dc.Ctx, key)
	return cmd.Val()
}
func (dc *DataCtx) HIncrBy(key string, field string, increment int64) (err error) {
	status := dc.Rds.HIncrBy(dc.Ctx, key, field, increment)
	return status.Err()
}
func (dc *DataCtx) HIncrByFloat(key string, field string, increment float64) (err error) {
	status := dc.Rds.HIncrByFloat(dc.Ctx, key, field, increment)
	return status.Err()
}
func (dc *DataCtx) HSetNX(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.HSetNX(dc.Ctx, key, field, bytes)
	return status.Err()
}

// golang version of python scan_iter
func (dc *DataCtx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = dc.Rds.Scan(dc.Ctx, cursor, match, count).Result()
	return keys, err
}

// get all keys that match the pattern, and return a map of key->value
func (dc *DataCtx) GetAll(match string, mapOut interface{}) (err error) {
	var (
		keys []string = []string{match}
		val  []byte
	)
	mapElem := reflect.TypeOf(mapOut)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Fatal("mapOut must be a map[string] struct/interface{}")
	}
	if keys, err = dc.Scan(match, 0, 1024*1024*1024); err != nil {
		return err
	}
	var result error
	structSupposed := mapElem.Elem()
	for _, key := range keys {
		if val, result = dc.Rds.Get(dc.Ctx, key).Bytes(); result != nil {
			err = result
			continue
		}
		obj := reflect.New(structSupposed).Interface()
		if msgpack.Unmarshal(val, obj) != nil {
			err = result
			continue
		} else {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(obj).Elem())
		}
	}
	return result
}

// set each key value of _map to redis string type key value
func (dc *DataCtx) SetAll(_map interface{}) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Fatal("mapOut must be a map[string] struct/interface{}")
	}
	//HSet each element of _map to redis
	var result error
	pipe := dc.Rds.Pipeline()
	for _, k := range reflect.ValueOf(_map).MapKeys() {
		v := reflect.ValueOf(_map).MapIndex(k)
		if bytes, err := msgpack.Marshal(v.Interface()); err != nil {
			result = err
		} else {
			pipe.Set(dc.Ctx, k.String(), bytes, -1)
		}
	}
	pipe.Exec(dc.Ctx)
	return result
}
