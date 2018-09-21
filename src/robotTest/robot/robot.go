package robot

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"lib/comm/method"
	"lib/comm/proto"
	"lib/comm/token"
	"lib/packet"

	"github.com/golang/protobuf/proto"
)

const (
	//运行状态
	ROBOT_STATUS_ERROR = "status_error"
	ROBOT_STATUS_INIT  = "status_init"
	ROBOT_STATUS_STOP  = "status_stop"
	ROBOT_STATUS_CONN  = "status_conn"
)

const (
	ROBOT_ACTION_CONN  = "conn"
	ROBOT_ACTION_LOGIN = "login"
)

type MsgResponse interface {
	GetBase() *lib_comm_proto.BaseResponse
}

type Robot struct {
	conn       net.Conn
	reqResLock sync.RWMutex
	reqResCh   map[uint32]chan []byte
	count      uint32
	needCmd    bool

	Uid              uint64
	Token string
	Status           string
	ActionStatistics map[string]*ActionStatistic
}

func NewRobot(uid uint64) *Robot {
	rb := &Robot{}
	rb.reqResLock = sync.RWMutex{}
	rb.reqResCh = make(map[uint32]chan []byte, 1024)

	rb.Uid = uid
	rb.Token = token.GenLoginToken(rb.Uid)
	rb.Status = ROBOT_STATUS_INIT
	rb.ActionStatistics = make(map[string]*ActionStatistic)
	rb.ActionStatistics[ROBOT_ACTION_CONN] = NewActionStatistic(ROBOT_ACTION_CONN)
	rb.ActionStatistics[ROBOT_ACTION_LOGIN] = NewActionStatistic(ROBOT_ACTION_LOGIN)
	return rb
}

func (rb *Robot) OnCommand() {
	rb.needCmd = true
}

func (rb *Robot) Conn(ip, port string) (err error) {
	startTime := time.Now().UnixNano()
	st := rb.ActionStatistics[ROBOT_ACTION_CONN]
	defer func() {
		if err == nil {
			useTime := time.Now().UnixNano() - startTime
			st.addActStatisticSuccess(useTime, 0, 0)
		} else {
			st.addActStatisticFail()
			rb.Status = ROBOT_STATUS_ERROR
		}
	}()

	if rb.Status != ROBOT_STATUS_INIT {
		return errors.New("robot is not init")
	}
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 3*time.Second)
	if err != nil {
		return err
	}
	rb.conn = conn
	rb.Status = ROBOT_STATUS_CONN
	go rb.getPacket()
	return nil
}

func (rb *Robot) DisConn() (err error) {
	defer func() {
		rb.Status = ROBOT_STATUS_STOP
	}()
	if rb.conn != nil {
		return rb.conn.Close()
	}
	return nil
}

func (rb *Robot) getPacket() {
	for {
		pkg, err := packet.GetPacket(rb.conn)
		if err != nil {
			return
		}
		if pkg.ProtocolType == packet.PROTOCOL_TYPE_COMMAND {
			if rb.needCmd {
				go rb.onCommandPacket(pkg)
			}
			continue
		}
		respMsg, err := packet.GetResponseMessageFromPacket(pkg)
		if err != nil || respMsg == nil {
			fmt.Printf("[ERROR] GetResponseMessageFromPacket fail , %v", err)
			continue
		}

		rb.reqResLock.RLock()
		if _, ok := rb.reqResCh[respMsg.Sequence]; ok {
			rb.reqResCh[respMsg.Sequence] <- respMsg.PbMsgBytes
		}
		rb.reqResLock.RUnlock()
	}
}

func (rb *Robot) onCommandPacket(pkg *packet.Packet) (err error) {
	// todo:处理推送消息
	return
}

func (rb *Robot) sendPacket(methodId uint32, pb interface{}, resp proto.Message) (sByteLen, rBytesLen int64, err error) {
	rb.reqResLock.Lock()
	ch := make(chan []byte)
	var seq uint32
	for {
		rb.count++
		seq = rb.count
		if _, ok := rb.reqResCh[seq]; ok {
			continue
		}
		break
	}
	rb.reqResCh[seq] = ch
	rb.reqResLock.Unlock()

	defer func() {
		rb.reqResLock.Lock()
		delete(rb.reqResCh, seq)
		rb.reqResLock.Unlock()
	}()

	var pbBytes []byte
	if pb == nil {
		pbBytes = []byte("hello")
	} else {
		pbBytes, err = proto.Marshal(pb.(proto.Message))
	}

	if err != nil {
		return 0, 0, err
	}

	pkgBytes, err := packet.MarshalRequestPacket(rb.Uid, methodId, seq, pbBytes)
	if err != nil {
		return 0, 0, err
	}
	_, err = rb.conn.Write(pkgBytes)
	if err != nil {
		close(ch)
		return 0, 0, errors.New("write fail")
	}

	timeout := time.NewTicker(time.Second * 10)
	var msgBytes []byte
	select {
	case msgBytes = <-ch:
		if resp != nil {
			err = proto.Unmarshal(msgBytes, resp)
		}

	case <-timeout.C:
		err = errors.New("req timeout")
	}
	return int64(len(pbBytes)), int64(len(msgBytes)), err
}

func (rb *Robot) checkResp(err error, name string, resp MsgResponse) error {
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New(fmt.Sprintf("%v resp is nil", name))
	}
	base := resp.GetBase()
	if base == nil {
		return errors.New(fmt.Sprintf("%v resp.base is nil", name))
	}

	if base.Ret != lib_comm_proto.RET_SYSTEM_SUCCESS {
		return errors.New(fmt.Sprintf("%v fail, ret:%v  detail_code:%v", name, base.Ret, base.LogicRet))
	}
	return nil
}

func (rb *Robot) Login(uid uint64) (err error) {
	startTime := time.Now().UnixNano()
	st := rb.ActionStatistics[ROBOT_ACTION_LOGIN]
	var sByteLen, rBytesLen int64
	defer func() {
		if err == nil {
			useTime := time.Now().UnixNano() - startTime
			st.addActStatisticSuccess(useTime, sByteLen, rBytesLen)
		} else {
			st.addActStatisticFail()
			rb.Status = ROBOT_STATUS_ERROR
		}
	}()

	req := &lib_comm_proto.LoginRequest{Base:&lib_comm_proto.BaseRequest{}, Token: rb.Token}
	resp := &lib_comm_proto.LoginResponse{}
	sByteLen, rBytesLen, err = rb.sendPacket(method.MethodId_login, req, resp)
	err = rb.checkResp(err, "login", resp)
	if err != nil {
		return err
	}
	return nil
}
