package test

import (
	"context"
	"fmt"
	"testing"

	"template/client"
)

func TestEcho(t *testing.T) {
	cli := client.NewTemplateClient()
	ctx := context.Background()
	rsp, err := cli.Echo(ctx, "echo test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}
