package micro

import (
	"context"
	"reflect"
	"strings"

	"lib/comm/errors"
	comm "lib/comm/proto"

	"github.com/golang/protobuf/proto"
	microErrors "github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
)

type BaseMessage interface {
	GetBase() *comm.BaseResponse
}

func ServerCommWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := fn(ctx, req, rsp)
		return commResponse(rsp.(proto.Message), err)
	}
}

func commResponse(rsp proto.Message, err error) error {
	if err == nil {
		setBaseResponse(rsp.(proto.Message), comm.RET_SYSTEM_SUCCESS, comm.LOGIC_RET_LOGIC_SUCCESS)
		return nil
	}

	// 将micro client.serverError 转为micro.Error并覆盖err
	errVal := reflect.ValueOf(err)
	if errVal.Kind() == reflect.String {
		err = microErrors.Parse(errVal.String())
	}

	switch err := err.(type) {
	case *microErrors.Error:
		if strings.Compare(err.Id, errors.SYSTEM_ERR_ID) == 0 {
			setBaseResponse(rsp, comm.RET_SYSTEM_FAIL, comm.LOGIC_RET(err.Code))
		} else if strings.Compare(err.Id, errors.LOGIC_ERR_ID) == 0 {
			setBaseResponse(rsp, comm.RET_LOGIC_FAIL, comm.LOGIC_RET(err.Code))
		} else {
			setBaseResponse(rsp, comm.RET_SYSTEM_FAIL, comm.LOGIC_RET(err.Code))
		}
	default:
		setBaseResponse(rsp, comm.RET_SYSTEM_FAIL, comm.LOGIC_RET_LOGIC_SYS_ERROR)
	}
	return nil
}

func setBaseResponse(rsp proto.Message, ret comm.RET, logicRet comm.LOGIC_RET) {
	switch rsp := rsp.(type) {
	case BaseMessage:
		base := rsp.GetBase()
		if base == nil {
			base = &comm.BaseResponse{}
			resp := reflect.ValueOf(rsp).Elem()
			resp.FieldByName("Base").Set(reflect.ValueOf(base))
		}
		base.Ret = ret
		base.LogicRet = logicRet
	default:
		// do nothing
	}
}
