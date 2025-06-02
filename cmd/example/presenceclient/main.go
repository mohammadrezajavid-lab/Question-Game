package main

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
)

func main() {

	client := presenceclient.NewClient(presenceclient.Config{
		Host:    "127.0.0.1",
		Port:    8086,
		Network: "tcp",
	})
	res, gErr := client.GetPresence(context.Background(), presenceparam.NewGetPresenceRequest([]uint{1, 2, 3, 4}))
	if gErr != nil {
		panic(gErr)
	}

	for _, item := range res.Items {
		fmt.Println("item", item.UserId, "  ", item.Timestamp)
	}
}
