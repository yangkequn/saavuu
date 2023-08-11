package data

import (
	"context"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type Ctx[k comparable, v any] struct {
	Ctx context.Context
	Rds *redis.Client
	Key string
}

func NewStruct[k comparable, v any]() *Ctx[k, v] {
	var Key string
	_type := reflect.TypeOf((*v)(nil))
	//take name of type v as key
	for _type.Kind() == reflect.Ptr || _type.Kind() == reflect.Slice {
		_type = _type.Elem()
	}
	Key = _type.Name()
	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	return &Ctx[k, v]{Ctx: context.Background(), Rds: config.Rds, Key: Key}
}

func New[k comparable, v any](Key string) *Ctx[k, v] {
	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	return &Ctx[k, v]{Ctx: context.Background(), Rds: config.Rds, Key: Key}
}
func (ctx *Ctx[k, v]) WithContext(c context.Context) *Ctx[k, v] {
	return &Ctx[k, v]{Ctx: c, Rds: ctx.Rds, Key: ctx.Key}
}

func (db *Ctx[k, v]) Time() (tm time.Time, err error) {
	cmd := db.Rds.Time(db.Ctx)
	return cmd.Result()
}
