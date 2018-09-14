package micro

import (
	"context"
	"log"
	"time"

	"github.com/micro/go-micro/client"
)

type TimeWrapper struct {
	client.Client
}

func (l *TimeWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	starTime := time.Now()
	err := l.Client.Call(ctx, req, rsp, opts...)
	endTime := time.Now()
	useTime := endTime.UnixNano() - starTime.UnixNano()
	if err != nil {
		log.Printf("[error] %s useTime: %d ms err: %v ", req.Method(), useTime/1000000, err)
	} else {
		log.Printf("[success] %s useTime: %d ms ", req.Method(), useTime/1000000)
	}
	return err
}

func timeWrapper(c client.Client) client.Client {
	return &TimeWrapper{c}
}