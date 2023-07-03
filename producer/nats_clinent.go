package producer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

const (
	NATSStreamingURL = "stanserver:4222"
	ClusterID        = "test-cluster"
	ClientID         = "test-publisher"
	Channel          = "testch"
)

type NatsStream struct {
	client stan.Conn
}

func NewNatsStream() *NatsStream {
	sc, err := stan.Connect(
		ClusterID, ClientID,
		stan.NatsURL(NATSStreamingURL),
	)
	if err != nil {
		log.Fatalf("error connection nats, %s", err.Error())
	}
	log.Println("nats connection successful")
	return &NatsStream{client: sc}
}

func (n *NatsStream) RunNatsSteaming(client *PostgresClient, cache *MemcacheClient) {
	_, err := n.client.Subscribe(
		Channel, func(m *stan.Msg) {
			var order Order
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = client.InsertOrder(order)
			if err != nil {
				log.Fatal(err)
				return
			}
			cache.Set(order.OrderUID, m.Data, time.Millisecond)
		}, stan.StartAtTimeDelta(time.Minute*10))
	if err != nil {
		log.Fatal(err)
	}
}
