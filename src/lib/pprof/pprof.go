package pprof

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func Init(name string, version string) {
	go func() {
		http.HandleFunc("/debug/pprof/info", func(resp http.ResponseWriter, req *http.Request) {
			resp.Write([]byte(fmt.Sprintf("{\"ServiceName\":\"%v\",\"ServiceVersion\":\"%v\"}", name, version)))
		})

		reStartTimes := 0
		for {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			addr := fmt.Sprintf("0.0.0.0:%v", 20101+r.Intn(99))
			log.Printf("profiling http [%v]  srv name is [%v]\n", addr, name)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				log.Printf("Get an error %s\n", err.Error())
			}
			reStartTimes++
			log.Printf("Restart Times is %v ,and now  get some error and try to restart after 1 sec\n", reStartTimes)
			time.Sleep(time.Second)
		}
	}()
}
