package micro

import "github.com/micro/go-micro/client"

func NewClient() client.Client {
	cli := client.NewClient(
		client.Wrap(clientTimeWrapper),
	)
	return cli
}