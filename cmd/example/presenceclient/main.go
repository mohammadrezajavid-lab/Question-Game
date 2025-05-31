package main

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := presenceclient.NewClient(conn)
	res, gErr := client.GetPresence(context.Background(), presenceparam.NewGetPresenceRequest([]uint{1, 2, 3, 4}))
	if gErr != nil {
		panic(gErr)
	}

	for _, item := range res.Items {
		fmt.Println("item", item.UserId, "  ", item.Timestamp)
	}
}
