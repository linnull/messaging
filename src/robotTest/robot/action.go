package robot

import "fmt"

func Conn(uid uint64) (err error) {
	rb := NewRobot(uid)
	err = rb.Conn(ip, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rb.DisConn()
	fmt.Println(rb.Status)
	rb.ActionStatistics[ROBOT_ACTION_CONN].PrintStatus()
	return
}

func Login(uid uint64) (err error) {
	rb := NewRobot(uid)
	err = rb.Conn(ip, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rb.DisConn()
	rb.ActionStatistics[ROBOT_ACTION_CONN].PrintStatus()

	err = rb.Login(uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rb.Status)
	rb.ActionStatistics[ROBOT_ACTION_LOGIN].PrintStatus()
	return
}
