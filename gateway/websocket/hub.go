package websocket

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

type Hub struct {
	clients             map[*Client]bool
	userClient          map[uint]*Client
	register            chan *Client
	unregister          chan *Client
	broadcast           chan BroadcastMessage
	quit                chan struct{}
	subscriber          broker.Subscriber
	numWorkers          uint
	broadcastBufferSize uint
}

type BroadcastMessage struct {
	UserIds []uint
	Message []byte
}

func NewHub(subscriber broker.Subscriber, numWorkers uint, broadcastBufferSize uint) *Hub {
	return &Hub{
		clients:             make(map[*Client]bool),
		userClient:          make(map[uint]*Client),
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		broadcast:           make(chan BroadcastMessage, broadcastBufferSize),
		quit:                make(chan struct{}),
		subscriber:          subscriber,
		numWorkers:          numWorkers,
		broadcastBufferSize: broadcastBufferSize,
	}
}

func (h *Hub) Run() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	go h.dispatcher(ctx)

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.userClient[client.userID] = client
			metrics.UserOnlineCounter.Inc()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userClient, client.userID)
				close(client.send)
				metrics.UserOnlineCounter.Dec()
			}

		case msg := <-h.broadcast:
			for _, uid := range msg.UserIds {
				if client, ok := h.userClient[uid]; ok {
					client.send <- msg.Message
				}
			}

		case <-h.quit:
			<-ctx.Done()
			return
		}
	}
}

func (h *Hub) dispatcher(ctx context.Context) {
	jobQueue, err := h.subscriber.SubscribeTopic(ctx, entity.GameSvcCreatedGameEvent)
	if err != nil {
		metrics.FailedSubscribeTopicCounter.With(prometheus.Labels{"topic": entity.GameSvcCreatedGameEvent}).Inc()
		logger.Fatal(err, "subscribe to topic failed", "topic", entity.GameSvcCreatedGameEvent)
	}

	var wg sync.WaitGroup
	for i := 0; i < int(h.numWorkers); i++ {
		wg.Add(1)
		go h.worker(ctx, jobQueue, &wg, i)
	}

	wg.Wait()
}

func (h *Hub) worker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("WebsocketGateway shutdown initiated", "worker_id", workerId)
			return
		case job, ok := <-jobs:
			if !ok {
				logger.Info("WebsocketGateway job channel closed", "worker_id", workerId)
				return
			}
			h.sendGameCreatedNotificationToClient(ctx, job)
		}
	}
}

func (h *Hub) sendGameCreatedNotificationToClient(ctx context.Context, payload string) {

	gameCreated := protobufencodedecode.DecodeGameSvcCreatedGameEvent(payload)
	type message struct {
		Event  string `json:"event"`
		GameId uint   `json:"game_id"`
	}
	newMsg := message{
		Event:  entity.GameSvcCreatedGameEvent,
		GameId: gameCreated.GameId,
	}
	msgPayload, _ := json.Marshal(newMsg)

	msg := BroadcastMessage{
		UserIds: gameCreated.PlayerIds,
		Message: msgPayload,
	}

	select {
	case h.broadcast <- msg:
	case <-ctx.Done():
		return
	}
}

// Close gracefully shuts down the hub by closing all client connections
// and stopping the hub's run loop.
func (h *Hub) Close() {
	var wg sync.WaitGroup
	for client := range h.clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			c.Close()
		}(client)
	}

	wg.Wait()
	metrics.UserOnlineCounter.Set(0)
	close(h.quit)
}
