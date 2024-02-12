package https

import (
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

var ErrBadCommand = errors.New("error bad command")

func (svcCtx *HttpContext) PostHandler() (ret interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn   map[string]interface{} = map[string]interface{}{}
		operation string
	)

	if operation, err = svcCtx.KeyFieldAtJwt(); err != nil {
		return "", err
	}
	if !permission.IsPermitted(permission.Post, svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	//db := &data.Ctx{Ctx: svcCtx.Ctx, Rds: config.Rds, Key: svcCtx.Key}
	db := data.New[interface{}, interface{}](svcCtx.Key)

	//service name is stored in svcCtx.Key
	switch svcCtx.Cmd {
	// all data that appears in the form or body is json format, will be stored in paramIn["JsonPack"]
	// this is used to support 3rd party api
	case "JSAPI":
		//convert query fields to JsonPack. but ignore K field(api name )
		svcCtx.Req.ParseForm()
		if len(svcCtx.Req.Form) > 0 {
			if paramIn["JsonPack"], err = msgpack.Marshal(svcCtx.Req.Form); err != nil {
				return nil, err
			}
		}
		return api.CallByHTTP(svcCtx.Key, paramIn)
	case "API":
		if MsgPack, _ := svcCtx.BodyBytes(); len(MsgPack) > 0 {
			paramIn["MsgPack"] = MsgPack
		}
		svcCtx.MergeJwtField(paramIn)
		return api.CallByHTTP(svcCtx.Key, paramIn)
	case "ZADD":
		var Score float64
		var obj interface{}
		if Score, err = strconv.ParseFloat(svcCtx.Req.FormValue("Score"), 64); err != nil {
			return "false", errors.New("parameter Score shoule be float")
		}
		//unmarshal msgpack
		if MsgPack, _ := svcCtx.BodyBytes(); len(MsgPack) == 0 {
			return "false", errors.New("missing MsgPack content")
		} else if err = msgpack.Unmarshal(MsgPack, &obj); err != nil {
			return "false", err
		}
		if err = db.ZAdd(redis.Z{Score: Score, Member: obj}); err != nil {
			return "false", err
		}
		return "true", nil
	default:
		err = ErrBadCommand
	}

	return ret, err
}
