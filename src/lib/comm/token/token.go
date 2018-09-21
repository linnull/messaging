package token

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	TOKEN_TYPE_LOGIN = iota
)

func GenLoginToken(uid uint64) string {
	return genToken(uid, TOKEN_TYPE_LOGIN)
}

func genToken(uid uint64, tokenType int) string {
	// todo:生成token
	return fmt.Sprintf("%d;%d", uid, tokenType)
}

func CheckLoginToken(token string) (uid uint64, err error) {
	// todo:检查token
	res := strings.Split(token, ";")
	if len(res) != 2 {
		return 0, errors.New("token is wrong")
	}
	tokenType, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, err
	}
	if tokenType != TOKEN_TYPE_LOGIN {
		return 0, errors.New("token type wrong")
	}
	return strconv.ParseUint(res[0], 10, 64)
}
