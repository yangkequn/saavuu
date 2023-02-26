package https

import (
	"errors"
	"strings"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

var ErrEmptyKeyOrField = errors.New("empty key or field")

func (svcCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		result    map[string]interface{} = map[string]interface{}{}
		bytes     []byte
		operation string = strings.ToLower(svcCtx.Cmd)
	)

	if strings.Contains(svcCtx.Field, "@") {
		if err := svcCtx.ParseJwtToken(); err != nil {
			return "false", err
		}
		if !permission.IsPutPermitted(svcCtx.Key, operation, &svcCtx.Field, svcCtx.jwtToken) {
			return "false", errors.New("permission denied")
		}
	} else {
		if !permission.IsPutPermitted(svcCtx.Key, operation, nil, nil) {
			return "false", errors.New("permission denied")
		}
	}

	switch svcCtx.Cmd {
	case "HSET":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", ErrEmptyKeyOrField
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := config.DataRds.HSet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	case "RPUSH":
		//error if empty Key or Field
		if svcCtx.Key == "" {
			return "false", ErrEmptyKeyOrField
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := config.DataRds.RPush(svcCtx.Ctx, svcCtx.Key, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	}
	return result, nil
}
