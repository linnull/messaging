package errors

import (
	"github.com/micro/go-micro/errors"

	comm "lib/comm/proto"
)

const SYSTEM_ERR_ID = "server_system_err"
const LOGIC_ERR_ID = "server_logic_err"

func NewSystemErr() error {
	return &errors.Error{
		Id: SYSTEM_ERR_ID,
	}
}

func NewLogicErr(logicRet comm.LOGIC_RET) error {
	return &errors.Error{
		Id:   LOGIC_ERR_ID,
		Code: int32(logicRet),
	}
}

var (
	SYSTEM_FAIL             = NewSystemErr()
	LOGIC_SYS_ERROR         = NewLogicErr(comm.LOGIC_RET_LOGIC_SYS_ERROR)
	LOGIC_LOGIN_TOKEN_ERROR = NewLogicErr(comm.LOGIC_RET_LOGIC_LOGIN_TOKEN_ERROR)
)
