package test

import (
	"context"
	"fmt"
	"testing"

	"logic/login/client"
)

func TestEcho(t *testing.T) {
	cli := client.NewLoginClient()
	ctx := context.Background()
	rsp, err := cli.Echo(ctx, "echo test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}
