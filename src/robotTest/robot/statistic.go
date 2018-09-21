package robot

import (
	"fmt"
)

type ActionStatistic struct {
	Name     string
	Run      int64
	Success  int64
	Fail     int64
	UseTime  []int64
	SumTime  int64
	AvgTime  int64
	SendByte int64
	RecByte  int64
}

func NewActionStatistic(Name string) *ActionStatistic {
	return &ActionStatistic{Name: Name, UseTime: []int64{}}
}

func (as *ActionStatistic) addActStatisticSuccess(useTime int64, sBytesLen, rBytesLen int64) {
	as.UseTime = append(as.UseTime, useTime)
	as.Run++
	as.Success++
	as.SumTime += useTime

	as.AvgTime = as.SumTime / as.Success
	as.RecByte += rBytesLen
	as.SendByte += sBytesLen
}

func (as *ActionStatistic) addActStatisticFail() {
	as.Run++
	as.Fail++
}

func (as *ActionStatistic) PrintStatus() {
	fmt.Printf("[%v] Run:%d, Success:%d, Fail:%d, SumTime:%vms, AvgTime:%vms, SendByte:%d, RecByte:%d\n", as.Name, as.Run, as.Success, as.Fail, float64(as.SumTime)/1000000, float64(as.AvgTime)/1000000, as.SendByte, as.RecByte)
}

func (as *ActionStatistic) PrintStatusWithNano() {
	fmt.Printf("[%v] Run:%d, Success:%d, Fail:%d, SumTime:%v, AvgTime:%v, SendByte:%d, RecByte:%d\n", as.Name, as.Run, as.Success, as.Fail, as.SumTime, as.AvgTime, as.SendByte, as.RecByte)
}
